package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Imranr2/DCUBE_API/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading env file")
		}
	}

	db := database.InitDB()
	fmt.Println(db)
}
