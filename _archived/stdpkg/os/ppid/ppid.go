package main

import (
	"log"
	"os"
	"os/exec"
)

// go build ppid.go
//
//	2018/11/02 17:44:29 parent: my pid is 45923
//	2018/11/02 17:44:29 child: my pid is 45924
//	2018/11/02 17:44:29 child: parent pid is 45923
func main() {
	if len(os.Args) < 2 {
		p, err := os.Executable()
		if err != nil {
			log.Fatalf("failed to get executable: %v", err)
		}
		cmd := exec.Command(p, "-child")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Start(); err != nil {
			log.Fatalf("failed to start a child process: %v", err)
		}

		log.Printf("parent: my pid is %d", os.Getpid())

		if err := cmd.Wait(); err != nil {
			log.Fatalf("failed to wait for child process: %v", err)
		}
		return
	}

	log.Printf("child: my pid is %d", os.Getpid())
	log.Printf("child: parent pid is %d", os.Getppid())
}
