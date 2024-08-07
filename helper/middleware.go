package helper

import (
	"github.com/go-chi/jwtauth"
	"net/http"
)

var (
	tokenAuth *jwtauth.JWTAuth
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("poke-secret"), nil)
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}

func TokenAuth(next http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)(next)
}

func RoleMiddleware(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
