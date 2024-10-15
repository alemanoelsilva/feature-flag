package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Person struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsAssigned bool   `json:"isAssigned"`
}

type PersonPaginationResponse struct {
	Items []Person `json:"items"`
	Total int      `json:"total"`
}

type Pagination struct {
	Page  int
	Limit int
	Next  bool
}

func (p *Person) GetPerson(featureFlagId int, pagination Pagination) ([]Person, Pagination) {
	// Start constructing the base URL
	var baseURL string
	if featureFlagId == 0 {
		baseURL = "http://localhost:9696/api/feature-flags/v1/people"
	} else {
		baseURL = "http://localhost:9696/api/feature-flags/v1/people/feature-flags/" + strconv.Itoa(featureFlagId)
	}

	if pagination.Page > 0 {
		baseURL += "?page=" + url.QueryEscape(strconv.Itoa(pagination.Page))
	}

	if pagination.Limit > 0 {
		baseURL += "&limit=" + url.QueryEscape(strconv.Itoa(pagination.Limit))
	}

	// Check if the name parameter is provided
	if p.Name != "" {
		// Encode the name parameter to handle special characters
		encodedName := url.QueryEscape(p.Name)

		// Append the query parameter to the base URL
		baseURL += "&name=" + encodedName
	}

	// Check if the isAssigned parameter is provided
	if p.IsAssigned {
		// Append the query parameter to the base URL
		baseURL += "&isAssigned=true"
	}

	fmt.Printf("URL %v", baseURL)
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

	var response PersonPaginationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error Unmarshal response body: %v", err)
	}

	var paginationResponse Pagination
	if pagination.Page == 0 {
		paginationResponse.Next = false
		return response.Items, paginationResponse
	}

	paginationResponse.Page = pagination.Page + 1
	paginationResponse.Limit = pagination.Limit
	paginationResponse.Next = response.Total > len(response.Items)

	return response.Items, paginationResponse
}
