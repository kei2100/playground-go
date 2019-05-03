package posfile

import (
	"os"
)

// PositionFile interface
type PositionFile interface {
	// FileInfo returns os.FileInfo
	FileInfo() os.FileInfo
	// Offset returns offset value
	Offset() int64
	// IncreaseOffset increases offset value
	IncreaseOffset(incr int)
	// Update updates fileInfo and offset
	Update(fileInfo os.FileInfo, offset int64)
}

// NewMemoryPositionFile creates a memoryPositionFile
func NewMemoryPositionFile(fileInfo os.FileInfo, offset int64) PositionFile {
	return &memoryPositionFile{fileInfo: fileInfo, offset: offset}
}

type memoryPositionFile struct {
	fileInfo os.FileInfo
	offset   int64
}

func (pf *memoryPositionFile) FileInfo() os.FileInfo {
	return pf.fileInfo
}

func (pf *memoryPositionFile) Offset() int64 {
	return pf.offset
}

func (pf *memoryPositionFile) IncreaseOffset(incr int) {
	pf.offset += int64(incr)
}

func (pf *memoryPositionFile) Update(fileInfo os.FileInfo, offset int64) {
	pf.fileInfo = fileInfo
	pf.offset = offset
}
