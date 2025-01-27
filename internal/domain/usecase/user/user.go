package user

import (
	"context"
	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	"github.com/golang-jwt/jwt"
	"time"
)

type Service interface {
	Create(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	Update(ctx context.Context, fieldOfUpdates map[string]any) error
	Delete(ctx context.Context, id string) error
}

type JWTService interface {
	Create(user entity.User, signingMethod jwt.SigningMethod, exp time.Duration, secret string) (string, error)
	Parse(token string, signingMethod jwt.SigningMethod, exp time.Duration, secret string) (entity.AccessToken, error)
	GetAccessToken(accessHead string) (string, error)
}

type userUsecase struct {
	service       Service
	jwtService    JWTService
	signingMethod jwt.SigningMethodHMAC
	secret        string
	expiration    time.Duration
	log           *logger.Logger
}

func New(
	service Service,
	jwtService JWTService,
	signingMethod jwt.SigningMethodHMAC,
	secret string,
	exp time.Duration,
	log *logger.Logger,
) *userUsecase {
	return &userUsecase{
		service:       service,
		jwtService:    jwtService,
		signingMethod: signingMethod,
		secret:        secret,
		expiration:    exp,
		log:           log,
	}
}
