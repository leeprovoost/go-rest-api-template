package models

import "fmt"

// Status is a custom response object we pass around the system and send back to the customer
// 404: Not found
// 500: Internal Server Error
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (s *Status) GoString() string {
	return fmt.Sprintf(`
{
	Status: %s,
	Message: %s,
}`,
		s.Status,
		s.Message,
	)
}
