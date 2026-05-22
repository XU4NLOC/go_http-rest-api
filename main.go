package main

import (
	"log"
	"net/http"

	"go_http-rest-api/db"
	"go_http-rest-api/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect()

	r := router.Setup()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
