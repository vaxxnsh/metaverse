package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, name, email, password string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) Create(ctx context.Context, name string, email string, password string) (*domain.User, error) {
	user, err := u.repo.Create(ctx, name, email, password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
