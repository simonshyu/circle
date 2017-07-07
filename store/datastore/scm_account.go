package datastore

import (
	"github.com/SimonXming/circle/model"
	"github.com/russross/meddler"
)

func (db *datastore) ScmAccountCreate(account *model.ScmAcount) error {
	return meddler.Insert(db, scmAccountTable, account)
}

const scmAccountTable = "scm_account"
