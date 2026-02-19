package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, email, name string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type userService struct {
	repository UserRepository
}

func NewUserService(r UserRepository) Service {
	return &userService{
		repository: r,
	}
}

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidName  = errors.New("invalid name")
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

func (s *userService) CreateUser(
	ctx context.Context,
	email string,
	name string,
) (*User, error) {

	if email == "" {
		return nil, ErrInvalidEmail
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	existing, err := s.repository.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrUserExists
	}

	u := &User{
		ID:        uuid.NewString(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repository.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, ErrUserNotFound
	}

	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
