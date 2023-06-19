package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
