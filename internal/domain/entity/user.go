package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
