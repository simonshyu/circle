package datastore

import (
	"github.com/russross/meddler"
	"github.com/simonshyu/circle/model"
)

func (db *datastore) SecretCreate(secret *model.Secret) error {
	return meddler.Insert(db, secretTable, secret)
}

func (db *datastore) SecretUpdate(secret *model.Secret) error {
	return meddler.Update(db, secretTable, secret)
}

const secretTable = "secrets"
