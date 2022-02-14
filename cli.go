package main

import (
	"fmt"
	"log"
	"strings"
)

func summarizeCLI(swap SwapTransaction) error {
	s, err := Summarize(swap)
	if err != nil {
		return err
	}
	fmt.Println("Swap summary:")

	ours, theirs := s.AmountSC.HumanString(), s.AmountSF.String()

	if s.ReceiveSF {
		theirs, ours = ours, theirs
	}

	fmt.Println("  You receive           ", ours)
	fmt.Println("  Counterparty receives ", theirs)
	if s.PayFee {
		fmt.Println("  You will also pay the 5 SC transaction fee.")
	}
	return nil
}

func createCLI(inStr, outStr string) {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		log.Fatal("Invalid swap: must specify one SC value and one SF value")
	}
	input, output := ParseCurrency(inStr), ParseCurrency(outStr)
	swap, err := CreateSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc accept", EncodeSwap(swap))
}

func acceptCLI(swapStr string) {
	swap, err := DecodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckAccept(swap); err != nil {
		log.Fatal(err)
	}
	summarizeCLI(swap)
	fmt.Print("Accept this swap? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	err = AcceptSwap(&swap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Swap accepted!")
	fmt.Println("ID:", swap.Transaction().ID())
	fmt.Println()
	fmt.Println("To proceed, ask your counterparty to run the following command:")
	fmt.Println()
	fmt.Println("    embc finish", EncodeSwap(swap))
}

func finishCLI(swapStr string) {
	swap, err := DecodeSwap(swapStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckFinish(swap); err != nil {
		log.Fatal(err)
	}
	summarizeCLI(swap)
	fmt.Print("Sign and broadcast this transaction? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		log.Fatal("Swap cancelled.")
	}
	if err := FinishSwap(&swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully broadcast swap transaction!")
	fmt.Println("ID:", swap.Transaction().ID())
}
