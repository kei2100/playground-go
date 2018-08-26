package repo

import (
	"database/sql"
	"log"

	"github.com/kei2100/playground-go/patterns/di/repo"
)

// Ping returns an implementation of the PingRepository
func Ping() repo.PingRepository {
	return &ping{db: db}
}

// ping implements the PingRepository
type ping struct {
	db *sql.DB
}

func (p *ping) Ping() bool {
	if err := p.db.Ping(); err != nil {
		log.Printf("repo: ping mysql returns %v", err)
		return false
	}
	return true
}
