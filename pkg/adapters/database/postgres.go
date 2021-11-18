package database

import (
	"context"
	"encoding/json"
	"time"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TaskGateway used for access to task database layer.
type TaskGateway struct {
	pool *pgxpool.Pool
}

// NewTaskGateway return task gateway implementation.
func NewTaskGateway(dsn string) (domain.Gateway, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	return &TaskGateway{pool: pool}, nil
}

//Create is for a task creation, returns a created task.
func (gw *TaskGateway) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	row := gw.pool.QueryRow(ctx, create, task.ID, task.ExecuteAt, task.Deadline, task.Payload, task.Meta)
	err := row.Scan(
		&task.ID,
		&task.ClaimID,
		&task.State,
		&task.ExecuteAt,
		&task.Deadline,
		&task.Payload,
		&task.Result,
		&task.Meta,
	)
	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "task_pkey" (SQLSTATE 23505)` {
			return nil, domain.Error{Code: domain.ErrDuplicateTask, Inner: err, Message: "task already set"}
		}
		return nil, err
	}
	return task, nil
}

// FindByID returns a task by id.
func (gw *TaskGateway) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task := &domain.Task{}
	row := gw.pool.QueryRow(ctx, findByID, id)
	err := row.Scan(
		&task.ID,
		&task.ClaimID,
		&task.State,
		&task.ExecuteAt,
		&task.Deadline,
		&task.Payload,
		&task.Result,
		&task.Meta,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.Error{Code: domain.ErrTaskNotFound, Message: "task not found"}
		}
		return nil, err
	}
	return task, nil
}

// ClaimPending locks and returns pending (or next-attempt failed) task.
func (gw *TaskGateway) ClaimPending(ctx context.Context, amount int) ([]*domain.Task, error) {
	tasks := make([]*domain.Task, 0)
	rows, err := gw.pool.Query(ctx, claimPending, domain.StatePending, domain.StateProcessing, amount)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.ClaimID,
			&task.State,
			&task.ExecuteAt,
			&task.Deadline,
			&task.Payload,
			&task.Result,
			&task.Meta,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		return nil, domain.Error{Code: domain.ErrNoPendingTasks, Message: "no pending tasks"}
	}
	return tasks, nil
}

// MarkAsSucceeded marks a task as succefully processed.
func (gw *TaskGateway) MarkAsSucceeded(ctx context.Context, id, claimID uuid.UUID, result map[string]interface{}) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	tag, err := gw.pool.Exec(ctx, markAsSucceeded, domain.StateSucceeded, id, claimID, string(resultJSON))
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.Error{Code: domain.ErrStaleResult, Message: "result is stale"}
	}
	return nil
}

// MarkAsFailed marks a task as failed.
func (gw *TaskGateway) MarkAsFailed(ctx context.Context, id, claimID uuid.UUID, reason string) error {
	tag, err := gw.pool.Exec(ctx, markAsFailed, domain.StateFailed, id, claimID, map[string]string{"failReason": reason})
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.Error{Code: domain.ErrStaleResult, Message: "result is stale"}
	}
	// TODO: check attempt, send to failure table and delete from task table?
	return nil
}
