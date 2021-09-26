package service

import "github.com/kei2100/playground-go/patterns/di/repo"

// Ping test
func Ping() bool {
	return repo.Ping.Ping()
}
