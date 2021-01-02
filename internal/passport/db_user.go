package passport

import (
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/palantir/stacktrace"
)

// Compile-time proof of interface implementation
var _ models.UserStorage = (*UserService)(nil)

// UserService will hold the connection and key db info
type UserService struct {
	UserList  map[int]models.User
	MaxUserID int
}

// NewUserService creates a new Carer Service with the system's database connection
func NewUserService(list map[int]models.User, count int) models.UserStorage {
	return &UserService{
		UserList:  list,
		MaxUserID: count,
	}
}

// ListUsers returns a list of JSON documents
func (service *UserService) ListUsers() ([]models.User, error) {
	var list []models.User
	for _, v := range service.UserList {
		list = append(list, v)
	}
	return list, nil
}

// GetUser returns a single JSON document
func (service *UserService) GetUser(i int) (models.User, error) {
	user, ok := service.UserList[i]
	if !ok {
		return models.User{}, stacktrace.NewError("Failure trying to retrieve user")
	}
	return user, nil
}

// AddUser adds a User JSON document, returns the JSON document with the generated id
func (service *UserService) AddUser(u models.User) (models.User, error) {
	service.MaxUserID = service.MaxUserID + 1
	u.ID = service.MaxUserID
	service.UserList[service.MaxUserID] = u
	return u, nil
}

// UpdateUser updates an existing user
func (service *UserService) UpdateUser(u models.User) (models.User, error) {
	id := u.ID
	_, ok := service.UserList[id]
	if !ok {
		return u, stacktrace.NewError("Failure trying to update user")
	}
	service.UserList[id] = u
	return service.UserList[id], nil
}

// DeleteUser deletes a user
func (service *UserService) DeleteUser(i int) error {
	_, ok := service.UserList[i]
	if !ok {
		return stacktrace.NewError("Failure trying to delete user")
	}
	delete(service.UserList, i)
	return nil
}

// CreateMockDataSet initialises a database for test purposes. It returns a list of User objects
// as well as the new max object count
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
