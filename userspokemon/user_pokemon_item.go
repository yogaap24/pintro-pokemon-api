package userspokemon

import (
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"net/http"
	"strconv"
	"time"
)

const PokeAPIURL = "https://pokeapi.co/api/v2/pokemon/"

type UserPokemon struct {
	Id         ulid.ULID
	UserId     ulid.ULID
	PokemonId  int
	Nickname   string
	CapturedAt time.Time
	Released   bool
}

func NewPokemon(userId ulid.ULID, pokemonId int, nickname string) (UserPokemon, error) {
	if nickname == "" {
		name, err := getPokemonName(pokemonId)
		if err != nil {
			return UserPokemon{}, err
		}

		nickname = name
	}

	id, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		return UserPokemon{}, err
	}

	return UserPokemon{
		Id:         id,
		UserId:     userId,
		PokemonId:  pokemonId,
		Nickname:   nickname,
		CapturedAt: time.Now(),
		Released:   false,
	}, nil
}

func ReleasePokemon(userPokemon *UserPokemon) error {
	if userPokemon.Released {
		return ErrPokemonAlreadyReleased
	}
	userPokemon.Released = true
	return nil
}

func UpdatePokemon(userPokemon *UserPokemon, nickname string) error {
	if userPokemon.Released {
		return ErrPokemonAlreadyReleased
	}
	userPokemon.Nickname = nickname
	return nil
}

func UnReleasePokemon(userPokemon *UserPokemon) error {
	if !userPokemon.Released {
		return ErrPokemonNotReleased
	}
	userPokemon.Released = false
	return nil
}

func getPokemonName(pokemonId int) (string, error) {
	resp, err := http.Get(PokeAPIURL + strconv.Itoa(pokemonId))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Name, nil
}
