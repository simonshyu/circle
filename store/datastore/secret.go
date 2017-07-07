package datastore

import (
	"github.com/SimonXming/circle/model"
	"github.com/russross/meddler"
)

func (db *datastore) SecretCreate(secret *model.Secret) error {
	return meddler.Insert(db, secretTable, secret)
}

const secretTable = "secrets"