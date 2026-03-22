package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type AdminService interface {
	Create(ctx context.Context, name, email, password string) (*domain.Admin, error)
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (a *adminService) Create(ctx context.Context, name string, email string, password string) (*domain.Admin, error) {
	admin, err := a.repo.Create(ctx, name, email, password)

	if err != nil {
		return nil, err
	}

	return admin, nil
}
