package apiserv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/freundallein/scheduler/pkg/utils/logging"
)

// adapt HTTP connection to ReadWriteCloser
type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

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
	rpcServer := rpc.NewServer()
	rpcServer.Register(&Scheduler{})
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/rpc/v0",
		func(w http.ResponseWriter, r *http.Request) {
			serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
			err := rpcServer.ServeRequest(serverCodec)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("err_while_serving_json_rpc")
				return
			}
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
func (svc *Service) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"addr": svc.httpserv.Addr,
	}).Info("api_svc_starting")
	return svc.httpserv.ListenAndServe()
}

// Shutdown provides graceful shutdown of the ops http server.
func (svc *Service) Shutdown(ctx context.Context) error {
	return svc.httpserv.Shutdown(ctx)
}
