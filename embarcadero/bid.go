package embarcadero

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/NebulousLabs/go-skynet"
	"gitlab.com/NebulousLabs/Sia/crypto"
	"gitlab.com/NebulousLabs/Sia/types"
	"gitlab.com/NebulousLabs/encoding"
)

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

	// broadcast setup transaction
	if err := siad.TransactionPoolRawPost(setupTxn, nil); err != nil {
		return Bid{}, err
	}

	return bid, nil
}

func uploadToSkynet(bid Bid) (string, error) {
	data := skynet.UploadData{"embarcaderobid": bytes.NewReader(encoding.Marshal(bid))}
	return skynet.Upload(data, skynet.DefaultUploadOptions)
}

func downloadFromSkynet(link string) (bid Bid, err error) {
	rc, err := skynet.Download(link, skynet.DefaultDownloadOptions)
	if err != nil {
		return
	}
	defer rc.Close()
	err = encoding.NewDecoder(rc, encoding.DefaultAllocLimit).Decode(&bid)
	return
}

func storeOnChain(bid Bid) (string, error) {
	txn := types.Transaction{
		ArbitraryData: [][]byte{append(BidPrefix[:], encoding.Marshal(bid)...)},
	}
	if err := siad.TransactionPoolRawPost(txn, nil); err != nil {
		return "", err
	}
	return bid.ID.String(), nil
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

func PlaceBid(inStr, outStr string, skynet, b64 bool) {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		log.Fatal("Invalid bid: must specify one SC value and one SF value")
	}
	input, output := parseCurrency(inStr), parseCurrency(outStr)
	bid, err := createBid(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case skynet:
		link, err := uploadToSkynet(bid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Bid created successfully.")
		fmt.Println("Share this link with your desired counterparty:")
		fmt.Println("   ", link)

	case b64:
		fmt.Println("Bid created successfully.")
		fmt.Println("Share this string with your desired counterparty:")
		fmt.Println(base64.StdEncoding.EncodeToString(encoding.Marshal(bid)))

	default:
		id, err := storeOnChain(bid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Bid created successfully.")
		fmt.Println("Your bid has been submitted for inclusion in the next block.")
		fmt.Println("When the bid appears on-chain, it will be listed in the 'bids' command.")
		fmt.Println("Your bid ID is:")
		fmt.Println("   ", id)
	}
}

func FillBid(mt *MarketTracker, bidStr string, skynet, b64 bool) {
	// load bid from specified source
	var bid Bid
	var err error
	switch {
	case skynet:
		bid, err = downloadFromSkynet(bidStr)

	case b64:
		var data []byte
		data, err = base64.StdEncoding.DecodeString(bidStr)
		if err == nil {
			err = encoding.Unmarshal(data, &bid)
		}

	default:
		var bids []Bid
		bids = mt.Bids()
		var matches int
		for _, b := range bids {
			if strings.HasPrefix(b.ID.String(), bidStr) {
				bid = b
				matches++
			}
		}
		if matches == 0 {
			err = errors.New("bid not found")
		} else if matches > 1 {
			err = errors.New("bid ID not unique; add more digits")
		}
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
	// sign and broadcast
	wspr, err := siad.WalletSignPost(fillTxn, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := siad.TransactionPoolRawPost(wspr.Transaction, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bid filled successfully.")
	fmt.Println("Transaction ID:", wspr.Transaction.ID())
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
