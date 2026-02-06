package passport

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
)

// Compile-time proof of interface implementation.
var _ models.UserStorage = (*UserService)(nil)

// UserService is an in-memory implementation of models.UserStorage.
type UserService struct {
	UserList  map[int]models.User
	MaxUserID int
}

// NewUserService creates a new UserService with the given data.
func NewUserService(list map[int]models.User, count int) models.UserStorage {
	return &UserService{
		UserList:  list,
		MaxUserID: count,
	}
}

// ListUsers returns all users sorted by ID.
func (s *UserService) ListUsers(_ context.Context) ([]models.User, error) {
	users := make([]models.User, 0, len(s.UserList))
	for _, v := range s.UserList {
		users = append(users, v)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users, nil
}

// GetUser returns a single user by ID.
func (s *UserService) GetUser(_ context.Context, id int) (models.User, error) {
	user, ok := s.UserList[id]
	if !ok {
		return models.User{}, fmt.Errorf("user %d not found", id)
	}
	return user, nil
}

// AddUser adds a new user with an auto-generated ID.
func (s *UserService) AddUser(_ context.Context, u models.User) (models.User, error) {
	s.MaxUserID++
	u.ID = s.MaxUserID
	s.UserList[s.MaxUserID] = u
	return u, nil
}

// UpdateUser replaces an existing user.
func (s *UserService) UpdateUser(_ context.Context, u models.User) (models.User, error) {
	if _, ok := s.UserList[u.ID]; !ok {
		return u, fmt.Errorf("user %d not found", u.ID)
	}
	s.UserList[u.ID] = u
	return s.UserList[u.ID], nil
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(_ context.Context, id int) error {
	if _, ok := s.UserList[id]; !ok {
		return fmt.Errorf("user %d not found", id)
	}
	delete(s.UserList, id)
	return nil
}

// CreateMockDataSet returns test data: a map of users and the max user ID.
func CreateMockDataSet() (map[int]models.User, int) {
	list := make(map[int]models.User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = models.User{
		ID:              0,
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "London",
	}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = models.User{
		ID:              1,
		FirstName:       "Jane",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Milton Keynes",
	}
	return list, len(list) - 1
}
