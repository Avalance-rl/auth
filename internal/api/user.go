package v1

import (
	"context"
	"github.com/avalance-rl/otiva-pkg/logger"
	authv1 "github.com/avalance-rl/otiva/proto/gen/avalance.auth.v1"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	usecase "github.com/avalance-rl/otiva/services/auth/internal/domain/usecase/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUsecase interface {
	Create(ctx context.Context, user usecase.CreateDTO) error
	Authenticate(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (entity.AccessToken, error)
}

type authServer struct {
	authv1.UnimplementedAuthServiceServer
	userUsecase UserUsecase
	log         *logger.Logger
}

func NewAuthServer(userUsecase UserUsecase, log *logger.Logger) authv1.AuthServiceServer {
	return &authServer{
		userUsecase: userUsecase,
		log:         log,
	}
}

func (s *authServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	user := usecase.CreateDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.userUsecase.Create(ctx, user); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// После успешной регистрации сразу аутентифицируем пользователя
	accessToken, err := s.userUsecase.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to authenticate user")
	}

	return &authv1.AuthResponse{
		Token: accessToken,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	accessToken, err := s.userUsecase.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		s.log.Error("failed to authenticate user", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	_, err = s.userUsecase.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	return &authv1.AuthResponse{
		Token: accessToken,
	}, nil
}

func (s *authServer) ValidateToken(ctx context.Context, req *authv1.TokenRequest) (*authv1.ValidationResponse, error) {
	token, err := s.userUsecase.ValidateToken(ctx, req.Token)
	if err != nil {
		s.log.Error("failed to validate token", zap.Error(err))
		return &authv1.ValidationResponse{
			Valid: false,
		}, nil
	}

	return &authv1.ValidationResponse{
		Valid:  token.Expired(),
		UserId: token.UserUUID.String(),
	}, nil
}
