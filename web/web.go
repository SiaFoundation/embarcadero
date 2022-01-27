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

func writeResponse(w http.ResponseWriter, r api.Response, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := make(map[string]interface{})
	response["status"] = r.Status
	response["message"] = r.Message
	response["data"] = data

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

type CreatePayload struct {
	InStr  string `json:"inStr"`
	OutStr string `json:"outStr"`
}

func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// inStr := ps.ByName("inStr")
	// outStr := ps.ByName("outStr")
	decoder := json.NewDecoder(r.Body)
	payload := CreatePayload{}
	err := decoder.Decode(&payload)

	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	fmt.Println(ps)

	response := api.Create(payload.InStr, payload.OutStr)
	data := make(map[string]interface{})

	writeResponse(w, response, data)
}

func Accept(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr := ps.ByName("swapStr")

	fmt.Println(swapStr)

	response := api.Accept(swapStr)
	data := make(map[string]interface{})

	writeResponse(w, response, data)
}

func Finish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr := ps.ByName("swapStr")

	fmt.Println(swapStr)

	response := api.Finish(swapStr)
	data := make(map[string]interface{})

	writeResponse(w, response, data)
}

func setup(apiPort string) {
	apiAddr := "localhost:" + apiPort
	router := httprouter.New()

	// API
	router.POST("/api/create", Create)
	router.POST("/api/accept", Accept)
	router.POST("/api/finish", Finish)

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
