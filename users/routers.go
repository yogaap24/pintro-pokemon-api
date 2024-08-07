package users

import (
	"context"
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

	r.Post("/login", loginHandler)

	r.Group(func(r chi.Router) {
		r.Use(helper.TokenAuth)
		r.Get("/profile", getProfileHandler)

		r.Group(func(r chi.Router) {
			r.Use(helper.RoleMiddleware(helper.RoleAdmin))
			r.Get("/", listUsersHandler)
			r.Post("/", createUserHandler)
			r.Get("/{id}", getUserHandler)
			r.Put("/{id}", updateUserHandler)
			r.Delete("/{id}", deleteUserHandler)
		})
	})

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

func loginHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	var j struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(req.Body).Decode(&j)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := authenticate(ctx, j.Username, j.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	claims := map[string]interface{}{
		"user_id": user.Id.String(),
		"role":    user.Role,
	}
	_, tokenString, err := helper.GetTokenAuth().Encode(claims)

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func listUsersHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	users, err := listUsers(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := findUser(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func getProfileHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	_, claims, _ := jwtauth.FromContext(ctx)
	currentUserId, userIdOk := claims["user_id"].(string)

	if !userIdOk {
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	currentUserIdUlid, err := ulid.Parse(currentUserId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("invalid user_id in token: %v", err))
		return
	}

	user, err := findUser(ctx, currentUserIdUlid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("error fetching user: %v", err))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var j struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(req.Body).Decode(&j)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := createUser(ctx, j.Username, j.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func updateUserHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var j struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = json.NewDecoder(req.Body).Decode(&j)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := updateUser(ctx, id, j.Username, j.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId := chi.URLParam(req, "id")

	id, err := ulid.Parse(userId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = deleteUser(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeMessage(w, http.StatusOK, "user deleted")
}
