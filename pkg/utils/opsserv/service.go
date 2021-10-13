package opsserv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Service used as an endpoint for operations management.
type Service struct {
	httpserv *http.Server
	Port     string
}

// New returns service instance
func New(opts ...Option) *Service {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}
	mux := http.NewServeMux()
	mux.Handle(
		"/ops/metrics",
		promhttp.Handler(),
	)
	mux.HandleFunc(
		"/ops/healthcheck",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	)
	addr := fmt.Sprintf("0.0.0.0:%s", svc.Port)
	svc.httpserv = &http.Server{
		Handler:           mux,
		Addr:              addr,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	return svc
}

// Run starts the ops http server.
func (svc *Service) Run() error {
	log.WithFields(log.Fields{
		"addr": svc.httpserv.Addr,
	}).Info("ops_svc_starting")
	return svc.httpserv.ListenAndServe()
}

// Shutdown provides graceful shutdown of the ops http server.
func (svc *Service) Shutdown(ctx context.Context) error {
	return svc.httpserv.Shutdown(ctx)
}
