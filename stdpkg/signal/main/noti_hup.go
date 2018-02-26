package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// just a test for 「nohupで実行しても、SIGHUPのハンドラを登録すればハンドラは実行される」
//
// (prepare)  $ go build -o bin/noti_hup stdpkg/signal/main/noti_hup.go
// # go run だとgoプロセスがHUPを受けて終了させてしまうので実行バイナリにしておく必要がある
//
// (patternA) $ nohup bin/noti_hup
// (patternB) $ nohup bin/noti_hup --hup=false
// (patternC) $ bin/noti_hup --hup=false
//
// (test)     $ kill -s USR1 $(pgrep -f noti_hup); kill -s HUP $(pgrep -f noti_hup)
//
// (resultA) 「received: usr1」,「received: hup」が表示される。プロセスは終了しない。
// (resultB) 「received: usr1」が表示される。プロセスは終了しない。
// (resultC) 「received: usr1」が表示される。プロセスは終了する。
//
func main() {
	hup := flag.Bool("hup", true, "handle HUP")
	usr1 := flag.Bool("usr1", true, "handle USR1")
	flag.Parse()

	handle := make([]os.Signal, 0)
	if *hup {
		handle = append(handle, syscall.SIGHUP)
	}
	if *usr1 {
		handle = append(handle, syscall.SIGUSR1)
	}

	c := make(chan os.Signal)
	signal.Notify(c, handle...)

	for {
		select {
		case s := <-c:
			log.Printf("received: %v", s.String())
		}
	}
}
