package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTP:", r.RequestURI)
	defer log.Println("exit ServeHTTP:", r.RequestURI)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		} else {
			log.Println("websocket.HandshakeError:", err)
		}
		return
	}
	c.socketWriter(r.Context(), ws)

}

func (c *Command) socketReader(ws *websocket.Conn) {
	defer fmt.Println("exit reader")
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			break
		}
		log.Println("message from client:", string(p))
		if bytes.Equal(p, []byte("stop")) {
			log.Println("signal from client to stop")
			break
		}
	}
}

func (c *Command) socketWriter(ctx context.Context, ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
		fmt.Println("exit from socketWriter")
	}()

	outCh := make(chan string)

	//try to exec command, send logs to outCh
	go func() {
		log.Println("try to start command:", c.pathName())
		defer fmt.Println("exit command:", c.pathName())

		// this commented block of code works only with some commands
		// maybe incorrect redirects or smth. tried different variations
		/*cmd := exec.CommandContext(ctx, c.pathName(), c.args...)
		if c.path != "" && c.path != "./" {
			cmd.Dir = c.path
		}
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}
		go commandLogScanner(ctx, cmdReader, outCh)
		if err := cmd.Start(); err != nil {
			log.Println(err)
			return
		}
		if err := cmd.Wait(); err != nil {
			log.Println(err)
		}*/

		// more low-level and complicated solution that works with all commands we need
		outr, outw, err := os.Pipe()
		if err != nil {
			internalError(ws, "stdout:", err)
			return
		}
		defer outr.Close()
		defer outw.Close()

		inr, inw, err := os.Pipe()
		if err != nil {
			internalError(ws, "stdin:", err)
			return
		}
		defer inr.Close()
		defer inw.Close()

		_, err = exec.LookPath(c.pathName())
		if err != nil {
			log.Println(err)
			return
		}

		cmdPath := c.ProcessName()
		if c.path == "" {
			if cmdPath, err = exec.LookPath(c.name); err != nil {
				log.Println(err)
				return
			}
		}
		args := append([]string{cmdPath}, c.args...)

		proc, err := os.StartProcess(args[0], args, &os.ProcAttr{
			Dir:   c.ProcessDir(),
			Files: []*os.File{inr, outw, outw},
		})
		if err != nil {
			internalError(ws, "start:", err)
			return
		}
		inr.Close()
		outw.Close()

		go commandLogScanner(ctx, outr, outCh)

		c.socketReader(ws)

		inw.Close()

		if err := proc.Signal(os.Interrupt); err != nil {
			log.Println("inter:", err)
		}

		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
			// A bigger bonk on the head.
			if err := proc.Signal(os.Kill); err != nil {
				log.Println("term:", err)
			}
		}

		if _, err := proc.Wait(); err != nil {
			log.Println("wait:", err)
		}
	}()

	// try to read command logs from outCh and send to websocket
	// periodically ping websocket connection
	for {
		select {
		case line := <-outCh /*tf.Lines*/ :
			p, err := json.Marshal(line)
			if err != nil {
				log.Println(err)
				break
			}
			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println("ws.SetWriteDeadline in line := <-outCh in command socketWriter:", err)
			}
			if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
				return
			}
		case <-pingTicker.C:
			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println("ws.SetWriteDeadline in pingTicker.C in command socketWriter:", err)
			}
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-ctx.Done():
			fmt.Println(ctx.Err(), "ctx.Done in writer")
			return
		}
	}
}

func commandLogScanner(ctx context.Context, rc io.Reader, out chan string) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println("exit commandLogScanner", ctx.Err())
			return
		case out <- scanner.Text():

			/*default:
			fmt.Println(scanner.Text())
			out <- scanner.Text()*/
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}
