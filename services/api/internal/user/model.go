package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/vaxxnsh/metaverse/api/internal/db"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func convertDBUsertoModel(dbUser db.User) *User {
	return &User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Role:      "user",
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
}
