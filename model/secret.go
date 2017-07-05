package model

// Secret represents a secret variable, such as a password or token.
// swagger:model registry
type Secret struct {
	ID     int64  `json:"id"              meddler:"secret_id,pk"`
	RepoID int64  `json:"-"               meddler:"secret_repo_id"`
	Name   string `json:"name"            meddler:"secret_name"`
	Value  string `json:"value,omitempty" meddler:"secret_value"`
}
