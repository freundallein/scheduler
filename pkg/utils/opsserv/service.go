package opsserv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Options - httpserv parameters
type Options struct {
	Port string
}

// Service used as an endpoint for operations management.
type Service struct {
	options *Options
	Srv     *http.Server
}

// New returns service instance
func New(options *Options) *Service {
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
	addr := fmt.Sprintf("0.0.0.0:%s", options.Port)
	srv := &http.Server{
		Handler:           mux,
		Addr:              addr,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	return &Service{
		options: options,
		Srv:     srv,
	}
}

// Run starts ops http server
func (svc *Service) Run() error {
	return svc.Srv.ListenAndServe()
}
