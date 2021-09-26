package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kei2100/playground-go/stdpkg/os/exec/rename_and_read_multi_proc/rrmp2/filelib"
)

var (
	logPath        = filepath.Join(".", "tmp", "test.log")
	rotatedLogPath = filepath.Join(".", "tmp", "test.1.log")
)

func main() {
	if os.Getenv("IS_CHILD") == "1" {
		child()
		return
	}

	os.Mkdir(filepath.Dir(logPath), 0755)
	os.Remove(logPath)
	os.Remove(rotatedLogPath)
	ioutil.WriteFile(logPath, []byte("test"), 0644)

	parent()

	// Windows
	// 2019/05/01 10:46:26 child: logfile opened. tmp\test.log
	// 2019/05/01 10:46:26 child: entries are...
	// 2019/05/01 10:46:26 child: tmp\test.log
	// 2019/05/01 10:46:27 parent: logfile renamed.
	// 2019/05/01 10:46:27 parent: entries are...
	// 2019/05/01 10:46:27 parent: tmp\test.1.log
	// 2019/05/01 10:46:29 child: logfile readed. tmp\test.log
	// 2019/05/01 10:46:29 child: content: testlog, err:<nil>
	// 2019/05/01 10:46:29 child: entries are...
	// 2019/05/01 10:46:29 child: tmp\test.1.log
	// 2019/05/01 10:46:29 child: logfile closed

	// OSX
	// 2019/05/01 10:44:45 child: logfile opened. tmp/test.log
	// 2019/05/01 10:44:45 child: entries are...
	// 2019/05/01 10:44:45 child: tmp/test.log
	// 2019/05/01 10:44:46 parent: logfile renamed.
	// 2019/05/01 10:44:46 parent: entries are...
	// 2019/05/01 10:44:46 parent: tmp/test.1.log
	// 2019/05/01 10:44:48 child: logfile readed. tmp/test.log
	// 2019/05/01 10:44:48 child: content: testlog, err:<nil>
	// 2019/05/01 10:44:48 child: entries are...
	// 2019/05/01 10:44:48 child: tmp/test.1.log
	// 2019/05/01 10:44:48 child: logfile closed
}

func child() {
	f, err := filelib.OpenLogFileToRead(logPath, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
		log.Println("child: logfile closed")
	}()

	log.Printf("child: logfile opened. %v", f.Name())
	lsLogDir("child")

	time.Sleep(3 * time.Second)

	content, err := ioutil.ReadAll(f)
	log.Printf("child: logfile readed. %v", f.Name())
	log.Printf("child: content: %v, err:%v", string(content), err)
	lsLogDir("child")
}

func parent() {
	com, err := os.Executable()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(com)
	cmd.Env = append(os.Environ(), "IS_CHILD=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()

	time.Sleep(time.Second)

	f.WriteString("log")
	f.Sync()
	f.Close()
	os.Rename(logPath, rotatedLogPath)

	log.Println("parent: logfile renamed.")
	lsLogDir("parent")

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
}

func lsLogDir(parentOrChild string) {
	m, err := filepath.Glob(filepath.Dir(logPath) + "/*")
	if err != nil {
		log.Printf("%v: glob error: %v", parentOrChild, err)
	}
	log.Printf("%v: entries are...", parentOrChild)
	for _, e := range m {
		log.Printf("%v: %v", parentOrChild, e)
	}
}
