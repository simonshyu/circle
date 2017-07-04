package datastore

import (
	// "fmt"
	"time"

	"github.com/SimonXming/circle/model"
	// "github.com/drone/drone/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) BuildCreate(build *model.Build) error {
	// TODO: get build number
	id := 1
	build.Number = id
	build.Created = time.Now().UTC().Unix()
	err := meddler.Insert(db, buildTable, build)
	if err != nil {
		return err
	}
	return nil
}

const buildTable = "builds"
