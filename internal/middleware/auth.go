// Package middleware provides HTTP middleware for authentication, authorization,
// and rate limiting.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kwasiga/secure-api/internal/auth"
)

type contextKey string

// claimsKey is the context key used to store parsed JWT claims.
const claimsKey contextKey = "claims"

// AuthMiddleware validates the Bearer token in the Authorization header.
// On success it attaches the parsed Claims to the request context.
// Returns 401 if the header is missing, malformed, or the token is invalid.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(header, "Bearer ")

			claims, err := auth.ValidateToken(tokenString, secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves the JWT Claims stored in the request context by AuthMiddleware.
// Returns false if the claims are absent (i.e. the route is not behind AuthMiddleware).
func GetClaims(r *http.Request) (*auth.Claims, bool) {
	claims, ok := r.Context().Value(claimsKey).(*auth.Claims)
	return claims, ok
}
