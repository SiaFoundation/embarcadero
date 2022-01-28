package api

import (
	"fmt"
	"strings"

	"go.sia.tech/embarcadero/core"
)

type Response struct {
	Status int
	Data   interface{}
}

func sendMessage(status int, message string) Response {
	data := make(map[string]interface{})
	data["message"] = message

	return Response{
		Status: status,
		Data:   data,
	}
}

func sendData(status int, data interface{}) Response {
	return Response{
		Status: status,
		Data:   data,
	}
}

func Create(inStr, outStr string) Response {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		return sendMessage(500, "Invalid swap: must specify one SC value and one SF value")
	}
	input, output := core.ParseCurrency(inStr), core.ParseCurrency(outStr)
	swap, err := core.CreateSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		return sendMessage(500, err.Error())
	}

	message := fmt.Sprintln("To proceed, ask your counterparty to run the following command:")
	message += fmt.Sprintln("")
	message += fmt.Sprintln("    embc accept", core.EncodeSwap(swap))

	return sendMessage(200, message)
}

func Accept(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return sendMessage(500, err.Error())
	}
	if err := core.CheckAccept(swap); err != nil {
		return sendMessage(500, err.Error())
	}

	err = core.AcceptSwap(&swap)
	if err != nil {
		return sendMessage(500, err.Error())
	}

	message := fmt.Sprintln("Swap accepted!")
	message += fmt.Sprintln("ID:", swap.Transaction().ID())
	message += fmt.Sprintln()
	message += fmt.Sprintln("To proceed, ask your counterparty to run the following command:")
	message += fmt.Sprintln()
	message += fmt.Sprintln("    embc finish", core.EncodeSwap(swap))

	return sendMessage(200, message)
}

func Finish(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return sendMessage(500, err.Error())
	}
	if err := core.CheckFinish(swap); err != nil {
		return sendMessage(500, err.Error())
	}
	err = core.FinishSwap(&swap)
	if err != nil {
		return sendMessage(500, err.Error())
	}
	if err := core.Siad.TransactionPoolRawPost(swap.Transaction(), nil); err != nil {
		return sendMessage(500, err.Error())
	}
	message := fmt.Sprintln("Successfully broadcast swap transaction!")
	message += fmt.Sprintln("ID:", swap.Transaction().ID())

	return sendMessage(200, message)
}

func Summarize(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return sendMessage(500, err.Error())
	}
	if err := core.CheckAccept(swap); err != nil {
		return sendMessage(500, err.Error())
	}

	message, err := core.Summarize(swap)

	if err != nil {
		return sendMessage(500, err.Error())
	}

	return sendMessage(200, message)
}

func Consensus() Response {
	c, err := core.Siad.ConsensusGet()

	if err != nil {
		return sendMessage(500, err.Error())
	}

	return sendData(200, c)
}

func Wallet() Response {
	c, err := core.Siad.WalletGet()

	if err != nil {
		return sendMessage(500, err.Error())
	}

	return sendData(200, c)
}
