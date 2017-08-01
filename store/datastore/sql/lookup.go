package sql

import (
	"github.com/simonshyu/circle/store/datastore/sql/mysql"
)

// Supported database drivers
const (
	DriverMysql = "mysql"
)

// Lookup returns the named sql statement compatible with
// the specified database driver.
func Lookup(driver string, name string) string {
	switch driver {
	case DriverMysql:
		return mysql.Lookup(name)
	default:
		return mysql.Lookup(name)
	}
}
