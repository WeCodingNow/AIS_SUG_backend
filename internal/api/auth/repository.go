package auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)

	// GetUserRole(ctx context.Context, user *models.User) (*models.Role, error)
	// GetUserWithRole(ctx context.Context,  string)
}
