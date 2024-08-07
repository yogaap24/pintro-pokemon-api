package userspokemon

import (
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"time"
)

func (u UserPokemon) MarshalJSON() ([]byte, error) {
	var j struct {
		Id         ulid.ULID `json:"id"`
		UserId     ulid.ULID `json:"user_id"`
		PokemonId  int       `json:"pokemon_id"`
		Nickname   string    `json:"nickname"`
		CapturedAt time.Time `json:"captured_at"`
		Released   bool      `json:"released"`
	}

	j.Id = u.Id
	j.UserId = u.UserId
	j.PokemonId = u.PokemonId
	j.Nickname = u.Nickname
	j.CapturedAt = u.CapturedAt
	j.Released = u.Released

	return json.Marshal(j)
}

func (u *UserPokemon) UnmarshalJSON(data []byte) error {
	var j struct {
		Id         ulid.ULID `json:"id"`
		UserId     ulid.ULID `json:"user_id"`
		PokemonId  int       `json:"pokemon_id"`
		Nickname   string    `json:"nickname"`
		CapturedAt string    `json:"captured_at"`
		Released   bool      `json:"released"`
	}

	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}

	CapturedAt, err := time.Parse(time.RFC3339, j.CapturedAt)
	if err != nil {
		return err
	}

	u.Id = j.Id
	u.UserId = j.UserId
	u.PokemonId = j.PokemonId
	u.Nickname = j.Nickname
	u.CapturedAt = CapturedAt
	u.Released = j.Released

	return nil
}
