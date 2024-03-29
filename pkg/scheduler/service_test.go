package scheduler

import (
	"context"
	"errors"
	"fmt"
	domain "github.com/freundallein/scheduler/pkg"
	"github.com/freundallein/scheduler/pkg/mock"
	"github.com/google/uuid"
	"testing"
)

var errExpected = errors.New("expected error")

func TestSet(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		task        *domain.Task
		expectedID  uuid.UUID
	}{
		{
			name:       "normal case",
			task:       &domain.Task{},
			expectedID: uuid.New(),
		},
		{
			name:        "error case",
			task:        &domain.Task{},
			expectedErr: errExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := New(
				&mock.Gateway{
					CreateFn: func(task *domain.Task) (*domain.Task, error) {
						if task != tt.task {
							t.Errorf("Expected: `%v`, got: `%v`", tt.task, task)
						}
						if task == nil || tt.expectedErr != nil {
							return nil, tt.expectedErr
						}
						task.ID = tt.expectedID
						task.State = domain.StatePending
						return task, nil
					},
				},
			)
			ctx := context.Background()
			observed, err := scheduler.Set(ctx, tt.task)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected `%v`, got: `%v`", tt.expectedErr, err)
			}
			if observed == nil {
				return
			}
			if observed.State != domain.StatePending {
				t.Errorf("Expected `%v`, got: `%v`", domain.StatePending, observed.State)
			}
			if observed.ID != tt.expectedID {
				t.Errorf("Expected `%v`, got: `%v`", tt.expectedID, observed.ID)
			}
		})
	}
}

func TestGet(t *testing.T) {
	expectedUUID := uuid.New()
	tests := []struct {
		name          string
		expectedErr   error
		task          *domain.Task
		expectedID    uuid.UUID
		expectedState domain.State
	}{
		{
			name: "normal case",
			task: &domain.Task{
				ID:    expectedUUID,
				State: domain.StatePending,
			},
			expectedID:    expectedUUID,
			expectedState: domain.StatePending,
		},
		{
			name:        "error case",
			expectedID:  expectedUUID,
			expectedErr: errExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := New(
				&mock.Gateway{
					FindByIDFn: func(id uuid.UUID) (*domain.Task, error) {
						if tt.expectedErr != nil {
							return nil, tt.expectedErr
						}
						if id != tt.expectedID {
							return nil, fmt.Errorf("task not found")
						}
						return tt.task, nil
					},
				},
			)
			ctx := context.Background()
			observed, err := scheduler.Get(ctx, tt.expectedID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected `%v`, got: `%v`", tt.expectedErr, err)
			}
			if observed == nil {
				return
			}
			if observed.State != tt.expectedState {
				t.Errorf("Expected `%v`, got: `%v`", tt.expectedState, observed.State)
			}
			if observed.ID != tt.expectedID {
				t.Errorf("Expected `%v`, got: `%v`", tt.expectedID, observed.ID)
			}
		})
	}
}
