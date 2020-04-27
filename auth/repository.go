package auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type UserRepository interface {
	// creates user and immediately assings it a role
	CreateUser(ctx context.Context, user *models.User, role *models.Role) error

	GetUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserRole(ctx context.Context, user *models.User) (*models.Role, error)

	// GetUserWithRole(ctx context.Context,  string)
	GetUserByName(ctx context.Context, username string) (*models.User, error)
}
