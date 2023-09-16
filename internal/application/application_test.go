package application

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/Imranr2/DCUBE_API/internal/session"
	"github.com/Imranr2/DCUBE_API/internal/urlshortener"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var users = []user.User{{
	ID: 1,
	Username: "test1",
	Password: "$2a$10$o4xsT2RBlIrK62FQkuPTcOs5NbPefWTz9pq4hU42UGZRopgCB2K4S",
}, {
	ID: 2,
	Username: "test2",
	Password: "$2a$10$tDDglbfPHHaBWYTq8mp1LutPsA/.Zz5Tfld0pwGaSXMIgMEU7kRKC",
}, {
	ID: 3,
	Username: "test4",
	Password: "$2a$10$tDDglbfPHHaBWYTq8mp1LutPsA/.Zz5Tfld0pwGaSXMIgMEU7kRKC",
}}

var urls = []urlshortener.ShortenedURL{{
	ID: 1,
	Original: "https://www.google.com/maps",
	Shortened: "dcu.be/test1",
	UserID: 1,
}, {
	ID: 2,
	Original: "https://www.youtube.com/test",
	Shortened: "dcu.be/youtubeTest",
	UserID: 1,
}, {
	ID: 3,
	Original: "https://www.google.com/maps",
	Shortened: "dcu.be/test4",
	UserID: 2,
}, {
	ID: 4,
	Original: "https://www.netflix.com/test",
	Shortened: "dcu.be/netflixTest",
	UserID: 2,
}}

func setup() (app *Application, db *gorm.DB) {
	dbName := "testDB"
	exec.Command("rm", "-f", dbName)

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&user.User{}, &urlshortener.ShortenedURL{})

	db.Create(users)
	db.Create(urls)

	app = &Application{}

	app.InitApp(db)
	return
}

func TestRegisterSuccess(t *testing.T) {
	app, db := setup()
	payload := []byte(`{"username":"test3", "password":"password"}`)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var newUser user.User
	err := db.First(&newUser, user.User{}).Error
	assert.Nil(t, err)
}

func TestRegisterFail(t *testing.T) {
	app, _ := setup()
	payload := []byte(`{"username":"test2", "password":"password"}`)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLoginSuccess(t *testing.T) {
	app, _ := setup()
	payload := []byte(`{"username":"test1", "password":"password1"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestLoginFail(t *testing.T) {
	app, _ := setup()
	payload := []byte(`{"username":"test2", "password":"wrongPassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetURLsSuccess(t *testing.T) {
	app, _ := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/url", nil)
	token, _ := session.GenerateToken(uint(1))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetURLsFail(t *testing.T) {
	app, _ := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 3)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/url", nil)
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	token, _ := session.GenerateToken(uint(3))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})
	
	resp = executeRequest(req, app)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCreateURLSuccess(t *testing.T) {
	app, db := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)
	payload := []byte(`{"original_url":"www.newurl.com"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/url", bytes.NewBuffer(payload))
	token, _ := session.GenerateToken(uint(1))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var url urlshortener.ShortenedURL
	err := db.Model(&urlshortener.ShortenedURL{}).First(&url, urlshortener.ShortenedURL{Original: "www.newurl.com"}).Error
	assert.Nil(t, err)
}

func TestCreateURLFail(t *testing.T) {
	app, _ := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)
	payload := []byte(`{"original_url":"www.newurl.com"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/url", bytes.NewBuffer(payload))

	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestDeleteURLSuccess(t *testing.T) {
	app, db := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, "/url/1", nil)
	token, _ := session.GenerateToken(uint(1))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})
	
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusOK, resp.Code)

	var url urlshortener.ShortenedURL
	err := db.Model(&urlshortener.ShortenedURL{}).First(&url, urlshortener.ShortenedURL{ID: 1}).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteURLFail(t *testing.T) {
	app, _ := setup()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, "/url/1", nil)

	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, "/url/200", nil)
	token, _ := session.GenerateToken(uint(1))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})

	resp = executeRequest(req, app)
	assert.Equal(t, http.StatusNotFound, resp.Code)

	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, "/url/3", nil)
	token, _ = session.GenerateToken(uint(1))
	req.AddCookie(&http.Cookie{Name: "token", Value: token.TokenString})

	resp = executeRequest(req, app)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func executeRequest(req *http.Request, app *Application) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, req)
	return recorder
}
