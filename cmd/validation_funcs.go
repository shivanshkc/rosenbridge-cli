package cmd

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
