package sql

import (
	"github.com/SimonXming/circle/store/datastore/sql/sqlite"
)

// Supported database drivers
const (
	DriverSqlite = "sqlite3"
	DriverMysql  = "mysql"
)

// Lookup returns the named sql statement compatible with
// the specified database driver.
func Lookup(driver string, name string) string {
	switch driver {
	case DriverMysql:
		return sqlite.Lookup(name)
	default:
		return sqlite.Lookup(name)
	}
}
