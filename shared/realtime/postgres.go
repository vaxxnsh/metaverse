package realtime

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vaxxnsh/shared/db"
)

var ErrSpaceNotFound = errors.New("space not found")

type StaticElement struct {
	ID                  string
	X, Y, Width, Height int32
	Static              bool
}
type Space struct {
	ID            string
	Width, Height int32
	Elements      []StaticElement
}

type PostgresProvider struct{ queries *db.Queries }

func NewPostgresProvider(pool *pgxpool.Pool) *PostgresProvider {
	return &PostgresProvider{queries: db.New(pool)}
}

func OpenPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func (p *PostgresProvider) GetRealtimeSpace(ctx context.Context, spaceID string) (Space, error) {
	var id pgtype.UUID
	if err := id.Scan(spaceID); err != nil || !id.Valid {
		return Space{}, ErrSpaceNotFound
	}
	row, err := p.queries.GetSpaceByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Space{}, ErrSpaceNotFound
	}
	if err != nil {
		return Space{}, err
	}
	elements, err := p.queries.GetSpaceElements(ctx, id)
	if err != nil {
		return Space{}, err
	}
	space := Space{ID: row.ID.String(), Width: row.Width, Height: row.Height, Elements: make([]StaticElement, 0, len(elements))}
	for _, element := range elements {
		space.Elements = append(space.Elements, StaticElement{ID: element.ElementID.String(), X: element.X, Y: element.Y, Width: element.ElementWidth, Height: element.ElementHeight, Static: element.Static})
	}
	return space, nil
}
