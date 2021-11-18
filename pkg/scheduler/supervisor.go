package scheduler

import (
	"context"
	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/prometheus/client_golang/prometheus"
	"time"

	domain "github.com/freundallein/scheduler/pkg"
)

// Supervisor implements a domain.Supervisor.
type Supervisor struct {
	taskGateway              domain.Gateway
	staleTasksDeletedCounter prometheus.Counter
}

// NewSupervisor returns a domain.Supervisor implementation.
func NewSupervisor(taskGateway domain.Gateway, opts ...SupervisorOption) *Supervisor {
	svc := &Supervisor{
		taskGateway: taskGateway,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// DeleteStaleTasks cleans storage from stale tasks.
func (svc *Supervisor) DeleteStaleTasks(ctx context.Context, staleHours int) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			rows, err := svc.taskGateway.DeleteStaleTasks(ctx, staleHours)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("supervisor_delete_stale_tasks_failure")
				return err
			}
			svc.staleTasksDeletedCounter.Add(float64(rows))
			log.WithFields(log.Fields{
				"rows": rows,
			}).Debug("supervisor_delete_stale_rows")
		}
	}
}
