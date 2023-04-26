package main

import (
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)
	defer nc.Drain()

	js, _ := nc.JetStream()

	addStream(js, "MSG")
}

func addStream(js nats.JetStreamContext, streamName string) {
	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"msg"},
		MaxMsgs:  100,
		MaxAge:   time.Hour,
	})
}