package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	list := make(map[int]User)
	list[0] = User{0, "John", "Doe", "31-12-1985", "London"}
	list[1] = User{1, "Jane", "Doe", "01-01-1992", "Milton Keynes"}
	db = &Database{list, 1}
	retCode := m.Run()
	os.Exit(retCode)
}

func TestList(t *testing.T) {
	list := db.List()
	count := len(list["users"])
	assert.Equal(t, 2, count, "There should be 2 items in the list.")
}

func TestGetSuccess(t *testing.T) {
	u, err := db.Get(0)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, u.Id, "they should be equal")
		assert.Equal(t, "John", u.FirstName, "they should be equal")
		assert.Equal(t, "Doe", u.LastName, "they should be equal")
		assert.Equal(t, "31-12-1985", u.DateOfBirth, "they should be equal")
		assert.Equal(t, "London", u.LocationOfBirth, "they should be equal")
	}
}

func TestGetFail(t *testing.T) {
	_, err := db.Get(10)
	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	u := User{
		FirstName:       "Apple",
		LastName:        "Jack",
		DateOfBirth:     "07-03-1972",
		LocationOfBirth: "Cambridge",
	}
	u = db.Add(u)
	// we should now have a user object with a database Id
	assert.Equal(t, 2, u.Id, "Expected database Id should be 2.")
	// we should now have 3 items in the list
	list := db.List()
	count := len(list["users"])
	assert.Equal(t, 3, count, "There should be 3 items in the list.")
}

func TestUpdateSuccess(t *testing.T) {
	u := User{
		Id:              0,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     "31-12-1985",
		LocationOfBirth: "Southend",
	}
	// check if there are no errors
	u2, err := db.Update(u)
	assert.Nil(t, err)
	// check returned user
	assert.Equal(t, 0, u2.Id, "they should be equal")
	assert.Equal(t, "John", u2.FirstName, "they should be equal")
	assert.Equal(t, "2 Doe", u2.LastName, "they should be equal")
	assert.Equal(t, "31-12-1985", u2.DateOfBirth, "they should be equal")
	assert.Equal(t, "Southend", u2.LocationOfBirth, "they should be equal")
}

func TestUpdateFail(t *testing.T) {
	u := User{
		Id:              20,
		FirstName:       "John",
		LastName:        "2 Doe",
		DateOfBirth:     "31-12-1985",
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
