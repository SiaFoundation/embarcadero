package cli

import (
	"fmt"
	"log"
	"strings"

	"go.sia.tech/embarcadero/core"
)

func Create(inStr, outStr string) {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		log.Fatal("Invalid swap: must specify one SC value and one SF value")
	}
	input, output := core.ParseCurrency(inStr), core.ParseCurrency(outStr)
	swap, err := core.CreateSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc accept", core.EncodeSwap(swap))
}

func Accept(swapStr string) {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := core.CheckAccept(swap); err != nil {
		log.Fatal(err)
	}
	core.Summarize(swap)
	fmt.Print("Accept this swap? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	err = core.AcceptSwap(&swap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Swap accepted!")
	fmt.Println("ID:", swap.Transaction().ID())
	fmt.Println()
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc finish", core.EncodeSwap(swap))
}

func Finish(swapStr string) {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := core.CheckFinish(swap); err != nil {
		log.Fatal(err)
	}
	core.Summarize(swap)
	fmt.Print("Sign and broadcast this transaction? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	err = core.FinishSwap(&swap)
	if err != nil {
		log.Fatal(err)
	}
	if err := core.Siad.TransactionPoolRawPost(swap.Transaction(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully broadcast swap transaction!")
	fmt.Println("ID:", swap.Transaction().ID())
}
