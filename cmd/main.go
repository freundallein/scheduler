package main

import (
	"context"
	"fmt"
	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	prometheusNamespace = "scheduler"
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

	tasksEnqueued := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "scheduler",
		Name:      "tasks_enqueued_total",
		Help:      "The total number of enqueued tasks.",
	})
	taskRequestPolled := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "scheduler",
		Name:      "tasks_polled_total",
		Help:      "The total number of task polling requests.",
	})
	tasksClaimed := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "worker",
		Name:      "tasks_claimed_total",
		Help:      "The total number of claimed tasks.",
	})
	tasksSucceeded := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "worker",
		Name:      "tasks_succeeded_total",
		Help:      "The total number of succeeded tasks.",
	})
	tasksFailed := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "worker",
		Name:      "tasks_failed_total",
		Help:      "The total number of failed tasks.",
	})
	service := scheduler.New(
		gateway,
		scheduler.WithTasksEnqueued(tasksEnqueued),
		scheduler.WithTaskRequestPolled(taskRequestPolled),
		scheduler.WithTasksClaimed(tasksClaimed),
		scheduler.WithTasksSucceeded(tasksSucceeded),
		scheduler.WithTasksFailed(tasksFailed),
	)

	staleTasksDeleted := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Subsystem: "supervisor",
		Name:      "stale_tasks_deleted",
		Help:      "The total number of deleted stale tasks.",
	})
	supervisor := scheduler.NewSupervisor(
		gateway,
		scheduler.WithStaleTasksDeleted(staleTasksDeleted),
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
