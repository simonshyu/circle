## 定义 agent 对于 server jsonrpc2.0 交互接口

```go
type Peer interface {
	// Next returns the next pipeline in the queue.
	Next(c context.Context) (*Pipeline, error)

	// Wait blocks until the pipeline is complete.
	Wait(c context.Context, id string) error

	// Init signals the pipeline is initialized.
	Init(c context.Context, id string, state State) error

	// Done signals the pipeline is complete.
	Done(c context.Context, id string, state State) error

	// (Not sure)Extend extends the pipeline deadline
	// Extend(c context.Context, id string) error

	// Update updates the pipeline state.
	Update(c context.Context, id string, state State) error

	// Upload uploads the pipeline artifact.
	Upload(c context.Context, id string, file *File) error

	// Log writes the pipeline log entry.
	Log(c context.Context, id string, line *Line) error
}
```

## 关于需要落库的数据

### store 目录定义与数据库的交互

* [] 具体数据库操作接口的实现
* [] 数据库 sql
* [] 支持多个数据库 backend

### model 目录定义与数据库表结构

* [] 定义数据库表结构

## 工作流程描述

* server 定义与 agent 交互的 基于 jsonrpc2 的 websocket 接口
    * agent 获取 job
    * 解析 job
    * 执行 job 
    * 更新 job 数据/状态
* server 定义供用户读取显示的数据接口
    * Config job pipeline 配置
    * Repo 代码仓库
    * Secret 密码或token
    * Build 代表一次构建行为
    * Task 代表 scheduled pipeline Task.(Task 被发送到 )
	* Logs 代表一次构建的日志
    * File 一次 Build 的 pipeline artifact.
	* Agents 代表构建节点信息
    * Proc 代表 build 发生的进程信息
	* 未完...
* 存在一个 Compiler 概念
	* 用于解析 yaml 文件
	* 经过 compiler 的解析，会生成一个 backend.Config 对象
	* backend.Config 包含 stages, network, volume, secret。更完整，抽象程度更高。
	* pipeline.frontend.yaml.Config => pipeline.backend.Config


## 待解决问题

1. 构建出的程序包如何上传到文件服务器？
	* 可以通过服务端提供提供一条上传shell命令模版，让用户在 shell 中调用。实现上传文件的功能。
2. 如何支持上传至 registry 时用户密码验证的问题？
	* agent 可以从 server 获取相应的 Secret 验证信息。
3. 日志的保存问题？
	* 构建过程中，通过 websocket 获取一个 log stream 获取日志。
	* 构建结束时，将日志上传至数据库。供以后读取。
4. 多节点共享拉取的代码库的问题？
	* 需要。
	* drone 是通过定义 docker volume，并且每个 step 共享同一个 volume 来实现的。
	* 根据 drone/agent/agent.go 可知，一个 agent 负责执行一次 build 的全部 stage.
	* http://docs.drone.io/workspace/
5. Event 的概念是什么用途？
	* Event 没有落库，看上去只是为了方便 log, 看来只是用在了 pubsub.Message 里。 
	* model.EventPush/model.EventPull/model.EventTag 等等是针对 webhook 所携带的 git 操作类型, 对于 build.build_event 字段

