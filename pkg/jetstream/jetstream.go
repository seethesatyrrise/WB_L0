package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/nats-io/nats.go"
	"http-nats-psql/pkg/orderModel"
	"os"
)

func RunJetStream(db *pg.DB) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	streamName := "MSG"

	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"msg"},
	})

	sub, err := js.SubscribeSync("msg", nats.AckExplicit())
	if err != nil {
		panic(err)
	}

	ch := make(chan struct{})
	defer close(ch)
	go getMessages(sub, db, ch)
	<-ch
}

func getMessages(sub *nats.Subscription, db *pg.DB, ch chan struct{}) {
	defer func() { ch <- struct{}{} }()
	for {
		msg, err := sub.NextMsgWithContext(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println("got message")

		msg.Ack()

		msgOut := &orderModel.Order{}
		err = json.Unmarshal(msg.Data, msgOut)
		if err != nil {
			panic(err)
		}

		insertOrder(msgOut, db)
	}
}

func insertOrder(order *orderModel.Order, db *pg.DB) error {
	var res interface{}
	_, err := db.Query(&res,
		`select insert_data(?);`,
		order)

	if err != nil {
		return err
	}

	return nil
}
