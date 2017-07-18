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
6. 拉取 git/svn 代码使用 golang 的哪个库？
	* 使用 plugin/git:latest 这个镜像

## 执行时的代码逻辑
0. 入口 pipeline.New(backend.Config).Run()
1. `pipeline.go` runtime.execAll(backend.Config.Stages.Steps)
2. `pipeline.go` runtime.exec(Steps.step)
3. `pipeline.go` runtime.engine.Exec(step) => dockerEngine.Exec(step)
4. `docker/convert.go` toConfig() step.Config => docker.container.Config
5. `docker/docker.go` ImagePull, ContainerCreate, ContainerStart

## git pull 的逻辑
1. step 定义了拉取代码所需的环境变量 WORKSPACE REMOTE_URL EMAIL PASSWORD 等
2. 定义一个负责拉取代码的镜像, 镜像运行需要以上环境变量
3. WORKSPACE 是一个可以与其他镜像共享的卷

## 表设计

* 核心表是 config
* config 应该关联了一个代码库 repo(可以多个 config 应该关联了同一个代码库)
* repo 应该关联了一个 secret
* 如果要实现和 gitlab 服务器的交互，需要一个 remote 和 remote.GitlabClient
* remote.GitlabClient 需要 host 和 secret 来与 gitlab 交互

场景一：
一个 dcos 用户登录，提供了若干个 scm 服务地址和用户名密码，scm 是对不同类型代码管理工具的抽象。
用户查看 scm 服务地址对应的代码库列表，选择其中的某个代码库，记录并转化为用于构建的 repo。
转化为 repo 的过程中，生成一条用于鉴权的 secret（默认 secret 类型为用户名密码，也可以新增其他类型的 secret）。
针对选中的每个 repo 建立一个构建相关的 config. (注意：config 包含了构建行为针对的branch，如果 hook 携带的 branch 的被 exclude，不会触发构建)

hook 触发的构建携带了需要 pull 的代码分支(tag,commit)，而用户点击构建时，需要 pull 的是上次构建的代码或者 default_branch

## circle 与 scm server 交互的接口

```golang

type Remote interface {
	// Login authenticates the session and returns the
	// remote user details.
	Login(w http.ResponseWriter, r *http.Request) (*model.User, error)

	// Auth authenticates the session and returns the remote user
	// login for the given token and secret
	Auth(token, secret string) (string, error)

	// Teams fetches a list of team memberships from the remote system.
	Teams(u *model.User) ([]*model.Team, error)

	// TeamPerm fetches the named organization permissions from
	// the remote system for the specified user.
	TeamPerm(u *model.User, org string) (*model.Perm, error)

	// Repo fetches the named repository from the remote system.
	Repo(u *model.User, owner, repo string) (*model.Repo, error)

	// Repos fetches a list of repos from the remote system.
	Repos(u *model.User) ([]*model.RepoLite, error)

	// Perm fetches the named repository permissions from
	// the remote system for the specified user.
	Perm(u *model.User, owner, repo string) (*model.Perm, error)

	// File fetches a file from the remote repository and returns in string
	// format.
	File(u *model.User, r *model.Repo, b *model.Build, f string) ([]byte, error)

	// FileRef fetches a file from the remote repository for the given ref
	// and returns in string format.
	FileRef(u *model.User, r *model.Repo, ref, f string) ([]byte, error)

	// Status sends the commit status to the remote system.
	// An example would be the GitHub pull request status.
	Status(u *model.User, r *model.Repo, b *model.Build, link string) error

	// Netrc returns a .netrc file that can be used to clone
	// private repositories from a remote system.
	Netrc(u *model.User, r *model.Repo) (*model.Netrc, error)

	// Activate activates a repository by creating the post-commit hook.
	Activate(u *model.User, r *model.Repo, link string) error

	// Deactivate deactivates a repository by removing all previously created
	// post-commit hooks matching the given link.
	Deactivate(u *model.User, r *model.Repo, link string) error

	// Hook parses the post-commit hook from the Request body and returns the
	// required data in a standard format.
	Hook(r *http.Request) (*model.Repo, *model.Build, error)
}

```

## 与 webhook 相关的问题

### 现状

* drone: webhook 事件关联哪个 commit 就去哪个 commit 取 .drone.yml 文件内容。所以，drone 可以支持任意版本。

* circle: 因为 pipeline.yml 文件不存放在代码仓库。所以，只能维持单个(或若干个)版本的 pipeline.yml，每个版本分别对应一些分支限制条件，只有 webhook 发生时的信息满足限制条件，才会触发这个 webhook 对应 commit id 的构建。pipeline.yml 里定义了关心的分支信息。

限制条件到这里只有：分支(pull, merge_request 等和 branch 相关的事件)


当任意分支发生某一类型事件时，与此事件绑定的 `hook_url` 将会被调用(携带代码库信息，分支信息和事件信息)。假设有如下 `hook_url`
```json
[
	"http://127.0.0.1:8000/scm-1/repo-1/hook",         # push，merge request
	"http://127.0.0.1:8000/scm-2/repo-1/hook",         # tag
]
```
1. 假设 master 分支发生 push 事件，`http://127.0.0.1:8000/scm-1/repo-1/hook` 触发。
服务器根据 scmID 和 repo 信息找到 repoID。因为 repo.allow_push，所以找到 repo 唯一对应的 config-1。假设 config 里定义了只有涉及到 develop 的变动才会构建。所以这次构建 skip。

2. 假设 master 分支发生 tag_push 事件，`http://127.0.0.1:8000/scm-2/repo-1/hook` 触发。
服务器根据 scmID 和 repo 信息找到 repoID。因为 repo.allow_tag。所以 `repo.config_list` 中每一个 config 都构建。


## 与 build 相关的问题

### drone:

* 默认第一次构建必须由 webhook 触发
* 用户手动触发的情况包括：
	* 复制某一次构建，并创建一次新的构建
	* 重新运行某一次构建

### circle:
* 允许 webhook 触发
* 用户手动触发的情况包括：
	* 新增构建类型 manual (意为: 用户手动触发 repo 默认分支 refs/heads/{DEFAULT_BRANCH} 的构建)
	* fork 一次构建
	* re-run 一次构建
