package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
)

type psqlUserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *psqlUserRepository {
	return &psqlUserRepository{
		queries: queries,
	}
}

func (r *psqlUserRepository) Create(
	ctx context.Context,
	id pgtype.UUID,
	name, email, password string,
) (db.User, error) {

	now := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}

	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	})
}
