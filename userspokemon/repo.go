package userspokemon

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func findUserPokemonById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (UserPokemon, error) {
	query := `SELECT id, user_id, pokemon_id, nickname, captured_at, released
				  FROM users_pokemons
			 WHERE id = $1;`

	row := tx.QueryRow(ctx, query, id)

	var userPokemon UserPokemon
	if err := row.Scan(&userPokemon.Id, &userPokemon.UserId, &userPokemon.PokemonId, &userPokemon.Nickname, &userPokemon.CapturedAt, &userPokemon.Released); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserPokemon{}, ErrorUserPokemonNotFound
		}
		return UserPokemon{}, err
	}

	return userPokemon, nil
}

func findUserPokemonByUserId(ctx context.Context, tx pgx.Tx, userId ulid.ULID) ([]UserPokemon, error) {
	query := `SELECT id, user_id, pokemon_id, nickname, captured_at, released
				  FROM users_pokemons
			 WHERE user_id = $1;`

	rows, err := tx.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userPokemons []UserPokemon
	for rows.Next() {
		var userPokemon UserPokemon
		if err := rows.Scan(&userPokemon.Id, &userPokemon.UserId, &userPokemon.PokemonId, &userPokemon.Nickname, &userPokemon.CapturedAt, &userPokemon.Released); err != nil {
			return nil, err
		}

		userPokemons = append(userPokemons, userPokemon)
	}

	return userPokemons, nil
}

func saveUserPokemon(ctx context.Context, tx pgx.Tx, userPokemon UserPokemon) error {
	query := `INSERT INTO users_pokemons (id, user_id, pokemon_id, nickname, captured_at, released)
                  VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (id) DO UPDATE SET
                  user_id = EXCLUDED.user_id,
                  pokemon_id = EXCLUDED.pokemon_id,
                  nickname = EXCLUDED.nickname,
                  captured_at = EXCLUDED.captured_at,
                  released = EXCLUDED.released;`

	_, err := tx.Exec(ctx, query, userPokemon.Id, userPokemon.UserId, userPokemon.PokemonId, userPokemon.Nickname, userPokemon.CapturedAt, userPokemon.Released)
	if err != nil {
		return err
	}

	return nil
}
