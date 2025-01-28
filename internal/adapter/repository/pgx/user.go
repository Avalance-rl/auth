package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/avalance-rl/otiva-pkg/database"
	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/avalance-rl/otiva/services/auth/internal/adapter/repository/pgx/model"
	"github.com/avalance-rl/otiva/services/auth/internal/domain/entity"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	pgxPool *database.PgxPool
	log     *logger.Logger
}

func NewProductRepository(pgxPool *database.PgxPool, log *logger.Logger) *userRepository {
	return &userRepository{
		pgxPool: pgxPool,
		log:     log,
	}
}

func (u *userRepository) Create(ctx context.Context, user entity.User) error {
	userRepo := model.User{}
	userRepo.ConvertFromEntity(user)

	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("users").
		Cols(
			"uuid",
			"email",
			"password",
			"created_at",
			"updated_at",
		).
		Values(
			userRepo.UUID,
			userRepo.Email,
			userRepo.Password,
			userRepo.CreatedAt,
			userRepo.UpdatedAt,
		)

	sql, args := sb.Build()

	err := u.pgxPool.Pool.QueryRow(ctx, sql, args...).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("user already exists")
			}
		}

		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (u *userRepository) FindByID(ctx context.Context, uuid string) (entity.User, error) {
	var userRepo model.User

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"uuid",
		"email",
		"password",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(sb.Equal("uuid", uuid))

	sql, args := sb.Build()

	row := u.pgxPool.Pool.QueryRow(ctx, sql, args...)
	err := row.Scan(
		&userRepo.UUID,
		&userRepo.Email,
		&userRepo.Password,
		&userRepo.CreatedAt,
		&userRepo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user with ID %s not found", uuid)
		}
		return entity.User{}, fmt.Errorf("failed to find user by ID: %w", err)
	}

	user := userRepo.ConvertToEntity()

	return user, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var userRepo model.User

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"uuid",
		"email",
		"password",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(sb.Equal("email", email))

	sql, args := sb.Build()

	row := u.pgxPool.Pool.QueryRow(ctx, sql, args...)
	err := row.Scan(
		&userRepo.UUID,
		&userRepo.Email,
		&userRepo.Password,
		&userRepo.CreatedAt,
		&userRepo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return entity.User{}, fmt.Errorf("failed to find user by ID: %w", err)
	}

	user := userRepo.ConvertToEntity()

	return user, nil
}

func (u *userRepository) Update(ctx context.Context, uuid string, fieldOfUpdates map[string]any) error {
	if len(fieldOfUpdates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	ub := sqlbuilder.NewUpdateBuilder()
	ub.SetFlavor(sqlbuilder.PostgreSQL)
	ub.Update("users")

	for field, value := range fieldOfUpdates {
		ub.Set(ub.Assign(field, value))
	}

	ub.Where(ub.Equal("uuid", uuid))

	sql, args := ub.Build()

	_, err := u.pgxPool.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("failed to update user: %v", pgErr.Message)
		}
		return fmt.Errorf("unexpected error during update: %w", err)
	}

	return nil
}

func (u *userRepository) Delete(ctx context.Context, uuid string) error {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.PostgreSQL)
	db.DeleteFrom("users").Where(db.Equal("uuid", uuid))

	sql, args := db.Build()

	commandTag, err := u.pgxPool.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with UUID %s not found", uuid)
	}

	return nil
}
