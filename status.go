package main

// Status is a custom response object we pass around the system and send back to the customer
// 404: Not found
// 500: Internal Server Error
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
