package service

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
)

type MapCreatorService interface {
	CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error)
	UpdateElementImage(ctx context.Context, elementId, imageUrl string) (*domain.Element, error)
	CreateAvatar(ctx context.Context, imageUrl, name string) (*domain.Avatar, error)
}

type mapCreatorService struct {
	repo repository.MapCreatorRepository
}

func NewMapCreatorService(repo repository.MapCreatorRepository) MapCreatorService {
	return &mapCreatorService{repo: repo}
}

func (s *mapCreatorService) CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error) {
	return s.repo.CreateElement(ctx, imageUrl, width, height, static)
}

func (s *mapCreatorService) UpdateElementImage(ctx context.Context, elementId, imageUrl string) (*domain.Element, error) {
	return s.repo.UpdateElementImage(ctx, elementId, imageUrl)
}

func (s *mapCreatorService) CreateAvatar(ctx context.Context, imageUrl, name string) (*domain.Avatar, error) {
	return s.repo.CreateAvatar(ctx, imageUrl, name)
}
