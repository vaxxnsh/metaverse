package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
	"github.com/vaxxnsh/metaverse/api/internal/utils"
)

type UserService interface {
	Create(ctx context.Context, name, email, password string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	PatchMetadata(ctx context.Context, userId, avatarId string) (*domain.User, error)
	FindUserByID(ctx context.Context, userId string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) Create(ctx context.Context, name string, email string, password string) (*domain.User, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.Create(ctx, name, email, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) PatchMetadata(ctx context.Context, userId, avatarId string) (*domain.User, error) {
	user, err := u.repo.PatchMetadata(ctx, userId, avatarId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) FindUserByID(ctx context.Context, userId string) (*domain.User, error) {
	user, err := u.repo.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
