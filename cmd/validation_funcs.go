package cmd

import (
	"net/url"
	"strings"
)

// checkBaseURL checks if the provided baseURL is valid.
func checkBaseURL(inputURL string) error {
	// Should not be empty.
	if inputURL == "" {
		return errURL
	}

	// Should start with a valid protocol.
	if !strings.HasPrefix(inputURL, "ws://") && !strings.HasPrefix(inputURL, "wss://") &&
		!strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		return errURL
	}

	// Ensure the URL is parsable before it's used in network requests.
	// This prevents panics or malformed requests in the API client layer.
	if _, err := url.Parse(inputURL); err != nil {
		return errURL
	}

	return nil
}

// checkClientID checks if the provided client ID is valid.
func checkClientID(clientID string) error {
	clientIDLen := len(clientID)
	// Validating the length of client ID.
	if clientIDLen < clientIDMinLen || clientIDLen > clientIDMaxLen {
		return errClientID
	}
	// Validating format of client ID.
	if !clientIDRegexp.MatchString(clientID) {
		return errClientID
	}
	return nil
}

// checkClientIDSlice checks if the provided slice of client IDs is valid.
func checkClientIDSlice(clientIDs []string) error {
	// Calling the checkClientID function upon every client ID in the slice.
	for _, id := range clientIDs {
		if err := checkClientID(id); err != nil {
			return err
		}
	}
	return nil
}
