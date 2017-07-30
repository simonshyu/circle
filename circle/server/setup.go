// +build !enterprise

package server

import (
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/store/datastore"
	"github.com/SimonXming/queue"
	"github.com/SimonXming/queue/gcp"

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

func setupGcpQueue(c *cli.Context, s store.Store) queue.Queue {
	project := "speedy-crane-157603"
	topic := "test"
	topicDone := "done"
	subscription := "reader"
	subscriptionDone := "done"
	tokenpath := "/Users/simon/Code/go/src/github.com/SimonXming/circle/google-platform-example-8e5a3d25c36f.json"
	q, err := gcp.New(
		gcp.WithProject(project),
		gcp.WithTopic(topic, topicDone),
		gcp.WithSubscription(subscription, subscriptionDone),
		gcp.WithServiceAccountToken(tokenpath),
	)
	if err != nil {
		println(err.Error())
	}
	return model.WithTaskStore(q, s)
}
