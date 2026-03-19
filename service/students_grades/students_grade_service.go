package studentsgrades

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
)

type ServiceStudentGrade struct {
	repo types.StudentsGradeStore
}

func NewServiceStudentsGrade(repo types.StudentsGradeStore) *ServiceStudentGrade {
	return &ServiceStudentGrade{
		repo: repo,
	}
}

func (s *ServiceStudentGrade) CreateNewStudentsGrade(ctx context.Context, payloads *types.StudentsGrade) error {
	return s.repo.CreateNewStudentsGrade(ctx, payloads)
}

func (s *ServiceStudentGrade) UpdateStudentsGrade(id uuid.UUID, ctx context.Context, payloads types.PayloadsStudentGradeUpdate) error {
	return s.repo.UpdateStudentsGrade(ctx, id, &payloads)
}

func (s *ServiceStudentGrade) DeleteStudentsGrade(id uuid.UUID, ctx context.Context) error {
	return s.repo.DeleteStudentsGrade(ctx, id)
}
