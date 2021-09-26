package defaults

import "time"

// Time returns the defaultV if the v is zero. returns v otherwise.
func Time(v, defaultV time.Time) time.Time {
	if v.IsZero() {
		return defaultV
	}
	return v
}
