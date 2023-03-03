package webterm

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type Webterm struct {
	cmd  *exec.Cmd
	tty  *os.File
	conn *websocket.Conn
}

func New(conn *websocket.Conn, command string, args []string) (*Webterm, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), "TERM=xterm")

	log.Info("Starting TTY")
	tty, err := pty.Start(cmd)
	if err != nil {
		log.Error("Not possible to start PTY", "err", err)
		return nil, err
	}

	wt := &Webterm{cmd, tty, conn}

	return wt, nil
}

func (wt *Webterm) ReadTTY() ([]byte, error) {
	buf := make([]byte, 1024)
	read, err := wt.tty.Read(buf)
	if err != nil {
		log.Error("Unable to read from PTY", "err", err)
		return nil, err
	}

	return buf[:read], nil
}

func (wt *Webterm) WriteWebsocket(message []byte) {
	wt.conn.WriteMessage(websocket.BinaryMessage, message)
}

func (wt *Webterm) ReadWebsocket() {
}

func (wt *Webterm) Close() {
	wt.cmd.Process.Kill()
	wt.cmd.Process.Wait()
	wt.tty.Close()
	wt.conn.Close()
}
