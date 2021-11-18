package apiserv

import (
	"context"
	"fmt"
	"github.com/freundallein/scheduler/pkg/scheduler"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/freundallein/scheduler/pkg/utils/logging"
)

// HttpConn adapts HTTP connection to ReadWriteCloser
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
	// HTTP Serve port
	Port string
	// Token is for request authorization
	Token string
	// WorkerToken is for worker's request authorization
	WorkerToken string
}

// New returns service instance
func New(service *scheduler.Service, opts ...Option) *Service {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}
	rpcServer := rpc.NewServer()
	rpcServer.Register(&Scheduler{
		svc: service,
	})
	rpcServer.Register(&Worker{
		svc: service,
	})
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/rpc/v0",
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Auth")
			if svc.Token != "" && auth != svc.Token {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - not authorized"))
				return
			}
			log.WithFields(log.Fields{
				"auth": auth,
			}).Debug("authentication_passed")
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
	mux.HandleFunc(
		"/worker/v0",
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Auth")
			if svc.WorkerToken != "" && auth != svc.WorkerToken {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - not authorized"))
				return
			}
			log.WithFields(log.Fields{
				"auth": auth,
			}).Debug("authentication_passed")
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
