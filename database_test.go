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
	var list map[string][]User = db.List()
	var count = len(list["users"])
	assert.Equal(t, count, 2, "they should be equal")
}

func TestGetSuccess(t *testing.T) {
	u, err := db.Get(0)
	if assert.Nil(t, err) {
		assert.Equal(t, u.Id, 0, "they should be equal")
		assert.Equal(t, u.FirstName, "John", "they should be equal")
		assert.Equal(t, u.LastName, "Doe", "they should be equal")
		assert.Equal(t, u.DateOfBirth, "31-12-1985", "they should be equal")
		assert.Equal(t, u.LocationOfBirth, "London", "they should be equal")
	}
}

func TestGetFail(t *testing.T) {
	_, err := db.Get(10)
	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	t.Errorf("Test not implemented.")
}

func TestDelete(t *testing.T) {
	t.Errorf("Test not implemented.")
}

func TestUpdate(t *testing.T) {
	t.Errorf("Test not implemented.")
}
