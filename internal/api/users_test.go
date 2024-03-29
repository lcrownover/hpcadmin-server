package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

func TestAPICreateUser(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a user, then get it back
	ur := UserRequest{
		Username:  "testapicreateuser",
		Email:     "testapicreateuser@localhost",
		FirstName: "TestAPI",
		LastName:  "CreateUser",
	}
	userReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
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
	if u.Username != ur.Username {
		t.Errorf("expected username %v got %v", ur.Username, u.Username)
	}
	if u.Email != ur.Email {
		t.Errorf("expected email %v got %v", ur.Email, u.Email)
	}
	if u.FirstName != ur.FirstName {
		t.Errorf("expected firstname %v got %v", ur.FirstName, u.FirstName)
	}
	if u.LastName != ur.LastName {
		t.Errorf("expected lastname %v got %v", ur.LastName, u.LastName)
	}
}

// TestGetAllUsers tests the GET /api/v1/users endpoint
// it creates a user, then gets all users and checks that the created user is in the list
func TestAPIGetAllUsers(t *testing.T) {
	// first we need to create a user, then get it back
	ur := UserRequest{
		Username:  "testapigetallusers",
		Email:     "testapigetallusers@localhost",
		FirstName: "TestAPI",
		LastName:  "GetAllUsers",
	}
	userReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
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
		if user.Username == ur.Username {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find user %v in the list of users", ur.Username)
	}
}

func TestAPIUpdateUser(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a user, then get it back
	ur := UserRequest{
		Username:  "testapiupdateuser",
		Email:     "testapiupdateuser@localhost",
		FirstName: "TestAPI",
		LastName:  "UpdateUser",
	}
	userReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
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
	if u.Username != ur.Username {
		t.Errorf("expected username %v got %v", ur.Username, u.Username)
	}
	if u.Email != ur.Email {
		t.Errorf("expected email %v got %v", ur.Email, u.Email)
	}
	if u.FirstName != ur.FirstName {
		t.Errorf("expected firstname %v got %v", ur.FirstName, u.FirstName)
	}
	if u.LastName != ur.LastName {
		t.Errorf("expected lastname %v got %v", ur.LastName, u.LastName)
	}

	// now update the user
	ur2 := UserRequest{
		Username:  "testapiupdateuser",
		Email:     "testapiupdateuser2@localhost",
		FirstName: "TestAPI2",
		LastName:  "UpdateUser2",
	}
	userReq2, err := json.Marshal(ur2)
	if err != nil {
		t.Fatal(err)
	}
	updateURL := fmt.Sprintf("http://localhost:3333/api/v1/users/%d", userResponse.Id)
	req, err = http.NewRequest("PUT", updateURL, bytes.NewBuffer([]byte(userReq2)))
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusCreated)
	}
	userResponse = UserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		t.Fatal(err)
	}
	u, err = data.GetUserById(th.DB, userResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != ur2.Username {
		t.Errorf("expected username %v got %v", ur2.Username, u.Username)
	}
	if u.Email != ur2.Email {
		t.Errorf("expected email %v got %v", ur2.Email, u.Email)
	}
	if u.FirstName != ur2.FirstName {
		t.Errorf("expected firstname %v got %v", ur2.FirstName, u.FirstName)
	}
	if u.LastName != ur2.LastName {
		t.Errorf("expected lastname %v got %v", ur2.LastName, u.LastName)
	}
}

func TestAPIDeleteUser(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a user, then delete it
	ur := UserRequest{
		Username:  "testapideleteuser",
		Email:     "testapideleteuser@localhost",
		FirstName: "TestAPI",
		LastName:  "deleteUser",
	}
	userReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
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
	if u.Username != ur.Username {
		t.Errorf("expected username %v got %v", ur.Username, u.Username)
	}
	if u.Email != ur.Email {
		t.Errorf("expected email %v got %v", ur.Email, u.Email)
	}
	if u.FirstName != ur.FirstName {
		t.Errorf("expected firstname %v got %v", ur.FirstName, u.FirstName)
	}
	if u.LastName != ur.LastName {
		t.Errorf("expected lastname %v got %v", ur.LastName, u.LastName)
	}

	// now delete the user
	deleteURL := fmt.Sprintf("http://localhost:3333/api/v1/users/%d", userResponse.Id)
	req, err = http.NewRequest("DELETE", deleteURL, bytes.NewBuffer([]byte(userReq)))
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}

	// now try to get the user again, should get an error
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
		t.Fatalf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}
	var usersResponse []UserResponse
	err = json.NewDecoder(resp.Body).Decode(&usersResponse)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the user is not in the list
	found := false
	for _, user := range usersResponse {
		if user.Username == ur.Username {
			found = true
		}
	}
	if found {
		t.Error("found user that should have been deleted")
	}
}
