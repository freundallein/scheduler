package main

import (
	"context"
	"fmt"
	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/oklog/run"
	"os"
	"os/signal"
	"syscall"

	"github.com/freundallein/scheduler/pkg/adapters/apiserv"
	"github.com/freundallein/scheduler/pkg/adapters/database"

	"github.com/freundallein/scheduler/pkg/scheduler"
	"github.com/freundallein/scheduler/pkg/utils"
	"github.com/freundallein/scheduler/pkg/utils/opsserv"
)

const (
	logLevelKey    = "LOG_LEVEL"
	opsPortKey     = "OPS_PORT"
	apiPortKey     = "API_PORT"
	databaseDSNKey = "DB_DSN"
	tokenKey       = "TOKEN"
	workerTokenKey = "WORKER_TOKEN"
	staleHoursKey  = "STALE_HOURS"
)

func main() {
	logLevel := utils.GetEnv(logLevelKey, "debug")
	log.Init("scheduler", logLevel)
	log.Info("init_service")
	apiPort := utils.GetEnv(apiPortKey, "8000")
	opsPort := utils.GetEnv(opsPortKey, "8001")
	token := utils.GetEnv(tokenKey, "token")
	workerToken := utils.GetEnv(workerTokenKey, "token")
	staleHours, err := utils.GetIntEnv(staleHoursKey, 24*7)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("stale_period_env_failure")
	}

	databaseDSN := utils.GetEnv(databaseDSNKey, "postgres://scheduler:scheduler@0.0.0.0:5432/scheduler")

	gateway, err := database.NewTaskGateway(databaseDSN)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("db_gateway_creation_failure")
		os.Exit(1)
	}

	service := scheduler.New(
		gateway,
	)
	supervisor := scheduler.NewSupervisor(
		gateway,
	)

	apiService := apiserv.New(
		service,
		apiserv.WithToken(token),
		apiserv.WithWorkerToken(workerToken),
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
				return fmt.Errorf("signal_recv %v", s)
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
	{
		g.Add(func() error {
			return supervisor.DeleteStaleTasks(ctx, staleHours)
		}, func(err error) {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("supervisor_interrupted")
		})
	}

	err = g.Run()
	log.WithFields(log.Fields{
		"err": err,
	}).Info("service_stopped")

}
