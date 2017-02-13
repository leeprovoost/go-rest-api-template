package models

import "fmt"

// Healthcheck will store information about its name and version
type Healthcheck struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (h *Healthcheck) GoString() string {
	return fmt.Sprintf(`
{
	AppName: %s,
	Version: %s,
}`,
		h.AppName,
		h.Version,
	)
}
