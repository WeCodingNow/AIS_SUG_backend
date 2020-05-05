package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/dgrijalva/jwt-go"
)

type AuthUseCase struct {
	userRepo   auth.UserRepository
	hashSalt   string
	signingKey []byte
	expireDate time.Duration
}

type Claims struct {
	jwt.StandardClaims
	// Username string `json:"username"`
	UserID int `json:"id"`
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		hashSalt:   hashSalt,
		signingKey: signingKey,
		expireDate: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthUseCase) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	hashedPwd := fmt.Sprintf("%x", pwd.Sum(nil))

	return a.userRepo.CreateUser(ctx, username, hashedPwd)

}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(
		accessToken, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return a.signingKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		user, err := a.userRepo.GetUserByID(ctx, claims.UserID)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, auth.ErrInvalidAccessToken
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	hashedPwd := fmt.Sprintf("%x", pwd.Sum(nil))
	user, err := a.userRepo.GetUser(ctx, username, hashedPwd)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(a.expireDate).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "ais_sug",
		},
		UserID: user.ID,
	})

	return token.SignedString(a.signingKey)
}
