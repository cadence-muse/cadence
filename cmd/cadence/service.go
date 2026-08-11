package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"cadence/pkg/common/log"
)

var errServiceStopped = errors.New("service stopped without errors")

func service(_ context.Context, config *config, logger log.Logger) error {
	router := mux.NewRouter()

	router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	})

	httpServer := &http.Server{
		Handler:           router,
		Addr:              config.ServeRESTAddress,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       time.Hour,
		WriteTimeout:      time.Hour,
	}
	logger.Info("Listening and serving...")
	err := httpServer.ListenAndServe()
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
