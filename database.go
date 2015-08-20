package main

import "errors"

type Database struct {
	UserList  map[int]User
	MaxUserId int
}

var db *Database

// List returns a list of JSON documents
func (db *Database) List() map[string][]User {
	var list []User = make([]User, 0)
	for _, v := range db.UserList {
		list = append(list, v)
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	return responseObject
}

// Retrieve a single JSON document
func (db *Database) Get(i int) (User, error) {
	user, ok := db.UserList[i]
	if ok {
		return user, nil
	} else {
		return user, errors.New("User does not exist")
	}
}

// Add a User JSON document, returns the JSON document with the generated id
func (db *Database) Add(u User) User {
	db.MaxUserId = db.MaxUserId + 1
	newUser := User{
		Id:              db.MaxUserId,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	db.UserList[db.MaxUserId] = newUser
	return newUser
}

// Update an existing user
func (db *Database) Update(u User) (User, error) {
	id := u.Id
	_, ok := db.UserList[id]
	if ok {
		db.UserList[id] = u
		return db.UserList[id], nil
	} else {
		return u, errors.New("User does not exist")
	}
}

// Delete a user
func (db *Database) Delete(i int) (bool, error) {
	_, ok := db.UserList[i]
	if ok {
		delete(db.UserList, i)
		return true, nil
	} else {
		return false, errors.New("Could not delete this user")
	}
}
