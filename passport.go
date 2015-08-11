package main

type Passport struct {
	Id           string `json:"id"`
	DateOfIssue  string `json:"date_of_issue"`
	DateOfExpiry string `json:"date_of_expiry"`
	Authority    string `json:"authority"`
	UserId       int    `json:"user_id"`
}
