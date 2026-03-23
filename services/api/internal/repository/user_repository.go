package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vaxxnsh/metaverse/api/internal/db"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, paswordHash string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	PatchMetadata(ctx context.Context, userId, avatarId string) (*domain.User, error)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return domain.DBUserToDomain(user), nil
}

func (u *userRepository) PatchMetadata(ctx context.Context, userId, avatarId string) (*domain.User, error) {
	userUUID := pgtype.UUID{}
	avatarUUID := pgtype.UUID{}

	if err := userUUID.Scan(userId); err != nil || !userUUID.Valid {
		return nil, err
	}

	if err := avatarUUID.Scan(avatarId); err != nil || !avatarUUID.Valid {
		return nil, err
	}

	_, err := u.queries.FindUserByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	avatar, err := u.queries.FindAvatarByID(ctx, avatarUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("avatar not found")
		}
		return nil, err
	}

	updatedUser, err := u.queries.PatchUserMetadata(ctx, db.PatchUserMetadataParams{
		ID:       userUUID,
		AvatarID: avatar.ID,
	})
	if err != nil {
		return nil, err
	}

	return domain.DBUserToDomain(updatedUser), nil
}
