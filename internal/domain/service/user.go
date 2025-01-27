package service

import (
	"context"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
)

type UserCreator interface {
	Create(ctx context.Context, user entity.User) error
}

type UserGetter interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type UserUpdater interface {
	Update(ctx context.Context, fieldOfUpdates map[string]any) error
}

type UserDeleter interface {
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userCreator UserCreator
	userGetter  UserGetter
	userUpdater UserUpdater
	userDeleter UserDeleter
}

func NewUserService(
	creator UserCreator,
	getter UserGetter,
	updater UserUpdater,
	deleter UserDeleter,
) *userService {
	return &userService{
		userCreator: creator,
		userGetter:  getter,
		userUpdater: updater,
		userDeleter: deleter,
	}
}

func (s *userService) Create(ctx context.Context, user entity.User) error {
	return s.userCreator.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (entity.User, error) {
	user, err := s.userGetter.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, fieldOfUpdates map[string]any) error {
	return s.userUpdater.Update(ctx, fieldOfUpdates)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userDeleter.Delete(ctx, id)
}
