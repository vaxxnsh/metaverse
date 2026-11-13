package ws

import "time"

// User is the authoritative live state for one room member.
type User struct {
	ID       string
	Position Position
	LastMove time.Time
}
