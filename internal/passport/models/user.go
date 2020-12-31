package models

import (
	"fmt"
	"time"
)

// User holds personal user information
type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	LocationOfBirth string    `json:"locationOfBirth"`
}

// UserStorage defines all the database operations
type UserStorage interface {
	ListUsers() ([]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) error
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (u *User) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	FirstName: %s,
	LastName: %s,
	DateOfBirth: %s,
	LocationOfBirth: %s,
}`,
		u.ID,
		u.FirstName,
		u.LastName,
		u.DateOfBirth,
		u.LocationOfBirth,
	)
}
