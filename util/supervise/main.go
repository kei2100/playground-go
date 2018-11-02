package main

import (
	"log"
	"os"
	"time"

	"github.com/kei2100/playground-go/util/supervise/supervisor"
)

func main() {
	if len(os.Args) < 2 {
		p, err := os.Executable()
		if err != nil {
			log.Fatalf("failed to get executable: %v", err)
		}
		sv := supervisor.NewSupervisor(p, "--child")
		if err := sv.Start(); err != nil {
			log.Println(err)
		}
		return
	}

	// child

	for {
		time.Sleep(time.Second)
	}
}
