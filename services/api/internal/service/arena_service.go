package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type ArenaService interface {
	GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error)
}

type arenaService struct {
	repo repository.ArenaRepository
}

func NewArenaService(repo repository.ArenaRepository) ArenaService {
	return &arenaService{repo: repo}
}

func (s *arenaService) GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error) {
	return s.repo.GetSpaceWithElements(ctx, spaceId)
}
