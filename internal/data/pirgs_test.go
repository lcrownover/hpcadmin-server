package data

import (
	"testing"
)

func TestCreatePirg(t *testing.T) {
	dbr, _ := NewDBRequest("localhost", 5432, "postgres", "postgres", "hpcadmin_test", true)
	db, err := NewDBConn(dbr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ur := UserRequest{
		Username:  "testcreatepirguser",
		Email:     "testcreatepirguser@localhost",
		FirstName: "Test",
		LastName:  "User",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	pr := PirgRequest{
		Name:     "testcreatepirg",
		OwnerId:  user.Id,
		AdminIds: []int{user.Id},
		UserIds:  []int{user.Id},
	}
	pirg, err := CreatePirg(db, &pr)
	if err != nil {
		t.Fatal(err)
	}
	if pirg.Name != pr.Name {
		t.Fatal("expected names to match")
	}
	if pirg.OwnerId != pr.OwnerId {
		t.Fatal("expected owner ids to match")
	}
}

func TestGetAllPirgs(t *testing.T) {
	dbr, _ := NewDBRequest("localhost", 5432, "postgres", "postgres", "hpcadmin_test", true)
	db, err := NewDBConn(dbr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ur := UserRequest{
		Username:  "testgetallpirgsuser",
		Email:     "testgetallpirgsuser@localhost",
		FirstName: "Test",
		LastName:  "User",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	pr1 := PirgRequest{
		Name:     "testcreatepirg1",
		OwnerId:  user.Id,
		AdminIds: []int{user.Id},
		UserIds:  []int{user.Id},
	}
	_, err = CreatePirg(db, &pr1)
	if err != nil {
		t.Fatal(err)
	}
	pirgs, err := GetAllPirgs(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(pirgs) < 1 {
		t.Fatal("expected 2 pirgs: ", len(pirgs))
	}
}
