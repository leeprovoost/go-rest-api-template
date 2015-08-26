package main

type User struct {
	Id              int    `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DateOfBirth     string `json:"dateOfBirth"`
	LocationOfBirth string `json:"locationOfBirth"`
}
