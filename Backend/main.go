package main

import (
	"database/sql"
	"encoding/json"

	// "time"
	"math/rand"

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
	ID              int    `json:"id"`
	Card_number     string `json:"card_number"`
	Cardholder_name string `json:"cardholder_name"`
	CVV            string `json:"cvv"`
	Expiry_month    int    `json:"expiry_month"`
	Expiry_year     int    `json:"expiry_year"`
	OTP   int    `json:"OTP"`
}

var db *sql.DB
func connect()  error {
	var err error
	db, err = sql.Open("mysql", "root:pall850@/acsservice")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
   // Test the database connection
	err = db.Ping()
	if err != nil {
		return  fmt.Errorf("failed to ping database: %v", err)
	}
	log.Println("Database connected successfully!")

    return  nil
}
func main() {
	// Connect to the database
	err := connect()
	if err != nil {
		log.Fatal(err)
	}
    // Close the database connection before the main function exits
	defer db.Close()
    // Initialize the router
	router := mux.NewRouter()

	router.HandleFunc("/process_payment", processPaymentHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
//access card data from the frontend
func processPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var card Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received card data: %+v\n", card)

	if card.Card_number == "" {
		http.Error(w, "Card holder name is required", http.StatusBadRequest)
		return
	}
	if card.Card_number == "" || len(card.Card_number) != 16 {
        http.Error(w,"Card number is required and must be 16 digits Please enter valid card number",http.StatusBadRequest)
		return
	}

	if(card.CVV=="" || len(card.CVV)!=3){
		http.Error(w,"CVV is required Please enter valid 3 digits cvv number",http.StatusBadRequest)
		return
	}

	if(card.Expiry_month==0||card.Expiry_year==0){
		http.Error(w,"Expiry month and year are required",http.StatusBadRequest)
		return
	}

	// Check if the card data exists in the database
	query := "SELECT * FROM card_information WHERE card_number = ? AND cardholder_name = ?"
	row := db.QueryRow(query, card.Card_number, card.Cardholder_name)
    var storedCard Card
	var otpValue sql.NullInt64
    err = row.Scan(
    // &id,
	&card.ID,
    &storedCard.Card_number,
    &storedCard.Cardholder_name,
    &storedCard.CVV,
    &storedCard.Expiry_month,
    &storedCard.Expiry_year,
    &otpValue,
)

	
	if err == sql.ErrNoRows {
		// Card data not found in the database
		log.Println("Error in card data founding:", err)

		http.Error(w, "Card data not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error querying the database:", err)
		http.Error(w, "Failed to query the database", http.StatusInternalServerError)
		return
	}
	log.Println("Card ID:", card.ID)

	// log.Println("storedcard",storedCard)
	// log.Println("Stored Card Number:", storedCard.Card_number)
	// log.Println("Stored Cardholder Name:", storedCard.Cardholder_name)
	// log.Println("Frontend Card Number:", card.Card_number)
	// log.Println("Frontend Cardholder Name:", card.Cardholder_name)
	// log.Println("Frontend CVV:", card.CVV)
	// log.Println("Stored CVV:", storedCard.CVV)
	// log.Println("Frontend Expiry Month:", card.Expiry_month)
	// log.Println("Stored Expiry Month:", storedCard.Expiry_month)
	// log.Println("Frontend Expiry Year:", card.Expiry_year)
	// log.Println("Stored Expiry Year:", storedCard.Expiry_year)

	// Compare the stored card data with the frontend data
	if card.CVV != storedCard.CVV || card.Expiry_month != storedCard.Expiry_month || card.Expiry_year != storedCard.Expiry_year {
		http.Error(w, "Card data does not match", http.StatusBadRequest)
		return
	}

	// Generate OTP
otp := generateOTP()

// Update the OTP in the database
updateQuery := "UPDATE card_information SET OTP = ? WHERE id = ?"
_, err = db.Exec(updateQuery, otp, storedCard.ID)
if err != nil {
    log.Println("Error updating OTP in the database:", err)
    http.Error(w, "Failed to update OTP in the database", http.StatusInternalServerError)
    return
}

log.Println("OTP",otp)
	// Send a response back to the frontend
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP added successfully"))
	w.Write([]byte("Payment processed successfully"))
}

// Function to generate a random OTP

func generateOTP() int {

	otp := rand.Intn(900000) + 100000
	return otp

}