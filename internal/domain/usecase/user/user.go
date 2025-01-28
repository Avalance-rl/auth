package user

import (
	"context"
	"fmt"
	"time"

	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, uuid string) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Update(ctx context.Context, uuid string, fieldOfUpdates map[string]any) error
	Delete(ctx context.Context, uuid string) error
}

type TokenService interface {
	Create(user entity.User) (string, error)
	Parse(token string) (entity.AccessToken, error)
}

type userUsecase struct {
	service       Service
	tokenService  TokenService
	signingMethod jwt.SigningMethodHMAC
	secret        string
	expiration    time.Duration
	log           *logger.Logger
}

func New(
	service Service,
	tokenService TokenService,
	signingMethod jwt.SigningMethodHMAC,
	secret string,
	exp time.Duration,
	log *logger.Logger,
) *userUsecase {
	return &userUsecase{
		service:       service,
		tokenService:  tokenService,
		signingMethod: signingMethod,
		secret:        secret,
		expiration:    exp,
		log:           log,
	}
}

func (u userUsecase) Create(ctx context.Context, user CreateDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	err = u.service.Create(ctx, user.ConvertToEntity())
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (u userUsecase) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := u.service.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := u.tokenService.Create(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u userUsecase) ValidateToken(ctx context.Context, token string) (entity.AccessToken, error) {
	tokenParsed, err := u.tokenService.Parse(token)
	if err != nil {
		return entity.AccessToken{}, err
	}

	return tokenParsed, nil
}

func (u userUsecase) GetByID(ctx context.Context, uuid string) (entity.User, error) {
	user, err := u.service.GetByID(ctx, uuid)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u userUsecase) Update(ctx context.Context, uuid string, fieldOfUpdates map[string]any) error {
	err := u.service.Update(ctx, uuid, fieldOfUpdates)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u userUsecase) Delete(ctx context.Context, uuid string) error {
	err := u.service.Delete(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
