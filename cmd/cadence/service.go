package main

import (
	"context"
	stderrors "errors"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/nightnoryu/go-kita/health"
	"github.com/nightnoryu/go-kita/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const shutdownTimeout = 10 * time.Second

var errServiceStopped = stderrors.New("service stopped without errors")

func service(ctx context.Context, config *config, logger log.Logger) error {
	router := mux.NewRouter()

	container, err := newDependencyContainer(ctx, config, logger, router)
	if err != nil {
		return errors.Wrap(err, "failed to initialize the dependency container")
	}
	defer func() {
		if closeErr := container.Close(); closeErr != nil {
			logger.Error(closeErr, "failed to close dependency container")
		}
	}()

	livenessHandler, err := health.NewLivenessHandler(health.LivenessConfig{})
	if err != nil {
		return errors.Wrap(err, "failed to create liveness handler")
	}
	readinessHandler, err := health.NewReadinessHandler(health.ReadinessConfig{
		Checks: []health.NamedCheck{
			{Name: "postgres", Check: container.checkDatabase},
			{Name: "redis", Check: container.checkRedis},
		},
		OnFailure: func(name string, err error) {
			logger.WithFields(log.Fields{"dependency": name}).Error(err, "readiness check failed")
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create readiness handler")
	}
	router.Handle("/healthz", livenessHandler).Methods(http.MethodGet)
	router.Handle("/readyz", readinessHandler).Methods(http.MethodGet)

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
