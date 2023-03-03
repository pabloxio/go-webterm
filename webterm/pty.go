package webterm

import (
	"sync"

	"github.com/charmbracelet/log"
)

func (wt *Webterm) TTYWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, err := wt.ReadTTY()
		if err != nil {
			log.Warn("Unable to read from PTY", "err", err)
			break
		}

		log.Info("Sending bytes from PTY to Websocket")
		wt.WriteWebsocket(msg)
	}

	log.Info("Stopping TTYWorker")
}
