package userspokemon

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrorUserPokemonNotFound  = errors.New("users pokemons not found")
	ErrPokemonCatchFailed     = errors.New("pokemon catch failed")
	ErrPokemonNotFound        = errors.New("pokemon not found")
	ErrPokemonAlreadyReleased = errors.New("pokemon already released")
	ErrPokemonNotReleased     = errors.New("pokemon not released")
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("Cannot assign nil pool")
	}

	pool = newPool

	return nil
}
