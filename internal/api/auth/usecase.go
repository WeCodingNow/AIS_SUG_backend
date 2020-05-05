package auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, username, password string) (*models.User, error)
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
