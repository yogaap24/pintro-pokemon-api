package users

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func findUserById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (User, error) {
	query := `SELECT id, username, password, role, created_at, updated_at, deleted_at 
				FROM users WHERE id = $1 
 			  AND deleted_at IS NULL`

	row := tx.QueryRow(ctx, query, id)

	var user User
	if err := row.Scan(
		&user.Id, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrorUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func findUserByUsernameAndPassword(ctx context.Context, tx pgx.Tx, username, password string) (User, error) {
	query := `SELECT id, username, password, role, created_at, updated_at, deleted_at 
				FROM users WHERE username = $1 AND password = $2
			  AND deleted_at IS NULL`

	row := tx.QueryRow(ctx, query, username, password)

	var user User
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrorUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func checkAdminExists(ctx context.Context, tx pgx.Tx) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin')`

	var exists bool
	if err := tx.QueryRow(ctx, query).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func saveUser(ctx context.Context, tx pgx.Tx, user User) error {
	query := `INSERT INTO users (id, username, password, role, created_at, updated_at, deleted_at) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (id) DO UPDATE SET
					username = $2,
					password = $3,
					updated_at = COALESCE(EXCLUDED.updated_at, users.updated_at),
					deleted_at = 
						CASE WHEN EXCLUDED.deleted_at IS NOT NULL 
					THEN EXCLUDED.deleted_at ELSE users.deleted_at END;`

	_, err := tx.Exec(ctx, query, user.Id, user.Username, user.Password, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}
