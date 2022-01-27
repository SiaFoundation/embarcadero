package main

import (
	"fmt"
	"log"

	"go.sia.tech/embarcadero/cli"
	"go.sia.tech/embarcadero/core"
	"go.sia.tech/embarcadero/web"
	"go.sia.tech/siad/node/api/client"
	"lukechampine.com/flagg"
)

var (
	rootUsage = `Usage:
    embc [flags] [action]

Actions:
	create        create a swap transaction
	accept        accept a swap transaction
	finish        sign + broadcast a swap transaction
	ui            start the web UI
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

	uiUsage = `Usage:
embc ui [apiPort]

Opens the UI in your browser. Runs the embarcadero API server in the background. Accepts a custom API port. For example:

	embc ui 8081
`
)

func main() {
	log.SetFlags(0)

	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)
	siadAddr := rootCmd.String("siad", "localhost:9980", "host:port that the siad API is running on")

	createCmd := flagg.New("create", createUsage)
	acceptCmd := flagg.New("accept", acceptUsage)
	finishCmd := flagg.New("finish", finishUsage)
	uiCmd := flagg.New("ui", uiUsage)

	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			{Cmd: createCmd},
			{Cmd: acceptCmd},
			{Cmd: finishCmd},
			{Cmd: uiCmd},
		},
	})
	args := cmd.Args()

	// initialize client
	opts, _ := client.DefaultOptions()
	opts.Address = *siadAddr
	core.Siad = client.New(opts)

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
		cli.Create(args[0], args[1])
	case acceptCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		cli.Accept(args[0])
	case finishCmd:
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		cli.Finish(args[0])
	case uiCmd:
		if len(args) != 0 {
			web.Serve(args[1])
		}
		web.Serve("9981")
	}
}
