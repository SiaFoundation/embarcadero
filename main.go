package main

import (
	"log"

	"go.sia.tech/siad/node/api/client"
	"lukechampine.com/flagg"
)

var (
	rootUsage = `Usage:
    embc [flags] [action]

Run 'embc' with no arguments to open a web UI in your browser.
Alternatively, use the actions below to conduct a swap via the CLI.

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
broadcast.
`
)

func main() {
	log.SetFlags(0)

	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)
	webAddr := rootCmd.String("addr", "localhost:8080", "HTTP service address")
	siadAddr := rootCmd.String("siad", "localhost:9980", "host:port that the siad API is running on")
	dev := rootCmd.Bool("dev", false, "run in dev mode")

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

	switch cmd {
	case rootCmd:
		serve(*webAddr, *dev)
	case createCmd:
		if len(args) != 2 {
			cmd.Usage()
			return
		}
		createCLI(args[0], args[1])
	case acceptCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		acceptCLI(args[0])
	case finishCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		finishCLI(args[0])
	}
}
