package repository

import (
	"context"

	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, paswordHash string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{queries: q}
}

func (u *userRepository) Create(ctx context.Context, name, email, paswordHash string) (*domain.User, error) {
	user, err := u.queries.CreateUser(ctx, db.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: paswordHash,
	})

	if err != nil {
		return nil, err
	}

	return domain.DBUserToDomain(user), nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.queries.FindUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return domain.DBUserToDomain(user), nil
}
