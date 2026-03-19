package absensis

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
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
