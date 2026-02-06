package status

// Response is a custom error response sent back to the client.
// It avoids leaking internal error details.
type Response struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}
