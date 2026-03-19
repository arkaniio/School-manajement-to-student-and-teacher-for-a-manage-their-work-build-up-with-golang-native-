package service

import (
	"context"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
)

type ServiceUser struct {
	repo types.UserStore
}

func NewServiceUser(repo types.UserStore) *ServiceUser {
	return &ServiceUser{repo: repo}
}

func (s *ServiceUser) CreateUser(ctx context.Context, user *types.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *ServiceUser) GetUserByID(id uuid.UUID, ctx context.Context) (*types.User, error) {
	return s.repo.GetUserById(id, ctx)
}

func (s *ServiceUser) GetUserByEmailAndUsername(email string, username string) (*types.User, error) {
	return s.repo.GetUserByEmailAndUsername(email, username)
}

func (s *ServiceUser) UpdateDataUser(id uuid.UUID, ctx context.Context, payload types.Update) error {
	return s.repo.UpdateDataUser(id, ctx, payload)
}
