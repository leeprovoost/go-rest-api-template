package models

import (
	"context"
	"time"
)

// User holds personal user information.
type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	LocationOfBirth string    `json:"locationOfBirth"`
}

// UserStorage defines all the database operations for users.
type UserStorage interface {
	ListUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id int) (User, error)
	AddUser(ctx context.Context, u User) (User, error)
	UpdateUser(ctx context.Context, u User) (User, error)
	DeleteUser(ctx context.Context, id int) error
}
