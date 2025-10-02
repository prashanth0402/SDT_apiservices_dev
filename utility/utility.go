package utility

// CheckError returns true if error exists and prints error code + message
func CheckError(err error) bool {
	return err != nil
}
