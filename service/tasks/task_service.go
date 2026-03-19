package tasks

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
)

type ServiceTask struct {
	repo types.TaskStore
}

func NewServiceTask(repo types.TaskStore) *ServiceTask {
	return &ServiceTask{repo: repo}
}

func (s *ServiceTask) CreateTask(ctx context.Context, task *types.Task) error {
	return s.repo.CreateNewTasks(ctx, task)
}

func (s *ServiceTask) GetTaskByID(id uuid.UUID, ctx context.Context) (*types.Task, error) {
	return s.repo.GetTaskById(id, ctx)
}

func (s *ServiceTask) GetAllTasks(ctx context.Context) ([]types.TaskWithStudents, error) {
	return s.repo.GetAllTaskIncludeStudents(ctx)
}

func (s *ServiceTask) UpdateTask(id uuid.UUID, ctx context.Context, payloads types.PayloadUpdate) error {
	return s.repo.UpdateTask(id, ctx, payloads)
}

func (s *ServiceTask) DeleteTask(id uuid.UUID, ctx context.Context) error {
	return s.repo.DeleteTask(id, ctx)
}
