package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/shared/db"
	"github.com/vaxxnsh/shared/domain"
)

type MetadataRepository interface {
	GetAvatars(ctx context.Context) ([]domain.Avatar, error)
	GetBulkUserAvatars(ctx context.Context, userIds []string) ([]domain.UserAvatar, error)
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

func (m *metadataRepository) GetBulkUserAvatars(ctx context.Context, userIds []string) ([]domain.UserAvatar, error) {
	uuids := make([]pgtype.UUID, 0, len(userIds))
	for _, id := range userIds {
		u := pgtype.UUID{}
		if err := u.Scan(id); err != nil || !u.Valid {
			continue
		}
		uuids = append(uuids, u)
	}

	rows, err := m.queries.GetBulkUserAvatars(ctx, uuids)
	if err != nil {
		return nil, err
	}

	return domain.MapBulkUserAvatarsToDomain(rows), nil
}
