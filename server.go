package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

func writeError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, err, code)
}

type createRequest struct {
	Offer   string `json:"offer"`
	Receive string `json:"receive"`
}

type createResponse struct {
	Swap SwapTransaction `json:"swap"`
}

func createHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cr createRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.Contains(cr.Offer, "SF") == strings.Contains(cr.Receive, "SF") {
		writeError(w, "Invalid swap: must specify one SC value and one SF value", http.StatusBadRequest)
		return
	}
	input, output := ParseCurrency(cr.Offer), ParseCurrency(cr.Receive)
	swap, err := CreateSwap(input, output, strings.Contains(cr.Offer, "SF"))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, createResponse{
		Swap: swap,
	})
}

type acceptRequest struct {
	Swap SwapTransaction `json:"swap"`
}

type acceptResponse struct {
	ID   string          `json:"id"`
	Swap SwapTransaction `json:"swap"`
}

func acceptHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ar acceptRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := CheckAccept(ar.Swap); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := AcceptSwap(&ar.Swap); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, acceptResponse{
		ID:   ar.Swap.Transaction().ID().String(),
		Swap: ar.Swap,
	})
}

type finishRequest struct {
	Swap SwapTransaction `json:"swap"`
}

type finishResponse struct {
	ID   string          `json:"id"`
	Swap SwapTransaction `json:"swap"`
}

func finishHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var fr finishRequest
	if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := CheckFinish(fr.Swap); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := FinishSwap(&fr.Swap); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, finishResponse{
		ID:   fr.Swap.Transaction().ID().String(),
		Swap: fr.Swap,
	})
}

type summarizeRequest struct {
	Swap SwapTransaction `json:"swap"`
}

type summarizeResponse struct {
	ID      string      `json:"id"`
	Summary SwapSummary `json:"summary"`
}

func summarizeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var fr summarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	summary, err := Summarize(fr.Swap)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, summarizeResponse{
		ID:      fr.Swap.Transaction().ID().String(),
		Summary: summary,
	})
}

func walletHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := siad.WalletGet()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c)
}

func consensusHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := siad.ConsensusGet()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c)
}

func serve(addr string, dev bool) {
	api := httprouter.New()
	api.POST("/api/create", createHandler)
	api.POST("/api/accept", acceptHandler)
	api.POST("/api/finish", finishHandler)
	api.POST("/api/summarize", summarizeHandler)
	api.GET("/api/wallet", walletHandler)
	api.GET("/api/consensus", consensusHandler)

	go func() {
		ui := buildUIHandler()
		err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				if r.Header.Get("Access-Control-Request-Method") != "" {
					w.Header().Set("Access-Control-Allow-Headers", "content-type")
					w.Header().Set("Access-Control-Allow-Methods", w.Header().Get("Allow"))
					// CORS, necessary for development only
					if dev {
						w.Header().Set("Access-Control-Allow-Origin", "*")
					}
				}
				w.WriteHeader(http.StatusNoContent)
				return
			} else if strings.HasPrefix(r.URL.Path, "/api") {
				api.ServeHTTP(w, r)
				return
			}
			ui.ServeHTTP(w, r)
		}))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on %v...", addr)

	if !dev {
		if err := open(addr); err != nil {
			log.Println("Warning: failed to automatically open web UI:", err)
			log.Println("Please navigate to the above URL in your browser.")
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	fmt.Println("Received interrupt, shutting down...")
}

func open(url string) error {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	switch runtime.GOOS {
	case "windows":
		_, err := exec.Command("cmd", "/c", "start", url).CombinedOutput()
		return err
	case "darwin":
		_, err := exec.Command("open", url).CombinedOutput()
		return err
	default: // linux, bsd
		_, err := exec.Command("xdg-open", url).CombinedOutput()
		return err
	}
}
