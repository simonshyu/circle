package datastore

import (
	"bytes"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store/datastore/sql"
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

const logTable = "logs"
