package main

import "time"

// Passport holds passport data
type Passport struct {
	ID           string    `json:"id"`
	DateOfIssue  time.Time `json:"dateOfIssue"`
	DateOfExpiry time.Time `json:"dateOfExpiry"`
	Authority    string    `json:"authority"`
	UserID       int       `json:"userId"`
}
