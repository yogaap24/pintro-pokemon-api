package pokemon

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"mda/helper"
	"net/http"
	"strconv"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(helper.TokenAuth)
	
	r.Get("/", listPokemonHandler)
	r.Get("/{id}", getPokemonHandler)
	r.Get("/{name}", getPokemonHandler)

	return r
}

func writeMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeMessage(w, status, err.Error())
}

func listPokemonHandler(w http.ResponseWriter, req *http.Request) {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	const (
		defaultLimit  = 10
		defaultOffset = 0
	)

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	pokemons, err := listPokemon(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(pokemons)
	if err != nil {
		return
	}
}

func getPokemonHandler(w http.ResponseWriter, req *http.Request) {
	idStr := chi.URLParam(req, "id")
	name := chi.URLParam(req, "name")

	var (
		pokemon interface{}
		err     error
	)

	if name != "" {
		pokemon, err = findPokemonByName(name)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		pokemon, err = findPokemonById(id)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(pokemon)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}
