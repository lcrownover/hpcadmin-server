// TODO(lcrown): table test for different database drivers?
package data

import (
	"testing"
)

func TestGetUserById(t *testing.T) { 
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
    ur := UserRequest{
        Username: "testgetuserbyid",
        Email: "testgetuserbyid@localhost",
        FirstName: "Test",
        LastName: "User",
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

func TestCreateUser(t *testing.T) { 
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
    ur := UserRequest{
        Username: "testcreateuser",
        Email: "testcreateuser@localhost",
        FirstName: "Test",
        LastName: "User",
    }
    user, err := CreateUser(db, &ur)
    if err != nil {
        t.Fatal(err)
    }
    if user.Id < 1 {
        t.Fatal("expected id to be greater than 0")
    }
    if user.Username != ur.Username {
        t.Fatal("expected usernames to match")
    }
    if user.Email != ur.Email {
        t.Fatal("expected emails to match")
    }
    if user.FirstName != ur.FirstName {
        t.Fatal("expected first names to match")
    }
    if user.LastName != ur.LastName {
        t.Fatal("expected last names to match")
    }
}

func TestUpdateUser(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
    ur := UserRequest{
        Username: "testupdateuser",
        Email: "testupdateuser@localhost",
        FirstName: "Test",
        LastName: "User",
    }
    user, err := CreateUser(db, &ur)
    if err != nil {
        t.Fatal(err)
    }
    updatedUr := UserRequest{
        Username: "testupdateuser2",
        Email: "testupdateuser2@localhost",
        FirstName: "Test2",
        LastName: "User2",
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
        t.Fatal("expected usernames to match")
    }
    if user2.Email != updatedUr.Email {
        t.Fatal("expected emails to match")
    }
    if user2.FirstName != updatedUr.FirstName {
        t.Fatal("expected first names to match")
    }
    if user2.LastName != updatedUr.LastName {
        t.Fatal("expected last names to match")
    }
}

func TestDeleteUser(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
    ur := UserRequest{
        Username: "testdeleteuser",
        Email: "testdeleteuser@localhost",
        FirstName: "Test",
        LastName: "User",
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

func TestGetAllUsers(t *testing.T) {
	dh := NewTestDataHandler()
	db := dh.DB
	defer db.Close()
    ur1 := UserRequest{
		Username:  "testgetallusers",
		Email:     "testgetallusers@localhost",
		FirstName: "Test",
		LastName:  "User",
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
