package token

import (
	"context"
	"time"

	"github.com/4aykovski/grpc_auth_sso/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct{}

func (m *Manager) GenerateJWTToken(
	ctx context.Context,
	user entity.User,
	app entity.App,
	tokenTTL time.Duration,
	secret string,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
