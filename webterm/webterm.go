package webterm

import (
	"errors"
	"io"
	"os"
	"os/exec"

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

func (wt *Webterm) ReadWebsocket() (*io.Reader, error) {
	messageType, reader, err := wt.conn.NextReader()
	if err != nil {
		return nil, err
	}

	if messageType != websocket.BinaryMessage {
		return nil, errors.New("Unexpected message type")
	}

	dataTypeBuf := make([]byte, 1)
	_, err = reader.Read(dataTypeBuf)
	if err != nil {
		return nil, err
	}

	return &reader, nil
}

func (wt *Webterm) Close() {
	wt.cmd.Process.Kill()
	wt.cmd.Process.Wait()
	wt.tty.Close()
	wt.conn.Close()
}
