// TODO(lcrown): table test for different database drivers?
package data

import (
	"testing"
)

func TestDataGetUserById(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
	ur := UserRequest{
		Username:  "testgetuserbyid",
		Email:     "testgetuserbyid@localhost",
		FirstName: "Test",
		LastName:  "User",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	if user.Id < 1 {
		t.Fatal("expected id to be greater than 0")
	}
	user2, err := GetUserById(db, user.Id)
	if err != nil {
		t.Fatal(err)
	}
	if user2.Id != user.Id {
		t.Fatal("expected user ids to match")
	}
	if user2.Username != user.Username {
		t.Fatal("expected usernames to match")
	}
}

func TestDataCreateUser(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
	ur := UserRequest{
		Username:  "testdatacreateuser",
		Email:     "testdatacreateuser@localhost",
		FirstName: "TestData",
		LastName:  "CreateUser",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	if user.Id < 1 {
		t.Fatal("expected id to be greater than 0")
	}
	if user.Username != ur.Username {
		t.Fatalf("expected username %v, got %v", ur.Username, user.Username)
	}
	if user.Email != ur.Email {
		t.Fatalf("expected email %v got %v", ur.Email, user.Email)
	}
	if user.FirstName != ur.FirstName {
		t.Fatalf("expected first name %v got %v", ur.FirstName, user.FirstName)
	}
	if user.LastName != ur.LastName {
		t.Fatalf("expected last name %v got %v", ur.LastName, user.LastName)
	}
}

func TestDataUpdateUser(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
	ur := UserRequest{
		Username:  "testdataupdateuser",
		Email:     "testdataupdateuser@localhost",
		FirstName: "TestData",
		LastName:  "UpdateUser",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	updatedUr := UserRequest{
		Username:  "testdataupdateuser2",
		Email:     "testdataupdateuser2@localhost",
		FirstName: "TestData2",
		LastName:  "UpdateUser2",
	}
	err = UpdateUser(db, user.Id, &updatedUr)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := GetUserById(db, user.Id)
	if err != nil {
		t.Fatal(err)
	}
	if user2.Username != updatedUr.Username {
		t.Fatalf("expected username %v got %v", updatedUr.Username, user2.Username)
	}
	if user2.Email != updatedUr.Email {
		t.Fatalf("expected email %v got %v", updatedUr.Email, user2.Email)
	}
	if user2.FirstName != updatedUr.FirstName {
		t.Fatalf("expected first name %v got %v", updatedUr.FirstName, user2.FirstName)
	}
	if user2.LastName != updatedUr.LastName {
		t.Fatalf("expected last name %v got %v", updatedUr.LastName, user2.LastName)
	}
}

func TestDataDeleteUser(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
	ur := UserRequest{
		Username:  "testdatadeleteuser",
		Email:     "testdatadeleteuser@localhost",
		FirstName: "TestData",
		LastName:  "DeleteUser",
	}
	user, err := CreateUser(db, &ur)
	if err != nil {
		t.Fatal(err)
	}
	err = DeleteUser(db, user.Id)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetUserById(db, user.Id)
	if err == nil {
		t.Fatal("expected error getting deleted user")
	}
}

func TestDataGetAllUsers(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
	ur1 := UserRequest{
		Username:  "testdatagetallusers",
		Email:     "testdatagetallusers@localhost",
		FirstName: "TestData",
		LastName:  "GetAllUsers",
	}
	_, err := CreateUser(db, &ur1)
	if err != nil {
		t.Fatal(err)
	}
	users, err := GetAllUsers(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) < 1 {
		t.Fatal("expected at least one user")
	}
}
