package server

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", ":8085", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
)

func Run(ctx context.Context) error {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./client_")))
	//each command has its own websocket handler
	for _, cmnd := range logCommands {
		log.Println(cmnd)
		mux.Handle(cmnd.pattern(), cmnd)
	}
	mux.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}

		cmdText := r.FormValue("cmd")
		if cmdText == "" {
			http.Error(w, "no command to exec", http.StatusBadRequest)
			return
		}
		log.Println("try to parse, cmdText:", cmdText)
		cmd := NewCommandFromString(cmdText)

		c := exec.CommandContext(r.Context(), "bash", "-c", cmd.String())
		c.Stdout = w
		c.Stderr = w
		if err := c.Run(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	//graceful shutdown, when signal from context
	go func() {
		<-ctx.Done()
		log.Println("Shutdown websocket server!")
		if err := server.Shutdown(ctx); err != nil {
			return
		}
	}()

	return server.ListenAndServe()
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)
