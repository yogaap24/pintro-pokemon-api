package users

import (
	"context"
	"github.com/oklog/ulid/v2"
)

func authenticate(ctx context.Context, username, password string) (User, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	user, err := findUserByUsernameAndPassword(ctx, tx, username, password)
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	return user, nil
}

func listUsers(ctx context.Context) (UserList, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return UserList{}, err
	}

	list, err := findAllUsers(ctx, tx)
	if err != nil {
		return UserList{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return UserList{}, err
	}

	return list, nil
}

func findUser(ctx context.Context, id ulid.ULID) (User, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}

	user, err := findUserById(ctx, tx, id)
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	return user, nil
}

func createUser(ctx context.Context, username, password string) (user User, err error) {
	userItem, err := NewUser(username, password)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveUser(ctx, tx, userItem)
	if err != nil {
		tx.Rollback(ctx)
		return
	}

	tx.Commit(ctx)

	return userItem, nil
}

func updateUser(ctx context.Context, id ulid.ULID, username, password string) (user User, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	user, err = findUserById(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return User{}, err
	}

	err = UpdateUser(&user, username, password)
	if err != nil {
		tx.Rollback(ctx)
		return User{}, err
	}

	err = saveUser(ctx, tx, user)
	if err != nil {
		tx.Rollback(ctx)
		return User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	return user, nil
}

func deleteUser(ctx context.Context, id ulid.ULID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	user, err := findUserById(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	DeleteUser(&user)

	err = saveUser(ctx, tx, user)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}
