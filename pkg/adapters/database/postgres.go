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
	query := `insert into task(id, execute_at, deadline, payload, meta) values ($1, $2, $3, $4, $5) returning *;`
	row := gw.pool.QueryRow(ctx, query, task.ID, task.ExecuteAt, task.Deadline, task.Payload, task.Meta)
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
	query := `select * from task where id=$1;`
	task := &domain.Task{}
	row := gw.pool.QueryRow(ctx, query, id)
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
	query := `
	with claimed_tasks as (
		select 
			id 
		from task 
		where 
			state <> 'succeeded' 
			and execute_at <= localtimestamp
		order by execute_at
		limit $1
		for update skip locked
	)
	update task 
	set 
		state = 'processing', 
		execute_at = localtimestamp + interval '1 minute',
		claim_id = uuid_generate_v4()
	from claimed_tasks
	where task.id = claimed_tasks.id
	returning task.*;
	`
	tasks := []*domain.Task{}
	rows, err := gw.pool.Query(ctx, query, amount)
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
	query := `
	update task
	set 
		state = 'succeeded',
		claim_id = null,
		result = $3
	where 
		id = $1
		and claim_id = $2;
	`
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	tag, err := gw.pool.Exec(ctx, query, id, claimID, string(resultJSON))
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
	query := `
	update task
	set 
		state = 'failed',
		claim_id = null,
		meta = meta::jsonb || $3 || CONCAT('{"attempts":', COALESCE(meta->>'attempts','0')::int + 1, '}')::jsonb
	where 
		id = $1
		and claim_id = $2;
	`
	tag, err := gw.pool.Exec(ctx, query, id, claimID, map[string]string{"failReason": reason})
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.Error{Code: domain.ErrStaleResult, Message: "result is stale"}
	}
	// TODO: check attempt, send to failure table and delete from task table?
	return nil
}
