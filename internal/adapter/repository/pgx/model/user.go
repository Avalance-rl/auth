package model

import (
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ConvertFromEntity(userEntity entity.User) {
	u.UUID = userEntity.UUID
	u.Email = userEntity.Email
	u.Password = userEntity.Password
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) ConvertToEntity() entity.User {
	return entity.User{
		UUID:     u.UUID,
		Username: u.Email,
		Email:    u.Email,
		Password: u.Password,
	}
}
