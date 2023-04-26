package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"http-nats-psql/internal/database"
	"http-nats-psql/internal/jetstream"
	"http-nats-psql/internal/rest"
	"http-nats-psql/internal/server"
	"http-nats-psql/internal/storage"
)

type App struct {
	db      *database.DB
	storage *storage.Storage
	js      *jetstream.Stream
	router  *gin.Engine
	server  *server.Server
	cfg     *AppConfig
}

func New(ctx context.Context) (app *App, err error) {
	app = &App{}
	app.cfg, err = newConfig()
	if err != nil {
		return nil, err
	}

	app.db, err = database.NewDatabaseConnection(&app.cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	app.storage = storage.NewStorage()
	if err := app.storage.RestoreData(ctx, app.db); err != nil {
		return nil, err
	}

	app.js, err = jetstream.NewJetStream(&app.cfg.JSConfig)
	if err != nil {
		return nil, err
	}

	rest := rest.NewRest(app.storage, app.db)

	router := gin.New()
	api := router.Group("/api")

	rest.Register(api)

	app.router = router

	app.server, err = server.NewServer(&app.cfg.ServerConfig, app.router.Handler())
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) (err error) {
	go app.js.GetMessages(ctx, app.storage, app.db)

	if err := app.server.Start(); err != nil {
		return err
	}

	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if app.server != nil {
			if err := app.server.Shutdown(gCtx); err != nil {
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if app.js != nil {
			app.js.Unsubscribe()
		}
		return nil
	})

	return g.Wait()
}
