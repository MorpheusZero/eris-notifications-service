package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	SetupServer()
}

type NotificationBody struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type EmailMessage struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func SetupServer() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/notify", notifyHandler)
	err := http.ListenAndServe(":9110", nil)
	if err != nil {
		panic(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "ok")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body NotificationBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	fmt.Println(body.Type)

	if body.Type != "" {
		if strings.ToLower(body.Type) == "email" {
			var email EmailMessage
			err := json.NewDecoder(strings.NewReader(body.Data)).Decode(&email)
			if err != nil {
				http.Error(w, "Failed to parse email content", http.StatusBadRequest)
				return
			}
			go emailHandler(email)
		} else {
			http.Error(w, "Invalid Type!", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid Request: Type is a required!", http.StatusBadRequest)
		return
	}

	// You can perform any processing or response sending here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func emailHandler(email EmailMessage) {
	fmt.Printf("Received Email Notification: %+v\n", email)
	// Connect to SES and send the email.
}
