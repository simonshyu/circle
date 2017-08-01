// +build !enterprise

package server

import (
	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/store"
	"github.com/simonshyu/circle/store/datastore"
	"github.com/simonshyu/queue"
	"github.com/simonshyu/queue/redis"

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

func setupRedisQueue(c *cli.Context, s store.Store) queue.Queue {
	redisConn, err := redis.New()
	if err != nil {
		return nil
	}
	return model.WithTaskStore(redisConn, s)
}
