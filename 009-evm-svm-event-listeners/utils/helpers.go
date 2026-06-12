package utils

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
