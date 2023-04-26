package main

import (
	"context"
	"fmt"
	"http-nats-psql/pkg/database"
	"http-nats-psql/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := database.GetConfig()
	if err != nil {
		panic(err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		fmt.Println("can't create server")
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	if err = srv.Start(); err != nil {
		fmt.Println("can't start server")
	}

	<-signalChan
	if err = srv.Shutdown(context.Background()); err != nil {
		fmt.Println("error shutting server")
	}
}
