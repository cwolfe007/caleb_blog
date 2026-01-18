package main

import (
	"crypto/subtle"
	"log"
	"net/http"
)

func main() {
	store := NewStore()
	handlers := NewHandlers(store)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /post/{id}", handlers.PostHandler)

	// Protected admin routes
	mux.HandleFunc("GET /admin", BasicAuth(handlers.AdminHandler))
	mux.HandleFunc("GET /admin/create", BasicAuth(handlers.CreateHandler))
	mux.HandleFunc("POST /admin/create", BasicAuth(handlers.CreateHandler))
	mux.HandleFunc("POST /admin/delete/{id}", BasicAuth(handlers.DeleteHandler))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte("writer")) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte("password123")) == 1

		if !usernameMatch || !passwordMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
