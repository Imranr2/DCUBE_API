package application

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/Imranr2/DCUBE_API/internal/urlshortener"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var users = []user.User{{
	ID: 1,
	Username: "test1",
	Password: "password1",
}, {
	ID: 2,
	Username: "test2",
	Password: "password2",
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

func TestRegister(t *testing.T) {
	app, db := setup()
	payload := []byte(`{"username":"test3", "password":"password"}`)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var newUser user.User
	err := db.First(&newUser, user.User{}).Error
	assert.Nil(t, err)
}

func executeRequest(req *http.Request, app *Application) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, req)
	return recorder
}

func checkResponseCode(t *testing.T, expected, actual uint) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

