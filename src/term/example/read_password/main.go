package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kei2100/playground-go/src/term"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	fmt.Printf("Please input password: ")
	// ReadPasswordContext を使うことで、パスワード入力待ち中に Ctrl+C などのシグナルでキャンセルできるようになる
	p, err := term.ReadPasswordContext(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("\n")
	fmt.Printf("Your input is %s\n", p)
	return nil
}
