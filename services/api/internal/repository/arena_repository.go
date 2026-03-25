package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type ArenaRepository interface {
	GetSpaceWithElements(ctx context.Context, spaceId string) (*domain.SpaceWithElements, error)
	AddElementToSpace(ctx context.Context, spaceId, elementId string, x, y int32) (*domain.SpaceElement, error)
	DeleteSpaceElement(ctx context.Context, spaceId string, x, y int32) error
	GetAllElements(ctx context.Context) ([]domain.ElementDetails, error)
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

func (r *arenaRepository) GetAllElements(ctx context.Context) ([]domain.ElementDetails, error) {
	rows, err := r.queries.GetAllElements(ctx)
	if err != nil {
		return nil, err
	}
	return domain.MapAllElementsToDomain(rows), nil
}

func (r *arenaRepository) DeleteSpaceElement(ctx context.Context, spaceId string, x, y int32) error {
	spaceUUID := pgtype.UUID{}
	if err := spaceUUID.Scan(spaceId); err != nil || !spaceUUID.Valid {
		return domain.ErrInvalidSpaceID
	}

	return r.queries.DeleteSpaceElement(ctx, db.DeleteSpaceElementParams{
		SpaceID: spaceUUID,
		X:       x,
		Y:       y,
	})
}

func (r *arenaRepository) AddElementToSpace(ctx context.Context, spaceId, elementId string, x, y int32) (*domain.SpaceElement, error) {
	spaceUUID := pgtype.UUID{}
	if err := spaceUUID.Scan(spaceId); err != nil || !spaceUUID.Valid {
		return nil, domain.ErrInvalidSpaceID
	}

	elementUUID := pgtype.UUID{}
	if err := elementUUID.Scan(elementId); err != nil || !elementUUID.Valid {
		return nil, domain.ErrInvalidElementID
	}

	se, err := r.queries.CreateSpaceElement(ctx, db.CreateSpaceElementParams{
		SpaceID:   spaceUUID,
		ElementID: elementUUID,
		X:         x,
		Y:         y,
	})
	if err != nil {
		return nil, err
	}

	return domain.DBSpaceElementToDomain(se), nil
}
