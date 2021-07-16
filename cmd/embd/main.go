package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"time"

	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/node/api/client"
	"gitlab.com/NebulousLabs/Sia/persist"
	"gitlab.com/NebulousLabs/Sia/types"
	"gitlab.com/NebulousLabs/encoding"

	"gitlab.com/NebulousLabs/embarcadero"
)

type marketTracker struct {
	height types.BlockHeight
	ccid   modules.ConsensusChangeID
	trades map[types.OutputID]embarcadero.Trade // identified by parent ID of bid input

	uscos       map[types.SiacoinOutputID]types.SiacoinOutput
	usfos       map[types.SiafundOutputID]types.SiafundOutput
	outputsBuf  *bufio.Writer
	outputsFile *os.File
	numDiffs    int // number of entries in outputsFile

	unsub   func() // from consensus
	syncing bool
	dir     string
	mu      sync.Mutex
}

func probablyFillTxn(txn types.Transaction) bool {
	if len(txn.SiafundOutputs) > 0 && len(txn.SiacoinOutputs) > 0 &&
		len(txn.SiafundInputs) > 0 && len(txn.SiacoinInputs) > 0 {
		for _, sig := range txn.TransactionSignatures {
			if !sig.CoveredFields.WholeTransaction {
				return true
			}
		}
	}
	return false
}

func extractBidTxn(height types.BlockHeight, txn types.Transaction) (types.Transaction, error) {
	// All we care about is that the transaction has a valid partial signature
	// that covers either one siacoin input and one siafund output, or one
	// siafund input and one siacoin output. (The transaction *shouldn't*
	// contain anything else, but it doesn't matter; we can extract just the
	// parts we care about without invalidating the signature.)
	validCoveredFields := func(cf types.CoveredFields) bool {
		if cf.WholeTransaction || len(cf.FileContracts)+len(cf.FileContractRevisions)+len(cf.StorageProofs)+len(cf.MinerFees)+len(cf.ArbitraryData)+len(cf.TransactionSignatures) > 0 {
			return false
		}
		scInput := len(cf.SiacoinInputs) == 1 && len(cf.SiafundOutputs) == 1 &&
			cf.SiacoinInputs[0] < uint64(len(txn.SiacoinInputs)) &&
			cf.SiafundOutputs[0] < uint64(len(txn.SiafundOutputs))
		sfInput := len(cf.SiafundInputs) == 1 && len(cf.SiacoinOutputs) == 1 &&
			cf.SiafundInputs[0] < uint64(len(txn.SiafundInputs)) &&
			cf.SiacoinOutputs[0] < uint64(len(txn.SiacoinOutputs))
		return sfInput || scInput
	}
	var keep types.Transaction
	for _, sig := range txn.TransactionSignatures {
		cf := sig.CoveredFields
		if !validCoveredFields(cf) {
			continue
		}
		if len(cf.SiacoinInputs) == 1 {
			keep.SiacoinInputs = []types.SiacoinInput{txn.SiacoinInputs[cf.SiacoinInputs[0]]}
			keep.SiafundOutputs = []types.SiafundOutput{txn.SiafundOutputs[cf.SiafundOutputs[0]]}
		} else {
			keep.SiafundInputs = []types.SiafundInput{txn.SiafundInputs[cf.SiafundInputs[0]]}
			keep.SiacoinOutputs = []types.SiacoinOutput{txn.SiacoinOutputs[cf.SiacoinOutputs[0]]}
		}
		keep.TransactionSignatures = []types.TransactionSignature{sig}
		break
	}
	return keep, keep.StandaloneValid(height)
}

func findBids(height types.BlockHeight, txns []types.Transaction) (newBids []embarcadero.Bid, filledBids []types.Transaction) {
	for _, txn := range txns {
		if probablyFillTxn(txn) {
			filledBids = append(filledBids, txn)
		}
		for _, arb := range txn.ArbitraryData {
			if !bytes.HasPrefix(arb, embarcadero.BidPrefix) {
				continue
			}
			var bid embarcadero.Bid
			if err := encoding.Unmarshal(bytes.TrimPrefix(arb, embarcadero.BidPrefix), &bid); err != nil {
				log.Println("failed to decode bid:", err)
				continue
			} else if bid.Transaction, err = extractBidTxn(height, bid.Transaction); err != nil {
				log.Println("invalid bid txn:", err)
				continue
			}
			if len(bid.Transaction.SiacoinInputs) == 1 {
				bid.ID = types.OutputID(bid.Transaction.SiacoinInputs[0].ParentID)
			} else {
				bid.ID = types.OutputID(bid.Transaction.SiafundInputs[0].ParentID)
			}
			bid.Height = height
			newBids = append(newBids, bid)
		}
	}
	return
}

func (mt *marketTracker) ProcessConsensusChange(cc modules.ConsensusChange) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if mt.ccid == modules.ConsensusChangeBeginning {
		mt.height-- // genesis block is height 0
	}
	mt.ccid = cc.ID

	// avoid writing to disk on every consensus change
	var shouldSave bool

	// update unspent outputs
	enc := encoding.NewEncoder(mt.outputsBuf)
	for _, diff := range cc.SiacoinOutputDiffs {
		if diff.Direction == modules.DiffApply {
			mt.uscos[diff.ID] = diff.SiacoinOutput
		} else {
			delete(mt.uscos, diff.ID)
		}
		enc.EncodeAll(true, diff)
		mt.numDiffs++
	}
	for _, diff := range cc.SiafundOutputDiffs {
		if diff.Direction == modules.DiffApply {
			mt.usfos[diff.ID] = diff.SiafundOutput
		} else {
			delete(mt.usfos, diff.ID)
		}
		enc.EncodeAll(false, diff)
		mt.numDiffs++
	}

	// process reverted blocks: delete bids that were placed, and un-fill bids
	// that were filled.
	for _, b := range cc.RevertedBlocks {
		mt.height--
		revertedBids, unfilledBids := findBids(mt.height, b.Transactions)
		shouldSave = shouldSave || len(revertedBids) > 0 || len(unfilledBids) > 0
		for _, bid := range revertedBids {
			delete(mt.trades, bid.ID)
		}
		for _, txn := range unfilledBids {
			// txn may have filled multiple bids
			var ids []types.OutputID
			for _, in := range txn.SiacoinInputs {
				ids = append(ids, types.OutputID(in.ParentID))
			}
			for _, in := range txn.SiafundInputs {
				ids = append(ids, types.OutputID(in.ParentID))
			}
			for _, id := range ids {
				if t, ok := mt.trades[id]; ok {
					t.Transaction = types.Transaction{}
					t.Height = 0
					mt.trades[id] = t
				}
			}
		}
	}
	// process applied blocks: create bids that were placed, and fill bids that
	// were filled.
	for _, b := range cc.AppliedBlocks {
		mt.height++
		appliedBids, filledBids := findBids(mt.height, b.Transactions)
		shouldSave = shouldSave || len(appliedBids) > 0 || len(filledBids) > 0
		for _, bid := range appliedBids {
			mt.trades[bid.ID] = embarcadero.Trade{Bid: bid}
		}
		for _, txn := range filledBids {
			// txn may have filled multiple bids
			var ids []types.OutputID
			for _, in := range txn.SiacoinInputs {
				ids = append(ids, types.OutputID(in.ParentID))
			}
			for _, in := range txn.SiafundInputs {
				ids = append(ids, types.OutputID(in.ParentID))
			}
			for _, id := range ids {
				if t, ok := mt.trades[id]; ok {
					t.Transaction = txn
					t.Height = mt.height
					mt.trades[id] = t
				}
			}
		}
	}

	// if a bid double-spends an output, mark it as invalid (but don't delete
	// it, because the double-spend might be reverted)
	for id, t := range mt.trades {
		if t.Height != 0 {
			continue // we only consider unfilled bids
		}
		if t.Bid.OfferingSF {
			sfo, ok := mt.usfos[types.SiafundOutputID(id)]
			t.Bid.Invalid = !ok || !sfo.Value.Equals(t.Bid.SF)
		} else {
			sco, ok := mt.uscos[types.SiacoinOutputID(id)]
			t.Bid.Invalid = !ok || !sco.Value.Equals(t.Bid.SC)
		}
		mt.trades[id] = t
	}

	// save at least once per day, or once per 10k blocks if we're syncing.
	//
	// NOTE: this may occasionally fail to fire if we process multiple blocks at
	// a time; but in steady-state, this should be rare.
	if mt.syncing {
		shouldSave = shouldSave || mt.height%10000 == 0
	} else {
		shouldSave = shouldSave || mt.height%144 == 0
	}

	if shouldSave {
		if err := mt.save(); err != nil {
			log.Fatal("Failed to save:", err)
		}
	}
}

func (mt *marketTracker) Bids() []embarcadero.Bid {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	bids := make([]embarcadero.Bid, 0, len(mt.trades))
	for _, trade := range mt.trades {
		if trade.Height == 0 && !trade.Bid.Invalid {
			bids = append(bids, trade.Bid)
		}
	}
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Height > bids[j].Height
	})
	return bids
}

func (mt *marketTracker) Trades() []embarcadero.Trade {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	trades := make([]embarcadero.Trade, 0, len(mt.trades))
	for _, trade := range mt.trades {
		if trade.Height != 0 {
			trades = append(trades, trade)
		}
	}
	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Height > trades[j].Height
	})
	return trades
}

var meta = persist.Metadata{
	Header:  "Embarcadero",
	Version: "0.1.0",
}

type persistData struct {
	Height          types.BlockHeight
	CCID            modules.ConsensusChangeID
	OutputsFileSize int64
	Trades          []embarcadero.Trade
}

func (mt *marketTracker) save() error {
	// flush and sync outputs file
	if err := mt.outputsBuf.Flush(); err != nil {
		return err
	} else if err := mt.outputsFile.Sync(); err != nil {
		return err
	}
	stat, err := mt.outputsFile.Stat()
	if err != nil {
		return err
	}
	// write JSON file atomically
	data := persistData{
		Height:          mt.height,
		CCID:            mt.ccid,
		OutputsFileSize: stat.Size(),
		Trades:          make([]embarcadero.Trade, 0, len(mt.trades)),
	}
	for _, t := range mt.trades {
		data.Trades = append(data.Trades, t)
	}
	err = persist.SaveJSON(meta, data, filepath.Join(mt.dir, "persist.json"))
	if err != nil {
		return err
	}
	return nil
}

func (mt *marketTracker) load() error {
	var data persistData
	if err := persist.LoadJSON(meta, &data, filepath.Join(mt.dir, "persist.json")); err != nil && !os.IsNotExist(err) {
		return err
	}
	mt.trades = make(map[types.OutputID]embarcadero.Trade)
	for _, t := range data.Trades {
		mt.trades[t.Bid.ID] = t
	}
	mt.height = data.Height
	mt.ccid = data.CCID

	// load outputs
	mt.uscos = make(map[types.SiacoinOutputID]types.SiacoinOutput)
	mt.usfos = make(map[types.SiafundOutputID]types.SiafundOutput)
	f, err := os.OpenFile(filepath.Join(mt.dir, "outputs.dat"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if err := f.Truncate(data.OutputsFileSize); err != nil {
		return err
	}
	br := bufio.NewReader(f)
	dec := encoding.NewDecoder(br, -1)
	for dec.Err() == nil {
		if isSC := dec.NextBool(); isSC {
			var diff modules.SiacoinOutputDiff
			if dec.Decode(&diff) == nil {
				if diff.Direction == modules.DiffApply {
					mt.uscos[diff.ID] = diff.SiacoinOutput
				} else {
					delete(mt.uscos, diff.ID)
				}
			}
		} else {
			var diff modules.SiafundOutputDiff
			if dec.Decode(&diff) == nil {
				if diff.Direction == modules.DiffApply {
					mt.usfos[diff.ID] = diff.SiafundOutput
				} else {
					delete(mt.usfos, diff.ID)
				}
			}
		}
	}
	if dec.Err() != io.EOF {
		return dec.Err()
	}

	mt.outputsFile = f
	mt.outputsBuf = bufio.NewWriter(f)
	return nil
}

func (mt *marketTracker) Close() error {
	mt.unsub()
	return mt.save()
}

func newMarketTracker(dir string) (*marketTracker, error) {
	mt := &marketTracker{dir: dir}
	if err := mt.load(); err != nil {
		return nil, err
	}
	return mt, nil
}

func main() {
	log.SetFlags(0)
	apiAddr := flag.String("a", "localhost:8080", "host:port to serve the embarcadero API on")
	siadAddr := flag.String("siad", "localhost:9980", "host:port that the siad API is listening on")
	dir := flag.String("d", ".", "directory where server state will be stored")
	flag.Parse()

	mt, err := start(*dir, *siadAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer mt.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	srv := serve(mt, *apiAddr)
	<-sigChan
	fmt.Println("Received interrupt, shutting down...")
	srv.Shutdown(context.Background())
}

func start(dir string, siadAddr string) (*marketTracker, error) {
	opts, _ := client.DefaultOptions()
	client := client.New(opts)
	client.Address = siadAddr
	mt, err := newMarketTracker(dir)
	if err != nil {
		return nil, err
	}
	errCh, unsub := client.ConsensusSetSubscribe(mt, mt.ccid, nil)
	mt.unsub = unsub
	mt.syncing = true
	for caughtUp := false; !caughtUp; {
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Println()
				return nil, err
			}
			caughtUp = true
		case <-time.After(time.Second):
			mt.mu.Lock()
			fmt.Printf("\rSynced to height %v...", mt.height)
			mt.mu.Unlock()
		}
	}
	fmt.Println()
	go func() {
		if err := <-errCh; err != nil {
			log.Fatalln("ConsensusSetSubscribe:", err)
		}
	}()
	// need to use a mutex here because ProcessConsensusChange may be called
	// concurrently
	mt.mu.Lock()
	mt.syncing = false
	mt.mu.Unlock()
	return mt, nil
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// encode nil slices as [] instead of null
	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice && val.Len() == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(v)
}

func serve(mt *marketTracker, APIaddr string) *http.Server {
	srv := &http.Server{
		Addr: APIaddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodGet {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			switch path.Base(req.URL.Path) {
			case "bids":
				writeJSON(w, mt.Bids())
			case "trades":
				writeJSON(w, mt.Trades())
			default:
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
		}),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on %v...", APIaddr)
	return srv
}
