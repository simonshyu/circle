// +build !enterprise

package server

import (
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/store/datastore"
	"github.com/SimonXming/queue"

	"github.com/urfave/cli"
)

func setupStore(c *cli.Context) store.Store {
	println(c.String("datasource"))
	return datastore.New(
		c.String("driver"),
		c.String("datasource"),
	)
}

func setupQueue(c *cli.Context, s store.Store) queue.Queue {
	return model.WithTaskStore(queue.New(), s)
}
