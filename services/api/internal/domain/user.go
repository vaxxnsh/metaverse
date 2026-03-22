package domain

import (
	"time"

	"github.com/vaxxnsh/metaverse/api/internal/db"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func DBUserToDomain(user db.User) *User {
	return &User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}
}
