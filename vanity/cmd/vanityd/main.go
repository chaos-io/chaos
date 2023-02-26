package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/chaos-io/chaos/core/logs"
)

func main() {
	if err := run(); err != nil {
		logs.Fatal("unexpected error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	srv := &http.Server{
		Addr:    flagRunAddr,
		Handler: newRouter(),
	}

	errChan := make(chan error, 1)
	go func() {
		logs.Infow("starting HTTP server", "addr", flagRunAddr)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("cannot run HTTP server: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		logs.Infow("shutting down gracefully", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	case err := <-errChan:
		return err
	}
}

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", handleIndex)
	r.Get("/*", handleGoGet)
	r.Get("/ping", handlePing)

	return r
}
