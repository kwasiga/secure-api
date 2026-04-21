package middleware

import (
	"net/http"
	"slices"
)

// RequireRole restricts a route to users whose JWT role claim matches one of the
// provided roles. Must be used after AuthMiddleware (claims must already be in context).
// Returns 401 if claims are missing, 403 if the role is not permitted.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !slices.Contains(roles, claims.Role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
