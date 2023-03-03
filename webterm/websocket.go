package webterm

import (
	"io"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

func (wt *Webterm) WebsocketWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		messageType, reader, err := wt.conn.NextReader()
		if err != nil {
			log.Warn("Unable to grab next reader", "err", err)
			wt.Close()
			break
		}

		if messageType != websocket.BinaryMessage {
			log.Warn("Unexpected meesage type")
			continue
		}

		dataTypeBuf := make([]byte, 1)
		_, err = reader.Read(dataTypeBuf)
		if err != nil {
			log.Warn("Unable to read MessateType from reader")
			break
		}

		log.Info("Copying bytes from Websocket to TTY")
		_, err = io.Copy(wt.tty, reader)
		if err != nil {
			log.Warn("Unable to copy from Websocket to TTY")
		}
	}

	log.Info("Stopping TTYWorker")
}
