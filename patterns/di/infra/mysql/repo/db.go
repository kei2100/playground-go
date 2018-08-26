package repo

import (
	"database/sql"
	"fmt"
	"net"

	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/kei2100/playground-go/util/sync"
)

var db *sql.DB
var initOnce sync.OnceOrError
var closeOnce sync.OnceOrError

// DB variables
var (
	DBUser = "root"
	DBPass = "root"
	DBHost = net.JoinHostPort("localhost", "3306")
)

// Init initialize connection to the database.
func Init() error {
	return initOnce.DoOrError(func() error {
		d, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", DBUser, DBPass, DBHost))
		if err != nil {
			return fmt.Errorf("repo: failed to create to mysql connection: %v", err)
		}
		db = d
		return nil
	})
}

// Close closes the database, releasing any open resources.
func Close() error {
	return closeOnce.DoOrError(func() error {
		if db == nil {
			return nil
		}
		return db.Close()
	})
}
