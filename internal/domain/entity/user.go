package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
