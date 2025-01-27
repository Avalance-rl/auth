package user

import (
	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/golang-jwt/jwt"
	"time"
)

type Service interface {
}

type JWTService interface {
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
