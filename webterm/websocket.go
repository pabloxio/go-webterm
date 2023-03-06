package webterm

import (
	"io"
	"sync"
)

func (wt *Webterm) WebsocketWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	defer wt.Close()

	for {
		reader, err := wt.ReadWebsocket()
		if err != nil {
			wt.logger.Warn("Unable to reader from Websocket", "err", err)
			break
		}

		wt.logger.Info("Copying bytes from Websocket to TTY")
		_, err = io.Copy(wt.tty, *reader)
		if err != nil {
			wt.logger.Warn("Unable to copy from Websocket to TTY")
			break
		}
	}

	wt.logger.Info("Stopping TTYWorker")
}
