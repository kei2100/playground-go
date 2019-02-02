package main

import (
	"log"
	"os"
	"os/exec"

	"bitbucket.org/avd/go-ipc/fifo"
)

func _main() {
	testData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	if os.Getenv("_WORKER") == "on" {
		log.Printf("pid %v: start worker", os.Getpid())

		wfifo, err := fifo.New("fifo", os.O_CREATE|os.O_WRONLY, 0666)
		panicIf(err)
		defer wfifo.Close()
		if written, err := wfifo.Write(testData); err != nil || written != len(testData) {
			panic("write")
		}
		return
	}

	log.Printf("pid %v: start master", os.Getpid())
	launchWorker()

	buff := make([]byte, len(testData))
	rfifo, err := fifo.New("fifo", os.O_CREATE|os.O_RDONLY, 0666)
	panicIf(err)
	defer rfifo.Close()

	if read, err := rfifo.Read(buff); err != nil || read != len(testData) {
		panic("read")
	}
	// ensure we've received valid data
	for i, b := range buff {
		println(b)
		if b != testData[i] {
			panic("wrong data")
		}
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func launchWorker() (*os.Process, error) {
	// launch worker
	bin, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(bin)
	cmd.Env = append(cmd.Env, "_WORKER=on")
	cmd.Env = append(cmd.Env, "TMPDIR="+os.TempDir())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd.Process, nil
}
