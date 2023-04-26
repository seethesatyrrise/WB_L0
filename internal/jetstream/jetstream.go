package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/internal"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/models"
	"http-nats-psql/internal/storage"
)

type Stream struct {
	conn      *nats.Conn
	stream    *nats.JetStreamContext
	subscr    *nats.Subscription
	subCh     chan *nats.Msg
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (stream *Stream) Unsubscribe() {
	stream.ctxCancel()
	stream.subscr.Unsubscribe()
	stream.conn.Drain()
}

func NewJetStream(jscfg *internal.JSConfig) *Stream {
	nc, err := nats.Connect(jscfg.URL)
	if err != nil {
		panic(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	ch := make(chan *nats.Msg)
	sub, err := nc.ChanSubscribe("msg", ch)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	newStream := &Stream{
		stream:    &js,
		subscr:    sub,
		subCh:     ch,
		conn:      nc,
		ctx:       ctx,
		ctxCancel: cancel,
	}

	return newStream
}

func GetMessages(stream *Stream, storage *storage.Storage, db *database.DB) {
	count := 1
	for {
		select {
		case <-stream.ctx.Done():
			return
		case msg := <-stream.subCh:
			fmt.Printf("got message #%d\n", count)
			count++

			order := &models.Order{}
			err := json.Unmarshal(msg.Data, order)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = db.InsertOrder(order)
			if err != nil {
				fmt.Println(err)
			}
			err = storage.InsertOrder(order)

			msg.AckSync()
		}
	}
}
