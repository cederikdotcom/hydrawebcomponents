package hydrawebcomponents

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

// RequireAuth wraps an HTTP handler with Bearer token authentication.
// Returns 401 JSON on failure (for API endpoints).
func (w *Web) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if w.IsAuthenticated(r) {
			next(wr, r)
			return
		}
		WriteError(wr, http.StatusUnauthorized, "unauthorized")
	}
}

// RequireWebAuth wraps an HTTP handler with cookie-based authentication.
// Redirects to /login on failure (for web UI pages).
func (w *Web) RequireWebAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if w.IsAuthenticated(r) {
			next(wr, r)
			return
		}
		http.Redirect(wr, r, "/login", http.StatusSeeOther)
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
