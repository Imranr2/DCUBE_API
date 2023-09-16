package application

import (
	"encoding/json"
	"net/http"
	"strconv"

	dcubeerrs "github.com/Imranr2/DCUBE_API/internal/errors"
	"github.com/Imranr2/DCUBE_API/internal/urlshortener"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var userManager user.UserManager
var urlShortenerManager urlshortener.URLShortenerManager

func InitApp(db *gorm.DB) {
	userManager = user.NewUserManager(db)
	urlShortenerManager = urlshortener.NewURLShortenerManager(db)
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

	respondWithJSON(w, http.StatusCreated, resp)
	return
}

func GetURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
		return
	}

	getRequest := urlshortener.GetRequest{UserID: userID}
	resp, err := urlShortenerManager.GetURL(getRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
	return
}

func CreateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
		return
	}

	var createRequest urlshortener.CreateRequest
	json.NewDecoder(r.Body).Decode(&createRequest)

	err := validateParams(createRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	createRequest.UserID = userID
	resp, err := urlShortenerManager.CreateURL(createRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)
	return
}

func DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
	}

	params := mux.Vars(r)
	urlID, ok := params["id"]

	if !ok {
		respondWithError(w, dcubeerrs.New(http.StatusBadRequest, "Missing URL ID"))
		return
	}

	u64, e := strconv.ParseUint(urlID, 10, 64)

	if e != nil {
		respondWithError(w, dcubeerrs.New(http.StatusBadRequest, "ID is not an unsigned integer"))
		return
	}

	var deleteRequest urlshortener.DeleteRequest
	deleteRequest.UserID = userID
	deleteRequest.ID = uint(u64)

	resp, err := urlShortenerManager.DeleteURL(deleteRequest)

	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
	return
}
