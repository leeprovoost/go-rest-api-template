package passport

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
)

// Compile-time proof of interface implementation.
var _ models.PassportStorage = (*PassportService)(nil)

// PassportService is an in-memory implementation of models.PassportStorage.
type PassportService struct {
	PassportList map[string]models.Passport
}

// NewPassportService creates a new PassportService with the given data.
func NewPassportService(list map[string]models.Passport) models.PassportStorage {
	return &PassportService{
		PassportList: list,
	}
}

// ListPassportsByUser returns all passports belonging to a user, sorted by ID.
func (s *PassportService) ListPassportsByUser(_ context.Context, userID int) ([]models.Passport, error) {
	var passports []models.Passport
	for _, p := range s.PassportList {
		if p.UserID == userID {
			passports = append(passports, p)
		}
	}
	if passports == nil {
		passports = []models.Passport{}
	}
	sort.Slice(passports, func(i, j int) bool {
		return passports[i].ID < passports[j].ID
	})
	return passports, nil
}

// GetPassport returns a single passport by ID.
func (s *PassportService) GetPassport(_ context.Context, id string) (models.Passport, error) {
	p, ok := s.PassportList[id]
	if !ok {
		return models.Passport{}, fmt.Errorf("passport %q not found", id)
	}
	return p, nil
}

// AddPassport stores a new passport. The client provides the passport ID.
func (s *PassportService) AddPassport(_ context.Context, p models.Passport) (models.Passport, error) {
	if _, exists := s.PassportList[p.ID]; exists {
		return models.Passport{}, fmt.Errorf("passport %q already exists", p.ID)
	}
	s.PassportList[p.ID] = p
	return p, nil
}

// UpdatePassport replaces an existing passport.
func (s *PassportService) UpdatePassport(_ context.Context, p models.Passport) (models.Passport, error) {
	if _, ok := s.PassportList[p.ID]; !ok {
		return p, fmt.Errorf("passport %q not found", p.ID)
	}
	s.PassportList[p.ID] = p
	return s.PassportList[p.ID], nil
}

// DeletePassport removes a passport by ID.
func (s *PassportService) DeletePassport(_ context.Context, id string) error {
	if _, ok := s.PassportList[id]; !ok {
		return fmt.Errorf("passport %q not found", id)
	}
	delete(s.PassportList, id)
	return nil
}

// CreateMockPassportDataSet returns test passport data.
func CreateMockPassportDataSet() map[string]models.Passport {
	list := make(map[string]models.Passport)
	doi, _ := time.Parse(time.RFC3339, "2020-01-15T00:00:00Z")
	doe, _ := time.Parse(time.RFC3339, "2030-01-15T00:00:00Z")
	list["012345678"] = models.Passport{
		ID:           "012345678",
		DateOfIssue:  doi,
		DateOfExpiry: doe,
		Authority:    "HMPO",
		UserID:       0,
	}
	doi, _ = time.Parse(time.RFC3339, "2019-06-01T00:00:00Z")
	doe, _ = time.Parse(time.RFC3339, "2029-06-01T00:00:00Z")
	list["987654321"] = models.Passport{
		ID:           "987654321",
		DateOfIssue:  doi,
		DateOfExpiry: doe,
		Authority:    "HMPO",
		UserID:       1,
	}
	return list
}
