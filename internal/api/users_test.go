package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

type testDataHandler struct {
	db *sql.DB
}

func newTestDataHandler() *testDataHandler {
	dbr := data.DBRequest{
		Host:       "localhost",
		Port:       5432,
		User:       "postgres",
		Password:   "postgres",
		DisableSSL: true,
	}
	db, err := data.NewDBConn(dbr)
	if err != nil {
		log.Fatal(err)
	}
	return &testDataHandler{
		db: db,
	}
}

func TestCreateUser(t *testing.T) {
	th := newTestDataHandler()

	// first we need to create a user, then get it back
	userReq := `{"username": "lcrownover", "email": "lcrownover@localhost", "firstname": "Lucas", "lastname": "Crownover"}`
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/users", bytes.NewBuffer([]byte(userReq)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	var userResp UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		t.Fatal(err)
	}

	u, err := data.GetUserById(th.db, userResp.Id)
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

func TestGetAllUsers(t *testing.T) {
	// Create a new request to the /users endpoint
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	tdh := newTestDataHandler()
	tuh := &UserHandler{tdh.db}

	// Create a new handler and call it's ServeHTTP method passing in the
	// the response recorder and the request object
	handler := http.HandlerFunc(tuh.GetAllUsers)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"username":"lcrownover","email":"lcrownover@localhost","firstname":"Lucas","lastname":"Crownover"}]`
	actual := rr.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}
