package main

// DataStorer defines all the database operations
type DataStorer interface {
	ListUsers() ([]User, error)
	GetUser(i int) (User, error)
	AddUser(u User) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(i int) error
}

// MockDB will hold the connection and key db info
type MockDB struct {
	UserList  map[int]User
	MaxUserID int
}
