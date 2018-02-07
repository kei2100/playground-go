package mysql

import (
	"database/sql"
	"fmt"
	"net"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnMaxLifetime(t *testing.T) {
	t.Skip("skip for automated test")

	// prepare
	// docker run -e MYSQL_ROOT_PASSWORD=pass -p 13306:3306 -d mysql:5.7

	ip := "root:pass"
	hp := net.JoinHostPort("0.0.0.0", "13306")
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/", ip, hp))
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

	if _, err := db.Query("SELECT 1"); err != nil {
		t.Fatal(err)
	}
}
