package health

// Check holds health check response data.
type Check struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}
