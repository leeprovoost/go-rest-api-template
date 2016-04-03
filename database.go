package main

import "errors"

// DataStorer defines all the database operations
type DataStorer interface {
	List() (map[string]User, error)
	Get(i int) (User, error)
	Add(u User) (User, error)
	Update(u User) (User, error)
	Delete(i int) (bool, error)
}

// Database will hold the connection and key db info
type Database struct {
	UserList  map[int]User
	MaxUserID int
}

var db *Database

// List returns a list of JSON documents
func (db *Database) List() (map[string][]User, error) {
	var list []User
	for _, v := range db.UserList {
		list = append(list, v)
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	return responseObject, nil
}

// Get a single JSON document
func (db *Database) Get(i int) (User, error) {
	user, ok := db.UserList[i]
	if !ok {
		return user, errors.New("user does not exist")
	}
	return user, nil
}

// Add a User JSON document, returns the JSON document with the generated id
func (db *Database) Add(u User) (User, error) {
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

// Update an existing user
func (db *Database) Update(u User) (User, error) {
	id := u.ID
	_, ok := db.UserList[id]
	if !ok {
		return u, errors.New("user does not exist")
	}
	db.UserList[id] = u
	return db.UserList[id], nil
}

// Delete a user
func (db *Database) Delete(i int) (bool, error) {
	_, ok := db.UserList[i]
	if !ok {
		return false, errors.New("could not delete this user")
	}
	delete(db.UserList, i)
	return true, nil
}
