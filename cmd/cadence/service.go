package main

import (
	"context"
	stderrors "errors"
	"io"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/nightnoryu/go-kita/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const shutdownTimeout = 10 * time.Second

var errServiceStopped = stderrors.New("service stopped without errors")

func service(ctx context.Context, config *config, logger log.Logger) error {
	router := mux.NewRouter()

	container, err := newDependencyContainer(config, logger, router)
	if err != nil {
		return errors.Wrap(err, "failed to initialize the dependency container")
	}
	defer func() {
		if closeErr := container.Close(); closeErr != nil {
			logger.Error(closeErr, "failed to close dependency container")
		}
	}()

	router.HandleFunc("/resilience/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	})

	router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, r *http.Request) {
		if readyErr := container.Ready(r.Context()); readyErr != nil {
			logger.Error(readyErr, "readiness check failed")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, http.StatusText(http.StatusServiceUnavailable))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	})

	router.Handle("/metrics", promhttp.HandlerFor(container.metrics.Registry, promhttp.HandlerOpts{}))

	httpServer := &http.Server{
		Handler:           router,
		Addr:              config.ServeRESTAddress,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       time.Hour,
		WriteTimeout:      time.Hour,
	}

	// Shutdown must use a fresh context; ctx is canceled by this point - hence the nolints
	go func() { //nolint:gosec
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil { //nolint:contextcheck
			logger.Error(shutdownErr, "failed to gracefully shut down HTTP server")
		}
	}()

	logger.Info("Listening and serving...")
	err = httpServer.ListenAndServe()
	return translateStopErr(err, errServiceStopped)
}

func translateStopErr(from, to error) error {
	switch {
	case errors.Is(from, http.ErrServerClosed):
		return to
	default:
		return from
	}
}
