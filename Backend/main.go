package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	// "time"

	// "regexp"
	"strings"
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
	OTP             int    `json:"OTP"`
	// Expiry          sql.NullTime `json:"expiry"`
}
// var otpValue sql.NullInt64


// type NullInt64 struct {
// 	sql.NullInt64
// }

// func (ni *NullInt64) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" {
// 		ni.Valid = false
// 		return nil
// 	}
// 	err := json.Unmarshal(data, &ni.Int64)
// 	if err == nil {
// 		ni.Valid = true
// 	}
// 	return err
// }

var storedCard Card
var db *sql.DB
var card Card
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
	router.HandleFunc("/match_otp",matchOTP).Methods("POST")
	router.HandleFunc("/resend_otp",resendOTP).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
//API for access card data from the frontend and match this with that database if it matches then generate OTP if not then 
//give a proper message.
func processPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var errorMessages []string//creating a slice
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
	//validation for card number 
if card.Card_number == "" {
	errorMessages = append(errorMessages, "Card number is required.")
}

if len(card.Card_number) != 16 {
	errorMessages = append(errorMessages, "Card number must be 16 digits.")
}

if strings.Contains(card.Card_number, "-") {
	errorMessages = append(errorMessages, "Card number cannot contain '-' character.")
}

if strings.Contains(card.Card_number, "+") {
	errorMessages = append(errorMessages, "Card number cannot contain '+' character.")
}

if strings.Contains(card.Card_number, "e") {
	errorMessages = append(errorMessages, "Card number cannot contain 'e' character.")
}

if strings.Contains(card.Card_number, " ") {
	errorMessages = append(errorMessages, "Card number cannot contain whitespace.")
}

if len(errorMessages) > 0 {
	errorString := strings.Join(errorMessages, " ")
	http.Error(w, errorString, http.StatusBadRequest)
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
    // Compare the stored card data with the frontend data
	if card.CVV != storedCard.CVV || card.Expiry_month != storedCard.Expiry_month || card.Expiry_year != storedCard.Expiry_year {
		http.Error(w, "Card data does not match", http.StatusBadRequest)
		return
	}
    // Generate OTP
    otp := generateOTP()
    // Update the OTP in the database
	// updateQuery := "UPDATE card_information SET OTP = ?, WHERE id = ?"
	updateQuery := "UPDATE card_information SET OTP = ? WHERE ID = ?"

	
	// expiry := time.Now().Add(1 * time.Minute)
	// _, err = db.Exec(updateQuery, otp,card.ID)
	_, err = db.Exec(updateQuery, otp, storedCard.ID)

    if err != nil {
    log.Println("Error updating OTP in the database:", err)
    http.Error(w, "Failed to update OTP in the database", http.StatusInternalServerError)
    return
}
    log.Println("OTP:",otp)
	// Schedule the deletion of OTP after the expiry time
	// time.Sleep(1 * time.Minute)
	// deleteOTPFromDatabase(db, card.ID)
	// fmt.Println("OTP deleted from the database")
	// Send a response back to the frontend
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP added successfully AND Payment processed successfully"))
}

// Function to generate a random OTP
func generateOTP() int {
    otp := rand.Intn(900000) + 100000
	return otp

}
// func deleteOTPFromDatabase(db *sql.DB, cardID int) {
// 	deleteQuery := "UPDATE card_information SET OTP = NULL, WHERE id = ?"
// 	_, err := db.Exec(deleteQuery, cardID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("OTP deleted from the database")
// }




//validation for card number of -+e and space character with regexp
// func containsInvalidCharsInCardNumber(cardNumber string) bool {
	// regex := regexp.MustCompile(`[-+e\s]`)
	// return regex.MatchString(cardNumber)
// }




// Match the OTP received from the frontend with the logger OTP
func matchOTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}
	query := "SELECT OTP FROM card_information WHERE id = ?"
	row := db.QueryRow(query, card.ID)
	var storedOTP sql.NullInt64
	err = row.Scan(&storedOTP)
	if err == sql.ErrNoRows {
		log.Println("No OTP found for the given ID:", card.ID)
		http.Error(w, "No OTP found", http.StatusNotFound)
		return
	}else if err != nil {
		log.Println("Error retrieving OTP from the database:", err)
		http.Error(w, "Failed to retrieve OTP from the database", http.StatusInternalServerError)
		return
	}

	if storedOTP.Valid && int(card.OTP)== int(storedOTP.Int64) {
		log.Println("OTP matched successfully")
	} else {
		log.Println("Invalid OTP provided")
		http.Error(w, "Invalid OTP", http.StatusBadRequest)
		return
	}

	w.Write([]byte("OTP matched successfully"))
}

//API for resend the OTP
func resendOTP(w http.ResponseWriter, r *http.Request) {
	
	var card Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Println("Error parsing JSON payload:", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	// Generate a new OTP
	otp := generateOTP()

	// Update the OTP in the database
	updateQuery := "UPDATE card_information SET OTP = ? WHERE ID = ?"
	_,err = db.Exec(updateQuery, otp, card.ID)
	if err != nil {
		log.Println("Error resending the OTP in the database:", err)
		http.Error(w, "Failed to resend OTP in the database", http.StatusInternalServerError)
		return
	}

	log.Println("card ID:", card.ID)
	log.Println("OTP:", otp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP resent successfully"))
}
