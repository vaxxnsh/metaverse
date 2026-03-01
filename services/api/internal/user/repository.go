package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
)

type createDBUserArgs struct {
	name     string
	email    string
	password string
	role     string
}

type Repository interface {
	Create(ctx context.Context, args createDBUserArgs) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
}

type repository struct {
	queries *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &repository{queries: q}
}

func (r *repository) Create(ctx context.Context, args createDBUserArgs) (db.User, error) {
	id := pgtype.UUID{}
	_ = id.Scan(uuid.New())

	now := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}

	createArgs := db.CreateUserParams{
		ID:        id,
		Name:      args.name,
		Email:     args.email,
		UpdatedAt: now,
		CreatedAt: now,
		Password:  args.password,
	}

	return r.queries.CreateUser(ctx, createArgs)
}

func (r *repository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
