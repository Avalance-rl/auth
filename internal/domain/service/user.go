package service

import (
	"context"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
)

type UserCreator interface {
	Create(ctx context.Context, user entity.User) error
}

type UserFinder interface {
	FindByID(ctx context.Context, uuid string) (entity.User, error)
}

type UserUpdater interface {
	Update(ctx context.Context, uuid string, fieldOfUpdates map[string]any) error
}

type UserDeleter interface {
	Delete(ctx context.Context, uuid string) error
}

type userService struct {
	userCreator UserCreator
	userFinder  UserFinder
	userUpdater UserUpdater
	userDeleter UserDeleter
}

func NewUserService(
	creator UserCreator,
	finder UserFinder,
	updater UserUpdater,
	deleter UserDeleter,
) *userService {
	return &userService{
		userCreator: creator,
		userFinder:  finder,
		userUpdater: updater,
		userDeleter: deleter,
	}
}

func (s *userService) Create(ctx context.Context, user entity.User) error {
	return s.userCreator.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, uuid string) (entity.User, error) {
	user, err := s.userFinder.FindByID(ctx, uuid)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, uuid string, fieldOfUpdates map[string]any) error {
	return s.userUpdater.Update(ctx, uuid, fieldOfUpdates)
}

func (s *userService) Delete(ctx context.Context, uuid string) error {
	return s.userDeleter.Delete(ctx, uuid)
}
