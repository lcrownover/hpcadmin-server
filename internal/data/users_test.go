// TODO(lcrown): table test for different database drivers?
package data

import (
	"database/sql"
	"testing"
)

func TestGetUserById(t *testing.T) { 
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=hpcadmin sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
    WipeDB(db)
	defer db.Close()
    ur := UserRequest{
        Username: "testuser",
        Email: "testuser@localhost",
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
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=hpcadmin sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
    WipeDB(db)
	defer db.Close()
    ur := UserRequest{
        Username: "testuser",
        Email: "testuser@localhost",
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

func TestGetAllUsers(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=hpcadmin sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
    WipeDB(db)
	defer db.Close()
	ur1 := UserRequest{
		Username:  "testuser",
		Email:     "testuser@localhost",
		FirstName: "Test",
		LastName:  "User",
	}
	ur2 := UserRequest{
		Username:  "testuser2",
		Email:     "testuser2@localhost",
		FirstName: "Test",
		LastName:  "User2",
	}
	_, err = CreateUser(db, &ur1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateUser(db, &ur2)
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
	if len(users) != 2 {
		t.Fatal("expected exactly two users")
	}
}
