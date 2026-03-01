package user

import (
	"context"
	"errors"
)

type Service interface {
	Signup(ctx context.Context, name, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Signup(ctx context.Context, name, email, password string) (*User, error) {
	dbUser, err := s.repo.Create(ctx, createDBUserArgs{
		name:     name,
		email:    email,
		password: password,
	})

	if err != nil {
		return nil, err
	}

	return convertDBUsertoModel(dbUser), nil
}

func (s *service) Login(ctx context.Context, email, password string) (*User, error) {
	dbUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return &User{}, errors.New("invalid credentials")
	}

	if dbUser.Password != password {
		return &User{}, errors.New("invalid credentials")
	}

	return convertDBUsertoModel(dbUser), nil
}
