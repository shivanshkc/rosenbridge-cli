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
