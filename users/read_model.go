package users

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
	"time"
)

type UserList struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

var emptyList UserList

func findAllUsers(ctx context.Context, tx pgx.Tx) (UserList, error) {
	var userCount int

	row := tx.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE deleted_at IS NULL")
	err := row.Scan(&userCount)

	if err != nil {
		return emptyList, err
	}

	if userCount == 0 {
		return emptyList, nil
	}

	users := make([]User, userCount)

	rows, err := tx.Query(
		ctx,
		"SELECT id, username, password, role, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL",
	)

	if err != nil {
		return emptyList, err
	}

	defer rows.Close()

	var i int

	for i = range users {
		var id ulid.ULID
		var username, password string
		var role string
		var createdAt time.Time
		var updatedAt, deletedAt null.Time

		if !rows.Next() {
			break
		}

		if err := rows.Scan(&id, &username, &password, &role, &createdAt, &updatedAt, &deletedAt); err != nil {
			return emptyList, err
		}

		users[i] = User{
			Id:        id,
			Username:  username,
			Password:  password,
			Role:      role,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		}
	}

	list := UserList{
		Users: users,
		Count: userCount,
	}

	return list, nil
}
