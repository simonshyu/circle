package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"sync"
	"time"

	"github.com/simonshyu/queue"
)

type entry struct {
	item     *queue.Task
	done     chan bool
	retry    int
	error    error
	deadline time.Time
}

type conn struct {
	sync.Mutex

	// opts   *Options
	client    *redis.Client
	running   map[string]*entry
	extension time.Duration
}

func New() (queue.Queue, error) {
	conn := new(conn)
	// 初始化 running
	conn.running = map[string]*entry{}
	conn.extension = time.Minute * 10
	conn.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return conn, nil
}

// Push pushes an task to the tail of this queue.
// 1. Task => Undone list(Redis)
func (c *conn) Push(ctx context.Context, task *queue.Task) error {
	out, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = c.client.LPush("channel-undone", out).Err()
	if err != nil {
		return err
	}
	return nil
}

// 2. Undone list(Redis) => Task
func (c *conn) Poll(ctx context.Context, f queue.Filter) (*queue.Task, error) {
	result, err := c.client.BLPop(0, "channel-undone").Result()
	if err != nil {
		return nil, err
	}

	task := new(queue.Task)
	err = json.Unmarshal([]byte(result[1]), task)
	if err != nil {
		return nil, err
	}
	c.running[task.ID] = &entry{
		item:     task,
		done:     make(chan bool),
		deadline: time.Now().Add(c.extension),
	}

	return task, nil
}

// Extend extends the deadline for a task.
func (c *conn) Extend(ctx context.Context, id string) error {
	println("Try extend deadline %v", 600)
	return nil
}

// Done signals the task is complete.
func (c *conn) Done(ctx context.Context, id string) error {
	return c.Error(ctx, id, nil)
}

// Error signals the task is complete with errors.
func (c *conn) Error(ctx context.Context, id string, err error) error {
	c.Lock()
	task, ok := c.running[id]
	if ok {
		task.error = err
		close(task.done)
		delete(c.running, id)
	}
	c.Unlock()
	return nil
}

// Evict removes a pending task from the queue.
func (c *conn) Evict(ctx context.Context, id string) error {
	return nil
}

// Wait waits until the task is complete.
// 3. Return error when task is done
func (c *conn) Wait(ctx context.Context, id string) error {
	c.Lock()
	task, ok := c.running[id]
	c.Unlock()
	if ok {
		select {
		case <-ctx.Done():
		case <-task.done:
			return task.error
		}
	}
	return nil
}

// Info returns internal queue information.
func (c *conn) Info(ctx context.Context) queue.InfoT {
	// TODO this will be different for gcp
	return queue.InfoT{}
}
