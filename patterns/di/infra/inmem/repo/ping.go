package repo

import "github.com/kei2100/playground-go/patterns/di/repo"

// Ping returns an implementation of the PingRepository
func Ping() repo.PingRepository {
	return &ping{}
}

// ping implements the PingRepository
type ping struct{}

func (p *ping) Ping() bool {
	return true
}
