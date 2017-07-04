package model

// swagger:model build
type Build struct {
	ID       int64 `json:"id"            meddler:"build_id,pk"`
	RepoID   int64 `json:"-"             meddler:"build_repo_id"`
	ConfigID int64 `json:"-"             meddler:"build_config_id"`
	Number   int   `json:"number"        meddler:"build_number"`
	Created  int64 `json:"created_at"    meddler:"build_created"`
	Started  int64 `json:"started_at"    meddler:"build_started"`
	Finished int64 `json:"finished_at"   meddler:"build_finished"`
}
