package middleware

import (
	"crypto/subtle"
	"net/http"
)

// AdminBasicAuth guards admin routes with HTTP Basic Auth. This is deliberately
// the simplest possible admin authentication for the MVP (prompt §3): no RBAC,
// no sessions, no tokens — just a shared operator username/password.
func AdminBasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			// Constant-time comparison to avoid leaking length/content via timing.
			userOK := subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1
			passOK := subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1
			if !ok || !userOK || !passOK {
				w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
