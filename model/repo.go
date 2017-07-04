package model

// swagger:model build
type Repo struct {
	ID       int64  `json:"id,omitempty"             meddler:"repo_id,pk"`
	Kind     string `json:"scm,omitempty"            meddler:"repo_scm"`
	Host     string `json:"host,omitempty"           meddler:"repo_host"`
	Clone    string `json:"clone_url,omitempty"      meddler:"repo_clone"`
	Branch   string `json:"default_branch,omitempty" meddler:"repo_branch"`
	Login    string `json:"login,omitempty"          meddler:"repo_login"`
	Password string `json:"password,omitempty"       meddler:"repo_password"`
}
