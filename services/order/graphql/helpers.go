package graphql

// stringValue converts a string pointer to a string value, returning empty string if nil
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
