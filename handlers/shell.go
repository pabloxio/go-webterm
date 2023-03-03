package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/pabloxio/go-webterm/webterm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebtermHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(fmt.Sprintf("Received connection from: %s", r.RemoteAddr))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Unable to ugrade HTTP connection:", "err", err)
		return
	}

	webterm, err := webterm.New(conn, "bash", []string{"-l"})
	if err != nil {
		log.Error("Unable to create Webterm", "err", err)
		return
	}
	defer webterm.Close()

	var wg sync.WaitGroup
	// pty -> xterm.js
	wg.Add(1)
	go webterm.TTYWorker(&wg)
	// xterm.js > pty
	wg.Add(1)
	go webterm.WebsocketWorker(&wg)

	log.Info("Waiting")
	wg.Wait()
	log.Info("Clossing connection")
}
