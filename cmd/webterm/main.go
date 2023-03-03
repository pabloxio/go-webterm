package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/charmbracelet/log"
	"github.com/pabloxio/go-webterm/handlers"
)

func main() {
	address := flag.String("address", "localhost", "Listening address")
	port := flag.Int("port", 8000, "Listening port")
	flag.Parse()

	http.HandleFunc("/healthz", handlers.HealthzHandler)
	http.HandleFunc("/", handlers.WebtermHandler)

	socket := fmt.Sprintf("%s:%d", *address, *port)

	log.Info(fmt.Sprintf("Listening on %s", socket))
	log.Fatal(http.ListenAndServe(socket, nil))
}
