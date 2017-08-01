package datastore

import (
	// "fmt"

	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/store/datastore/sql"
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

func (db *datastore) ProcList(build *model.Build) ([]*model.Proc, error) {
	stmt := sql.Lookup(db.driver, "proc-find-build")
	data := []*model.Proc{}
	err := meddler.QueryAll(db, &data, stmt, build.ID)
	return data, err
}

func (db *datastore) ProcLoad(id int64) (*model.Proc, error) {
	stmt := sql.Lookup(db.driver, "proc-find-id")
	proc := new(model.Proc)
	err := meddler.QueryRow(db, proc, stmt, id)
	return proc, err
}

func (db *datastore) ProcChild(build *model.Build, pid int, child string) (*model.Proc, error) {
	stmt := sql.Lookup(db.driver, "proc-find-build-ppid")
	proc := new(model.Proc)
	err := meddler.QueryRow(db, proc, stmt, build.ID, pid, child)
	return proc, err
}

func (db *datastore) ProcUpdate(proc *model.Proc) error {
	return meddler.Update(db, procsTable, proc)
}

func (db *datastore) ProcClear(build *model.Build) error {
	//stmt1 := sql.Lookup(db.driver, "files-delete-build")
	stmt := sql.Lookup(db.driver, "proc-delete-build")
	_, err := db.Exec(stmt, build.ID)
	return err
}

const procsTable = "procs"
