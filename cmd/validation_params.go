package cmd

import (
	"fmt"
	"regexp"
)

// Validation params.
const (
	clientIDMinLen = 1
	clientIDMaxLen = 100
)

// Validation params that can't be Go constants.
var (
	clientIDRegexp = regexp.MustCompile("^[a-zA-Z0-9-@._]*$")
)

// All validation errors.
var (
	errClientID = fmt.Errorf("a client id should be between %d and %d chars, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())
)
