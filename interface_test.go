package main

import "testing"

// TestDoStructsSatisfyInterface is a helper test function that just validates
// whether our data structs are satisfying our DataStorer struct
func TestDoStructsSatisfyInterface(t *testing.T) {
	var _ DataStorer = (*MockDB)(nil)
}
