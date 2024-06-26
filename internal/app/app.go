package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/quolpr/booking/internal/transaction"
)

type App struct {
	ServiceProvider *serviceProvider
	config          *config
}

func NewApp(ctx context.Context) *App {
	config := newConfig()

	return &App{
		ServiceProvider: newServiceProvider(ctx, config),
		config:          config,
	}
}

func (app *App) Logger() *slog.Logger {
	return app.ServiceProvider.Logger
}

func (app *App) ServeHTTP(ctx context.Context) error {
	logger := app.ServiceProvider.Logger

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      newRoutes(app.ServiceProvider),
		BaseContext: func(listener net.Listener) context.Context {
			ctx := context.Background()
			ctx = transaction.WithTransactionManager(ctx, app.ServiceProvider.TrManager)
			return ctx
		},
	}

	serverErrorCh := make(chan error)
	go func() {
		defer close(serverErrorCh)

		logger.Info("Server started", "port", app.config.port)
		err := httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server closing", "error", err)
		}

		select {
		case serverErrorCh <- err:
			return
		default:
			return
		}
	}()

	select {
	case err := <-serverErrorCh:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		//nolint: contextcheck
		return httpServer.Shutdown(ctx)
	}
}
