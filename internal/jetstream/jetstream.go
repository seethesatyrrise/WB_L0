package jetstream

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
)

type Stream struct {
	conn   *nats.Conn
	stream *nats.JetStreamContext
	subscr *nats.Subscription
	subCh  chan *nats.Msg
}

func (s *Stream) Unsubscribe() {
	s.subscr.Unsubscribe()
	s.conn.Drain()
}

func NewJetStream(jscfg *JSConfig) (*Stream, error) {
	nc, err := nats.Connect(jscfg.URL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	ch := make(chan *nats.Msg)
	sub, err := nc.ChanSubscribe("msg", ch)
	if err != nil {
		return nil, err
	}

	newStream := &Stream{
		stream: &js,
		subscr: sub,
		subCh:  ch,
		conn:   nc,
	}

	return newStream, nil
}

func (s *Stream) GetMessages(ctx context.Context, storage *storage.Storage, db *database.DB) {
	count := 1
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.subCh:
			fmt.Printf("got message #%d\n", count)
			count++

			id, err := db.InsertOrder(ctx, msg.Data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = storage.InsertOrder(id, msg.Data)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if err := msg.AckSync(); err != nil {
				continue
			}
		}
	}
}
