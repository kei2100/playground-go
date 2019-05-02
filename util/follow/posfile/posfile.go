package posfile

import (
	"os"
)

// PositionFile interface
type PositionFile interface {
	// Update updates fileInfo and offset
	Update(fileInfo os.FileInfo, offset int64)
	// Offset returns offset value
	Offset() int64
	// IncreaseOffset increases offset value
	IncreaseOffset(incr int)
}

// NewMemoryPositionFile creates a memoryPositionFile
func NewMemoryPositionFile(fileInfo os.FileInfo, offset int64) PositionFile {
	return &memoryPositionFile{fileInfo: fileInfo, offset: offset}
}

type memoryPositionFile struct {
	fileInfo os.FileInfo
	offset   int64
}

func (pf *memoryPositionFile) Update(fileInfo os.FileInfo, offset int64) {
	pf.fileInfo = fileInfo
	pf.offset = offset
}

func (pf *memoryPositionFile) Offset() int64 {
	return pf.offset
}

func (pf *memoryPositionFile) IncreaseOffset(incr int) {
	pf.offset += int64(incr)
}
