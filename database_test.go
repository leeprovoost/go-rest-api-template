package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	list := make(map[int]User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = User{0, "John", "Doe", dt, "London"}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = User{1, "Jane", "Doe", dt, "Milton Keynes"}
	db = &Database{list, 1}
	retCode := m.Run()
	os.Exit(retCode)
}

func TestList(t *testing.T) {
	list, _ := db.List()
	count := len(list["users"])
	assert.Equal(t, 2, count, "There should be 2 items in the list.")
}

func TestGetSuccess(t *testing.T) {
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u, err := db.Get(0)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, u.ID, "they should be equal")
		assert.Equal(t, "John", u.FirstName, "they should be equal")
		assert.Equal(t, "Doe", u.LastName, "they should be equal")
		assert.Equal(t, dt, u.DateOfBirth, "they should be equal")
		assert.Equal(t, "London", u.LocationOfBirth, "they should be equal")
	}
}

func TestGetFail(t *testing.T) {
	_, err := db.Get(10)
	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	dt, _ := time.Parse(time.RFC3339, "1972-03-07T00:00:00Z")
	u := User{
		FirstName:       "Apple",
		LastName:        "Jack",
		DateOfBirth:     dt,
		LocationOfBirth: "Cambridge",
	}
	u, _ = db.Add(u)
	// we should now have a user object with a database Id
	assert.Equal(t, 2, u.ID, "Expected database Id should be 2.")
	// we should now have 3 items in the list
	list, _ := db.List()
	count := len(list["users"])
	assert.Equal(t, 3, count, "There should be 3 items in the list.")
}

func TestUpdateSuccess(t *testing.T) {
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := User{
		ID:              0,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	// check if there are no errors
	u2, err := db.Update(u)
	assert.Nil(t, err)
	// check returned user
	assert.Equal(t, 0, u2.ID, "they should be equal")
	assert.Equal(t, "John", u2.FirstName, "they should be equal")
	assert.Equal(t, "2 Doe", u2.LastName, "they should be equal")
	assert.Equal(t, dt, u2.DateOfBirth, "they should be equal")
	assert.Equal(t, "Southend", u2.LocationOfBirth, "they should be equal")
}

func TestUpdateFail(t *testing.T) {
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	u := User{
		ID:              20,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Southend",
	}
	_, err := db.Update(u)
	assert.NotNil(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	ok, err := db.Delete(1)
	assert.Equal(t, true, ok, "they should be equal")
	assert.Nil(t, err)
}

func TestDeleteFail(t *testing.T) {
	ok, err := db.Delete(10)
	assert.Equal(t, false, ok, "they should be equal")
	assert.NotNil(t, err)
}
