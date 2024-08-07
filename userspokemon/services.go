package userspokemon

import (
	"context"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"math/rand"
	"mda/helper"
	"time"
)

func catchPokemon(ctx context.Context, userId ulid.ULID, pokemonId int, nickname string) (UserPokemon, error) {
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() > 0.5 {
		return UserPokemon{}, ErrPokemonCatchFailed
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return UserPokemon{}, err
	}

	userPokemon, err := NewPokemon(userId, pokemonId, nickname)
	if err != nil {
		return UserPokemon{}, err
	}

	err = saveUserPokemon(ctx, tx, userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return UserPokemon{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return UserPokemon{}, err
	}

	return userPokemon, nil
}

func listUserPokemons(ctx context.Context, userId ulid.ULID) ([]UserPokemon, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	userPokemons, err := findUserPokemonByUserId(ctx, tx, userId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return userPokemons, nil
}

func releasePokemon(ctx context.Context, userPokemonId ulid.ULID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	userPokemon, err := findUserPokemonById(ctx, tx, userPokemonId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if userPokemon.Released {
		tx.Rollback(ctx)
		return ErrPokemonAlreadyReleased
	}

	primeGen := helper.NewPrimeGenerator()
	primeNumber, primeErr := primeGen.GetUniquePrime()

	if primeErr != nil {
		tx.Rollback(ctx)
		return primeErr
	}

	if err := helper.IsPrime(primeNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = ReleasePokemon(&userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = saveUserPokemon(ctx, tx, userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	helper.DefaultAttempts++
	if helper.DefaultAttempts == helper.DefaultThreshold {
		helper.DefaultThreshold += 6
	}

	return nil
}

func unReleasePokemon(ctx context.Context, userPokemonId ulid.ULID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	userPokemon, err := findUserPokemonById(ctx, tx, userPokemonId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if !userPokemon.Released {
		tx.Rollback(ctx)
		return ErrPokemonNotReleased
	}

	err = UnReleasePokemon(&userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = saveUserPokemon(ctx, tx, userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func updatePokemon(ctx context.Context, userPokemonId ulid.ULID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			log.Error().Msgf("Recovered in updatePokemon: %v", r)
		}
	}()

	userPokemon, err := findUserPokemonById(ctx, tx, userPokemonId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if userPokemon.Released {
		tx.Rollback(ctx)
		return ErrPokemonAlreadyReleased
	}

	fibValue := helper.GetNextFibonacciValue(userPokemon.Id.String())
	newNickname := helper.GenerateNickName(userPokemon.Nickname, fibValue)

	err = UpdatePokemon(&userPokemon, newNickname)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = saveUserPokemon(ctx, tx, userPokemon)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
