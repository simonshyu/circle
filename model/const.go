package model

const (
	RepoGit = "git"
)

const (
	EventPush   = "push"
	EventPull   = "pull_request"
	EventTag    = "tag"
	EventManual = "manual"
	// EventDeploy = "deployment"
)

const (
	StatusSkipped  = "skipped"
	StatusPending  = "pending"
	StatusRunning  = "running"
	StatusSuccess  = "success"
	StatusFailure  = "failure"
	StatusKilled   = "killed"
	StatusError    = "error"
	StatusBlocked  = "blocked"
	StatusDeclined = "declined"
)

// example: gitlab project have these three project type.
const (
	VisibilityPublic   = "public"
	VisibilityPrivate  = "private"
	VisibilityInternal = "internal"
)
