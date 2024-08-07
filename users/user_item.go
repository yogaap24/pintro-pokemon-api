package users

import (
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
	"time"
)

type User struct {
	Id        ulid.ULID
	Username  string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt null.Time
	DeletedAt null.Time
}

func NewUser(username, password string) (User, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		return User{}, err
	}

	return User{
		Id:        id,
		Username:  username,
		Password:  password,
		Role:      "user",
		CreatedAt: time.Now(),
	}, nil
}

func UpdateUser(user *User, username, password string) error {
	user.Username = username
	user.Password = password
	user.UpdatedAt = null.TimeFrom(time.Now())

	return nil
}

func DeleteUser(user *User) {
	user.DeletedAt = null.TimeFrom(time.Now())
}
