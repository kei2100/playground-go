package os

import "os"

// GetenvOrDefault retrieves the value of the environment variable named by the key.
// If the variable is not present, returns the default value
func GetenvOrDefault(key, defaultV string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultV
}
