package datastore

import (
	// gosql "database/sql"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) ConfigCreate(config *model.Config) error {
	return meddler.Insert(db, "config", config)
}

func (db *datastore) ConfigLoad(id int64) (*model.Config, error) {
	stmt := sql.Lookup(db.driver, "config-find-id")
	conf := new(model.Config)
	err := meddler.QueryRow(db, conf, stmt, id)
	return conf, err
}
