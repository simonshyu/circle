package model

// swagger:model build
type Build struct {
	ID       int64   `json:"id"            meddler:"build_id,pk"`
	RepoID   int64   `json:"-"             meddler:"build_repo_id"`
	ConfigID int64   `json:"-"             meddler:"build_config_id"`
	Number   int     `json:"number"        meddler:"build_number"`
	Event    string  `json:"event"         meddler:"build_event"`
	Status   string  `json:"status"        meddler:"build_status"`
	Error    string  `json:"error"         meddler:"build_error"`
	Created  int64   `json:"created_at"    meddler:"build_created"`
	Started  int64   `json:"started_at"    meddler:"build_started"`
	Enqueued int64   `json:"enqueued_at"   meddler:"build_enqueued"`
	Finished int64   `json:"finished_at"   meddler:"build_finished"`
	Commit   string  `json:"commit"        meddler:"build_commit"`
	Branch   string  `json:"branch"        meddler:"build_branch"`
	Ref      string  `json:"ref"           meddler:"build_ref"`
	Refspec  string  `json:"refspec"       meddler:"build_refspec"`
	Remote   string  `json:"remote"        meddler:"build_remote"`
	Procs    []*Proc `json:"procs,omitempty" meddler:"-"`
}
