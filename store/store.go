package store

import (
	"io"

	"github.com/SimonXming/circle/model"

	"golang.org/x/net/context"
)

type Store interface {
	ConfigCreate(*model.Config) error
}
