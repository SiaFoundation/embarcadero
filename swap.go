package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"reflect"
	"strings"

	"go.sia.tech/siad/crypto"
	"go.sia.tech/siad/node/api/client"
	"go.sia.tech/siad/types"
)

var siad *client.Client

var minerFee = types.SiacoinPrecision.Mul64(5)

// A SwapTransaction is a transaction that swaps Siacoin for Siafunds between
// two parties.
type SwapTransaction struct {
	SiacoinInputs  []types.SiacoinInput         `json:"siacoinInputs"`
	SiafundInputs  []types.SiafundInput         `json:"siafundInputs"`
	SiacoinOutputs []types.SiacoinOutput        `json:"siacoinOutputs"`
	SiafundOutputs []types.SiafundOutput        `json:"siafundOutputs"`
	Signatures     []types.TransactionSignature `json:"signatures"`
}

// A SwapSummary details the amount of Siacoins and Siafunds received and spent
// by a party during a swap.
type SwapSummary struct {
	ReceiveSF bool           `json:"receiveSF"`
	ReceiveSC bool           `json:"receiveSC"`
	AmountSF  types.Currency `json:"amountSF"`
	AmountSC  types.Currency `json:"amountSC"`
	MinerFee  types.Currency `json:"minerFee"`
	Stage     int            `json:"stage"`
}

// Transaction converts the swap transaction into a full transaction.
func (swap *SwapTransaction) Transaction() types.Transaction {
	return types.Transaction{
		SiacoinInputs:         swap.SiacoinInputs,
		SiafundInputs:         swap.SiafundInputs,
		SiacoinOutputs:        swap.SiacoinOutputs,
		SiafundOutputs:        swap.SiafundOutputs,
		MinerFees:             []types.Currency{minerFee},
		TransactionSignatures: swap.Signatures,
	}
}

// ParseCurrency parses a suffixed Siacoin or Siafund string into a currency
// value.
func ParseCurrency(amount string) types.Currency {
	amount = strings.TrimSpace(amount)
	if strings.HasSuffix(amount, "SF") || strings.HasSuffix(amount, "H") {
		i, ok := new(big.Int).SetString(strings.TrimRight(amount, "SFH"), 10)
		if !ok {
			log.Fatal("Invalid currency")
		}
		return types.NewCurrency(i)
	}

	units := []string{"pS", "nS", "uS", "mS", "SC", "KS", "MS", "GS", "TS"}
	for i, unit := range units {
		if strings.HasSuffix(amount, unit) {
			value := strings.TrimSpace(strings.TrimSuffix(amount, unit))
			r, ok := new(big.Rat).SetString(value)
			if !ok {
				log.Fatal("Invalid currency")
			}
			exp := 24 + 3*(int64(i)-4)
			mag := new(big.Int).Exp(big.NewInt(10), big.NewInt(exp), nil)
			r.Mul(r, new(big.Rat).SetInt(mag))
			if !r.IsInt() {
				log.Fatal("Currency must be an integer")
			}
			return types.NewCurrency(r.Num())
		}
	}
	log.Fatal("Must specify units of currency")
	return types.Currency{}
}

func encodeJSON(w io.Writer, v interface{}) {
	// encode nil slices as [] instead of null
	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice && val.Len() == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func addSC(swap *SwapTransaction, amount types.Currency) error {
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return fmt.Errorf("failed to get unspent outputs: %w", err)
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiacoinOutput {
			wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
			if err != nil {
				return fmt.Errorf("failed to get address %v unlock conditions: %w", u.UnlockHash, err)
			}
			swap.SiacoinInputs = append(swap.SiacoinInputs, types.SiacoinInput{
				ParentID:         types.SiacoinOutputID(u.ID),
				UnlockConditions: wucg.UnlockConditions,
			})
			inputSum = inputSum.Add(u.Value)
			if inputSum.Cmp(amount) >= 0 {
				break
			}
		}
	}
	if inputSum.Cmp(amount) < 0 {
		return errors.New("insufficient funds")
	}
	// add a change output, if necessary
	if !inputSum.Equals(amount) {
		wag, err := siad.WalletAddressGet()
		if err != nil {
			return fmt.Errorf("failed to get change output address: %w", err)
		}
		swap.SiacoinOutputs = append(swap.SiacoinOutputs, types.SiacoinOutput{
			UnlockHash: wag.Address,
			Value:      inputSum.Sub(amount),
		})
	}
	return nil
}

func addSF(swap *SwapTransaction, amount types.Currency) error {
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return fmt.Errorf("failed to get wallet unspent outputs: %w", err)
	}
	wag, err := siad.WalletAddressGet()
	if err != nil {
		return fmt.Errorf("failed to get wallet address: %w", err)
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiafundOutput {
			wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
			if err != nil {
				return fmt.Errorf("failed to get address %v unlock conditions: %w", u.UnlockHash, err)
			}
			swap.SiafundInputs = append(swap.SiafundInputs, types.SiafundInput{
				ParentID:         types.SiafundOutputID(u.ID),
				UnlockConditions: wucg.UnlockConditions,
				ClaimUnlockHash:  wag.Address,
			})
			inputSum = inputSum.Add(u.Value)
			if inputSum.Cmp(amount) >= 0 {
				break
			}
		}
	}
	if inputSum.Cmp(amount) < 0 {
		return errors.New("insufficient funds")
	}
	// add a change output, if necessary
	if !inputSum.Equals(amount) {
		swap.SiafundOutputs = append(swap.SiafundOutputs, types.SiafundOutput{
			UnlockHash: wag.Address,
			Value:      inputSum.Sub(amount),
		})
	}
	return nil
}

func signSC(swap *SwapTransaction) error {
	var toSign []crypto.Hash
	for _, sci := range swap.SiacoinInputs {
		swap.Signatures = append(swap.Signatures, types.TransactionSignature{
			ParentID:       crypto.Hash(sci.ParentID),
			PublicKeyIndex: 0,
			CoveredFields:  types.FullCoveredFields,
		})
		toSign = append(toSign, crypto.Hash(sci.ParentID))
	}
	txn := swap.Transaction()
	wspr, err := siad.WalletSignPost(txn, toSign)
	swap.Signatures = wspr.Transaction.TransactionSignatures
	return err
}

func signSF(swap *SwapTransaction) error {
	var toSign []crypto.Hash
	for _, sfi := range swap.SiafundInputs {
		swap.Signatures = append(swap.Signatures, types.TransactionSignature{
			ParentID:       crypto.Hash(sfi.ParentID),
			PublicKeyIndex: 0,
			CoveredFields:  types.FullCoveredFields,
		})
		toSign = append(toSign, crypto.Hash(sfi.ParentID))
	}
	txn := swap.Transaction()
	wspr, err := siad.WalletSignPost(txn, toSign)
	swap.Signatures = wspr.Transaction.TransactionSignatures
	return err
}

// CreateSwap creates a new SwapTransaction swapping the input amount for the
// output amount.
func CreateSwap(inputAmount, outputAmount types.Currency, offeringSF bool) (SwapTransaction, error) {
	wag, err := siad.WalletAddressGet()
	if err != nil {
		return SwapTransaction{}, err
	}
	var swap SwapTransaction
	if offeringSF {
		swap.SiacoinOutputs = append(swap.SiacoinOutputs, types.SiacoinOutput{
			Value:      outputAmount,
			UnlockHash: wag.Address,
		})
		swap.SiafundOutputs = append(swap.SiafundOutputs, types.SiafundOutput{
			Value:      inputAmount,
			UnlockHash: types.UnlockHash{}, // to be filled in by counterparty
		})
		if err := addSF(&swap, inputAmount); err != nil {
			return SwapTransaction{}, fmt.Errorf("failed to add siafunds to swap transaction: %w", err)
		}
	} else {
		swap.SiacoinOutputs = append(swap.SiacoinOutputs, types.SiacoinOutput{
			Value:      inputAmount,
			UnlockHash: types.UnlockHash{}, // to be filled in by counterparty
		})
		swap.SiafundOutputs = append(swap.SiafundOutputs, types.SiafundOutput{
			Value:      outputAmount,
			UnlockHash: wag.Address,
		})
		// the party that contributes SC is responsible for paying the miner fee
		if err := addSC(&swap, inputAmount.Add(minerFee)); err != nil {
			return SwapTransaction{}, fmt.Errorf("failed to add siacoins to swap transaction: %w", err)
		}
	}
	return swap, nil
}

// CheckAccept checks that the counterparty's swap transaction is valid
func CheckAccept(swap SwapTransaction) error {
	if len(swap.SiacoinInputs) == 0 && len(swap.SiafundInputs) == 0 {
		return errors.New("transaction has no inputs")
	} else if len(swap.SiacoinInputs) > 0 && len(swap.SiafundInputs) > 0 {
		return errors.New("only one set of inputs should be provided")
	} else if len(swap.SiacoinOutputs) == 0 && len(swap.SiafundOutputs) == 0 {
		return errors.New("transaction has no outputs")
	} else if swap.SiacoinOutputs[0].UnlockHash == (types.UnlockHash{}) && swap.SiafundOutputs[0].UnlockHash == (types.UnlockHash{}) {
		return errors.New("one output address should be left unspecified")
	} else if len(swap.Signatures) > 0 {
		return errors.New("transaction should not have any signatures yet")
	}
	return nil
}

// AcceptSwap accepts and signs a swap transaction.
func AcceptSwap(swap *SwapTransaction) error {
	wag, err := siad.WalletAddressGet()
	if err != nil {
		return fmt.Errorf("failed to get wallet address: %w", err)
	} else if len(swap.SiacoinInputs) == 0 {
		swap.SiafundOutputs[0].UnlockHash = wag.Address
		if err := addSC(swap, swap.SiacoinOutputs[0].Value.Add(minerFee)); err != nil {
			return fmt.Errorf("failed to add siacoin inputs: %w", err)
		}
		return signSC(swap)
	}
	swap.SiacoinOutputs[0].UnlockHash = wag.Address
	if err := addSF(swap, swap.SiafundOutputs[0].Value); err != nil {
		return fmt.Errorf("failed to add siafund inputs: %w", err)
	}
	return signSF(swap)
}

// CheckFinish checks that the accepted swap transaction is valid
func CheckFinish(swap SwapTransaction) error {
	if len(swap.SiacoinInputs) == 0 || len(swap.SiafundInputs) == 0 {
		return errors.New("transaction is missing inputs")
	} else if len(swap.SiacoinOutputs) == 0 || len(swap.SiafundOutputs) == 0 {
		return errors.New("transaction is missing outputs")
	} else if swap.SiacoinOutputs[0].UnlockHash == (types.UnlockHash{}) || swap.SiafundOutputs[0].UnlockHash == (types.UnlockHash{}) {
		return errors.New("one or both swap output addresses have been left unspecified")
	} else if len(swap.Signatures) == 0 {
		return errors.New("transaction is missing counterparty signatures")
	}

	wag, err := siad.WalletAddressesGet()
	if err != nil {
		return fmt.Errorf("failed to get wallet addresses: %w", err)
	}
	belongsToUs := make(map[types.UnlockHash]bool)
	for _, addr := range wag.Addresses {
		belongsToUs[addr] = true
	}

	var haveSCSignature bool
	for _, sci := range swap.SiacoinInputs {
		if crypto.Hash(sci.ParentID) == swap.Signatures[0].ParentID {
			haveSCSignature = true
			break
		}
	}
	if haveSCSignature {
		// all of the SF inputs should belong to us
		for _, sfi := range swap.SiafundInputs {
			if !belongsToUs[sfi.UnlockConditions.UnlockHash()] {
				return errors.New("counterparty added an SF input that does not belong to us")
			}
		}
		// none of the SC inputs should belong to us
		for _, sci := range swap.SiacoinInputs {
			if belongsToUs[sci.UnlockConditions.UnlockHash()] {
				return errors.New("counterparty added an SC input that belongs to us")
			}
		}
		// all of the SF change outputs should belong to us
		for _, sfo := range swap.SiafundOutputs[1:] {
			if !belongsToUs[sfo.UnlockHash] {
				return errors.New("counterparty added an SF output that does not belong to us")
			}
		}
		// the SC output should belong to us
		if !belongsToUs[swap.SiacoinOutputs[0].UnlockHash] {
			return errors.New("the SC output address does not belong to us")
		}
	} else {
		// all of the SC inputs should belong to us
		for _, sci := range swap.SiacoinInputs {
			if !belongsToUs[sci.UnlockConditions.UnlockHash()] {
				return errors.New("counterparty added an SC input that does not belong to us")
			}
		}
		// none of the SF inputs should belong to us
		for _, sfi := range swap.SiafundInputs {
			if belongsToUs[sfi.UnlockConditions.UnlockHash()] {
				return errors.New("counterparty added an SF input that belongs to us")
			}
		}
		// all of the SC change outputs should belong to us
		for _, sco := range swap.SiacoinOutputs[1:] {
			if !belongsToUs[sco.UnlockHash] {
				return errors.New("counterparty added an SC output that does not belong to us")
			}
		}
		// the SF output should belong to us
		if !belongsToUs[swap.SiafundOutputs[0].UnlockHash] {
			return errors.New("the SF output address does not belong to us")
		}
	}
	return nil
}

// FinishSwap signs and broadcasts an accepted swap transaction.
func FinishSwap(swap *SwapTransaction) error {
	var haveSCSignatures bool
	for _, sci := range swap.SiacoinInputs {
		if crypto.Hash(sci.ParentID) == swap.Signatures[0].ParentID {
			haveSCSignatures = true
			break
		}
	}
	var err error
	if haveSCSignatures {
		err = signSF(swap)
	} else {
		err = signSC(swap)
	}
	if err != nil {
		return fmt.Errorf("failed to sign swap transaction: %w", err)
	}
	return siad.TransactionPoolRawPost(swap.Transaction(), nil)
}

// Summarize returns a summary of the swap.
func Summarize(swap SwapTransaction) (s SwapSummary, err error) {
	wag, err := siad.WalletAddressesGet()
	if err != nil {
		return SwapSummary{}, fmt.Errorf("failed to get wallet addresses: %w", err)
	}

	for _, addr := range wag.Addresses {
		s.ReceiveSC = swap.SiacoinOutputs[0].UnlockHash == addr
		s.ReceiveSF = swap.SiafundOutputs[0].UnlockHash == addr
		if s.ReceiveSC || s.ReceiveSF {
			break
		}
	}
	s.AmountSC = swap.SiacoinOutputs[0].Value
	s.AmountSF = swap.SiafundOutputs[0].Value
	s.MinerFee = minerFee
	s.Stage = 0
	if s.ReceiveSC || s.ReceiveSF {
		s.Stage++
		if len(swap.Signatures) > 0 {
			s.Stage++
			receiveSCHasSigned := swap.Signatures[0].ParentID == crypto.Hash(swap.SiacoinInputs[0].ParentID)
			receiveSFHasSigned := swap.Signatures[0].ParentID == crypto.Hash(swap.SiafundInputs[0].ParentID)
			if s.ReceiveSC && receiveSCHasSigned {
				s.Stage++
			}
			if s.ReceiveSF && receiveSFHasSigned {
				s.Stage++
			}
		}
	}
	return
}
