package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/4aykovski/grpc_auth_sso/internal/adapters/repository"
	"github.com/4aykovski/grpc_auth_sso/internal/entity"
)

type userRepository interface {
	SaveUser(ctx context.Context, user entity.User) (int64, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
}

type adminRepository interface {
	GetAdmin(ctx context.Context, userId int) (entity.Admin, error)
}

type appRepository interface {
	GetApp(ctx context.Context, appID int) (entity.App, error)
}

type tokenManager interface {
	GenerateJWTToken(
		ctx context.Context,
		user entity.User,
		app entity.App,
		tokenTTL time.Duration,
		secret string,
	) (string, error)
}

type secretManager interface {
	GetSecret(ctx context.Context, appID int) (string, error)
}

type hasher interface {
	Hash(password string) (string, error)
	Check(password string, hash string) bool
}

type Service struct {
	log *slog.Logger

	userRepo  userRepository
	appRepo   appRepository
	adminRepo adminRepository

	tokenManager  tokenManager
	secretManager secretManager
	hasher        hasher

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid appId")
	ErrInvalidUserId      = errors.New("invalid userId")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// New creates new auth Service
func New(
	log *slog.Logger,
	userRepo userRepository,
	appRepo appRepository,
	adminRepo adminRepository,
	tokenManager tokenManager,
	secretManager secretManager,
	hasher hasher,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Service {
	return &Service{
		log:             log,
		userRepo:        userRepo,
		appRepo:         appRepo,
		adminRepo:       adminRepo,
		tokenManager:    tokenManager,
		secretManager:   secretManager,
		hasher:          hasher,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type LoginDTO struct {
	Email    string
	Password string
	AppId    int
}

// Login checks if user with given credentials exists in the system
//
// If user exists, but password is incorrect, returns error ErrInvalidCredentials
// If user doesn't exist, returns error ErrInvalidCredentials
// If app doesn't exist, returns error ErrInvalidAppId
func (s *Service) Login(ctx context.Context, dto LoginDTO) (string, error) {
	user, err := s.userRepo.GetUser(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
		}

		return "", fmt.Errorf("can't login user: %w", err)
	}
	s.log.Debug("user", slog.Int("user", int(user.ID)))

	if ok := s.hasher.Check(dto.Password, user.PasswordHash); !ok {
		return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
	}

	app, err := s.appRepo.GetApp(ctx, dto.AppId)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			return "", fmt.Errorf("can't login user: %w", ErrInvalidAppId)
		}

		return "", fmt.Errorf("can't login user: %w", err)
	}
	s.log.Debug("app", slog.String("app", app.Name), slog.Int("appId", app.ID))

	secret, err := s.secretManager.GetSecret(ctx, dto.AppId)
	if err != nil {
		return "", fmt.Errorf("can't login user: %w", err)
	}
	s.log.Debug("secret", slog.String("secret", secret))

	token, err := s.tokenManager.GenerateJWTToken(
		ctx,
		user,
		app,
		s.accessTokenTTL,
		secret,
	)
	if err != nil {
		return "", fmt.Errorf("can't login user: %w", err)
	}

	return token, nil
}

type RegisterDTO struct {
	Email    string
	Password string
}

// Register creates new user in the system
//
// If user with the same email already exists, returns error ErrUserAlreadyExists
func (s *Service) Register(ctx context.Context, dto RegisterDTO) (int64, error) {
	passHash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return -1, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entity.User{
		Email:        dto.Email,
		PasswordHash: passHash,
	}

	id, err := s.userRepo.SaveUser(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return -1, fmt.Errorf("failed to save user: %w", ErrUserAlreadyExists)
		}

		return -1, fmt.Errorf("failed to save user: %w", err)
	}

	return id, nil
}

type IsAdminDTO struct {
	UserId int
}

// IsAdmin checks if user is admin
//
// If user doesn't exist, returns error ErrInvalidUserId
func (s *Service) IsAdmin(ctx context.Context, dto IsAdminDTO) (bool, error) {
	_, err := s.adminRepo.GetAdmin(ctx, dto.UserId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return false, fmt.Errorf("failed to check if user is admin: %w", ErrInvalidUserId)
		}

		return false, fmt.Errorf("failed to check if user is admin: %w", err)
	}

	return true, nil
}
