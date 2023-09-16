package session

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const invalidUserId uint = 0

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

type Session struct {
	TokenString    string
	ExpirationTime time.Time
}

func getToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")

	if err != nil {
		return "", err
	}

	return c.Value, err
}

func GenerateToken(id uint) (Session, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return Session{}, err
	}

	return Session{
		TokenString:    tokenString,
		ExpirationTime: expirationTime,
	}, nil
}

func VerifyToken(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		return invalidUserId, err
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return invalidUserId, err
	}

	if !tkn.Valid {
		return invalidUserId, errors.New("Token has expired")
	}

	return claims.ID, nil
}
