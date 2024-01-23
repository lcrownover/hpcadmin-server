package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

func TestAPICreatePirg(t *testing.T) {
	th := NewTestDataHandler()
	// TODO(lcrown): create a user in the beginning of all these, then fix all the tests

	// first we need to create a pirg, then get it back
	ur := PirgRequest{
		Name:  "testapicreatepirg",
		Email:     "testapicreatepirg@localhost",
		FirstName: "TestAPI",
		LastName:  "CreatePirg",
	}
	pirgReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/pirgs", bytes.NewBuffer([]byte(pirgReq)))
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
	var pirgResponse PirgResponse
	err = json.NewDecoder(resp.Body).Decode(&pirgResponse)
	if err != nil {
		t.Fatal(err)
	}

	u, err := data.GetPirgById(th.DB, pirgResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Pirgname != ur.Pirgname {
		t.Errorf("expected pirgname %v got %v", ur.Pirgname, u.Pirgname)
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

// TestGetAllPirgs tests the GET /api/v1/pirgs endpoint
// it creates a pirg, then gets all pirgs and checks that the created pirg is in the list
func TestAPIGetAllPirgs(t *testing.T) {
	// first we need to create a pirg, then get it back
	ur := PirgRequest{
		Pirgname:  "testapigetallpirgs",
		Email:     "testapigetallpirgs@localhost",
		FirstName: "TestAPI",
		LastName:  "GetAllPirgs",
	}
	pirgReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/pirgs", bytes.NewBuffer([]byte(pirgReq)))
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

	// now get all pirgs
	req, err = http.NewRequest("GET", "http://localhost:3333/api/v1/pirgs", nil)
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
	var pirgsResponse []PirgResponse
	err = json.NewDecoder(resp.Body).Decode(&pirgsResponse)
	if err != nil {
		t.Fatal(err)
	}

	// check if the pirgsResponse has more than 1 pirg
	if len(pirgsResponse) < 1 {
		t.Errorf("expected more than 1 pirg, got %v", len(pirgsResponse))
	}

	// check if the pirg we created is in the list
	found := false
	for _, pirg := range pirgsResponse {
		if pirg.Pirgname == ur.Pirgname {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find pirg %v in the list of pirgs", ur.Pirgname)
	}
}

func TestAPIUpdatePirg(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a pirg, then get it back
	ur := PirgRequest{
		Pirgname:  "testapiupdatepirg",
		Email:     "testapiupdatepirg@localhost",
		FirstName: "TestAPI",
		LastName:  "UpdatePirg",
	}
	pirgReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/pirgs", bytes.NewBuffer([]byte(pirgReq)))
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
	var pirgResponse PirgResponse
	err = json.NewDecoder(resp.Body).Decode(&pirgResponse)
	if err != nil {
		t.Fatal(err)
	}

	u, err := data.GetPirgById(th.DB, pirgResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Pirgname != ur.Pirgname {
		t.Errorf("expected pirgname %v got %v", ur.Pirgname, u.Pirgname)
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

	// now update the pirg
	ur2 := PirgRequest{
		Pirgname:  "testapiupdatepirg",
		Email:     "testapiupdatepirg2@localhost",
		FirstName: "TestAPI2",
		LastName:  "UpdatePirg2",
	}
	pirgReq2, err := json.Marshal(ur2)
	if err != nil {
		t.Fatal(err)
	}
	updateURL := fmt.Sprintf("http://localhost:3333/api/v1/pirgs/%d", pirgResponse.Id)
	req, err = http.NewRequest("PUT", updateURL, bytes.NewBuffer([]byte(pirgReq2)))
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
	pirgResponse = PirgResponse{}
	err = json.NewDecoder(resp.Body).Decode(&pirgResponse)
	if err != nil {
		t.Fatal(err)
	}
	u, err = data.GetPirgById(th.DB, pirgResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Pirgname != ur2.Pirgname {
		t.Errorf("expected pirgname %v got %v", ur2.Pirgname, u.Pirgname)
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

func TestAPIDeletePirg(t *testing.T) {
	th := NewTestDataHandler()

	// first we need to create a pirg, then delete it
	ur := PirgRequest{
		Pirgname:  "testapideletepirg",
		Email:     "testapideletepirg@localhost",
		FirstName: "TestAPI",
		LastName:  "deletePirg",
	}
	pirgReq, err := json.Marshal(ur)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:3333/api/v1/pirgs", bytes.NewBuffer([]byte(pirgReq)))
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
	var pirgResponse PirgResponse
	err = json.NewDecoder(resp.Body).Decode(&pirgResponse)
	if err != nil {
		t.Fatal(err)
	}

	u, err := data.GetPirgById(th.DB, pirgResponse.Id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Pirgname != ur.Pirgname {
		t.Errorf("expected pirgname %v got %v", ur.Pirgname, u.Pirgname)
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

	// now delete the pirg
	deleteURL := fmt.Sprintf("http://localhost:3333/api/v1/pirgs/%d", pirgResponse.Id)
	req, err = http.NewRequest("DELETE", deleteURL, bytes.NewBuffer([]byte(pirgReq)))
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

	// now try to get the pirg again, should get an error
	req, err = http.NewRequest("GET", "http://localhost:3333/api/v1/pirgs", nil)
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
	var pirgsResponse []PirgResponse
	err = json.NewDecoder(resp.Body).Decode(&pirgsResponse)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the pirg is not in the list
	found := false
	for _, pirg := range pirgsResponse {
		if pirg.Pirgname == ur.Pirgname {
			found = true
		}
	}
	if found {
		t.Error("found pirg that should have been deleted")
	}
}
