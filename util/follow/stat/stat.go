package stat

import "os"

// Stat returns the named FileStat
func Stat(file *os.File) (*FileStat, error) {
	return stat(file)
}

// SameFile reports whether st1 and st2 are the same
func SameFile(st1, st2 *FileStat) bool {
	return st1.sameFile(st2)
}
