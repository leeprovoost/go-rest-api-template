package main

import (
	"os"
	"testing"
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
	if count != 2 {
		t.Errorf("Expected 2 elements in the list, only counting %v.", count)
	}
}

func TestGet(t *testing.T) {
	t.Errorf("Test not implemented.")
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
