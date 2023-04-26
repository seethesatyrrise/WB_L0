package server

import (
	"context"
	"fmt"
	"http-nats-psql/internal"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/jetstream"
	"http-nats-psql/internal/rest"
	"http-nats-psql/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	httpServer *http.Server
	stream     *jetstream.Stream
}

func NewServer(dbcfg *internal.DBConfig, jscfg *internal.JSConfig) (*server, error) {
	db, err := database.NewDatabaseConnection(dbcfg)
	if err != nil {
		return nil, err
	}

	storage := storage.NewStorage()
	storage.RestoreData(db)

	ordersStream := jetstream.NewJetStream(jscfg)

	go jetstream.GetMessages(ordersStream, storage, db)

	ordersRest := rest.NewRest(storage, db)

	router := gin.New()
	api := router.Group("/api")

	ordersRest.Register(api)

	return &server{
		httpServer: &http.Server{
			Addr:    ":" + dbcfg.ServerPort,
			Handler: router,
		},
		stream: ordersStream,
	}, nil
}

func (s *server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	fmt.Println("server started")
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	if s.stream != nil {
		s.stream.Unsubscribe()
	}

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	fmt.Println("server stopped")
	return nil
}
