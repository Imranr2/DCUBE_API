package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Imranr2/DCUBE_API/internal/application"
	"github.com/Imranr2/DCUBE_API/internal/session"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func run(r *mux.Router) {
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "X-Requested-With", "Content-Type", "Authorization", "Accept"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", handlers.CORS(credentials, headers, methods, origins)(r)))
}

func InitServer() {
	r := mux.NewRouter()
	initRoutes(r)
	run(r)
}

func initRoutes(r *mux.Router) {
	r.Use(commonMiddleware)
	r.HandleFunc("/login", application.Login).Methods(http.MethodPost)
	r.HandleFunc("/register", application.Register).Methods(http.MethodPost)

	api := r.PathPrefix("/url").Subrouter()
	api.Use(tokenValidatorMiddleware)
	api.HandleFunc("", application.GetURLs).Methods(http.MethodGet)
	api.HandleFunc("", application.CreateURL).Methods(http.MethodPost)
	api.HandleFunc("/{id}", application.DeleteURL).Methods(http.MethodDelete)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func tokenValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := session.VerifyToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext((ctx)))
	})
}

