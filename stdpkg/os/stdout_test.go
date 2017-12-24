package os

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func ExampleStdoutPrintln() {
	fmt.Println("test")

	// Output:
	// test
}

func ExampleStdoutUseStdout() {
	os.Stdout.Write([]byte("test"))

	// Output:
	// test
}

func ExampleStdoutUseFD() {
	o := os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	if o == nil {
		log.Fatal("file is nil")
	}
	o.Write([]byte("test"))

	// "test"は出力されるが、Output: では補足できない。

	// Output:
	//
}
