package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type SpaceService interface {
	CreateSpace(ctx context.Context, name string, width, height int32, mapId string) (*domain.Space, error)
}

type spaceService struct {
	repo repository.SpaceRepository
}

func NewSpaceService(repo repository.SpaceRepository) SpaceService {
	return &spaceService{repo: repo}
}

func (s *spaceService) CreateSpace(ctx context.Context, name string, width, height int32, mapId string) (*domain.Space, error) {
	return s.repo.CreateSpace(ctx, name, width, height, mapId)
}
