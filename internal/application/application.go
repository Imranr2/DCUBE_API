package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	dcubeerrs "github.com/Imranr2/DCUBE_API/internal/errors"
	"github.com/Imranr2/DCUBE_API/internal/urlshortener"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"github.com/go-playground/validator"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var userManager user.UserManager
var urlShortenerManager urlshortener.URLShortenerManager

type Application struct {
	router *mux.Router
}

func (app *Application) InitApp(db *gorm.DB) {
	app.initManagers(db)
	app.router = mux.NewRouter()
	app.initRoutes()
}

func (app *Application) Run() {
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{
		"Access-Control-Allow-Headers",
		"X-Requested-With",
		"Content-Type",
		"Authorization",
		"Accept",
	})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	url := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(url, handlers.CORS(credentials, headers, methods, origins)(app.router)))
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest user.Request
	json.NewDecoder(r.Body).Decode(&loginRequest)

	err := app.validateParams(loginRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	resp, err := userManager.Login(loginRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest user.Request
	json.NewDecoder(r.Body).Decode(&registerRequest)

	err := app.validateParams(registerRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	resp, err := userManager.Register(registerRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	app.respondWithJSON(w, http.StatusCreated, resp)
}

func (app *Application) GetURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		app.respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
		return
	}

	getRequest := urlshortener.GetRequest{UserID: userID}
	resp, err := urlShortenerManager.GetURL(getRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *Application) CreateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		app.respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
		return
	}

	var createRequest urlshortener.CreateRequest
	json.NewDecoder(r.Body).Decode(&createRequest)

	err := app.validateParams(createRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	createRequest.UserID = userID
	resp, err := urlShortenerManager.CreateURL(createRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	app.respondWithJSON(w, http.StatusCreated, resp)
}

func (app *Application) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)

	if !ok {
		app.respondWithError(w, dcubeerrs.New(http.StatusInternalServerError, "Invalid user id"))
	}

	params := mux.Vars(r)
	urlID, ok := params["id"]

	if !ok {
		app.respondWithError(w, dcubeerrs.New(http.StatusBadRequest, "Missing URL ID"))
		return
	}

	u64, e := strconv.ParseUint(urlID, 10, 64)

	if e != nil {
		app.respondWithError(w, dcubeerrs.New(http.StatusBadRequest, "ID is not an unsigned integer"))
		return
	}

	var deleteRequest urlshortener.DeleteRequest
	deleteRequest.UserID = userID
	deleteRequest.ID = uint(u64)

	resp, err := urlShortenerManager.DeleteURL(deleteRequest)

	if err != nil {
		app.respondWithError(w, err)
		return
	}

	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *Application) initManagers(db *gorm.DB) {
	userManager = user.NewUserManager(db)
	urlShortenerManager = urlshortener.NewURLShortenerManager(db)
}

func (app *Application) initRoutes() {
	app.router.Use(commonMiddleware)
	app.router.HandleFunc("/register", app.Register).Methods(http.MethodPost)

	api := app.router.PathPrefix("/url").Subrouter()
	api.Use(tokenValidatorMiddleware)
	api.Use(setCookieMiddleware)
	api.HandleFunc("", app.GetURLs).Methods(http.MethodGet)
	api.HandleFunc("", app.CreateURL).Methods(http.MethodPost)
	api.HandleFunc("/{id}", app.DeleteURL).Methods(http.MethodDelete)

	login := app.router.Path("/login").Subrouter()
	login.Use(setCookieMiddleware)
	login.HandleFunc("", app.Login).Methods(http.MethodPost)
}

func (app *Application) validateParams(s interface{}) dcubeerrs.Error {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		return dcubeerrs.New(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (app *Application) respondWithError(w http.ResponseWriter, err dcubeerrs.Error) {
	app.respondWithJSON(w, err.StatusCode(), map[string]string{"error": err.Message()})
}

func (app *Application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
