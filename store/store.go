package store

import (
	// "io"

	"github.com/SimonXming/circle/model"
	// "golang.org/x/net/context"
)

type Store interface {
	ConfigCreate(*model.Config) error
	ConfigLoad(int64) (*model.Config, error)

	RepoCreate(*model.Repo) error

	SecretCreate(*model.Secret) error

	BuildCreate(*model.Build) error
}
