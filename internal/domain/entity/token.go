package entity

import (
	"time"

	"github.com/google/uuid"
)

type userRole string

const (
	userRoleUser  userRole = "user"
	userRoleAdmin userRole = "admin"
)

type AccessToken struct {
	UserUUID  uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
	UserRole  userRole
}

func (at *AccessToken) Expired() bool {
	return at.ExpiresAt.Before(time.Now())
}

func (at *AccessToken) UserRoleToString() string {
	return string(at.UserRole)
}

func (at *AccessToken) SetUserRoleFromString(role string) {
	at.UserRole = userRole(role)
}
