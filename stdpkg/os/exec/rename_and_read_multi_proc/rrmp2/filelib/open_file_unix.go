// +build darwin freebsd linux

package filelib

import "os"

// OpenLogFileToRead calls os.OpenFile
func OpenLogFileToRead(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	return os.OpenFile(name, flag, perm)
}
