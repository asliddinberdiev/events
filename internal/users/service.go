package users

import (
	"context"

	"github.com/asliddinberdiev/events/internal/common"
	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req RegisterPayload) (*User, error) {
	return s.repo.Create(ctx, req)
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetByTelegramID(ctx context.Context, telegram_id string) (*User, error) {
	return s.repo.GetByTelegramID(ctx, telegram_id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) GetAll(ctx context.Context, params common.SearchRequest) (*common.AllResponse, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *UserService) Update(ctx context.Context, req UpdateUser) error {
	return s.repo.Update(ctx, req)
}

func (s *UserService) Delete(ctx context.Context, user_id uuid.UUID) error {
	return s.repo.Delete(ctx, user_id)
}
