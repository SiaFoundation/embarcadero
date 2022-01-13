package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/node/api/client"
	"gitlab.com/NebulousLabs/Sia/types"
	"gitlab.com/NebulousLabs/encoding"
	"lukechampine.com/flagg"
)

var (
	rootUsage = `Usage:
    embc [flags] [action]

Actions:
    bids          list open bids
    trades        list completed trades
    place         place a bid
    fill          fill an open bid
`
	bidsUsage = `Usage:
embc bids

Lists open bids.
`
	tradesUsage = `Usage:
embc trades

Lists completed trades.
`
	placeUsage = `Usage:
embc place offer ask

Places a bid, offering SC or SF for the opposite. For example:

	embc place 7MS 2SF
	
places a bid that offers your 7 MS for the counterparty's 2 SF.  Outputs an
base64-encoded bid containing an unsigned transaction.
`
	fillUsage = `Usage:
embc fill bid

Attempts to fill the provided bid.  Expects a base64-encoded string.  Outputs a
partially signed base64-encoded transaction.
`

	completeUsage = `Usage:
embc complete txn

Signs the base64-encoded transaction and broadcasts it.
`
)

func main() {
	log.SetFlags(0)

	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)
	siadAddr := rootCmd.String("siad", "localhost:9980", "host:port that the siad API is running on")

	placeCmd := flagg.New("place", placeUsage)
	fillCmd := flagg.New("fill", fillUsage)
	completeCmd := flagg.New("complete", completeUsage)

	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			{Cmd: placeCmd},
			{Cmd: fillCmd},
			{Cmd: completeCmd},
		},
	})
	args := cmd.Args()

	// initialize clients
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
	case placeCmd:
		if len(args) != 2 {
			cmd.Usage()
			return
		}
		placeBid(args[0], args[1])
	case fillCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		fillBid(args[0])
	case completeCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		completeTransaction(args[0])
	}
}

var siad *client.Client

// A Bid is an unfilled trade.
type Bid struct {
	Transaction types.Transaction
	ID          types.OutputID
	Height      types.BlockHeight
	SF, SC      types.Currency
	OfferingSF  bool
	Invalid     bool `json:"omitempty"`
}

func createBid(inputAmount, outputAmount types.Currency, offeringSF bool) (Bid, error) {
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return Bid{}, err
	}

	var setupTxn types.Transaction
	bid := Bid{
		OfferingSF: offeringSF,
	}
	if bid.OfferingSF {
		bid.SF, bid.SC = inputAmount, outputAmount

		wag, err := siad.WalletAddressGet()
		if err != nil {
			return Bid{}, err
		}

		// construct setup transaction
		var inputSum types.Currency
		for _, u := range wug.Outputs {
			if u.FundType == types.SpecifierSiafundOutput {
				wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
				if err != nil {
					return Bid{}, err
				}
				setupTxn.SiafundInputs = append(setupTxn.SiafundInputs, types.SiafundInput{
					ParentID:         types.SiafundOutputID(u.ID),
					UnlockConditions: wucg.UnlockConditions,
					ClaimUnlockHash:  wag.Address,
				})
				setupTxn.TransactionSignatures = append(setupTxn.TransactionSignatures, types.TransactionSignature{
					ParentID:      crypto.Hash(types.SiacoinOutputID(u.ID)),
					CoveredFields: types.FullCoveredFields,
				})
				inputSum = inputSum.Add(u.Value)
				if inputSum.Cmp(inputAmount) >= 0 {
					break
				}
			}
		}
		if inputSum.Cmp(inputAmount) < 0 {
			return Bid{}, errors.New("insufficient funds")
		}
		wag2, err := siad.WalletAddressGet()
		if err != nil {
			return Bid{}, err
		}
		setupTxn.SiafundOutputs = []types.SiafundOutput{{
			UnlockHash: wag2.Address,
			Value:      inputAmount,
		}}
		if !inputSum.Equals(inputAmount) {
			// add change output
			wag3, err := siad.WalletAddressGet()
			if err != nil {
				return Bid{}, err
			}
			setupTxn.SiafundOutputs = append(setupTxn.SiafundOutputs, types.SiafundOutput{
				UnlockHash: wag3.Address,
				Value:      inputSum.Sub(inputAmount),
			})
		}

		// construct bid transaction
		wucg, err := siad.WalletUnlockConditionsGet(wag2.Address)
		if err != nil {
			return Bid{}, err
		}
		wag4, err := siad.WalletAddressGet()
		if err != nil {
			return Bid{}, err
		}
		bid.Transaction = types.Transaction{
			SiafundInputs: []types.SiafundInput{{
				ParentID:         setupTxn.SiafundOutputID(0),
				UnlockConditions: wucg.UnlockConditions,
			}},
			SiacoinOutputs: []types.SiacoinOutput{{
				UnlockHash: wag4.Address,
				Value:      outputAmount,
			}},
			TransactionSignatures: []types.TransactionSignature{{
				ParentID:       crypto.Hash(setupTxn.SiafundOutputID(0)),
				PublicKeyIndex: 0,
				Timelock:       0,
				CoveredFields: types.CoveredFields{
					SiafundInputs:  []uint64{0},
					SiacoinOutputs: []uint64{0},
				},
			}},
		}
		bid.ID = types.OutputID(setupTxn.SiafundOutputID(0))
	} else {
		bid.SC, bid.SF = inputAmount, outputAmount
		// construct setup transaction
		var inputSum types.Currency
		for _, u := range wug.Outputs {
			if u.FundType == types.SpecifierSiacoinOutput {
				wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
				if err != nil {
					return Bid{}, err
				}
				setupTxn.SiacoinInputs = append(setupTxn.SiacoinInputs, types.SiacoinInput{
					ParentID:         types.SiacoinOutputID(u.ID),
					UnlockConditions: wucg.UnlockConditions,
				})
				setupTxn.TransactionSignatures = append(setupTxn.TransactionSignatures, types.TransactionSignature{
					ParentID:      crypto.Hash(types.SiacoinOutputID(u.ID)),
					CoveredFields: types.FullCoveredFields,
				})
				inputSum = inputSum.Add(u.Value)
				if inputSum.Cmp(inputAmount) >= 0 {
					break
				}
			}
		}
		if inputSum.Cmp(inputAmount) < 0 {
			return Bid{}, errors.New("insufficient funds")
		}
		wag2, err := siad.WalletAddressGet()
		if err != nil {
			return Bid{}, err
		}
		setupTxn.SiacoinOutputs = []types.SiacoinOutput{{
			UnlockHash: wag2.Address,
			Value:      inputAmount,
		}}
		if !inputSum.Equals(inputAmount) {
			// add change output
			wag3, err := siad.WalletAddressGet()
			if err != nil {
				return Bid{}, err
			}
			setupTxn.SiacoinOutputs = append(setupTxn.SiacoinOutputs, types.SiacoinOutput{
				UnlockHash: wag3.Address,
				Value:      inputSum.Sub(inputAmount),
			})
		}

		// construct bid transaction
		wucg, err := siad.WalletUnlockConditionsGet(wag2.Address)
		if err != nil {
			return Bid{}, err
		}
		wag4, err := siad.WalletAddressGet()
		if err != nil {
			return Bid{}, err
		}
		bid.Transaction = types.Transaction{
			SiacoinInputs: []types.SiacoinInput{{
				ParentID:         setupTxn.SiacoinOutputID(0),
				UnlockConditions: wucg.UnlockConditions,
			}},
			SiafundOutputs: []types.SiafundOutput{{
				UnlockHash: wag4.Address,
				Value:      outputAmount,
			}},
			TransactionSignatures: []types.TransactionSignature{{
				ParentID:       crypto.Hash(setupTxn.SiacoinOutputID(0)),
				PublicKeyIndex: 0,
				Timelock:       0,
				CoveredFields: types.CoveredFields{
					SiacoinInputs:  []uint64{0},
					SiafundOutputs: []uint64{0},
				},
			}},
		}
		bid.ID = types.OutputID(setupTxn.SiacoinOutputID(0))
	}

	// sign
	wspr, err := siad.WalletSignPost(setupTxn, nil)
	if err != nil {
		return Bid{}, err
	}
	setupTxn = wspr.Transaction
	wspr, err = siad.WalletSignPost(bid.Transaction, nil)
	if err != nil {
		return Bid{}, err
	}
	bid.Transaction = wspr.Transaction

	return bid, nil
}

func fillBidTxn(bid Bid) (types.Transaction, error) {
	wug, err := siad.WalletUnspentGet()
	if err != nil {
		return types.Transaction{}, err
	}

	fillTxn := bid.Transaction
	if bid.OfferingSF {
		// fill in siacoin input(s) and siafund output
		var inputSum types.Currency
		for _, u := range wug.Outputs {
			if u.FundType == types.SpecifierSiacoinOutput {
				wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
				if err != nil {
					return types.Transaction{}, err
				}
				fillTxn.SiacoinInputs = append(fillTxn.SiacoinInputs, types.SiacoinInput{
					ParentID:         types.SiacoinOutputID(u.ID),
					UnlockConditions: wucg.UnlockConditions,
				})
				fillTxn.TransactionSignatures = append(fillTxn.TransactionSignatures, types.TransactionSignature{
					ParentID:       crypto.Hash(u.ID),
					PublicKeyIndex: 0,
					CoveredFields:  types.FullCoveredFields,
				})
				inputSum = inputSum.Add(u.Value)
				if inputSum.Cmp(bid.SC) >= 0 {
					break
				}
			}
		}
		if inputSum.Cmp(bid.SC) < 0 {
			return types.Transaction{}, errors.New("insufficient funds")
		}

		if !inputSum.Equals(bid.SC) {
			// add change output
			wag, err := siad.WalletAddressGet()
			if err != nil {
				return types.Transaction{}, err
			}
			fillTxn.SiacoinOutputs = append(fillTxn.SiacoinOutputs, types.SiacoinOutput{
				UnlockHash: wag.Address,
				Value:      inputSum.Sub(bid.SC),
			})
		}
		wag, err := siad.WalletAddressGet()
		if err != nil {
			return types.Transaction{}, err
		}
		fillTxn.SiafundOutputs = []types.SiafundOutput{{
			UnlockHash: wag.Address,
			Value:      bid.SF,
		}}
	} else {
		cuhwag, err := siad.WalletAddressGet()
		if err != nil {
			return types.Transaction{}, err
		}
		// fill in siafund input(s) and siacoin output
		var inputSum types.Currency
		for _, u := range wug.Outputs {
			if u.FundType == types.SpecifierSiafundOutput {
				wucg, err := siad.WalletUnlockConditionsGet(u.UnlockHash)
				if err != nil {
					return types.Transaction{}, err
				}
				fillTxn.SiafundInputs = append(fillTxn.SiafundInputs, types.SiafundInput{
					ParentID:         types.SiafundOutputID(u.ID),
					UnlockConditions: wucg.UnlockConditions,
					ClaimUnlockHash:  cuhwag.Address,
				})
				fillTxn.TransactionSignatures = append(fillTxn.TransactionSignatures, types.TransactionSignature{
					ParentID:       crypto.Hash(u.ID),
					PublicKeyIndex: 0,
					CoveredFields:  types.FullCoveredFields,
				})
				inputSum = inputSum.Add(u.Value)
				if inputSum.Cmp(bid.SF) >= 0 {
					break
				}
			}
		}
		if inputSum.Cmp(bid.SF) < 0 {
			return types.Transaction{}, errors.New("insufficient funds")
		}

		if !inputSum.Equals(bid.SF) {
			// add change output
			wag, err := siad.WalletAddressGet()
			if err != nil {
				return types.Transaction{}, err
			}
			fillTxn.SiafundOutputs = append(fillTxn.SiafundOutputs, types.SiafundOutput{
				UnlockHash: wag.Address,
				Value:      inputSum.Sub(bid.SF),
			})
		}
		wag, err := siad.WalletAddressGet()
		if err != nil {
			return types.Transaction{}, err
		}
		fillTxn.SiacoinOutputs = []types.SiacoinOutput{{
			UnlockHash: wag.Address,
			Value:      bid.SC,
		}}
	}
	return fillTxn, nil
}

func placeBid(inStr, outStr string) {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		log.Fatal("Invalid bid: must specify one SC value and one SF value")
	}
	input, output := parseCurrency(inStr), parseCurrency(outStr)
	bid, err := createBid(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bid created successfully.")
	fmt.Println("Share this string with your desired counterparty:")
	fmt.Println(base64.StdEncoding.EncodeToString(encoding.Marshal(bid)))
}

func fillBid(bidStr string) {
	// load bid from specified source
	var bid Bid
	var data []byte
	data, err := base64.StdEncoding.DecodeString(bidStr)
	if err == nil {
		err = encoding.Unmarshal(data, &bid)
	}

	if err != nil {
		log.Fatal(err)
	}
	// display bid details and require confirmation
	fmt.Println("Bid details:")
	var theirs, yours string
	if bid.OfferingSF {
		theirs = fmt.Sprintf("%v SF", bid.SF)
		yours = bid.SC.HumanString()
	} else {
		theirs = bid.SC.HumanString()
		yours = fmt.Sprintf("%v SF", bid.SF)
	}
	fmt.Printf("Counterparty wants to trade their %v for your %v.\n", theirs, yours)
	fmt.Print("Accept? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Trade cancelled.")
	}

	fillTxn, err := fillBidTxn(bid)
	if err != nil {
		log.Fatal(err)
	}

	// sign
	wspr, err := siad.WalletSignPost(fillTxn, nil)
	if err != nil {
		log.Fatal(err)
	}

	encoded, err := json.Marshal(wspr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bid filled successfully.")
	fmt.Println("Transaction:", base64.StdEncoding.EncodeToString(encoded))
}

func completeTransaction(txnStr string) {
	decoded, err := base64.StdEncoding.DecodeString(txnStr)
	if err != nil {
		log.Fatal(err)
	}

	var txn types.Transaction
	if err := json.Unmarshal(decoded, &txn); err != nil {
		log.Fatal(err)
	}

	// sign and broadcast
	signed, err := siad.WalletSignPost(txn, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := siad.TransactionPoolRawPost(signed.Transaction, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully broadcasted trade")
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
