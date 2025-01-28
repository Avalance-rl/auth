package user

import (
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
)

type CreateDTO struct {
	Email    string
	Password string
}

func (u CreateDTO) ConvertToEntity() entity.User {
	return entity.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
