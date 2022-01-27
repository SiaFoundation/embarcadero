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
)

type response struct {
	status int
	data   map[string]interface{}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// encode nil slices as [] instead of null
	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice && val.Len() == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(v)
}

func postCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	inStr := ps.ByName("inStr")
	outStr := ps.ByName("outStr")

	Create(inStr, outStr)

	response := response{
		status: 200,
		data: map[string]interface{}{
			"status": "success",
		},
	}

	writeJSON(w, response)
}

func postAccept(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr := ps.ByName("swapStr")

	Accept(swapStr)

	response := response{
		status: 200,
		data: map[string]interface{}{
			"status": "success",
		},
	}

	writeJSON(w, response)
}

func postFinish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	swapStr := ps.ByName("swapStr")

	Finish(swapStr)

	response := response{
		status: 200,
		data: map[string]interface{}{
			"status": "success",
		},
	}

	writeJSON(w, response)
}

func getPing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Success")

	response := response{
		status: 200,
		data: map[string]interface{}{
			"status": "success",
		},
	}

	writeJSON(w, response)
}

func setup(apiPort string) {
	apiAddr := "locahost:" + apiPort
	router := httprouter.New()

	// API
	router.POST("/api/create", postCreate)
	router.POST("/api/accept", postAccept)
	router.POST("/api/finish", postFinish)
	router.GET("/api/ping", getPing)

	// UI
	handlerUI := buildUIHandler()
	router.NotFound = handlerUI

	go func() {
		if err := http.ListenAndServe(apiAddr, router); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on %v...", apiAddr)
}

func serve(apiAddr string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	setup(apiAddr)

	<-sigChan
	fmt.Println("Received interrupt, shutting down...")
}
