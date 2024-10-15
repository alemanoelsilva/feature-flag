package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type AssignmentRequest struct {
	FeatureFlagID int `json:"featureFlagId"`
	PersonID      int `json:"personId"`
}

func (ar *AssignmentRequest) ApplyAssignment() {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(ar)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	// Create a new request using http.NewRequest
	req, err := http.NewRequest("POST", "http://localhost:9696/api/feature-flags/v1/assignments", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Cookie header
	req.Header.Set("Cookie", COOKIE)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusCreated {
		log.Println("Assignment applied successfully.")
	} else {
		log.Printf("Failed to apply assignment: %s\n", resp.Status)
	}
}

func (ar *AssignmentRequest) DeleteAssignment() {
	// Convert the request body to JSON
	jsonData, err := json.Marshal(ar)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	// Create a new request using http.NewRequest
	req, err := http.NewRequest("DELETE", "http://localhost:9696/api/feature-flags/v1/assignments", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Cookie header
	req.Header.Set("Cookie", COOKIE)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusCreated {
		log.Println("Assignment removed successfully.")
	} else {
		log.Printf("Failed to delete assignment: %s\n", resp.Status)
	}
}
