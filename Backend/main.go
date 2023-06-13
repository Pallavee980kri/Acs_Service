package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)
type Card struct {
	CardNumber     int    `json:"card_number"`
	CardholderName string `json:"cardholder_name"`
	CVV            int    `json:"cvv"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	RandomNumber   int    `json:"random_number"`
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
    // Define the API endpoint to save the card details
	router.HandleFunc("/random-number", func(w http.ResponseWriter, r *http.Request) {
		createRandomNumber(w, r, db)
	}).Methods("POST")
    router.HandleFunc("/random-check",randomCheck).Methods("GET")
    // Start the server
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
func createRandomNumber(w http.ResponseWriter,r *http.Request, db *sql.DB) {
	var card Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
    // Create a new random number generator with a specific seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    // Generate a random number using the new generator
    randomNum := rng.Intn(900000) + 100000
    fmt.Println(randomNum)
    query := "INSERT INTO card_information (card_number, cardholder_name, cvv, expiry_month, expiry_year, random_number) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, card.CardNumber, card.CardholderName, card.CVV, card.ExpiryMonth, card.ExpiryYear, randomNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    // Respond with the card details and random number
	card.RandomNumber = randomNum
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func randomCheck(w http.ResponseWriter,r*http.Request){
fmt.Println("hello")
}