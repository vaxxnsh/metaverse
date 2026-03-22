package service

import (
	"context"
	"errors"

	"github.com/vaxxnsh/metaverse/api/internal/utils"
)

type AuthService interface {
	RegisterUser(ctx context.Context, name, email, password string) (string, error)
	RegisterAdmin(ctx context.Context, name, email, password string) (string, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	LoginAdmin(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userService  UserService
	adminService AdminService
}

func NewAuthService(u UserService, a AdminService) AuthService {
	return &authService{
		userService:  u,
		adminService: a,
	}
}

func (a *authService) RegisterUser(ctx context.Context, name, email, password string) (string, error) {
	user, err := a.userService.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", errors.New("user with this email already exits")
	}

	newUser, err := a.userService.Create(ctx, name, email, password)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := a.userService.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID)
}

func (a *authService) RegisterAdmin(ctx context.Context, name, email, password string) (string, error) {
	admin, err := a.adminService.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if admin != nil {
		return "", errors.New("user with this email already exits")
	}

	newAdmin, err := a.userService.Create(ctx, name, email, password)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(newAdmin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) LoginAdmin(ctx context.Context, email, password string) (string, error) {
	admin, err := a.userService.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := utils.CheckPassword(admin.Password, password); err != nil {
		return "", err
	}

	return utils.GenerateToken(admin.ID)
}
