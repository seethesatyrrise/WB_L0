package main

import (
	"context"
	"fmt"
	"http-nats-psql/internal"
	"http-nats-psql/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbcfg, err := internal.GetDBConfig()
	if err != nil {
		panic(err)
	}
	jscfg := internal.GetJSConfig()

	srv, err := server.NewServer(dbcfg, jscfg)
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
