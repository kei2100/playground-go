package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	logPath        = filepath.Join(".", "test.log")
	rotatedLogPath = filepath.Join(".", "test.1.log")
)

func main() {
	panicIf(os.Remove(logPath))
	panicIf(os.Remove(rotatedLogPath))
	panicIf(ioutil.WriteFile(logPath, []byte{}, 0644))

	if os.Getenv("IS_CHILD") == "1" {
		child()
		return
	}
	parent()
}

func child() {
	lg, err := os.OpenFile(logPath, os.O_RDONLY, 0)
	panicIf(err)
	defer lg.Close()

	time.Sleep(3 * time.Second)

	eb, err := ioutil.ReadAll(lg)
	panicIf(err)
	fmt.Println(string(eb))
}

func parent() {
	com, err := os.Executable()
	panicIf(err)
	cmd := exec.Command(com)
	cmd.Env = append(os.Environ(), "IS_CHILD=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	panicIf(cmd.Start())

	time.Sleep(time.Second)

	lg, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	panicIf(err)

	lg.WriteString("log")
	panicIf(lg.Close())
	// windowsだと「The process cannot access the file because it is being used by another process.」
	panicIf(os.Rename(logPath, rotatedLogPath))

	panicIf(cmd.Wait())
	return
}

func panicIf(err error) {
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}
