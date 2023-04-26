package jetstream

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
	"http-nats-psql/internal/utils"
	"time"
)

type Stream struct {
	conn   *nats.Conn
	stream nats.JetStreamContext
	subscr *nats.Subscription
	subCh  chan *nats.Msg
}

func (s *Stream) Drain() {
	s.conn.Drain()
}

func NewJetStream(jscfg *JSConfig) (*Stream, error) {
	nc, err := nats.Connect(jscfg.URL, nats.Name("http-nats-pgsql"))
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.ConsumerInfo(jscfg.StreamName, jscfg.ConsumerName)
	if err != nil {
		js.AddConsumer(jscfg.StreamName, &nats.ConsumerConfig{
			Durable:        jscfg.ConsumerName,
			DeliverSubject: jscfg.SubjectName,
			DeliverPolicy:  nats.DeliverAllPolicy,
			AckPolicy:      nats.AckExplicitPolicy,
			AckWait:        time.Second * 10,
		})
	}

	ch := make(chan *nats.Msg, 5)
	sub, err := js.ChanSubscribe(jscfg.SubjectName, ch, nats.Durable(jscfg.ConsumerName), nats.AckExplicit())

	if err != nil {
		return nil, err
	}

	newStream := &Stream{
		stream: js,
		subscr: sub,
		subCh:  ch,
		conn:   nc,
	}

	return newStream, nil
}

func (s *Stream) GetMessages(ctx context.Context, storage *storage.Storage, db *database.DB, v *utils.Validator) {
	defer close(s.subCh)
	count := 1
	for {

		select {
		case <-ctx.Done():
			return
		case msg := <-s.subCh:
			fmt.Printf("got message #%d\n", count)
			count++

			valid, err := v.Validate(msg.Data)
			if err != nil {
				if err := msg.Ack(); err != nil {
					utils.Logger.Error(err.Error())
					continue
				}
				utils.Logger.Error(err.Error())
				continue
			}
			if !valid {
				utils.Logger.Error("invalid message")
				continue
			}

			id, err := db.InsertOrder(ctx, msg.Data)
			if err != nil {
				utils.Logger.Error(err.Error())
				continue
			}
			err = storage.InsertOrder(id, msg.Data)
			if err != nil {
				utils.Logger.Error(err.Error())
				continue
			}

			if err := msg.Ack(); err != nil {
				utils.Logger.Error(err.Error())
			}
		}
	}
}
