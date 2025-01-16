package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/cyber/test-project/clients"
	"github.com/cyber/test-project/config"
	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/services"
)

type Application struct {
	Config       config.Config
	server       *http.Server
	shutdownOnce sync.Once
}

func InitApplication(configPath string) (*Application, error) {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	err = logging.Init(cfg.Logging)
	if err != nil {
		return nil, err
	}

	return &Application{
		Config: cfg,
	}, nil
}

func (app *Application) Start() error {
	log.Println("Starting application")

	baseHttpClient := http.Client{}

	asanaClientOptions := clients.ClientOptions{
		ServiceName: "asana",
		BaseClient:  &baseHttpClient,
		BaseURL:     app.Config.Asana.BaseURL,
	}
	asanaClient := clients.NewAsanaClient(asanaClientOptions)
	asanaService := services.NewAsanaUsersService(asanaClient, app.Config.Asana.AccessToken)

	routerConfig := RouterConfig{
		AsanaService: asanaService,
	}

	router, err := NewRouter(routerConfig)
	if err != nil {
		return err
	}

	app.server = &http.Server{
		Addr:    app.Config.Http.Addr,
		Handler: router,
	}

	listener, err := net.Listen("tcp", app.server.Addr)
	if err == nil {
		logging.Logger.Info("service started at", zap.String("address", app.server.Addr))
		err = app.server.Serve(listener)
	}

	if errors.Is(err, http.ErrServerClosed) {
		logging.Logger.Info("HTTP service stopped")
		return nil
	}

	if err != nil {
		logging.Logger.Error("HTTP service failed", zap.Error(err))
	}

	return err
}

func (app *Application) Shutdown() {
	app.shutdownOnce.Do(func() {
		defer func() {
			err := logging.Logger.Sync()
			if err != nil {
				logging.Logger.Error("error calling logger.Sync()", zap.Error(err))
			}
		}()

		c := make(chan struct{})
		go func() {
			app.stop()
			close(c)
		}()

		select {
		case <-c:
			logging.Logger.Info("application shutdown successfully complete")
		case <-time.After(app.Config.ShutdownTimeout):
			logging.Logger.Info("could not shutdown application in", zap.Duration("shutdown_timeout", app.Config.ShutdownTimeout))
		}
	})
}

func (app *Application) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), app.Config.ShutdownTimeout)

	logging.WithLogger(ctx, logging.Logger)

	defer cancel()
}
