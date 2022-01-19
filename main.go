package main

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
	"lukechampine.com/flagg"
)

var (
	rootUsage = `Usage:
    embc [flags] [action]

Actions:
	create        create a swap transaction
	accept        accept a swap transaction
	finish        sign + broadcast a swap transaction
`
	createUsage = `Usage:
embc create [ours] [theirs]

Creates a transaction that swaps SC for SF, or vice versa. For example:

	embc create 7MS 2SF
	
creates a transaction that swaps your 7 MS for the counterparty's 2 SF.
The transaction is unsigned, and only contains inputs from your wallet.
The counterparty must add their own inputs with 'embc accept' before the
transaction can be signed and broadcast.
`
	acceptUsage = `Usage:
embc accept [txn]

Displays a proposed swap transaction. If you accept the proposal, your inputs
will be added to complete the swap. The resulting transaction must be returned
to the original party and countersigned with 'embc finish' before it is valid
and ready for broadcasting.
`

	finishUsage = `Usage:
embc finish [txn]

Displays a proposed swap transaction. If you accept the proposal, your
signatures will be added, finalizing the transaction. The transaction is then
broadcasted.
`
)

var siad *client.Client

func main() {
	log.SetFlags(0)

	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)
	siadAddr := rootCmd.String("siad", "localhost:9980", "host:port that the siad API is running on")

	createCmd := flagg.New("create", createUsage)
	acceptCmd := flagg.New("accept", acceptUsage)
	finishCmd := flagg.New("finish", finishUsage)

	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			{Cmd: createCmd},
			{Cmd: acceptCmd},
			{Cmd: finishCmd},
		},
	})
	args := cmd.Args()

	// initialize client
	opts, _ := client.DefaultOptions()
	opts.Address = *siadAddr
	siad = client.New(opts)

	// handle command
	switch cmd {
	case rootCmd:
		if len(args) > 0 {
			cmd.Usage()
			return
		}
		fmt.Println("embc v0.1.0")
	case createCmd:
		if len(args) != 2 {
			cmd.Usage()
			return
		}
		create(args[0], args[1])
	case acceptCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		accept(args[0])
	case finishCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		finish(args[0])
	}
}

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

func create(inStr, outStr string) {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		log.Fatal("Invalid swap: must specify one SC value and one SF value")
	}
	input, output := parseCurrency(inStr), parseCurrency(outStr)
	swap, err := createSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc accept", encodeSwap(swap))
}

func accept(swapStr string) {
	swap, err := decodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := checkAccept(swap); err != nil {
		log.Fatal(err)
	}
	summarize(swap)
	fmt.Print("Accept this swap? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	err = acceptSwap(&swap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Swap accepted!")
	fmt.Println("ID:", swap.Transaction().ID())
	fmt.Println()
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc finish", encodeSwap(swap))
}

func finish(swapStr string) {
	swap, err := decodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := checkFinish(swap); err != nil {
		log.Fatal(err)
	}
	summarize(swap)
	fmt.Print("Sign and broadcast this transaction? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	err = finishSwap(&swap)
	if err != nil {
		log.Fatal(err)
	}
	if err := siad.TransactionPoolRawPost(swap.Transaction(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully broadcast swap transaction!")
	fmt.Println("ID:", swap.Transaction().ID())
}

func summarize(swap swapTransaction) {
	wag, err := siad.WalletAddressesGet()
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

func parseCurrency(amount string) types.Currency {
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

func encodeSwap(swap swapTransaction) string {
	return base64.StdEncoding.EncodeToString(encoding.Marshal(swap))
}

func decodeSwap(s string) (swapTransaction, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return swapTransaction{}, err
	}
	var swap swapTransaction
	err = encoding.Unmarshal(data, &swap)
	return swap, err
}

func addSC(swap *swapTransaction, amount types.Currency) error {
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return err
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiacoinOutput {
			wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
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
		wag, err := siad.WalletAddressGet()
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
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return err
	}
	wag, err := siad.WalletAddressGet()
	if err != nil {
		return err
	}
	var inputSum types.Currency
	for _, u := range wug.Outputs {
		if u.FundType == types.SpecifierSiafundOutput {
			wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
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
	wspr, err := siad.WalletSignPost(txn, toSign)
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
	wspr, err := siad.WalletSignPost(txn, toSign)
	swap.Signatures = wspr.Transaction.TransactionSignatures
	return err
}

func createSwap(inputAmount, outputAmount types.Currency, offeringSF bool) (swapTransaction, error) {
	wag, err := siad.WalletAddressGet()
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

func checkAccept(swap swapTransaction) error {
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

func acceptSwap(swap *swapTransaction) error {
	wag, err := siad.WalletAddressGet()
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

func checkFinish(swap swapTransaction) error {
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

func finishSwap(swap *swapTransaction) error {
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
