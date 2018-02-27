package os

import "os"

// GetenvOrDefault call os.Getenv by the key, and returns the value.
// Which will be the default value if the variable is not present.
func GetenvOrDefault(key, defaultV string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		v = defaultV
	}
	return v
}
