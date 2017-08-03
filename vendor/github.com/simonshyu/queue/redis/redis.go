package redis

/*
如果有多个 server 实例的需求，暂定不同 server 用不同的 running queue 来保存 running 信息；
同时不同 server 用不同的 pendding queue 来保存和获取 pendding 信息
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"

	"github.com/simonshyu/queue"
)

const POP_TIMEOUT = 0 // 0 == Blocking forever

type entry struct {
	Item     *queue.Task `json:"item"`
	done     chan bool   `json:"done"`
	Retry    int         `json:"retry"`
	Error    error       `json:"error"`
	Deadline time.Time   `json:"deadline"`
}

type conn struct {
	sync.Mutex

	opts      *Options
	client    *redis.Client
	running   map[string]*entry
	extension time.Duration
}

func New(opts ...Option) (queue.Queue, error) {
	conn := new(conn)
	// 初始化 running
	conn.running = map[string]*entry{}
	conn.extension = time.Minute * 10
	// conn.extension = time.Second * 10

	conn.opts = new(Options)
	conn.opts.hostIdentity = "host01"
	conn.opts.penddingQueueName = "pendding-queue"
	for _, opt := range opts {
		// 初始化 options
		opt(conn.opts)
	}

	conn.client = redis.NewClient(&redis.Options{
		Addr:     conn.opts.addr,
		Password: conn.opts.password,
		DB:       conn.opts.db,
	})

	runningQueueKey := fmt.Sprintf("%s:running:queue", conn.opts.hostIdentity)

	allRunning, err := conn.client.HGetAll(runningQueueKey).Result()
	if err != nil {
		return nil, err
	}
	for _, taskRaw := range allRunning {
		runningTask := new(entry)
		err = json.Unmarshal([]byte(taskRaw), runningTask)
		if err != nil {
			return nil, err
		}
		runningTask.done = make(chan bool)
		conn.running[runningTask.Item.ID] = runningTask
	}
	return conn, nil
}

// Push pushes an task to the tail of this queue.
// 1. Task => Undone list(Redis)
func (c *conn) Push(ctx context.Context, task *queue.Task) error {
	taskRaw, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = c.client.LPush(c.opts.penddingQueueName, taskRaw).Err()
	if err != nil {
		return err
	}
	go c.tracking()
	return nil
}

/*
现在的问题是什么？
running list 的的任务状态保存在内存里，如果 drone server 重启，所有正在运行的任务状态将会丢失。

queue 的状态包括什么？
当前还没有被 worker 接受的 task (已保存在 redis)
当前正在运行的 task (需要保存在 redis)

agent 需要做的事情
运行完成的 task (需要保存在 redis)
运行结束时，每个 proc 的 log (需要保存在 redis)
运行结束时，每个 proc 的 log (需要保存在 redis)


如果 redis server 重启呢，会发生什么？
因为 running list 在内存里，且与 redis server 的交互目前只发生在从 pendding 队列里获取任务。redis server 重启时，可以设计成让 client 获取任务时阻塞。
而正在运行的任务，由于保存在内存里，是不会受到 redis master 掉线的影响的。

首先解决的问题是，drone server 重启的问题。
由于需要保持的状态有：
哪个任务正在运行和任务的状态信息(entry.item, entry.retry, entry.error, entry.deadline)

当 drone server 重启时，所有 agent 正在运行的任务都将会继续运行，但在处理以下行为时会出现问题：
0. agent.Next() 只涉及 queue 的接口调用
1. agent.Wait() 只涉及 queue 的接口调用, 会引起 running task not found 错误
2. agent.Extend() 只涉及 queue 的接口调用, 会引起 running task not found 错误
3. agent.Init() 只涉及数据库更新
4. agent.Upload() 只涉及数据库更新
5. agent.Update() 只涉及数据库更新
6. agent.Done() 同时涉及数据库更新和 queue 的接口调用, 因为 task 不在 running list 里什么都不会发生

期望的行为是:
0. agent.Next() 阻塞，直到 server 重连成功
1. agent.Wait() 重启之后会读取 running list 到内存里，因此会正常处理
2. agent.Extend() 重启之后会读取 running list 到内存里，因此会正常处理
3. agent.Init()
4. agent.Upload()
5. agent.Update()
6. agent.Done() 重启之后会读取 running list 到内存里，因此会正常处理

但是，由于 agent.Wait() 函数在 server 重启后不会正常返回(比如 agent 对于 break signal 就不会正常处理),
但是当任务运行结束时，对于 agent.Done() 的调用需要正常处理，因此，需要
*/

// 2. Undone list(Redis) => Task
func (c *conn) Poll(ctx context.Context, f queue.Filter) (*queue.Task, error) {
	result, err := c.client.BRPop(POP_TIMEOUT, c.opts.penddingQueueName).Result()
	if err != nil {
		return nil, err
	}
	taskRawData := result[1]

	task := new(queue.Task)
	err = json.Unmarshal([]byte(taskRawData), task)
	if err != nil {
		return nil, err
	}
	c.running[task.ID] = &entry{
		Item:     task,
		done:     make(chan bool),
		Deadline: time.Now().Add(c.extension),
	}

	taskRaw, err := json.Marshal(c.running[task.ID])
	if err != nil {
		return nil, err
	}
	runningQueueKey := fmt.Sprintf("%s:running:queue", c.opts.hostIdentity)
	runningTaskKey := fmt.Sprintf("running:task:%s", task.ID)
	err = c.client.HSet(runningQueueKey, runningTaskKey, taskRaw).Err()
	if err != nil {
		return nil, err
	}

	go c.tracking()
	return task, nil
}

// Extend extends the deadline for a task.
func (c *conn) Extend(ctx context.Context, id string) error {
	c.Lock()
	defer c.Unlock()

	task, ok := c.running[id]
	if ok {
		task.Deadline = time.Now().Add(c.extension)
		return nil
	}
	return queue.ErrNotFound
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
		task.Error = err
		close(task.done)
		delete(c.running, id)
		c.deleteTaskFromRunningQueue(id)
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
			return task.Error
		}
	}
	return nil
}

// Info returns internal queue information.
func (c *conn) Info(ctx context.Context) queue.InfoT {
	c.Lock()
	stats := queue.InfoT{}
	stats.Stats.Running = len(c.running)
	penddingLength, _ := c.client.LLen(c.opts.penddingQueueName).Result()
	stats.Stats.Pending = int(penddingLength)
	for _, entry := range c.running {
		stats.Running = append(stats.Running, entry.Item)
	}
	c.Unlock()
	return stats
}

// every call this method will checking if task.deadline is arrived.
func (c *conn) tracking() {
	c.Lock()
	defer c.Unlock()

	// TODO(bradrydzewski) move this to a helper function
	// push items to the front of the queue if the item expires.
	for id, task := range c.running {
		if time.Now().After(task.Deadline) {
			taskRaw, err := json.Marshal(task.Item)
			if err != nil {
				log.Printf("re-added to pending queue error: %v \n", err)
			}
			err = c.client.RPush(c.opts.penddingQueueName, taskRaw).Err()
			if err != nil {
				log.Printf("re-added to pending queue error: %v \n", err)
			}

			close(task.done)
			delete(c.running, id)
			c.deleteTaskFromRunningQueue(id)
		}
	}
}

func (c *conn) deleteTaskFromRunningQueue(taskID string) {
	runningQueueKey := fmt.Sprintf("%s:running:queue", c.opts.hostIdentity)
	runningTaskKey := fmt.Sprintf("running:task:%s", taskID)
	err := c.client.HDel(runningQueueKey, runningTaskKey).Err()
	if err != nil {
		log.Printf("queue: delete %s key %s error: %v\n", runningQueueKey, runningTaskKey, err)
	}
}
