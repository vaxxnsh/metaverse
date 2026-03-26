package service

import (
	"context"

	"github.com/vaxxnsh/shared/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type AdminService interface {
	Create(ctx context.Context, name, email, password string) (*domain.Admin, error)
	FindByEmail(ctx context.Context, email string) (*domain.Admin, error)
	PatchMetadata(ctx context.Context, adminId, avatarId string) (*domain.Admin, error)
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

func (a *adminService) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	admin, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (a *adminService) PatchMetadata(ctx context.Context, adminId, avatarId string) (*domain.Admin, error) {
	admin, err := a.repo.PatchMetadata(ctx, adminId, avatarId)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
