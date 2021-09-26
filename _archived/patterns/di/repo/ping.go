package repo

// PingRepository is an interface for the ping test to the repository
type PingRepository interface {
	Ping() bool
}

// Ping holds a reference to the implementation of the PingRepository
var Ping PingRepository
