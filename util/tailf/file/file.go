package file

import "os"

func Open(name string) (*os.File, error) {
	return openFile(name, os.O_RDONLY, 0)
}
