package users

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrorUserNotFound = errors.New("user not found")
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("Cannot assign nil pool")
	}

	pool = newPool

	return nil
}
