package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

func TestCreateUser(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a user, then get it back
	userReq := `{"username": "lcrownover", "email": "lcrownover@localhost", "firstname": "Lucas", "lastname": "Crownover"}`
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/users", bytes.NewBuffer([]byte(userReq)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "testkey1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusCreated)
	}
	var userResponse UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		t.Fatal(err)
	}

	u, err := data.GetUserById(th.DB, userResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "lcrownover" {
		t.Errorf("expected username to be lcrownover, got %v", u.Username)
	}
	if u.Email != "lcrownover@localhost" {
		t.Errorf("expected email to be lcrownover@localhost, got %v", u.Email)
	}
	if u.FirstName != "Lucas" {
		t.Errorf("expected firstname to be Lucas, got %v", u.FirstName)
	}
	if u.LastName != "Crownover" {
		t.Errorf("expected lastname to be Crownover, got %v", u.LastName)
	}
}

// TestGetAllUsers tests the GET /api/v1/users endpoint
// it creates a user, then gets all users and checks that the created user is in the list
func TestGetAllUsers(t *testing.T) {
	// first we need to create a user, then get it back
	userReq := `{"username": "testgetall", "email": "testgetall@localhost", "firstname": "test", "lastname": "getall"}`
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/users", bytes.NewBuffer([]byte(userReq)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "testkey1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusCreated)
	}

	// now get all users
	req, err = http.NewRequest("GET", "http://localhost:3333/api/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "testkey1")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusCreated)
	}
	var usersResponse []UserResponse
	err = json.NewDecoder(resp.Body).Decode(&usersResponse)
	if err != nil {
		t.Fatal(err)
	}

	// check if the usersResponse has more than 1 user
	if len(usersResponse) < 1 {
		t.Errorf("expected more than 1 user, got %v", len(usersResponse))
	}

	// check if the user we created is in the list
	found := false
	for _, user := range usersResponse {
		if user.Username == "testgetall" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find user testgetall in the list of users")
	}
}
