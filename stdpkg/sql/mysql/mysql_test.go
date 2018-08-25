package mysql

import (
	"database/sql"
	"fmt"
	"net"
	"testing"
	"time"

	"context"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kei2100/playground-go/util/wait"
)

// prepare
// docker run -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 -d mysql:5.7
var (
	idPass   = "root:root"
	hostPort = net.JoinHostPort("0.0.0.0", "3306")
)

func TestConnMaxLifetime(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/", idPass, hostPort))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	//db.SetConnMaxLifetime(1)

	if _, err := db.Exec("SET wait_timeout=1"); err != nil {
		t.Fatal(err)
	}

	fmt.Println("sleep...")
	time.Sleep(2 * time.Second)

	// wait_timeout指定により既存の接続がタイムアウトしたときは、再作成してリトライする挙動をする。
	// [mysql] 2018/02/13 00:12:31 packets.go:33: unexpected EOFのようなstderr出力がある。
	//
	// db.SetConnMaxLifetime(1) にした場合はクライアント側でも古いコネクションを捨てるので、stderr出力はない。
	if _, err := db.Query("SELECT 1"); err != nil {
		t.Fatal(err)
	}
}

func TestUseContext(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/", idPass, hostPort))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := db.Query("SELECT sleep(3)"); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Millisecond)

	go func() {
		defer wg.Done()
		_, err := db.QueryContext(ctx, "SELECT 1")
		if err == nil {
			t.Errorf("go no error, want an error")
		}
		fmt.Printf("%T: %s", err, err)
	}()

	wait.WGroup(&wg, 10*time.Second)
}
