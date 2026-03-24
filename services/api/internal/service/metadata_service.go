package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type MetaDataService interface {
	GetAvatars(ctx context.Context, userId string) ([]domain.Avatar, error)
}

type metadataService struct {
	repo        repository.MetadataRepository
	userService UserService
}

func NewMetadataService(repo repository.MetadataRepository, userService UserService) MetaDataService {
	return &metadataService{
		repo:        repo,
		userService: userService,
	}
}

func (m *metadataService) GetAvatars(ctx context.Context, userId string) ([]domain.Avatar, error) {
	user, err := m.userService.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return m.repo.GetAvatars(ctx)
}
