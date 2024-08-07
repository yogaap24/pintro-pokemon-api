package pokemon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const PokeAPIURL = "https://pokeapi.co/api/v2/pokemon/"

const (
	defaultLimit  = 10
	defaultOffset = 0
)

type PokemonSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonListResponse struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []PokemonSummary `json:"results"`
}

type Pokemon struct {
	Id                     int           `json:"id"`
	Name                   string        `json:"name"`
	Order                  int           `json:"order"`
	Height                 int           `json:"height"`
	Weight                 int           `json:"weight"`
	BaseExperience         int           `json:"base_experience"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	HeldItems              []interface{} `json:"held_items"`
	Abilities              []interface{} `json:"abilities"`
	Forms                  []interface{} `json:"forms"`
	GameIndices            []interface{} `json:"game_indices"`
	Moves                  []interface{} `json:"moves"`
	Species                interface{}   `json:"species"`
	Sprites                interface{}   `json:"sprites"`
	Stats                  []interface{} `json:"stats"`
	Types                  []interface{} `json:"types"`
}

func extractId(url string) int {
	var id int
	_, err := fmt.Sscanf(url, PokeAPIURL+"%d/", &id)
	if err != nil {
		return 0
	}
	return id
}

func listPokemon(limit, offset int) (*PokemonListResponse, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	if offset < 0 {
		offset = defaultOffset
	}

	resp, err := http.Get(fmt.Sprintf("%s?limit=%d&offset=%d", PokeAPIURL, limit, offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list PokemonListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	for i := range list.Results {
		list.Results[i].Id = extractId(list.Results[i].URL)
	}

	return &list, nil
}

func findPokemonById(id int) (*Pokemon, error) {
	resp, err := http.Get(fmt.Sprintf("%s%d/", PokeAPIURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}

func findPokemonByName(name string) (*Pokemon, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s/", PokeAPIURL, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pokemon Pokemon

	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}
