package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TODO: how to receive it from the web
const COOKIE = "sess=eyJwYXNzcG9ydCI6eyJ1c2VyIjoie1widXNlcklkXCI6MjkyLFwicGVyc29uSWRcIjozNTIxNyxcInVzZXJFbWFpbFwiOlwiYWxleGFuZHJlLnNpbHZhQG1pbmRlcmEuY29tXCJ9In19; sess.sig=39qe2IvY5JL9Q7gKt3gasaRSGtI"

type Person struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsAssigned bool   `json:"isAssigned"`
}

type FeatureFlag struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	IsGlobal       bool   `json:"isGlobal"`
	ExpirationDate string `json:"expirationDate"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	Person         Person `json:"person"`
}

type PaginationResponse struct {
	Items []FeatureFlag `json:"items"`
	Total int           `json:"total"`
}

type FeatureFlagRequest struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	IsGlobal       bool   `json:"isGlobal"`
	ExpirationDate string `json:"expirationDate"`
}

type APIError struct {
	Error string `json:"error"`
}

type ErrorResponse struct {
	IsError bool
	Field   string
	Message string
}

func (ff *FeatureFlag) GetFeatureFlag() []FeatureFlag {
	// Start constructing the base URL
	baseURL := "http://localhost:9696/api/feature-flags/v1/feature-flags?limit=100"

	// Check if the name parameter is provided
	if ff.ID != 0 {
		// Encode the name parameter to handle special characters
		encodedID := url.QueryEscape(strconv.Itoa(ff.ID))

		// Append the query parameter to the base URL
		baseURL += "&id=" + encodedID
	}

	// Check if the name parameter is provided
	if ff.Name != "" {
		// Encode the name parameter to handle special characters
		encodedName := url.QueryEscape(ff.Name)

		// Append the query parameter to the base URL
		baseURL += "&name=" + encodedName
	}

	// Check if the isActive parameter is provided
	if ff.IsActive {
		// Append the query parameter to the base URL
		baseURL += "&isActive=true"
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatalf("Error request: %v", err)
	}

	// Set the Cookie header
	req.Header.Set("Cookie", COOKIE)

	// Use the default HTTP client to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	// Ensure the response body is closed after we are done reading it
	defer resp.Body.Close()

	// Check for successful response (Status Code 200)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var response PaginationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error Unmarshal response body: %v", err)
	}

	return response.Items
}

func (ff *FeatureFlag) CreateFeatureFlag() ErrorResponse {
	// Convert the request body to JSON
	featureFlag := FeatureFlagRequest{
		Name:           ff.Name,
		Description:    ff.Description,
		IsActive:       ff.IsActive,
		ExpirationDate: ff.ExpirationDate,
	}
	jsonData, err := json.Marshal(featureFlag)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	// Create a new request using http.NewRequest
	req, err := http.NewRequest("POST", "http://localhost:9696/api/feature-flags/v1/feature-flags", bytes.NewBuffer(jsonData))
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
		log.Println("Feature flag created successfully.")
	} else {
		log.Printf("Failed to create feature flag: %s\n", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var apiError APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		log.Fatalf("Error Unmarshal response body error: %v", err)
	}

	var response ErrorResponse
	if resp.StatusCode == http.StatusBadRequest {
		response.IsError = true
		response.Field = strings.Split(apiError.Error, "|")[0]
		response.Message = strings.Split(apiError.Error, "|")[1]
	}

	if resp.StatusCode == http.StatusConflict {
		response.IsError = true
		response.Field = "Page"
		response.Message = apiError.Error
	}

	return response
}

func (ff *FeatureFlag) UpdateFeatureFlag() ErrorResponse {
	// Convert the request body to JSON
	featureFlag := FeatureFlagRequest{
		Description:    ff.Description,
		IsActive:       ff.IsActive,
		IsGlobal:       ff.IsGlobal,
		ExpirationDate: ff.ExpirationDate,
	}
	jsonData, err := json.Marshal(featureFlag)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	url := fmt.Sprintf("http://localhost:9696/api/feature-flags/v1/feature-flags/%d", ff.ID)

	// Create a new request using http.NewRequest
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
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
	if resp.StatusCode == http.StatusOK {
		log.Println("Feature flag updated successfully.")
	} else {
		log.Printf("Failed to update feature flag: %s\n", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var apiError APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		log.Fatalf("Error Unmarshal response body error: %v", err)
	}

	var response ErrorResponse
	if resp.StatusCode == http.StatusBadRequest {
		response.IsError = true
		response.Field = strings.Split(apiError.Error, "|")[0]
		response.Message = strings.Split(apiError.Error, "|")[1]
	}

	if resp.StatusCode == http.StatusConflict {
		response.IsError = true
		response.Field = "Page"
		response.Message = apiError.Error
	}

	return response
}
