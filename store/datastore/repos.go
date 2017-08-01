package datastore

import (
	// "fmt"

	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/store/datastore/sql"
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

func (db *datastore) GetRepoScmName(scmID int64, name string) (*model.Repo, error) {
	stmt := sql.Lookup(db.driver, "repo-find-scm-id-fullname")
	repo := new(model.Repo)
	err := meddler.QueryRow(db, repo, stmt, scmID, name)
	return repo, err
}

const repoTable = "repos"
