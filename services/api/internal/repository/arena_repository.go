package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type ArenaRepository interface {
	GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error)
}

type arenaRepository struct {
	queries *db.Queries
}

func NewArenaRepository(q *db.Queries) ArenaRepository {
	return &arenaRepository{queries: q}
}

func (r *arenaRepository) GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error) {
	spaceUUID := pgtype.UUID{}
	if err := spaceUUID.Scan(spaceId); err != nil || !spaceUUID.Valid {
		return nil, domain.ErrInvalidSpaceID
	}

	space, err := r.queries.GetSpaceByID(ctx, spaceUUID)
	if err != nil {
		return nil, domain.ErrSpaceNotFound
	}

	elements, err := r.queries.GetSpaceElements(ctx, spaceUUID)
	if err != nil {
		return nil, err
	}

	result := domain.DBSpaceElementsToArena(space, elements)
	return &result, nil
}
