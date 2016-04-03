package main

import "errors"

// DataStorer defines all the database operations
type DataStorer interface {
	ListUsers() (map[string][]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) (bool, error)
}

// MockDB will hold the connection and key db info
type MockDB struct {
	UserList  map[int]User
	MaxUserID int
}

// ListUsers returns a list of JSON documents
func (db *MockDB) ListUsers() (map[string][]User, error) {
	var list []User
	for _, v := range db.UserList {
		list = append(list, v)
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	return responseObject, nil
}

// GetUser returns a single JSON document
func (db *MockDB) GetUser(i int) (User, error) {
	user, ok := db.UserList[i]
	if !ok {
		return user, errors.New("user does not exist")
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
		return u, errors.New("user does not exist")
	}
	db.UserList[id] = u
	return db.UserList[id], nil
}

// DeleteUser deletes a user
func (db *MockDB) DeleteUser(i int) (bool, error) {
	_, ok := db.UserList[i]
	if !ok {
		return false, errors.New("could not delete this user")
	}
	delete(db.UserList, i)
	return true, nil
}
