package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheshenia/nanostarter/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Server terminated!")
		cancel()
	}()

	if err := server.Run(ctx); err != nil {
		log.Println(err)
	}
}
