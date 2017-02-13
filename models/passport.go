package models

import (
	"fmt"
	"time"
)

// Passport holds passport data
type Passport struct {
	ID           string    `json:"id"`
	DateOfIssue  time.Time `json:"dateOfIssue"`
	DateOfExpiry time.Time `json:"dateOfExpiry"`
	Authority    string    `json:"authority"`
	UserID       int       `json:"userId"`
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (p *Passport) GoString() string {
	return fmt.Sprintf(`
{
	ID: %s,
	DateOfIssue: %s,
	DateOfExpiry: %s,
	Authority: %s,
	UserID: %d,
}`,
		p.ID,
		p.DateOfIssue,
		p.DateOfExpiry,
		p.Authority,
		p.UserID,
	)
}
