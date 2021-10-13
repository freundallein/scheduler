package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/oklog/run"

	"github.com/freundallein/scheduler/pkg/adapters/apiserv"

	"github.com/freundallein/scheduler/pkg/utils"
	"github.com/freundallein/scheduler/pkg/utils/opsserv"
)

const (
	logLevelKey = "LOG_LEVEL"
	opsPortKey  = "OPS_PORT"
	apiPortKey  = "API_PORT"
	pgDSNKey    = "PG_DSN"
)

func main() {
	logLevel := utils.GetEnv(logLevelKey, "debug")
	log.Init("scheduler", logLevel)
	log.Info("init_service")
	apiPort := utils.GetEnv(apiPortKey, "8000")
	opsPort := utils.GetEnv(opsPortKey, "8001")

	apiService := apiserv.New(
		apiserv.WithPort(apiPort),
	)
	opsService := opsserv.New(
		opsserv.WithPort(opsPort),
	)

	ctx, cancel := context.WithCancel(context.Background())
	var g run.Group
	{
		g.Add(func() error {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				return nil
			case s := <-sig:
				return fmt.Errorf("sig_recv %v", s)
			}
		}, func(err error) {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("service_interrupted")
			cancel()
		})
	}
	{
		g.Add(func() error {
			return opsService.Run(ctx)
		}, func(err error) {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("ops_svc_interrupted")
			if err := opsService.Shutdown(ctx); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Info("ops_svc_shutdown")
			}
		})
	}
	{
		g.Add(func() error {
			return apiService.Run(ctx)
		}, func(err error) {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("api_svc_interrupted")
			if err := apiService.Shutdown(ctx); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Info("api_svc_shutdown")
			}
		})
	}

	err := g.Run()
	log.WithFields(log.Fields{
		"err": err,
	}).Info("service_stopped")

}