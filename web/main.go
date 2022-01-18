package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"

	"github.com/SiaFoundation/embarcadero/embarcadero"
	"github.com/julienschmidt/httprouter"
)

var (
	// FYI: ../ is not allowed in the embed directive, so dist will need to be copied here.
	//go:embed dist/*
	Assets   embed.FS
	buildDir = "dist"
	mt       *embarcadero.MarketTracker
)

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

// API routes

func bids(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	writeJSON(w, mt.Bids())
}

func trades(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	writeJSON(w, mt.Trades())
}

func fillBid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	skynet, _ := strconv.ParseBool(ps.ByName("skynet"))
	b64, _ := strconv.ParseBool(ps.ByName("b64"))
	embarcadero.FillBid(mt, ps.ByName("bidStr"), skynet, b64)
	fmt.Fprint(w, "Success")
}

func placeBid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	skynet, _ := strconv.ParseBool(ps.ByName("skynet"))
	b64, _ := strconv.ParseBool(ps.ByName("b64"))
	embarcadero.PlaceBid(ps.ByName("inStr"), ps.ByName("outStr"), skynet, b64)
	fmt.Fprint(w, "Success")
}

// Static file server for UI

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func buildUIHandler() http.Handler {
	defaultPath := path.Join(buildDir, "index.html")

	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(buildDir, name)
		f, err := Assets.Open(assetPath)

		if os.IsNotExist(err) {
			// Fallback to index.html
			return Assets.Open(defaultPath)
		}

		return f, err
	})

	return http.FileServer(http.FS(handler))
}

// Server

func Serve(_mt *embarcadero.MarketTracker, apiAddr string) {
	mt = _mt
	router := httprouter.New()

	// API
	router.GET("/api/bids", bids)
	router.POST("/api/bids/fill", fillBid)
	router.POST("/api/bids/place", placeBid)
	router.GET("/api/trades", trades)

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
