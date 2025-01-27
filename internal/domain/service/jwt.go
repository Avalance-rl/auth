package service

import (
	"fmt"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type token struct {
	signingMethod jwt.SigningMethod
	expiration    time.Duration
	secret        string
}

func New(
	signingMethod jwt.SigningMethod,
	expiration time.Duration,
	secret string,
) *token {
	return &token{
		signingMethod: signingMethod,
		expiration:    expiration,
		secret:        secret,
	}
}

func (t *token) Create(
	user entity.User,
	signingMethod jwt.SigningMethod,
	exp time.Duration,
	secret string,
) (string, error) {
	now := time.Now()
	accessPayload := jwt.MapClaims{
		"sub": user.UUID.String(),
		"iat": now.Unix(),
		"exp": now.Add(exp).Unix(),
	}
	accessToken := jwt.NewWithClaims(signingMethod, accessPayload)

	jwtSecretKey := []byte(secret)
	accessTokenString, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func (t *token) Parse(
	token string,
	signingMethod jwt.SigningMethod,
	exp time.Duration,
	secret string,
) (entity.AccessToken, error) {
	claim := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != signingMethod.Alg() {
			return nil, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}

	sub, ok := claim["sub"].(string)
	if !ok {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}

	role, ok := claim["role"].(string)
	if !ok {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}

	unixExpiresAt, ok := claim["exp"].(float64)
	if !ok {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}
	expiresAt := time.Unix(int64(unixExpiresAt), 0)

	unixIssuedAt, ok := claim["iat"].(float64)
	if !ok {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}

	issuedAt := time.Unix(int64(unixIssuedAt), 0)

	accessToken := entity.AccessToken{}
	accessToken.SetUserRoleFromString(role)
	userUUID, err := uuid.Parse(sub)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}
	accessToken.UserUUID = userUUID
	accessToken.ExpiresAt = expiresAt
	accessToken.IssuedAt = issuedAt

	return accessToken, nil
}

func (t *token) GetAccessToken(accessHead string) (string, error) {
	const bearer = "Bearer "
	if !strings.HasPrefix(accessHead, bearer) {
		return "", fmt.Errorf("ErrInvalidToken.NewWithNoMessage()")
	}
	accessToken := strings.TrimPrefix(accessHead, bearer)

	return accessToken, nil
}
