package webterm

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type Webterm struct {
	cmd    *exec.Cmd
	tty    *os.File
	conn   *websocket.Conn
	logger log.Logger
}

func New(conn *websocket.Conn, command string, args []string, logger log.Logger) (*Webterm, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), "TERM=xterm")

	logger.Info("Starting TTY")
	tty, err := pty.Start(cmd)
	if err != nil {
		logger.Error("Not possible to start PTY", "err", err)
		return nil, err
	}

	wt := &Webterm{cmd, tty, conn, logger}

	return wt, nil
}

func (wt *Webterm) ReadTTY() ([]byte, error) {
	buf := make([]byte, 1024)
	read, err := wt.tty.Read(buf)
	if err != nil {
		wt.logger.Error("Unable to read from PTY", "err", err)
		return nil, err
	}

	return buf[:read], nil
}

func (wt *Webterm) WriteWebsocket(message []byte) {
	wt.conn.WriteMessage(websocket.BinaryMessage, message)
}

func (wt *Webterm) ReadWebsocket() ([]byte, error) {
	wsMessageType, message, err := wt.conn.ReadMessage()
	// wsMessageType, reader, err := wt.conn.NextReader()
	if err != nil {
		return nil, err
	}

	if wsMessageType != websocket.BinaryMessage {
		return nil, errors.New("Unexpected message type")
	}

	return message, nil
}

func (wt *Webterm) Close() {
	wt.cmd.Process.Kill()
	wt.cmd.Process.Wait()
	wt.tty.Close()
	wt.conn.Close()
}

func (wt *Webterm) TTYWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	defer wt.Close()

	for {
		msg, err := wt.ReadTTY()
		if err != nil {
			wt.logger.Warn("Unable to read from PTY", "err", err)
			break
		}

		wt.logger.Debug("Sending bytes from PTY to Websocket")
		wt.WriteWebsocket(msg)
	}

	wt.logger.Info("Stopping TTYWorker")
}

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
