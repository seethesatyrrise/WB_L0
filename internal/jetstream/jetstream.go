package jetstream

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/storage"
	"http-nats-psql/internal/utils"
)

type Stream struct {
	conn   *nats.Conn
	stream nats.JetStreamContext
	subscr *nats.Subscription
	subCh  chan *nats.Msg
}

func (s *Stream) Unsubscribe() {
	s.subscr.Unsubscribe()
	//s.stream.DeleteConsumer("MSG", "server")
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

	js.ConsumerInfo(jscfg.StreamName, jscfg.ConsumerName)
	//js.AddConsumer(streamName, &nats.ConsumerConfig{
	//	Durable:       consumerName,
	//	DeliverPolicy: nats.DeliverAllPolicy,
	//	//AckPolicy: nats.AckExplicitPolicy,
	//})

	ch := make(chan *nats.Msg)
	sub, err := nc.ChanSubscribe(jscfg.SubjectName, ch)
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

			if err := msg.AckSync(); err != nil {
				continue
			}
		}
	}
}
