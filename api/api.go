package api

import (
	"fmt"
	"strings"

	"go.sia.tech/embarcadero/core"
)

type Response struct {
	Status  int
	Message string
}

func Create(inStr, outStr string) Response {
	if strings.Contains(inStr, "SF") == strings.Contains(outStr, "SF") {
		return Response{
			Status:  500,
			Message: "Invalid swap: must specify one SC value and one SF value",
		}
	}
	input, output := core.ParseCurrency(inStr), core.ParseCurrency(outStr)
	swap, err := core.CreateSwap(input, output, strings.Contains(inStr, "SF"))
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}

	message := fmt.Sprintln("To proceed, ask your counterparty to run the following command:")
	message += fmt.Sprintln("")
	message += fmt.Sprintln("    embc accept", core.EncodeSwap(swap))

	return Response{
		Status:  200,
		Message: message,
	}
}

func Summarize(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	if err := core.CheckAccept(swap); err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}

	message, err := core.Summarize(swap)

	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}

	return Response{
		Status:  200,
		Message: message,
	}
}

func Accept(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	if err := core.CheckAccept(swap); err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}

	err = core.AcceptSwap(&swap)
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}

	message := fmt.Sprintln("Swap accepted!")
	message += fmt.Sprintln("ID:", swap.Transaction().ID())
	message += fmt.Sprintln()
	message += fmt.Sprintln("To proceed, ask your counterparty to run the following command:")
	message += fmt.Sprintln()
	message += fmt.Sprintln("    embc finish", core.EncodeSwap(swap))

	return Response{
		Status:  200,
		Message: message,
	}
}

func Finish(swapStr string) Response {
	swap, err := core.DecodeSwap(swapStr)
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	if err := core.CheckFinish(swap); err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	err = core.FinishSwap(&swap)
	if err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	if err := core.Siad.TransactionPoolRawPost(swap.Transaction(), nil); err != nil {
		return Response{
			Status:  500,
			Message: err.Error(),
		}
	}
	message := fmt.Sprintln("Successfully broadcast swap transaction!")
	message += fmt.Sprintln("ID:", swap.Transaction().ID())

	return Response{
		Status:  200,
		Message: message,
	}
}
