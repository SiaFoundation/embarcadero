package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"text/tabwriter"

	"gitlab.com/NebulousLabs/Sia/types"
	"lukechampine.com/flagg"

	"github.com/SiaFoundation/embarcadero/embarcadero"
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
embc place [flags] offer ask

Places a bid, offering SC or SF for the opposite. For example:

	embc place 7MS 2SF
	
places a bid that offers your 7 MS for the counterparty's 2 SF. By default,
the bid is stored on the Sia blockchain, which makes the bid public and
incurs a transaction fee. Use the --skynet or --b64 flags to create a private
bid with no fee.
`
	fillUsage = `Usage:
embc fill [flags] bid

Attempts to fill the provided bid, which can take a number of forms; by
default, it expects one of the bid IDs returned by the bids command. If the
--skynet flag is set, it expects a Skynet link. If the --b64 flag is set, it
expects a base64-encoded string.

Once filled, the completed trade transaction is signed and broadcast. This
will incur a transaction fee.
`
)

func main() {
	log.SetFlags(0)

	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)
	addr := rootCmd.String("a", "http://localhost:8080", "host:port that the embarcadero API is running on")

	bidsCmd := flagg.New("bids", bidsUsage)
	tradesCmd := flagg.New("trades", tradesUsage)
	placeCmd := flagg.New("place", placeUsage)
	fillCmd := flagg.New("fill", fillUsage)

	var skynet, b64 bool
	placeCmd.BoolVar(&skynet, "skynet", false, "place bid via Skynet instead of storing it on the blockchain")
	placeCmd.BoolVar(&b64, "base64", false, "place bid via base64 strings instead of storing it on the blockchain")
	fillCmd.BoolVar(&skynet, "skynet", false, "fill bid via Skynet instead of storing it on the blockchain")
	fillCmd.BoolVar(&b64, "base64", false, "fill bid via base64 strings instead of storing it on the blockchain")

	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			{Cmd: bidsCmd},
			{Cmd: tradesCmd},
			{Cmd: placeCmd},
			{Cmd: fillCmd},
		},
	})
	args := cmd.Args()

	if skynet && b64 {
		log.Fatal("Cannot use both skynet and base64; choose one!")
	}

	// initialize clients
	embd = embarcadero.NewClient(*addr)

	// handle command
	switch cmd {
	case rootCmd:
		if len(args) > 0 {
			cmd.Usage()
			return
		}
		fmt.Println("embc v0.1.0")
	case bidsCmd:
		if len(args) > 0 {
			cmd.Usage()
			return
		}
		viewBids()
	case tradesCmd:
		if len(args) > 0 {
			cmd.Usage()
			return
		}
		viewTrades()
	case placeCmd:
		if len(args) != 2 {
			cmd.Usage()
			return
		}
		embd.PlaceBid(args[0], args[1], skynet, b64)
	case fillCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		embd.FillBid(args[0], skynet, b64)
	}
}

var (
	embd *embarcadero.Client
)

func viewBids() {
	bids, err := embd.Bids()
	if err != nil {
		log.Fatal(err)
	}
	if len(bids) == 0 {
		fmt.Println("No bids")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHeight\tBid\tAsk")
	for _, b := range bids {
		if b.OfferingSF {
			fmt.Fprintf(w, "%v\t%v\t%v SF\t%v\n", b.ID.String()[:8], b.Height, b.SF, formatCurrency(b.SC))
		} else {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v SF\n", b.ID.String()[:8], b.Height, formatCurrency(b.SC), b.SF)
		}
	}
	w.Flush()
}

func viewTrades() {
	trades, err := embd.Trades()
	if err != nil {
		log.Fatal(err)
	}
	if len(trades) == 0 {
		fmt.Println("No trades")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Bid Height\tFill Height\tPut\tAsk")
	for _, t := range trades {
		if t.Bid.OfferingSF {
			fmt.Fprintf(w, "%v\t%v\t%v SF\t%v\n", t.Bid.Height, t.Height, t.Bid.SF, formatCurrency(t.Bid.SC))
		} else {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v SF\n", t.Bid.Height, t.Height, formatCurrency(t.Bid.SC), t.Bid.SF)
		}
	}
	w.Flush()
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
