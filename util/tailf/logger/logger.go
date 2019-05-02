package logger

import (
	"log"
	"os"
	"sync"
)

var (
	logger   Logger = log.New(os.Stderr, "", log.LstdFlags)
	loggerMu sync.RWMutex
)

type Logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func Set(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger = l
}

func Print(v ...interface{}) {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	logger.Print(v...)
}

func Println(v ...interface{}) {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	logger.Println(v...)
}

func Printf(format string, v ...interface{}) {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	logger.Printf(format, v...)
}
