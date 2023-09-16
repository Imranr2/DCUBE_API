package session

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const invalidUserID uint = 0
const validDuration = 5 * time.Minute

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
	expirationTime := time.Now().Add(validDuration)

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
	token, err := getToken(r)

	if err != nil {
		return invalidUserID, err
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return invalidUserID, err
	}

	if !tkn.Valid {
		return invalidUserID, errors.New("token has expired")
	}

	return claims.ID, nil
}
