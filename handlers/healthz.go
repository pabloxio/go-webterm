package handlers

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	log.Info(fmt.Sprintf("Healthz 200 OK"))
}
