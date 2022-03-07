package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	stageToDescription = []string{
		"Waiting for you to accept",
		"Waiting for counterparty to accept",
		"Waiting for counterparty to finish",
		"Waiting for you to finish",
		"Swap transaction complete"}
)

func encodeSwapFile(s SwapTransaction) (name string) {
	txnID := s.Transaction().ID().String()
	name = fmt.Sprintf("embc_txn_%s.json", txnID[0:6])
	f, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	encodeJSON(f, s)
	return
}

func decodeSwapFile(filePath string) (swap SwapTransaction, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&swap); err != nil {
		fmt.Printf("%s", swap.Transaction().ID().String())
	}

	return
}

func printSummary(swap SwapTransaction) error {
	s, err := Summarize(swap)
	if err != nil {
		return err
	}
	ours, theirs := s.AmountSC.HumanString(), s.AmountSF.String()+" SF"
	if s.ReceiveSF {
		theirs, ours = ours, theirs
	}
	fmt.Println("Swap summary:")
	fmt.Println("  You receive           ", ours)
	fmt.Println("  Counterparty receives ", theirs)
	fmt.Println("  Stage                 ", stageToDescription[s.Stage])
	if s.ReceiveSF {
		fmt.Println()
		fmt.Println("  You will also pay the 5 SC transaction fee.")
	}
	return nil
}

func printTransaction(swap SwapTransaction) error {
	s, err := Summarize(swap)
	if err != nil {
		return err
	}
	nextFilePath := encodeSwapFile(swap)
	fmt.Println("Transaction:")
	fmt.Println("  ID:   ", swap.Transaction().ID())
	fmt.Println("  File: ", nextFilePath)
	fmt.Println()
	if s.Stage > 3 {
		return nil
	}
	command := "accept"
	if s.Stage > 1 {
		command = "finish"
	}
	fmt.Println("To proceed, send your counterparty the transaction file and ask them to run the following command:")
	fmt.Println()
	fmt.Println("  embc", command, nextFilePath)
	fmt.Println()
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
	printSummary(swap)
	fmt.Println()
	printTransaction(swap)
}

func acceptCLI(filePath string) {
	swap, err := decodeSwapFile(filePath)
	printSummary(swap)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckAccept(swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Printf("Accept this swap? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	fmt.Println()
	if !strings.EqualFold(resp, "y") {
		log.Fatal("  Swap cancelled.")
	} else if err = AcceptSwap(&swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Swap accepted!")
	fmt.Println()
	printTransaction(swap)
}

func finishCLI(filePath string) {
	swap, err := decodeSwapFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckFinish(swap); err != nil {
		log.Fatal(err)
	}
	printSummary(swap)
	fmt.Println()
	fmt.Printf("Sign and broadcast this transaction? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	fmt.Println()
	if !strings.EqualFold(resp, "y") {
		log.Fatal("  Swap cancelled.")
	} else if err := FinishSwap(&swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Successfully broadcast swap transaction!")
	fmt.Println()
	printTransaction(swap)
}
