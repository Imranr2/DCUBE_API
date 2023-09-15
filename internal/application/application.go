package application

import (
	"encoding/json"
	"net/http"

	dcubeerrs "github.com/Imranr2/DCUBE_API/internal/errors"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

var userManager user.UserManager

func InitApp(db *gorm.DB) {
	userManager = user.NewUserManager(db)
}

func validateParams(s interface{}) dcubeerrs.Error {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		return dcubeerrs.New(http.StatusBadRequest, err.Error())
	}
	return nil
}

func respondWithError(w http.ResponseWriter, err dcubeerrs.Error) {
	respondWithJSON(w, err.StatusCode(), map[string]string{"error": err.Message()})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest user.Request
	json.NewDecoder(r.Body).Decode(&loginRequest)

	err := validateParams(loginRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	resp, err := userManager.Login(loginRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
	return
}

func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest user.Request
	json.NewDecoder(r.Body).Decode(&registerRequest)

	err := validateParams(registerRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	resp, err := userManager.Register(registerRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
	return
}

func GetURLs(w http.ResponseWriter, r *http.Request) {
	return
}

func CreateURL(w http.ResponseWriter, r *http.Request) {
	return
}

func DeleteURL(w http.ResponseWriter, r *http.Request) {
	return
}
