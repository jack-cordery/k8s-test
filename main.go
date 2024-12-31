package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Response struct {
	Message string `json:"message"`
	Status string `json:"status"`
}

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	
	
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)


	log.Println("Some things ", dbHost, dbName, dbPassword, dbUser, dbURI)


	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", dbURI)
		dbStatus := "connected"
		if err != nil {
			dbStatus = "error"
			log.Println("error", err)
			err = json.NewEncoder(w).Encode(Response{
				Message: "not connected",
				Status: dbStatus,
			})
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		defer db.Close()
		log.Println("connected")

		err = db.Ping()

		if err != nil {
			dbStatus = "error"
			log.Println("Another error", err)
			err = json.NewEncoder(w).Encode(Response{
				Message: "not connected",
				Status: dbStatus,
			})
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		response := Response{
			Message: "API is running",
			Status: dbStatus,
		}

		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			log.Fatal(err)
		}
	})

	log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}