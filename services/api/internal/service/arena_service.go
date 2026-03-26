package service

import (
	"context"

	"github.com/vaxxnsh/shared/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type ArenaService interface {
	GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error)
	AddElementToSpace(ctx context.Context, spaceId, elementId string, x, y int32) (*domain.SpaceElement, error)
	DeleteSpaceElement(ctx context.Context, spaceId string, x, y int32) error
	GetAllElements(ctx context.Context) ([]domain.ElementDetails, error)
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

func (s *arenaService) AddElementToSpace(ctx context.Context, spaceId, elementId string, x, y int32) (*domain.SpaceElement, error) {
	return s.repo.AddElementToSpace(ctx, spaceId, elementId, x, y)
}

func (s *arenaService) DeleteSpaceElement(ctx context.Context, spaceId string, x, y int32) error {
	return s.repo.DeleteSpaceElement(ctx, spaceId, x, y)
}

func (s *arenaService) GetAllElements(ctx context.Context) ([]domain.ElementDetails, error) {
	return s.repo.GetAllElements(ctx)
}
