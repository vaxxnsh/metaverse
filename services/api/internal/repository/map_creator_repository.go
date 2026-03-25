package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type MapCreatorRepository interface {
	CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error)
	UpdateElementImage(ctx context.Context, elementId, imageUrl string) (*domain.Element, error)
	CreateAvatar(ctx context.Context, imageUrl, name string) (*domain.Avatar, error)
}

type mapCreatorRepository struct {
	queries *db.Queries
}

func NewMapCreatorRepository(q *db.Queries) MapCreatorRepository {
	return &mapCreatorRepository{queries: q}
}

func (r *mapCreatorRepository) CreateAvatar(ctx context.Context, imageUrl, name string) (*domain.Avatar, error) {
	avatar, err := r.queries.CreateAvatar(ctx, db.CreateAvatarParams{
		ImageUrl: imageUrl,
		Name:     pgtype.Text{String: name, Valid: name != ""},
	})
	if err != nil {
		return nil, err
	}

	return domain.DBAvatarToDomain(avatar), nil
}

func (r *mapCreatorRepository) UpdateElementImage(ctx context.Context, elementId, imageUrl string) (*domain.Element, error) {
	elementUUID := pgtype.UUID{}
	if err := elementUUID.Scan(elementId); err != nil || !elementUUID.Valid {
		return nil, domain.ErrInvalidElementID
	}

	element, err := r.queries.UpdateElementImage(ctx, db.UpdateElementImageParams{
		ID:       elementUUID,
		ImageUrl: imageUrl,
	})
	if err != nil {
		return nil, err
	}

	return domain.DBElementToDomain(element), nil
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
