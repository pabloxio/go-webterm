package webterm

import (
	"encoding/json"
	"sync"

	"github.com/creack/pty"
)

func (wt *Webterm) WebsocketWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	defer wt.Close()

	for {
		message, err := wt.ReadWebsocket()
		if err != nil {
			wt.logger.Warn("Unable to reader from Websocket", "err", err)
			break
		}

		switch message[0] {
		case 0:
			wt.logger.Debug("Copying bytes from Websocket to TTY")
			_, err = wt.tty.Write(message[1:])
			if err != nil {
				wt.logger.Warn("Unable to copy from Websocket to TTY")
				break
			}
		case 1:
			var ptyWinSize pty.Winsize
			err = json.Unmarshal(message[1:], &ptyWinSize)
			if err != nil {
				wt.logger.Error("Unable to read resize:", "err", err)
			}

			wt.logger.Info("TTY Resize", ptyWinSize)
			err = pty.Setsize(wt.tty, &ptyWinSize)
			if err != nil {
				wt.logger.Error("Unable to set TTY size:", "err", err)
			}
		default:
			wt.logger.Warn("Unknow data type")
		}
	}

	wt.logger.Info("Stopping TTYWorker")
}
