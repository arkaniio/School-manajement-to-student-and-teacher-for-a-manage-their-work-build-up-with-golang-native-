package students

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
)

type ServiceStudens struct {
	repo types.StudentStore
}

func NewHandlerService(repo types.StudentStore) *ServiceStudens {
	return &ServiceStudens{
		repo: repo,
	}
}

func (s *ServiceStudens) CreateNewStudent(ctx context.Context, students *types.Student) error {
	return s.repo.CreateNewStudent(ctx, students)
}

func (s *ServiceStudens) UpdateAsStudent(id uuid.UUID, ctx context.Context, payloads *types.UpdateAsStudent) error {
	return s.UpdateAsStudent(id, ctx, payloads)
}

func (s *ServiceStudens) DeleteStudents(id uuid.UUID, ctx context.Context) error {
	return s.repo.DeleteStudents(id, ctx)
}

func (s *ServiceStudens) GetAllStudents(ctx context.Context) ([]types.Student, error) {
	return s.repo.GetAllStudents(ctx)
}
