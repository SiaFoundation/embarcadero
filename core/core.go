package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"gitlab.com/NebulousLabs/encoding"
	"go.sia.tech/siad/crypto"
	"go.sia.tech/siad/node/api/client"
	"go.sia.tech/siad/types"
)

var Siad *client.Client

var minerFee = types.SiacoinPrecision.Mul64(5)

type swapTransaction struct {
	SiacoinInputs  []types.SiacoinInput
	SiafundInputs  []types.SiafundInput
	SiacoinOutputs []types.SiacoinOutput
	SiafundOutputs []types.SiafundOutput
	Signatures     []types.TransactionSignature
}

func (swap *swapTransaction) Transaction() types.Transaction {
	return types.Transaction{
		SiacoinInputs:         swap.SiacoinInputs,
		SiafundInputs:         swap.SiafundInputs,
		SiacoinOutputs:        swap.SiacoinOutputs,
		SiafundOutputs:        swap.SiafundOutputs,
		MinerFees:             []types.Currency{minerFee},
		TransactionSignatures: swap.Signatures,
	}
}

func Summarize(swap swapTransaction) {
	wag, err := Siad.WalletAddressesGet()
	if err != nil {
		log.Fatal(err)
	}
	var receiveSC bool
	for _, addr := range wag.Addresses {
		if swap.SiacoinOutputs[0].UnlockHash == addr {
			receiveSC = true
			break
		}
	}
	ours := swap.SiacoinOutputs[0].Value.HumanString()
	theirs := swap.SiafundOutputs[0].Value.String() + " SF"
	if !receiveSC {
		ours, theirs = theirs, ours
	}
	fmt.Println("Swap summary:")
	fmt.Println("  You receive           ", ours)
	fmt.Println("  Counterparty receives ", theirs)
	if !receiveSC {
		fmt.Println("  You will also pay the 5 SC transaction fee.")
	}
}

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

func formatCurrency(c types.Currency) string {
	pico := types.SiacoinPrecision.Div64(1e12)
	if c.Cmp(pico) < 0 {
		return c.String() + " H"
	}
	mag := pico
	unit := ""
	for _, unit = range []string{"pS", "nS", "uS", "mS", "SC", "KS", "MS", "GS", "TS"} {
		if c.Cmp(mag.Mul64(1e3)) < 0 {
			break
		} else if unit != "TS" {
			mag = mag.Mul64(1e3)
		}
	}
	num := new(big.Rat).SetInt(c.Big())
	denom := new(big.Rat).SetInt(mag.Big())
	res, _ := new(big.Rat).Mul(num, denom.Inv(denom)).Float64()
	return fmt.Sprintf("%.4g %s", res, unit)
}

func EncodeSwap(swap swapTransaction) string {
	return base64.StdEncoding.EncodeToString(encoding.Marshal(swap))
}

func DecodeSwap(s string) (swapTransaction, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return swapTransaction{}, err
	}
	var swap swapTransaction
	err = encoding.Unmarshal(data, &swap)
	return swap, err
}

func addSC(swap *swapTransaction, amount types.Currency) error {
	wug, err := Siad.WalletUnspentGet()
	if err != nil {
		return err
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiacoinOutput {
			wucg, err := Siad.WalletUnlockConditionsGet(u.UnlockHash)
			if err != nil {
				return err
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
		wag, err := Siad.WalletAddressGet()
		if err != nil {
			return err
		}
		swap.SiacoinOutputs = append(swap.SiacoinOutputs, types.SiacoinOutput{
			UnlockHash: wag.Address,
			Value:      inputSum.Sub(amount),
		})
	}
	return nil
}

func addSF(swap *swapTransaction, amount types.Currency) error {
	wug, err := Siad.WalletUnspentGet()
	if err != nil {
		return err
	}
	wag, err := Siad.WalletAddressGet()
	if err != nil {
		return err
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiafundOutput {
			wucg, err := Siad.WalletUnlockConditionsGet(u.UnlockHash)
			if err != nil {
				return err
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

func signSC(swap *swapTransaction) error {
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
	wspr, err := Siad.WalletSignPost(txn, toSign)
	swap.Signatures = wspr.Transaction.TransactionSignatures
	return err
}

func signSF(swap *swapTransaction) error {
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
	wspr, err := Siad.WalletSignPost(txn, toSign)
	swap.Signatures = wspr.Transaction.TransactionSignatures
	return err
}

func CreateSwap(inputAmount, outputAmount types.Currency, offeringSF bool) (swapTransaction, error) {
	wag, err := Siad.WalletAddressGet()
	if err != nil {
		return swapTransaction{}, err
	}
	var swap swapTransaction
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
			return swapTransaction{}, err
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
			return swapTransaction{}, err
		}
	}
	return swap, nil
}

func CheckAccept(swap swapTransaction) error {
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

func AcceptSwap(swap *swapTransaction) error {
	wag, err := Siad.WalletAddressGet()
	if err != nil {
		return err
	}
	if len(swap.SiacoinInputs) == 0 {
		swap.SiafundOutputs[0].UnlockHash = wag.Address
		if err := addSC(swap, swap.SiacoinOutputs[0].Value.Add(minerFee)); err != nil {
			return err
		}
		return signSC(swap)
	} else {
		swap.SiacoinOutputs[0].UnlockHash = wag.Address
		if err := addSF(swap, swap.SiafundOutputs[0].Value); err != nil {
			return err
		}
		return signSF(swap)
	}
}

func CheckFinish(swap swapTransaction) error {
	if len(swap.SiacoinInputs) == 0 || len(swap.SiafundInputs) == 0 {
		return errors.New("transaction is missing inputs")
	} else if len(swap.SiacoinOutputs) == 0 || len(swap.SiafundOutputs) == 0 {
		return errors.New("transaction is missing outputs")
	} else if swap.SiacoinOutputs[0].UnlockHash == (types.UnlockHash{}) || swap.SiafundOutputs[0].UnlockHash == (types.UnlockHash{}) {
		return errors.New("one or both swap output addresses have been left unspecified")
	} else if len(swap.Signatures) == 0 {
		return errors.New("transaction is missing counterparty signatures")
	}

	wag, err := Siad.WalletAddressesGet()
	if err != nil {
		return err
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

func FinishSwap(swap *swapTransaction) error {
	var haveSCSignatures bool
	for _, sci := range swap.SiacoinInputs {
		if crypto.Hash(sci.ParentID) == swap.Signatures[0].ParentID {
			haveSCSignatures = true
			break
		}
	}
	if haveSCSignatures {
		return signSF(swap)
	} else {
		return signSC(swap)
	}
}
