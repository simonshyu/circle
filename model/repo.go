package model

type RepoLite struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Clone    string `json:"clone_url"`
	Branch   string `json:"default_branch"`
}

// swagger:model build
type Repo struct {
	ID          int64  `json:"id,omitempty"               meddler:"repo_id,pk"`
	ScmId       int64  `json:"-"                          meddler:"repo_scm_id"`
	Owner       string `json:"owner,omitempty"            meddler:"repo_owner"`
	FullName    string `json:"full_name"                  meddler:"repo_full_name"`
	Link        string `json:"link_url,omitempty"         meddler:"repo_link"`
	Name        string `json:"name,omitempty"             meddler:"repo_name"`
	Clone       string `json:"clone_url,omitempty"        meddler:"repo_clone"`
	Branch      string `json:"default_branch,omitempty"   meddler:"repo_branch"`
	Timeout     int64  `json:"timeout,omitempty"          meddler:"repo_timeout"`
	IsPrivate   bool   `json:"private,omitempty"          meddler:"repo_private"`
	AllowPull   bool   `json:"allow_pr"                   meddler:"repo_allow_pr"`
	AllowPush   bool   `json:"allow_push"                 meddler:"repo_allow_push"`
	AllowTag    bool   `json:"allow_tags"                 meddler:"repo_allow_tags"`
	AllowManual bool   `json:"allow_manual"               meddler:"repo_allow_manual"`
	Counter     int    `json:"last_build"                 meddler:"repo_counter"`
	Hash        string `json:"-"                          meddler:"repo_hash"`
}

/*
Field Explanation:
	Link: repo 的关联地址
	Timeout: repo 代码库拉取的超时时间
	AllowManual: 允许手动触发构建
	IsPrivate: 判断是否需要 netrc 文件

暂时不开启字段：
	Visibility  string `json:"visibility"               meddler:"repo_visibility"`

	IsTrusted   bool   `json:"trusted"                  meddler:"repo_trusted"`
*/
