package main

type Database struct {
	UserList  map[int]User
	MaxUserId int
}

var db *Database
