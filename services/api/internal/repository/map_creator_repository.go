package repository

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type MapCreatorRepository interface {
	CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error)
}

type mapCreatorRepository struct {
	queries *db.Queries
}

func NewMapCreatorRepository(q *db.Queries) MapCreatorRepository {
	return &mapCreatorRepository{queries: q}
}

func (r *mapCreatorRepository) CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error) {
	element, err := r.queries.CreateElement(ctx, db.CreateElementParams{
		ImageUrl: imageUrl,
		Width:    width,
		Height:   height,
		Static:   static,
	})
	if err != nil {
		return nil, err
	}

	return domain.DBElementToDomain(element), nil
}
