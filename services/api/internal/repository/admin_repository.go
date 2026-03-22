package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type AdminRepository interface {
	Create(ctx context.Context, name, email, paswordHash string) (*domain.Admin, error)
	FindByEmail(ctx context.Context, email string) (*domain.Admin, error)
}

type adminRepository struct {
	queries *db.Queries
}

func NewAdminRepository(q *db.Queries) AdminRepository {
	return &adminRepository{queries: q}
}

func (a *adminRepository) Create(ctx context.Context, name, email, paswordHash string) (*domain.Admin, error) {
	admin, err := a.queries.CreateAdmin(ctx, db.CreateAdminParams{
		Name:     name,
		Email:    email,
		Password: paswordHash,
	})

	if err != nil {
		return nil, err
	}

	return domain.DBAdminToDomain(admin), nil
}

func (a *adminRepository) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	admin, err := a.queries.FindAdminByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return domain.DBAdminToDomain(admin), nil
}
