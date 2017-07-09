package datastore

import (
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) ScmAccountCreate(account *model.ScmAccount) error {
	return meddler.Insert(db, scmAccountTable, account)
}

func (db *datastore) ScmAccountList() ([]*model.ScmAccount, error) {
	stmt := sql.Lookup(db.driver, "scm-account-list")
	data := []*model.ScmAccount{}
	err := meddler.QueryAll(db, &data, stmt)
	return data, err
}

func (db *datastore) ScmAccountLoad(id int64) (*model.ScmAccount, error) {
	stmt := sql.Lookup(db.driver, "scm-account-find-id")
	account := new(model.ScmAccount)
	err := meddler.QueryRow(db, account, stmt, id)
	return account, err
}

const scmAccountTable = "scm_account"
