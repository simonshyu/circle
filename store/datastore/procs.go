package datastore

import (
	// "fmt"

	"github.com/SimonXming/circle/model"
	// "github.com/SimonXming/circle/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) ProcCreate(procs []*model.Proc) error {
	for _, proc := range procs {
		if err := meddler.Insert(db, procsTable, proc); err != nil {
			return err
		}
	}
	return nil
}

const procsTable = "procs"
