package users

import (
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
	"time"
)

func (u User) MarshalJSON() ([]byte, error) {
	var j struct {
		Id        ulid.ULID  `json:"id"`
		Username  string     `json:"username"`
		Password  string     `json:"password"`
		Role      string     `json:"role"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		DeletedAt *time.Time `json:"deleted_at,omitempty"`
	}

	j.Id = u.Id
	j.Username = u.Username
	j.Password = u.Password
	j.Role = u.Role
	j.CreatedAt = u.CreatedAt
	j.UpdatedAt = u.UpdatedAt.Ptr()
	j.DeletedAt = u.DeletedAt.Ptr()

	return json.Marshal(j)
}

func parseNullStringToNullTime(s null.String) (t null.Time) {
	if !s.Valid {
		return
	}

	ts, err := time.Parse(time.RFC3339, s.String)

	if err != nil {
		return
	}

	return null.TimeFrom(ts)
}

func (u *User) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID   `json:"id"`
		Username  string      `json:"username"`
		Password  string      `json:"password"`
		Role      string      `json:"role"`
		CreatedAt string      `json:"created_at"`
		UpdatedAt null.String `json:"updated_at"`
		DeletedAt null.String `json:"deleted_at"`
	}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return err
	}

	CreatedAt, err := time.Parse(time.RFC3339, j.CreatedAt)
	if err != nil {
		return err
	}

	UpdatedAt := parseNullStringToNullTime(j.UpdatedAt)
	DeletedAt := parseNullStringToNullTime(j.DeletedAt)

	u = &User{
		Id:        j.Id,
		Username:  j.Username,
		Password:  j.Password,
		Role:      j.Role,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
		DeletedAt: DeletedAt,
	}

	return nil
}
