package main

import "github.com/palantir/stacktrace"

// DataStorer defines all the database operations
type DataStorer interface {
	ListUsers() ([]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) error
}

// MockDB will hold the connection and key db info
type MockDB struct {
	UserList  map[int]User
	MaxUserID int
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
		return User{}, stacktrace.NewError("Failure trying to retrieve user")
	}
	return user, nil
}

// AddUser adds a User JSON document, returns the JSON document with the generated id
func (db *MockDB) AddUser(u User) (User, error) {
	db.MaxUserID = db.MaxUserID + 1
	u.ID = db.MaxUserID
	db.UserList[db.MaxUserID] = u
	return u, nil
}

// UpdateUser updates an existing user
func (db *MockDB) UpdateUser(u User) (User, error) {
	id := u.ID
	_, ok := db.UserList[id]
	if !ok {
		return u, stacktrace.NewError("Failure trying to update user")
	}
	db.UserList[id] = u
	return db.UserList[id], nil
}

// DeleteUser deletes a user
func (db *MockDB) DeleteUser(i int) error {
	_, ok := db.UserList[i]
	if !ok {
		return stacktrace.NewError("Failure trying to delete user")
	}
	delete(db.UserList, i)
	return nil
}
