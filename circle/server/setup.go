// +build !enterprise

package server

import (
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/store/datastore"

	"github.com/urfave/cli"
)

func setupStore(c *cli.Context) store.Store {
	println(c.String("datasource"))
	return datastore.New(
		c.String("driver"),
		c.String("datasource"),
	)
}
