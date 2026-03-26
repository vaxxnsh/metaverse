package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAdminNotFound  = errors.New("admin not found")
	ErrAvatarNotFound = errors.New("avatar not found")
	ErrInvalidMapID   = errors.New("invalid map id")
	ErrMapNotFound    = errors.New("map not found")
	ErrInvalidSpaceID   = errors.New("invalid space id")
	ErrSpaceNotFound    = errors.New("space not found")
	ErrInvalidElementID = errors.New("invalid element id")
)
