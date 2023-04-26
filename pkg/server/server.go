package server

import (
	"context"
	"fmt"
	"http-nats-psql/pkg/database"
	"http-nats-psql/pkg/jetstream"
	"http-nats-psql/pkg/repo"
	"http-nats-psql/pkg/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	httpServer *http.Server
	stream     *jetstream.Stream
}

func NewServer(cfg *database.Config) (*server, error) {
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		return nil, err
	}

	stream := jetstream.NewJetStream(db)

	ordersRepo := repo.NewRepo(db)
	ordersRest := rest.NewRest(ordersRepo)

	router := gin.New()
	api := router.Group("/api")

	ordersRest.Register(api)

	return &server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: router,
		},
		stream: stream,
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
