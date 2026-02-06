package models

import (
	"context"
	"time"
)

// Passport holds passport data.
type Passport struct {
	ID           string    `json:"id"`
	DateOfIssue  time.Time `json:"dateOfIssue"`
	DateOfExpiry time.Time `json:"dateOfExpiry"`
	Authority    string    `json:"authority"`
	UserID       int       `json:"userId"`
}

// PassportStorage defines all the database operations for passports.
type PassportStorage interface {
	ListPassportsByUser(ctx context.Context, userID int) ([]Passport, error)
	GetPassport(ctx context.Context, id string) (Passport, error)
	AddPassport(ctx context.Context, p Passport) (Passport, error)
	UpdatePassport(ctx context.Context, p Passport) (Passport, error)
	DeletePassport(ctx context.Context, id string) error
}
