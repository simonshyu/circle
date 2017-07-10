package datastore

import (
	// "fmt"

	"github.com/SimonXming/circle/model"
	"github.com/russross/meddler"
)

func (db *datastore) RepoCreate(repo *model.Repo) error {
	return meddler.Insert(db, repoTable, repo)
}

func (db *datastore) GetRepoName(name string) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.QueryRow(db, repo, rebind(repoNameQuery), name)
	return repo, err
}

const repoTable = "repos"

const repoNameQuery = `
SELECT *
FROM repos
WHERE repo_full_name = ?
LIMIT 1;
`
