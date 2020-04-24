package auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByName(ctx context.Context, username string) (*models.User, error)
}
