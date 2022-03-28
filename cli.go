package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var stageToDescription = []string{
	1: "Waiting for you to accept",
	2: "Waiting for counterparty to accept",
	3: "Waiting for you to finish",
	4: "Waiting for counterparty to finish",
	5: "Swap transaction pending",
	6: "Swap transaction confirmed",
}

func encodeSwapFile(s SwapTransaction) (string, error) {
	txnID := s.transaction().ID()
	f, err := os.Create(fmt.Sprintf("embc_txn_%x.json", txnID[:4]))
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := encodeJSON(f, s); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func decodeSwapFile(filePath string) (swap SwapTransaction, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return SwapTransaction{}, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&swap)
	return
}

func waitingOnCounterpartyOrComplete(s SwapSummary) bool {
	for _, stage := range []int{2, 4, 5, 6} {
		if stage == s.Stage {
			return true
		}
	}
	return false
}

func printSummary(s SwapSummary) error {
	ours, theirs := s.AmountSC.HumanString(), s.AmountSF.String()+" SF"
	if s.ReceiveSF {
		theirs, ours = ours, theirs
	}
	fmt.Println("Swap summary:")
	fmt.Println("  You receive:           ", ours)
	fmt.Println("  Counterparty receives: ", theirs)
	fmt.Println("  Status:                ", stageToDescription[s.Stage])
	if s.ReceiveSF {
		fmt.Println()
		fmt.Println("  You will also pay the 5 SC transaction fee.")
	}
	return nil
}

func printTransaction(swap SwapTransaction) error {
	s, err := summarize(swap)
	if err != nil {
		return err
	}
	nextFilePath, err := encodeSwapFile(swap)
	if err != nil {
		return err
	}
	fmt.Println("Transaction:")
	fmt.Println("  ID:   ", swap.transaction().ID())
	fmt.Println("  File: ", nextFilePath)
	fmt.Println()
	if s.Stage > 4 {
		return nil
	}
	command := "accept"
	if s.Stage > 2 {
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
	input, output := parseCurrency(inStr), parseCurrency(outStr)
	swap, err := createSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		log.Fatal(err)
	}
	sum, err := summarize(swap)
	if err != nil {
		log.Fatal(err)
	}
	printSummary(sum)
	fmt.Println()
	printTransaction(swap)
}

func acceptCLI(filePath string) {
	swap, err := decodeSwapFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	sum, err := summarize(swap)
	if err != nil {
		log.Fatal(err)
	}
	printSummary(sum)
	if waitingOnCounterpartyOrComplete(sum) {
		return
	}
	if err := checkAccept(swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Printf("Accept this swap? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	fmt.Println()
	if !strings.EqualFold(resp, "y") {
		log.Fatal("  Swap cancelled.")
	} else if err = acceptSwap(&swap); err != nil {
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
	if err := checkFinish(swap, false); err != nil {
		log.Fatal(err)
	}
	sum, err := summarize(swap)
	if err != nil {
		log.Fatal(err)
	}
	printSummary(sum)
	if waitingOnCounterpartyOrComplete(sum) {
		return
	}
	fmt.Println()
	fmt.Printf("Sign and broadcast this transaction? [y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	fmt.Println()
	if !strings.EqualFold(resp, "y") {
		log.Fatal("  Swap cancelled.")
	} else if err := finishSwap(&swap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Successfully broadcast swap transaction!")
	fmt.Println()
	printTransaction(swap)
}
