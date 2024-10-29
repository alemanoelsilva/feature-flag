package auth

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AuthUser struct {
	UserID    int      `json:"userId"`
	PersonID  int      `json:"personId"`
	UserEmail string   `json:"userEmail"`
	Roles     []string `json:"roles"`
	RolesID   []int    `json:"rolesIds"`
}

type AuthUserResponse struct {
	UserID    int
	PersonID  int
	UserEmail string
	Roles     []string
	RolesID   []int
	IsAdmin   bool
}

func getAuth(cookie string) (AuthUser, error) {
	// Start constructing the base URL
	baseURL := "http://localhost:9101/auth/session"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Cookie header
	req.Header.Set("Cookie", cookie)

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

	var authInfo AuthUser
	if err := json.Unmarshal(body, &authInfo); err != nil {
		log.Fatalf("Error Unmarshal response body: %v", err)
	}

	return authInfo, nil
}

func GetAuthInfo(cookie string) (AuthUserResponse, error) {
	// auth, err := getAuth(cookie)
	// if err != nil {
	// 	return AuthUserResponse{}, err
	// }

	// isAdmin := false
	// for _, role := range auth.Roles {
	// 	// TODO: just a info, only users with role FEATURE_FLAG is allowed to use this api. It ensures only the development team is able to make any change on it
	// 	if role == "FEATURE_FLAG" {
	// 		isAdmin = true
	// 	}
	// }

	// authInfo := AuthUserResponse{
	// 	UserID: auth.UserID,
	// 	// PersonID:  auth.PersonID,
	// 	// TODO: while testing, returns 1
	// 	PersonID:  1,
	// 	UserEmail: auth.UserEmail,
	// 	Roles:     auth.Roles,
	// 	RolesID:   auth.RolesID,
	// 	// IsAdmin:   isAdmin,
	// 	// TODO: while testing, returns true
	// 	IsAdmin: true,
	// }

	authInfo := AuthUserResponse{
		UserID:    1,
		PersonID:  1,
		UserEmail: "email",
		Roles:     []string{"ADMIN"},
		RolesID:   []int{1},
		IsAdmin:   true,
	}

	return authInfo, nil
}
