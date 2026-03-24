package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type SpaceService interface {
	CreateSpace(ctx context.Context, creatorId, name string, width, height int32, mapId string) (*domain.Space, error)
	DeleteSpace(ctx context.Context, spaceId string) error
	GetSpacesByCreator(ctx context.Context, creatorId string) ([]domain.SpaceSummary, error)
}

type spaceService struct {
	repo repository.SpaceRepository
}

func NewSpaceService(repo repository.SpaceRepository) SpaceService {
	return &spaceService{repo: repo}
}

func (s *spaceService) CreateSpace(ctx context.Context, creatorId, name string, width, height int32, mapId string) (*domain.Space, error) {
	return s.repo.CreateSpace(ctx, creatorId, name, width, height, mapId)
}

func (s *spaceService) DeleteSpace(ctx context.Context, spaceId string) error {
	return s.repo.DeleteSpace(ctx, spaceId)
}

func (s *spaceService) GetSpacesByCreator(ctx context.Context, creatorId string) ([]domain.SpaceSummary, error) {
	return s.repo.GetSpacesByCreator(ctx, creatorId)
}
