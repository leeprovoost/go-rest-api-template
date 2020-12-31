package passport

import (
	"testing"
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	list, _ := appEnv.UserStore.ListUsers()
	count := len(list)
	assert.Equal(t, 2, count, "There should be 2 items in the list.")
}

func TestGetUserSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u, err := appEnv.UserStore.GetUser(0)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, u.ID, "they should be equal")
		assert.Equal(t, "John", u.FirstName, "they should be equal")
		assert.Equal(t, "Doe", u.LastName, "they should be equal")
		assert.Equal(t, dt, u.DateOfBirth, "they should be equal")
		assert.Equal(t, "London", u.LocationOfBirth, "they should be equal")
	}
}

func TestGetUserFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	_, err := appEnv.UserStore.GetUser(10)
	assert.NotNil(t, err)
}

func TestAddUser(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	dt, _ := time.Parse(time.RFC3339, "1972-03-07T00:00:00Z")
	u := models.User{
		FirstName:       "Apple",
		LastName:        "Jack",
		DateOfBirth:     dt,
		LocationOfBirth: "Cambridge",
	}
	u, _ = appEnv.UserStore.AddUser(u)
	// we should now have a user object with a database Id
	assert.Equal(t, 2, u.ID, "Expected database Id should be 2.")
	// we should now have 3 items in the list
	list, _ := appEnv.UserStore.ListUsers()
	count := len(list)
	assert.Equal(t, 3, count, "There should be 3 items in the list.")
}

func TestUpdateUserSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := models.User{
		ID:              0,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	// check if there are no errors
	u2, err := appEnv.UserStore.UpdateUser(u)
	assert.Nil(t, err)
	// check returned user
	assert.Equal(t, 0, u2.ID, "they should be equal")
	assert.Equal(t, "John", u2.FirstName, "they should be equal")
	assert.Equal(t, "2 Doe", u2.LastName, "they should be equal")
	assert.Equal(t, dt, u2.DateOfBirth, "they should be equal")
	assert.Equal(t, "Southend", u2.LocationOfBirth, "they should be equal")
}

func TestUpdateUserFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := models.User{
		ID:              20,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	_, err := appEnv.UserStore.UpdateUser(u)
	assert.NotNil(t, err)
}

func TestDeleteUserSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.UserStore.DeleteUser(1)
	assert.Nil(t, err)
}

func TestDeleteUserFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.UserStore.DeleteUser(10)
	assert.NotNil(t, err)
}
