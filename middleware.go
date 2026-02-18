package hydrawebcomponents

import (
	"log"
	"net/http"
	"time"

	"github.com/cederikdotcom/hydraapi"
)

// RequireAuth wraps an HTTP handler with Bearer token authentication.
// Returns 401 JSON on failure (for API endpoints).
func (w *Web) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return w.auth.RequireAuth(next)
}

// RequireWebAuth wraps an HTTP handler with cookie-based authentication.
// Redirects to /login on failure (for web UI pages).
func (w *Web) RequireWebAuth(next http.HandlerFunc) http.HandlerFunc {
	return w.auth.RequireWebAuth("/login", next)
}

// LogRequest is HTTP middleware that logs each request's method, path, and duration.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	hydraapi.WriteJSON(w, status, v)
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	hydraapi.WriteError(w, status, msg)
}
