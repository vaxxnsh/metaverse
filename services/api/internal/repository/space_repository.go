package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type SpaceRepository interface {
	CreateSpace(ctx context.Context, creatorId, name string, width, height int32, mapId string) (*domain.Space, error)
	DeleteSpace(ctx context.Context, spaceId string) error
	GetSpacesByCreator(ctx context.Context, creatorId string) ([]domain.SpaceSummary, error)
}

type spaceRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewSpaceRepository(pool *pgxpool.Pool, q *db.Queries) SpaceRepository {
	return &spaceRepository{pool: pool, queries: q}
}

func (s *spaceRepository) CreateSpace(ctx context.Context, creatorId, name string, width, height int32, mapId string) (*domain.Space, error) {
	creatorUUID := pgtype.UUID{}
	if err := creatorUUID.Scan(creatorId); err != nil || !creatorUUID.Valid {
		return nil, domain.ErrUserNotFound
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	space, err := qtx.CreateSpace(ctx, db.CreateSpaceParams{
		CreatorID: creatorUUID,
		Name:      name,
		Width:     width,
		Height:    height,
	})
	if err != nil {
		return nil, err
	}

	if mapId != "" {
		mapUUID := pgtype.UUID{}
		if err := mapUUID.Scan(mapId); err != nil || !mapUUID.Valid {
			return nil, domain.ErrInvalidMapID
		}

		mapElements, err := qtx.GetMapElementsByMapID(ctx, mapUUID)
		if err != nil {
			return nil, domain.ErrMapNotFound
		}

		if len(mapElements) > 0 {
			elementIDs := make([]pgtype.UUID, 0, len(mapElements))
			xs := make([]int32, 0, len(mapElements))
			ys := make([]int32, 0, len(mapElements))

			for _, me := range mapElements {
				elementIDs = append(elementIDs, me.ElementID)
				xs = append(xs, me.X)
				ys = append(ys, me.Y)
			}

			if len(elementIDs) > 0 {
				err = qtx.BulkCreateSpaceElements(ctx, db.BulkCreateSpaceElementsParams{
					SpaceID:    space.ID,
					ElementIds: elementIDs,
					Xs:         xs,
					Ys:         ys,
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return domain.DBSpaceToDomain(space), nil
}

func (s *spaceRepository) DeleteSpace(ctx context.Context, spaceId string) error {
	spaceUUID := pgtype.UUID{}
	if err := spaceUUID.Scan(spaceId); err != nil || !spaceUUID.Valid {
		return domain.ErrInvalidSpaceID
	}

	return s.queries.DeleteSpace(ctx, spaceUUID)
}

func (s *spaceRepository) GetSpacesByCreator(ctx context.Context, creatorId string) ([]domain.SpaceSummary, error) {
	creatorUUID := pgtype.UUID{}
	if err := creatorUUID.Scan(creatorId); err != nil || !creatorUUID.Valid {
		return nil, domain.ErrUserNotFound
	}

	rows, err := s.queries.GetSpacesByCreator(ctx, creatorUUID)
	if err != nil {
		return nil, err
	}

	return domain.MapSpacesToSummary(rows), nil
}
