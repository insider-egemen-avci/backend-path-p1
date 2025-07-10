package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"insider-egemen-avci/backend-path-p1/internal/models"
	"insider-egemen-avci/backend-path-p1/internal/services"
)

type contextKey string

const userContextKey = contextKey("user")

func (h *APIHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		userID, err := validateDummyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := h.userService.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "User not found or error fetching user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}

func validateDummyToken(token string) (int64, error) {
	var id int64
	_, err := fmt.Sscan(token, &id)
	return id, err
}

func (h *APIHandler) RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "User not found in context", http.StatusInternalServerError)
				return
			}

			err := services.CheckRole(user, requiredRoles...)
			if err != nil {
				http.Error(w, "Forbidden: "+err.Error(), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
