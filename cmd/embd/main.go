package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/SiaFoundation/embarcadero/embarcadero"
	"github.com/SiaFoundation/embarcadero/web"
)

func main() {
	log.SetFlags(0)
	apiAddr := flag.String("p", "localhost:8080", "host:port to serve the embarcadero API on")
	siadAddr := flag.String("siad", "localhost:9980", "host:port that the siad API is listening on")
	dir := flag.String("d", ".", "directory where server state will be stored")
	flag.Parse()

	mt, err := embarcadero.Start(*dir, *siadAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer mt.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	web.Serve(mt, *apiAddr)
	<-sigChan
	fmt.Println("Received interrupt, shutting down...")
}
