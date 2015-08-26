package main

type Passport struct {
	Id           string `json:"id"`
	DateOfIssue  string `json:"dateOfIssue"`
	DateOfExpiry string `json:"dateOfExpiry"`
	Authority    string `json:"authority"`
	UserId       int    `json:"userId"`
}
