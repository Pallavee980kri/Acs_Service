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

var storedCard Card
var db *sql.DB
var card Card

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

		errorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}
	log.Printf("Received card data: %+v\n", card)
	if card.Cardholder_name == "" {

		errorMessagesResponse(w, r, "Card holder name is required")
		return
	}
	//validation for card number
	if card.Card_number == "" {

		errorMessagesResponse(w, r, "Card number is required.")
		return
	}

	if len(card.Card_number) != 16 {

		errorMessagesResponse(w, r, "Card number must be 16 digits.")
		return
	}

	if strings.Contains(card.Card_number, "-") {
		errorMessagesResponse(w, r, "Card number cannot contain '-' character.")
		return
	}

	if strings.Contains(card.Card_number, "+") {

		errorMessagesResponse(w, r, "Card number cannot contain '+' character.")
		return
	}

	if strings.Contains(card.Card_number, "e") {

		errorMessagesResponse(w, r, "Card number cannot contain 'e' character.")
		return
	}

	if strings.Contains(card.Card_number, " ") {

		errorMessagesResponse(w, r, "Card number cannot contain whitespace.")
		return
	}

	// if len(errorMessages) > 0 {
	// 	errorString := strings.Join(errorMessages, " ")
	// 	http.Error(w, errorString, http.StatusBadRequest)
	// 	return
	// }

	if card.CVV == "" || len(card.CVV) != 3 {

		errorMessagesResponse(w, r, "CVV is required Please enter valid 3 digits cvv number")
		return

	}
	
	// Check if the card data exists in the database
	query := "SELECT * FROM card_information WHERE card_number = ? AND cardholder_name = ?"
	row := db.QueryRow(query, card.Card_number, card.Cardholder_name)
	err = row.Scan(
		// &id,
		&storedCard.ID,
		&storedCard.Card_number,
		&storedCard.Cardholder_name,
		&storedCard.CVV,
		&storedCard.Expiry_month,
		&storedCard.Expiry_year,
		// &otpValue,
		&storedCard.OTP,
		// &storedCard.Expiry,
		&storedCard.Count,
	)
	if err == sql.ErrNoRows {

		log.Println("Error in card data founding:", err)
		errorMessagesResponse(w, r, "Card Data Not Found")
		return
	} else if err != nil {
		log.Println("Error querying the database:", err)
		errorMessagesResponse(w, r, "Failed to query the database")
		return
	}
	// Compare the stored card data with the frontend data
	if card.CVV != storedCard.CVV || card.Expiry_month != storedCard.Expiry_month || card.Expiry_year != storedCard.Expiry_year {
		errorMessagesResponse(w, r, "Card data does not match")
		return
	}
	// Generate OTP
	otp := generateOTP()
	updateQuery := "UPDATE card_information SET OTP = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, otp, storedCard.ID)

	if err != nil {
		log.Println("Error updating OTP in the database:", err)
		errorMessagesResponse(w, r, "Failed to update OTP in the database")
		return
	}
	log.Println("OTP:", otp)
	w.WriteHeader(http.StatusOK)
	//sending the response into json format
	successMessageResponse(w,r,"OTP added successfully!")
   //this code is for timer of 15 seconds of deleting the otp after some secnd from db
	go func() {
		newTimer := time.NewTimer(120 * time.Second)
	
		// Wait for the timer to expire
		<-newTimer.C
	
		// Perform the update in the background
		go func() {
			queryForUpdateOTP := "UPDATE card_information SET OTP = 0 WHERE Card_number = ?"
			_, err := db.Exec(queryForUpdateOTP, card.Card_number)
			if err != nil {
				log.Println("Error updating OTP:", err)
				return
			}
			log.Println("OTP deleted successfully")
		}()
	}()
	
	
}

// Function to generate a random OTP
func generateOTP() int {
	otp := rand.Intn(900000) + 100000
	return otp

}
func matchOTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		errorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}

	query := "SELECT OTP, count FROM card_information WHERE Card_number = ?"
	row := db.QueryRow(query, card.Card_number)
	var storedOTP sql.NullInt64
	var count int
	err = row.Scan(&storedOTP, &count)
	if err == sql.ErrNoRows {
		log.Println("No OTP found for the given card_number:", card.Card_number)
		errorMessagesResponse(w, r, "No OTP found")
		return
	} else if err != nil {
		log.Println("Error retrieving OTP from the database:", err)
		errorMessagesResponse(w, r, "Failed to retrieve OTP from the database")
		return
	}

	if storedOTP.Valid && int(card.OTP) == int(storedOTP.Int64) {

		count = 0
		// Update the count in the database
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)
			errorMessagesResponse(w, r, "Failed to update OTP count")
			return
		}

		log.Println("OTP matched successfully. Count:", count)
	} else {
		if count >= 3 {
			log.Println("OTP matched maximum number of times")

			errorMessagesResponse(w, r, "You have reached maximum attemps to submit OTP!")
			return
		}
		count++
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)

			errorMessagesResponse(w, r, "Failed to update OTP count")
			return
		}
		log.Println("Invalid OTP provided")
		errorMessagesResponse(w, r, "Invalid OTP")
		return
	}

	successMessageResponse(w,r,"OTP matched successfully")
}

// API for resend the OTP
func resendOTP(w http.ResponseWriter, r *http.Request) {

	var card Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)

		errorMessagesResponse(w, r, "Failed to parse JSON payload")
		return
	}

	// Generate a new OTP
	otp := generateOTP()

	// Update the OTP in the database
	updateQuery := "UPDATE card_information SET OTP = ? WHERE Card_number = ?"
	_, err = db.Exec(updateQuery, otp, card.Card_number)
	if err != nil {
		log.Println("Error resending the OTP in the database:", err)
		errorMessagesResponse(w, r, "Failed to resend OTP in the database")
		return
	}
	log.Println("OTP:", otp)
	w.WriteHeader(http.StatusOK)
	successMessageResponse(w,r,"OTP resent successfully")
}

func errorMessagesResponse(w http.ResponseWriter, r *http.Request, msg string) {
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

func successMessageResponse(w http.ResponseWriter,r *http.Request,msg string){
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}