package posfile

import (
	"encoding/gob"
	"log"
	"os"
)

// PositionFile interface
type PositionFile interface {
	// Close closes this PositionFile
	Close() error
	// FileInfo returns os.FileInfo
	// TODO instead of os.FileInfo
	FileInfo() os.FileInfo
	// Offset returns offset value
	Offset() int64
	// IncreaseOffset increases offset value
	IncreaseOffset(incr int)
	// Update updates fileInfo and offset
	Update(fileInfo os.FileInfo, offset int64)
}

type entry struct {
	fileInfo os.FileInfo
	offset   int64
}

func (pf *entry) FileInfo() os.FileInfo {
	return pf.fileInfo
}

func (pf *entry) Offset() int64 {
	return pf.offset
}

func (pf *entry) IncreaseOffset(incr int) {
	pf.offset += int64(incr)
}

func (pf *entry) Update(fileInfo os.FileInfo, offset int64) {
	pf.fileInfo = fileInfo
	pf.offset = offset
}

// Open opens named PositionFile
func Open(name string) (PositionFile, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var ent entry
	if fi.Size() == 0 {
		return &file{f: f, entry: ent}, nil
	}
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&ent); err != nil {
		return nil, err
	}
	return &file{f: f, entry: ent}, nil
}

type file struct {
	f *os.File
	entry
}

func (f *file) Close() error {
	return f.f.Close()
}

func (f *file) IncreaseOffset(incr int) {
	// TODO
	f.Update(f.FileInfo(), f.Offset()+int64(incr))
}

func (f *file) Update(fileInfo os.FileInfo, offset int64) {
	f.entry.Update(fileInfo, offset)
	if _, err := f.f.Seek(0, 0); err != nil {
		log.Print(err)
		// TODO
	}
	enc := gob.NewEncoder(f.f)
	if err := enc.Encode(&f.entry); err != nil {
		log.Print(err)
		// TODO
	}
}

// InMemory creates a inMemory PositionFile
func InMemory(fileInfo os.FileInfo, offset int64) PositionFile {
	return &inMemory{entry{fileInfo: fileInfo, offset: offset}}
}

type inMemory struct {
	entry
}

func (pf *inMemory) Close() error {
	return nil
}
