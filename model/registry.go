package model

import "errors"

var (
	errRegistryAddressInvalid  = errors.New("Invalid Registry Address")
	errRegistryUsernameInvalid = errors.New("Invalid Registry Username")
	errRegistryPasswordInvalid = errors.New("Invalid Registry Password")
)

// Registry represents a docker registry with credentials.
// swagger:model registry
type Registry struct {
	ID       int64  `json:"id"       meddler:"registry_id,pk"`
	RepoID   int64  `json:"-"        meddler:"registry_repo_id"`
	Address  string `json:"address"  meddler:"registry_addr"`
	Username string `json:"username" meddler:"registry_username"`
	Password string `json:"password" meddler:"registry_password"`
	Email    string `json:"email"    meddler:"registry_email"`
	Token    string `json:"token"    meddler:"registry_token"`
}
