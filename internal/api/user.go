package v1

import (
	"context"
	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/avalance-rl/otiva/services/auth/internal/api/dto"
	usecase "github.com/avalance-rl/otiva/services/auth/internal/domain/usecase/user"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type UserUsecase interface {
	Create(ctx context.Context, user usecase.CreateDTO) error
	Authenticate(ctx context.Context, email, password string) (string, error)
}

type userHandler struct {
	userUsecase UserUsecase
	log         *logger.Logger
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewUserHandler(userUsecase UserUsecase, log *logger.Logger) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
		log:         log,
	}
}

func (u *userHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /user", u.register)
}

func (u *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := usecase.CreateDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	err := u.userUsecase.Create(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (u *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := u.userUsecase.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := dto.UserLoginResp{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
