package pgx

import (
	"github.com/avalance-rl/otiva-pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

func NewProductRepository(db *pgxpool.Pool, log *logger.Logger) *userRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}
