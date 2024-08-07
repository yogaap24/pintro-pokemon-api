package userspokemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/oklog/ulid/v2"
	"mda/helper"
	"net/http"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(helper.TokenAuth)
	r.Get("/", listUserPokemonsHandler)
	r.Post("/", catchPokemonHandler)
	r.Put("/released/{id}", releasePokemonHandler)
	r.Put("/unreleased/{id}", unReleasePokemonHandler)
	r.Put("/rename/{id}", renamePokemonHandler)

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

func listUserPokemonsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	_, claims, _ := jwtauth.FromContext(ctx)
	IdUser, ok := claims["user_id"].(string)
	if !ok || IdUser == "" {
		writeError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userId, err := ulid.Parse(IdUser)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	userPokemons, err := listUserPokemons(ctx, userId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userPokemons)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func catchPokemonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var j struct {
		PokemonId int `json:"pokemon_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&j)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	_, claims, _ := jwtauth.FromContext(ctx)
	IdUser, ok := claims["user_id"].(string)
	if !ok || IdUser == "" {
		writeError(w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	userId, err := ulid.Parse(IdUser)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	userPokemon, err := catchPokemon(ctx, userId, j.PokemonId, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Message: "Pokemon caught successfully",
		Data:    userPokemon,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func releasePokemonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userPokemonId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userPokemonId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = releasePokemon(ctx, id)

	if errors.Is(err, ErrPokemonNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}

	if errors.Is(err, ErrPokemonAlreadyReleased) {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		if primeErr, ok := err.(*helper.PrimeError); ok {
			writeError(w, http.StatusBadRequest, fmt.Errorf("Pokemon release failed: %v", primeErr))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Pokemon released successfully",
	})

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func unReleasePokemonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userPokemonId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userPokemonId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = unReleasePokemon(ctx, id)

	if errors.Is(err, ErrPokemonNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}

	if errors.Is(err, ErrPokemonNotReleased) {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Pokemon catch back successfully",
	})

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func renamePokemonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userPokemonId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userPokemonId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = updatePokemon(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPokemonNotFound) {
			writeError(w, http.StatusNotFound, err)
		} else if errors.Is(err, ErrPokemonAlreadyReleased) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Pokemon renamed successfully",
	})

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
