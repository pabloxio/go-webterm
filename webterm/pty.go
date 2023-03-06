package webterm

import "sync"

func (wt *Webterm) TTYWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	defer wt.Close()

	for {
		msg, err := wt.ReadTTY()
		if err != nil {
			wt.logger.Warn("Unable to read from PTY", "err", err)
			break
		}

		wt.logger.Info("Sending bytes from PTY to Websocket")
		wt.WriteWebsocket(msg)
	}

	wt.logger.Info("Stopping TTYWorker")
}
