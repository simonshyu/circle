package datastore

import (
	"bytes"
	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/store/datastore/sql"
	"github.com/russross/meddler"
	"io"
	"io/ioutil"
)

func (db *datastore) LogFind(proc *model.Proc) (io.ReadCloser, error) {
	stmt := sql.Lookup(db.driver, "log-find-job-id")
	log := new(model.LogData)
	err := meddler.QueryRow(db, log, stmt, proc.ID)
	buf := bytes.NewBuffer(log.Data)
	return ioutil.NopCloser(buf), err
}

func (db *datastore) LogSave(proc *model.Proc, r io.Reader) error {
	stmt := sql.Lookup(db.driver, "log-find-job-id")
	log := new(model.LogData)
	err := meddler.QueryRow(db, log, stmt, proc.ID)
	if err != nil {
		log = &model.LogData{ProcID: proc.ID}
	}
	log.Data, _ = ioutil.ReadAll(r)
	return meddler.Save(db, logTable, log)
}

const logTable = "logs"
