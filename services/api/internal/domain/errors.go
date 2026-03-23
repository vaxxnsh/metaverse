package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAdminNotFound  = errors.New("admin not found")
	ErrAvatarNotFound = errors.New("avatar not found")
)
