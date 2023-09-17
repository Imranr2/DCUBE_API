package application

import (
	"context"
	"net/http"

	"github.com/Imranr2/DCUBE_API/internal/session"
)

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

func setCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("user_id").(uint)

		newToken, err := session.GenerateToken(userID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    newToken.TokenString,
			Expires:  newToken.ExpirationTime,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		next.ServeHTTP(w, r)
	})
}
