package datastore

import (
	"github.com/russross/meddler"
	"github.com/simonshyu/circle/model"
	// "github.com/simonshyu/circle/store/datastore/sql"
)

func (db *datastore) RegistryCreate(registry *model.Registry) error {
	return meddler.Insert(db, registryTable, registry)
}

func (db *datastore) RegistryUpdate(registry *model.Registry) error {
	return meddler.Update(db, registryTable, registry)
}

const registryTable = "registry"
