package entity

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
