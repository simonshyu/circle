package model

// swagger:model build
type Repo struct {
	ID        int64  `json:"id,omitempty"             meddler:"repo_id,pk"`
	ScmId     int64  `json:"-"                        meddler:"repo_scm_id"`
	Owner     string `json:"owner,omitempty"          meddler:"repo_owner"`
	FullName  string `json:"full_name"                meddler:"repo_full_name"`
	Name      string `json:"name,omitempty"           meddler:"repo_name"`
	Clone     string `json:"clone_url,omitempty"      meddler:"repo_clone"`
	Branch    string `json:"default_branch,omitempty" meddler:"repo_branch"`
	AllowPush bool   `json:"allow_push"               meddler:"repo_allow_push"`
	Hash      string `json:"-"                        meddler:"repo_hash"`
}
