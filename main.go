package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	handler "go_digital_sign/handlers"
	"go_digital_sign/utils"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to MongoDB
	client, err := utils.GetMongoClient()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	// Set up the collection in the handler
	db := client.Database(os.Getenv("MONGO_DB"))
	signCollection := db.Collection(os.Getenv("MONGO_COLLECTION"))
	handler.SetSignCollection(signCollection)

	r := mux.NewRouter()
	r.HandleFunc("/api/signdocument", handler.SignDocument).Methods("POST")
	r.HandleFunc("/api/signatureverify", handler.VerifySignature).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started at", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
