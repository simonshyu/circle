package datastore

import (
	// "fmt"

	"github.com/SimonXming/circle/model"
	"github.com/russross/meddler"
)

func (db *datastore) RepoCreate(repo *model.Repo) error {
	return meddler.Insert(db, repoTable, repo)
}

const repoTable = "repos"
