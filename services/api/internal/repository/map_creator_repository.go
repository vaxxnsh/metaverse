package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type DefaultElement struct {
	ElementId string
	X         int32
	Y         int32
}

type MapCreatorRepository interface {
	CreateElement(ctx context.Context, imageUrl string, width, height int32, static bool) (*domain.Element, error)
	UpdateElementImage(ctx context.Context, elementId, imageUrl string) (*domain.Element, error)
	CreateAvatar(ctx context.Context, imageUrl, name string) (*domain.Avatar, error)
	CreateMap(ctx context.Context, name string, width, height int32, defaultElements []DefaultElement) (*domain.Map, error)
}

type mapCreatorRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewMapCreatorRepository(pool *pgxpool.Pool, q *db.Queries) MapCreatorRepository {
	return &mapCreatorRepository{pool: pool, queries: q}
}

func (r *mapCreatorRepository) CreateMap(ctx context.Context, name string, width, height int32, defaultElements []DefaultElement) (*domain.Map, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	m, err := qtx.CreateMap(ctx, db.CreateMapParams{
		Name:   name,
		Width:  width,
		Height: height,
	})
	if err != nil {
		return nil, err
	}

	if len(defaultElements) > 0 {
		elementIDs := make([]pgtype.UUID, 0, len(defaultElements))
		xs := make([]int32, 0, len(defaultElements))
		ys := make([]int32, 0, len(defaultElements))

		for _, e := range defaultElements {
			elementUUID := pgtype.UUID{}
			if err := elementUUID.Scan(e.ElementId); err != nil || !elementUUID.Valid {
				return nil, domain.ErrInvalidElementID
			}
			elementIDs = append(elementIDs, elementUUID)
			xs = append(xs, e.X)
			ys = append(ys, e.Y)
		}

		if err := qtx.BulkCreateMapElements(ctx, db.BulkCreateMapElementsParams{
			MapID:      m.ID,
			ElementIds: elementIDs,
			Xs:         xs,
			Ys:         ys,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return domain.DBMapToDomain(m), nil
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
