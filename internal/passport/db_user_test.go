package passport

import (
	"context"
	"testing"
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUsers(t *testing.T) {
	srv := NewTestServer()
	list, _ := srv.userStore.ListUsers(context.Background())
	assert.Equal(t, 2, len(list), "there should be 2 items in the list")
}

func TestGetUserSuccess(t *testing.T) {
	srv := NewTestServer()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u, err := srv.userStore.GetUser(context.Background(), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, u.ID)
	assert.Equal(t, "John", u.FirstName)
	assert.Equal(t, "Doe", u.LastName)
	assert.Equal(t, dt, u.DateOfBirth)
	assert.Equal(t, "London", u.LocationOfBirth)
}

func TestGetUserFail(t *testing.T) {
	srv := NewTestServer()
	_, err := srv.userStore.GetUser(context.Background(), 10)
	assert.Error(t, err)
}

func TestAddUser(t *testing.T) {
	srv := NewTestServer()
	dt, _ := time.Parse(time.RFC3339, "1972-03-07T00:00:00Z")
	u := models.User{
		FirstName:       "Apple",
		LastName:        "Jack",
		DateOfBirth:     dt,
		LocationOfBirth: "Cambridge",
	}
	u, _ = srv.userStore.AddUser(context.Background(), u)
	assert.Equal(t, 2, u.ID, "expected database ID should be 2")

	list, _ := srv.userStore.ListUsers(context.Background())
	assert.Equal(t, 3, len(list), "there should be 3 items in the list")
}

func TestUpdateUserSuccess(t *testing.T) {
	srv := NewTestServer()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := models.User{
		ID:              0,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	u2, err := srv.userStore.UpdateUser(context.Background(), u)
	require.NoError(t, err)
	assert.Equal(t, 0, u2.ID)
	assert.Equal(t, "John", u2.FirstName)
	assert.Equal(t, "2 Doe", u2.LastName)
	assert.Equal(t, dt, u2.DateOfBirth)
	assert.Equal(t, "Southend", u2.LocationOfBirth)
}

func TestUpdateUserFail(t *testing.T) {
	srv := NewTestServer()
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := models.User{
		ID:              20,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	_, err := srv.userStore.UpdateUser(context.Background(), u)
	assert.Error(t, err)
}

func TestDeleteUserSuccess(t *testing.T) {
	srv := NewTestServer()
	err := srv.userStore.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
}

func TestDeleteUserFail(t *testing.T) {
	srv := NewTestServer()
	err := srv.userStore.DeleteUser(context.Background(), 10)
	assert.Error(t, err)
}
