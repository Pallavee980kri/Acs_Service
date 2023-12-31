package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	// "github.com/Pallavee980kri/Acs_Service/config"
	// "github.com/Pallavee980kri/Acs_Service/structType"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Card struct {
	ID              int    `json:"id"`
	Card_number     string `json:"card_number"`
	Cardholder_name string `json:"cardholder_name"`
	CVV             string `json:"cvv"`
	Expiry_month    int    `json:"expiry_month"`
	Expiry_year     int    `json:"expiry_year"`
	OTP             int    `json:"OTP"`
	Count           int    `json:"count"`
}

type carddetails struct {
	Card_number string `json:"card_number"`
}

var storedCard Card
var db *sql.DB
var card Card
var cancelTimer = make(chan struct{})

func connect() error {

	var err error
	db, err = sql.Open("mysql", "root:pall850@/acsservice")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	// Test the database connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	log.Println("Database connected successfully!")
	return nil
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
	router.HandleFunc("/match_otp", matchOTP).Methods("POST")
	router.HandleFunc("/resend_otp", resendOTP).Methods("POST")
	// handling cors error
	http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router))
}

// API for access card data from the frontend and match this with that database if it matches then generate OTP if not then
// give a proper message.
func processPaymentHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		ErrorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}
	// log.Printf("Received card data: %+v\n", card)
	if card.Cardholder_name == "" {
		ErrorMessagesResponse(w, r, "Card holder name is required")
		return
	}
	if card.Card_number == "" {
		ErrorMessagesResponse(w, r, "Card number is required.")
		return
	}
	if strings.Contains(card.Card_number, ".") {
		ErrorMessagesResponse(w, r, "Card number cannot contain '.' character")
		return
	}
	if len(card.Card_number) != 16 {
		ErrorMessagesResponse(w, r, "Card number must be 16 digits.")
		return
	}

	if strings.Contains(card.Card_number, "-") {
		ErrorMessagesResponse(w, r, "Card number cannot contain '-' character.")
		return
	}

	if strings.Contains(card.Card_number, "+") {
		ErrorMessagesResponse(w, r, "Card number cannot contain '+' character.")
		return
	}

	if strings.Contains(card.Card_number, "e") {
		ErrorMessagesResponse(w, r, "Card number cannot contain 'e' character.")
		return
	}

	if strings.Contains(card.Card_number, " ") {
		ErrorMessagesResponse(w, r, "Card number cannot contain whitespace.")
		return
	}
	if card.CVV == "" {
		ErrorMessagesResponse(w, r, "CVV is required")
		return

	}
	if len(card.CVV) != 3 {
		ErrorMessagesResponse(w, r, "Please enter valid 3 digits cvv number")
		return
	}
	if strings.Contains(card.CVV, ".") {
		ErrorMessagesResponse(w, r, "CVV cannot contain '.' character")
		return
	}
	if strings.Contains(card.CVV, "-") {
		ErrorMessagesResponse(w, r, "Card number cannot contain '-' character.")
		return
	}

	if strings.Contains(card.CVV, "+") {
		ErrorMessagesResponse(w, r, "Card number cannot contain '+' character.")
		return
	}

	if strings.Contains(card.CVV, "e") {
		ErrorMessagesResponse(w, r, "Card number cannot contain 'e' character.")
		return
	}

	if strings.Contains(card.CVV, " ") {
		ErrorMessagesResponse(w, r, "Card number cannot contain whitespace.")
		return
	}

	// Check if the card data exists in the database
	query := "SELECT * FROM card_information WHERE card_number = ? "
	row := db.QueryRow(query, card.Card_number)
	err = row.Scan(
		&storedCard.ID,
		&storedCard.Card_number,
		&storedCard.Cardholder_name,
		&storedCard.CVV,
		&storedCard.Expiry_month,
		&storedCard.Expiry_year,
		&storedCard.OTP,
		&storedCard.Count,
	)
	if storedCard.Card_number == "" {
		ErrorMessagesResponse(w, r, "Card Number Does Not Match!")
		return
	}

	if card.Cardholder_name != storedCard.Cardholder_name {
		ErrorMessagesResponse(w, r, "Card Holder Name Does Not Match!")
		return
	}
	if card.CVV != storedCard.CVV {
		ErrorMessagesResponse(w, r, "CVV Does Not Match!")
		return
	}
	if card.Expiry_month != storedCard.Expiry_month {
		ErrorMessagesResponse(w, r, "Expiry Month Does Not Match!")
		return
	}
	if card.Expiry_year != storedCard.Expiry_year {
		ErrorMessagesResponse(w, r, "Expiry year Does Not Match!")
		return
	}

	if err != nil {
		log.Println("Error querying the database:", err)
		ErrorMessagesResponse(w, r, "Failed to query the database")
		return
	}
	otp := generateOTP()
	updateQuery := "UPDATE card_information SET OTP = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, otp, storedCard.ID)

	if err != nil {
		log.Println("Error updating OTP in the database:", err)
		ErrorMessagesResponse(w, r, "Failed to update OTP in the database")
		return
	}
	log.Println("OTP:", otp)
	w.WriteHeader(http.StatusOK)
	successMessageResponse(w, r, "OTP added successfully!")
	//this code is for timer of 15 seconds of deleting the otp after some secnd from db
	go func() {
		select {
		case <-time.After(1 * time.Minute):
			go func() {
				queryForUpdateOTP := "UPDATE card_information SET OTP = 0 WHERE Card_number = ?"
				_, err := db.Exec(queryForUpdateOTP, card.Card_number)
				if err != nil {
					log.Println("Error updating OTP:", err)
					return
				}

				log.Println("OTP deleted successfully")

			}()
		case <-cancelTimer:
		}
	}()

}

// submit OTP
func matchOTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		ErrorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}
	query := "SELECT OTP, count FROM card_information WHERE Card_number = ?"
	row := db.QueryRow(query, card.Card_number)
	var storedOTP int
	var count int
	err = row.Scan(&storedOTP, &count)
	if err == sql.ErrNoRows {
		log.Println("No OTP found for the given card_number:", card.Card_number)
		ErrorMessagesResponse(w, r, "No OTP found")
		return
	} else if err != nil {
		log.Println("Error retrieving OTP from the database:", err)
		ErrorMessagesResponse(w, r, "Failed to retrieve OTP from the database")
		return
	}
	if storedOTP == 0 && card.OTP == storedOTP {
		count = 0
		// Update the count in the database
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error Updating In OTP Count:", err)
			ErrorMessagesResponse(w, r, "Failed to update OTP count")
			return
		}

		log.Println("OTP Matched Successfully.")
	} else if storedOTP != 0 && card.OTP == storedOTP {
		count = 0
		// Update the count in the database
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)
			ErrorMessagesResponse(w, r, "Failed to update OTP count")
			return
		}
		log.Println("OTP matched successfully.")
	} else {
		if count >= 3 {
			log.Println("OTP matched maximum number of times")

			ErrorMessagesResponse(w, r, "You have reached maximum attemps to submit OTP Please try again !")
			return
		}
		count++
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)

			ErrorMessagesResponse(w, r, "Failed to update OTP count")
			return
		}
		log.Println("Invalid OTP provided")
		ErrorMessagesResponse(w, r, "Invalid OTP")
		return
	}

	successMessageResponse(w, r, "OTP matched successfully")
}

// API for resend the OTP
func resendOTP(w http.ResponseWriter, r *http.Request) {
	cancelTimer <- struct{}{}
	log.Println("hello otp checked")
	var carddata carddetails
	err := json.NewDecoder(r.Body).Decode(&carddata)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		ErrorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}
	otp := generateOTP()
	// Update the OTP in the database
	updateQuery := "UPDATE card_information SET OTP = ? WHERE card_number = ?"
	_, err = db.Exec(updateQuery, otp, carddata.Card_number)
	log.Println("hello otp")
	if err != nil {
		log.Println("Error resending the OTP in the database:", err)
		ErrorMessagesResponse(w, r, "Failed to resend OTP in the database")
		return
	}
	log.Println("OTP:", otp)
	w.WriteHeader(http.StatusOK)
	successMessageResponse(w, r, "OTP resent successfully")
	go func() {
		select {
		case <-time.After(1 * time.Minute):
			go func() {
				queryForUpdateOTP := "UPDATE card_information SET OTP = 0 WHERE Card_number = ?"
				_, err := db.Exec(queryForUpdateOTP, card.Card_number)
				if err != nil {
					log.Println("Error updating OTP:", err)
					return
				}

				log.Println("OTP deleted successfully")

			}()
		case <-cancelTimer:
		}
	}()

}

// Function to generate a random OTP
// from[100000-900000]to (0-n)
func generateOTP() int {
	otp := rand.Intn(900000) + 100000
	return otp
}

// sending error msg in json format
func ErrorMessagesResponse(w http.ResponseWriter, r *http.Request, msg string) {
	statusCode := http.StatusNotFound
	w.WriteHeader(statusCode)
	// Creating the error response message
	errorResponse := map[string]string{
		"error": msg,
	}
	// Marshal the error response into JSON
	responseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println("Failed to marshal error response:", err)
		return
	}
	// Set the response content type
	w.Header().Set("Content-Type", "application/json")
	// Send the JSON response
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Println("Failed to send response:", err)
	}

}

// sending success message in json format
func successMessageResponse(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
