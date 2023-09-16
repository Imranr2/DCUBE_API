package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Imranr2/DCUBE_API/internal/urlshortener"
	"github.com/Imranr2/DCUBE_API/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDatabaseURL() (databaseURL string) {
	if os.Getenv("ENV") == "PROD" {
		databaseURL = os.Getenv("DATABASE_URL")
	} else {
		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("DATABASE_USERNAME"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_NET"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
		)
	}

	return
}

func InitDB() (db *gorm.DB) {
	db, err := gorm.Open(postgres.Open(getDatabaseURL()), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	err = db.AutoMigrate(
		&user.User{},
		&urlshortener.ShortenedURL{},
	)

	if err != nil {
		log.Fatal("Unable to perform migration")
	}

	return
}
