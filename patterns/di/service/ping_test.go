package service

import (
	"log"
	"testing"

	imrepo "github.com/kei2100/playground-go/patterns/di/infra/inmem/repo"
	mrepo "github.com/kei2100/playground-go/patterns/di/infra/mysql/repo"
	"github.com/kei2100/playground-go/patterns/di/repo"
)

func TestPing_mysql(t *testing.T) {
	if err := mrepo.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := mrepo.Close(); err != nil {
			log.Printf("repo.Close returns an error: %v", err)
		}
	}()

	repo.Ping = mrepo.Ping()
	if ok := Ping(); !ok {
		t.Error("Ping test failure")
	}
}

func TestPing_inmem(t *testing.T) {
	repo.Ping = imrepo.Ping()
	if ok := Ping(); !ok {
		t.Error("Ping test failure")
	}
}
