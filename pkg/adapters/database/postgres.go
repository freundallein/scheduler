package database

import (
	"context"
	"time"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskGateway struct {
	pool *pgxpool.Pool
}

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
func (gw *TaskGateway) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `insert into task(id, executed_at, deadline, payload) values ($1, $2, $3, $4) returning *;`
	row := gw.pool.QueryRow(ctx, query, task.ID, task.ExecuteAt, task.Deadline, task.Payload)
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
		return nil, err
	}
	return task, nil
}

func (gw *TaskGateway) FindByID(ctx context.Context, id string) (*domain.Task, error) {
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
		return nil, err
	}
	return task, nil
}

func (gw *TaskGateway) ClaimPending(ctx context.Context, amount int) ([]*domain.Task, error) {
	return nil, nil
}
func (gw *TaskGateway) MarkAsSucceeded(ctx context.Context, id, claimID, result string) error {
	return nil
}
func (gw *TaskGateway) MarkAsFailed(ctx context.Context, id, claimID, reason string) error {
	return nil
}
