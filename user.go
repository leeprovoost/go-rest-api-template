package main

import (
	"errors"
	"time"

	"github.com/palantir/stacktrace"
)

// User holds personal user information
type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	LocationOfBirth string    `json:"locationOfBirth"`
}

// ListUsers returns a list of JSON documents
func (db *MockDB) ListUsers() ([]User, error) {
	var list []User
	for _, v := range db.UserList {
		list = append(list, v)
	}
	return list, nil
}

// GetUser returns a single JSON document
func (db *MockDB) GetUser(i int) (User, error) {
	user, ok := db.UserList[i]
	if !ok {
		err := errors.New("user does not exist")
		return user, stacktrace.Propagate(err, "Failure trying to retrieve user")
	}
	return user, nil
}

// AddUser adds a User JSON document, returns the JSON document with the generated id
func (db *MockDB) AddUser(u User) (User, error) {
	db.MaxUserID = db.MaxUserID + 1
	newUser := User{
		ID:              db.MaxUserID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	db.UserList[db.MaxUserID] = newUser
	return newUser, nil
}

// UpdateUser updates an existing user
func (db *MockDB) UpdateUser(u User) (User, error) {
	id := u.ID
	_, ok := db.UserList[id]
	if !ok {
		err := errors.New("user does not exist")
		return u, stacktrace.Propagate(err, "Failure trying to update user")
	}
	db.UserList[id] = u
	return db.UserList[id], nil
}

// DeleteUser deletes a user
func (db *MockDB) DeleteUser(i int) error {
	_, ok := db.UserList[i]
	if !ok {
		err := errors.New("user does not exist")
		return stacktrace.Propagate(err, "Failure trying to delete user")
	}
	delete(db.UserList, i)
	return nil
}
