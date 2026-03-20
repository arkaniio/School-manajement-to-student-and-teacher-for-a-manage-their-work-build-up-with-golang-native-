package absensis

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
)

type ServiceAbsensi struct {
	repo types.AbsensiStore
}

func NewServiceAbsensi(repo types.AbsensiStore) *ServiceAbsensi {
	return &ServiceAbsensi{
		repo: repo,
	}
}

func (s *ServiceAbsensi) GetWeeklyStats(ctx context.Context) (*types.AbsensiStats, error) {
	return s.repo.GetWeeklyStats(ctx)
}

func (s *ServiceAbsensi) GetMonthlyStats(ctx context.Context) (*types.AbsensiStats, error) {
	return s.repo.GetMonthlyStats(ctx)
}

func (s *ServiceAbsensi) GetAllAbsensiWithStudents(ctx context.Context) ([]types.AbsensiWithStudent, error) {
	return s.repo.GetAllAbsensiWithStudents(ctx)
}

func (s *ServiceAbsensi) UpdateAbsensiById(id uuid.UUID, ctx context.Context, payloads types.PayloadAbsensisUpdate) error {
	return s.repo.UpdateAbsensiById(id, ctx, payloads)
}

func (s *ServiceAbsensi) DeleteAbsensisById(id uuid.UUID, ctx context.Context) error {
	return s.repo.DeleteAbsensisById(id, ctx)
}
