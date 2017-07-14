package datastore

import (
	// "fmt"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) RepoCreate(repo *model.Repo) error {
	return meddler.Insert(db, repoTable, repo)
}

func (db *datastore) RepoList() ([]*model.Repo, error) {
	stmt := sql.Lookup(db.driver, "repo-list")
	data := []*model.Repo{}
	err := meddler.QueryAll(db, &data, stmt)
	return data, err
}

func (db *datastore) RepoFind(scm *model.ScmAccount) ([]*model.Repo, error) {
	stmt := sql.Lookup(db.driver, "repo-find-scm-id")
	data := []*model.Repo{}
	err := meddler.QueryAll(db, &data, stmt, scm.ID)
	return data, err
}

func (db *datastore) RepoLoad(id int64) (*model.Repo, error) {
	stmt := sql.Lookup(db.driver, "repo-find-id")
	repo := new(model.Repo)
	err := meddler.QueryRow(db, repo, stmt, id)
	return repo, err
}

func (db *datastore) GetRepoName(name string) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.QueryRow(db, repo, rebind(repoNameQuery), name)
	return repo, err
}

func (db *datastore) GetRepoScmName(scmID int64, name string) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.QueryRow(db, repo, rebind(repoScmIDNameQuery), scmID, name)
	return repo, err
}

const repoTable = "repos"

const repoNameQuery = `
SELECT *
FROM repos
WHERE repo_full_name = ?
LIMIT 1;
`
const repoScmIDNameQuery = `
SELECT *
FROM repos
WHERE repo_scm_id    = ?
  AND repo_full_name = ?
LIMIT 1;
`
