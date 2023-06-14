package main

import (
	"database/sql"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"log"

	// "math/rand"
	"net/http"
	// "time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)
type Card struct {
	CardNumber     string `json:"cardNumber"`
	CardholderName string `json:"cardHolderName"`
	CVV            string `json:"cvv"`
	ExpiryMonth    int    `json:"expiryMonth"`
	ExpiryYear     int    `json:"expiryYear"`
	RandomNumber   int    `json:"randomNumber"`
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:pall850@/acsservice")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
   // Test the database connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
    fmt.Println("Database connected successfully!")
    return db, nil
}
func main() {
	// Connect to the database
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
    // Close the database connection before the main function exits
	defer db.Close()
    // Initialize the router
	router := mux.NewRouter()

	router.HandleFunc("/process_payment", processPaymentHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func processPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var card Card
	fmt.Println(card)
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received card data: %+v\n", card)
	if card.CardholderName == "" {
		http.Error(w, "Card holder name is required", http.StatusBadRequest)
		return
	}
	if(card.CardNumber==""){
		http.Error(w,"Card number is required",http.StatusBadRequest)
	}

	if(card.CVV==""){
		http.Error(w,"CVV is required",http.StatusBadRequest)
	}

	if(card.ExpiryMonth==0||card.ExpiryYear==0){
		http.Error(w,"Expiry month and year are required",http.StatusBadRequest)
	}
	// Send a response back to the frontend
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment processed successfully"))
}