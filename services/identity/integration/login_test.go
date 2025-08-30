package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type LoginResponse struct {
	Data   LoginData `json:"data"`
	Errors []any     `json:"errors"`
}

type LoginData struct {
	Login LoginPayload `json:"login"`
}

type LoginPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	User         User   `json:"user"`
}

type User struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func TestLoginMutationShape(t *testing.T) {
	query := `mutation($email:String!,$password:String!){ login(input:{email:$email,password:$password}){ accessToken expiresIn user { email } } }`
	vars := map[string]any{"email": "admin@example.com", "password": "Admin123!"}
	body, _ := json.Marshal(map[string]any{"query": query, "variables": vars})
	resp, err := http.Post("http://localhost:8001/graphql", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func TestLoginMutationComplete(t *testing.T) {
	query := `mutation($email:String!,$password:String!){
		login(input:{email:$email,password:$password}){
			accessToken
			refreshToken
			expiresIn
			user {
				email
				username
				firstName
				lastName
			}
		}
	}`
	vars := map[string]any{"email": "admin@example.com", "password": "Admin123!"}
	body, _ := json.Marshal(map[string]any{"query": query, "variables": vars})

	resp, err := http.Post("http://localhost:8001/graphql", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	if result.Data.Login.AccessToken == "" {
		t.Error("expected accessToken to be non-empty")
	}

	if result.Data.Login.ExpiresIn <= 0 {
		t.Error("expected expiresIn to be positive")
	}

	if result.Data.Login.User.Email != "admin@example.com" {
		t.Errorf("expected email admin@example.com, got %s", result.Data.Login.User.Email)
	}

	if result.Data.Login.User.Username != "admin" {
		t.Errorf("expected username admin, got %s", result.Data.Login.User.Username)
	}
}
