package repository

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type MetadataRepository interface {
	GetAvatars(ctx context.Context) ([]domain.Avatar, error)
}

type metadataRepository struct {
	queries *db.Queries
}

func NewMetadataRepository(q *db.Queries) MetadataRepository {
	return &metadataRepository{
		queries: q,
	}
}

func (m *metadataRepository) GetAvatars(ctx context.Context) ([]domain.Avatar, error) {
	avatars, err := m.queries.GetAvatars(ctx)
	if err != nil {
		return nil, err
	}

	return domain.MapAvatarsToDomain(avatars), nil
}
