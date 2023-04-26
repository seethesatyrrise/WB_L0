package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/pkg/orderModel"
	"http-nats-psql/pkg/storage"
	"os"
)

type Stream struct {
	connection   *nats.Conn
	stream       *nats.JetStreamContext
	subscription *nats.Subscription
}

func (stream *Stream) Unsubscribe() {
	stream.subscription.Unsubscribe()
	stream.connection.Drain()
}

func NewJetStream(storage *storage.Storage) *Stream {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	sub, err := js.SubscribeSync("msg")
	if err != nil {
		panic(err)
	}

	newStream := &Stream{stream: &js,
		subscription: sub,
		connection:   nc,
	}

	go getMessages(newStream, storage)

	return newStream
}

func getMessages(stream *Stream, storage *storage.Storage) {
	for count := 1; ; count++ {
		msg, err := stream.subscription.NextMsgWithContext(context.Background())
		if err != nil {
			fmt.Println("no subscription")
			return
		}
		fmt.Println("got message #", count)

		msg.AckSync()

		msgOut := &orderModel.Order{}
		err = json.Unmarshal(msg.Data, msgOut)
		if err != nil {
			panic(err)
		}

		storage.InsertOrder(msgOut)
	}
}
