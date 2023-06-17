package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"time"

	// "regexp"
	"strings"
	// "time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
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
	Count           int    `json:"count"`
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
// var jsonResp []byte
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
	//  // Define your request handler
	//  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //     // Your request handling logic here
    // })

    // Wrap the handler with the CORS middleware
    // corsHandler := handleCORS(handler)

    // Register the handler with the router
    // router.Handle("/", corsHandler)
	http.ListenAndServe(":8000",
	
handlers.CORS(
	handlers.AllowedOrigins([]string{"*"}),
	handlers.AllowedMethods([]string{"GET","POST"}),
	handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"}),
)(router))
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
    // if(card.Expiry_month==0||card.Expiry_year==0){
	// 	http.Error(w,"Expiry month and year are required",http.StatusBadRequest)
	// 	return
	// }
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
		// Card data not found in the database
		log.Println("Error in card data founding:", err)
        // http.Error(w, "Card data not found", http.StatusNotFound)
		statusCode := http.StatusNotFound // Use the desired status code
	w.WriteHeader(statusCode)

	// Create the error response
	errorResponse := map[string]string{
		"error": "card data not found", // Replace with your desired error message
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
	// coalesceQuery:="SELECT COALESCE(OTP, 0) AS OTP FROM card_information"
	updateQuery := "UPDATE card_information SET OTP = ? WHERE ID = ?"

	
	// expiry := time.Now().Add(1 * time.Minute)
	// _, err = db.Exec(updateQuery, otp,card.ID)
	// _, err1:= db.Exec(coalesceQuery)
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
	// w.Write([]byte("OTP added successfully AND Payment processed successfully"))
	// json.NewEncoder(w).Encode("OTP added success")
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "OTP added successfully"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	newtimer := time.NewTimer(15 * time.Second)
  
    // Notifying the channel
    <-newtimer.C
	queryForUpdateOTP := "UPDATE card_information SET OTP = 0 WHERE Card_number = ?"
	_, err = db.Exec(queryForUpdateOTP, card.Card_number)
	if err != nil {
		log.Println("Error updating OTP:", err)
		return
	}
	
	log.Println("OTP deleted successfully")
// 	queryForUpdateOTP := "UPDATE card_information SET OTP = ? WHERE Card_number = ?"
// stmt, err := db.Prepare(queryForUpdateOTP)
// if err != nil {
//     log.Println("Error preparing statement:", err)
//     return
// }
// defer stmt.Close()

// _, err = stmt.Exec(nil, card.Card_number)
// if err != nil {
//     log.Println("Error updating OTP:", err)
//     return
// }

// log.Println("OTP deleted successfully")

	

	
	 

	
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
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	query := "SELECT OTP, count FROM card_information WHERE Card_number = ?"
	row := db.QueryRow(query, card.Card_number)
	var storedOTP sql.NullInt64
	var count int
	err = row.Scan(&storedOTP, &count)
	if err == sql.ErrNoRows {
		log.Println("No OTP found for the given card_number:", card.Card_number)
		http.Error(w, "No OTP found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error retrieving OTP from the database:", err)
		http.Error(w, "Failed to retrieve OTP from the database", http.StatusInternalServerError)
		return
	}

	if storedOTP.Valid && int(card.OTP) == int(storedOTP.Int64) {
		

		count=0
		// Update the count in the database
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)
			http.Error(w, "Failed to update OTP count", http.StatusInternalServerError)
			return
		}

		log.Println("OTP matched successfully. Count:", count)
	} else {
		if count >= 3 {
			log.Println("OTP matched maximum number of times")
			http.Error(w, "OTP matched maximum number of times", http.StatusForbidden)
			return
		}
		count++
		updateQuery := "UPDATE card_information SET count = ? WHERE Card_number = ?"
		_, err := db.Exec(updateQuery, count, card.Card_number)
		if err != nil {
			log.Println("Error updating OTP count:", err)
			http.Error(w, "Failed to update OTP count", http.StatusInternalServerError)
			return
		}
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
	updateQuery := "UPDATE card_information SET OTP = ? WHERE Card_number = ?"
	_,err = db.Exec(updateQuery, otp, card.Card_number)
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


