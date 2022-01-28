package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"go.sia.tech/embarcadero/api"
	"go.sia.tech/embarcadero/static"
)

type CreatePayload struct {
	InStr  string `json:"inStr"`
	OutStr string `json:"outStr"`
}

type SwapPayload struct {
	SwapStr string `json:"swapStr"`
}

func writeResponse(w http.ResponseWriter, r api.Response) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := make(map[string]interface{})
	response["status"] = r.Status
	response["data"] = r.Data

	fmt.Println(response)

	writeJSON(w, response)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := make(map[string]interface{})
	response["status"] = 500
	response["message"] = err.Error()

	fmt.Println(response)

	writeJSON(w, response)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// encode nil slices as [] instead of null
	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice && val.Len() == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func getSwapStr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (string, error) {
	decoder := json.NewDecoder(r.Body)
	payload := SwapPayload{}
	err := decoder.Decode(&payload)

	if err != nil {
		return "", err
	}

	return payload.SwapStr, nil
}

func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	payload := CreatePayload{}
	err := decoder.Decode(&payload)

	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	fmt.Printf("Create\n\tinStr: %s\n\toutStr: %s\n\n", payload.InStr, payload.OutStr)

	response := api.Create(payload.InStr, payload.OutStr)

	writeResponse(w, response)
}

func Accept(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr, err := getSwapStr(w, r, ps)

	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	fmt.Printf("Accept\n\tswapStr: %s\n\n", swapStr)
	response := api.Accept(swapStr)
	writeResponse(w, response)
}

func Finish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr, err := getSwapStr(w, r, ps)

	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	fmt.Printf("Finish\n\tswapStr: %s\n\n", swapStr)
	response := api.Finish(swapStr)
	writeResponse(w, response)
}

func Summarize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr, err := getSwapStr(w, r, ps)

	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	fmt.Printf("Summarize\n\tswapStr: %s\n\n", swapStr)
	response := api.Summarize(swapStr)
	writeResponse(w, response)
}

func Consensus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := api.Consensus()

	fmt.Printf("Consensus\n\n")
	writeResponse(w, response)
}

func Wallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := api.Wallet()

	fmt.Printf("Wallet\n\n")
	writeResponse(w, response)
}

func setup(apiPort string) {
	apiAddr := "localhost:" + apiPort
	router := httprouter.New()

	// API
	router.POST("/api/create", Create)
	router.POST("/api/accept", Accept)
	router.POST("/api/finish", Finish)
	router.POST("/api/summarize", Summarize)
	router.GET("/api/consensus", Consensus)
	router.GET("/api/wallet", Wallet)

	// UI
	handlerUI := static.BuildUIHandler()
	router.NotFound = handlerUI

	// CORS, for development only
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Headers", "content-type")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	go func() {
		if err := http.ListenAndServe(apiAddr, router); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on %v...", apiAddr)
}

func Serve(apiAddr string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	setup(apiAddr)

	<-sigChan
	fmt.Println("Received interrupt, shutting down...")
}
