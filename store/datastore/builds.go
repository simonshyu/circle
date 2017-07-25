package datastore

import (
	"fmt"
	"time"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store/datastore/sql"
	"github.com/russross/meddler"
)

func (db *datastore) BuildCreate(build *model.Build, procs ...*model.Proc) error {
	id, err := db.incrementRepoRetry(build.RepoID)
	if err != nil {
		return err
	}
	build.Number = id
	build.Created = time.Now().UTC().Unix()
	build.Enqueued = build.Created
	err = meddler.Insert(db, buildTable, build)
	if err != nil {
		return err
	}
	for _, proc := range procs {
		proc.BuildID = build.ID
		err = meddler.Insert(db, "procs", proc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *datastore) BuildFind(repo *model.Repo) ([]*model.Build, error) {
	stmt := sql.Lookup(db.driver, "build-find-repo-id")
	builds := []*model.Build{}
	err := meddler.QueryAll(db, &builds, stmt, repo.ID)
	return builds, err
}

func (db *datastore) BuildLoad(id int64) (*model.Build, error) {
	stmt := sql.Lookup(db.driver, "build-find-id")
	build := new(model.Build)
	err := meddler.QueryRow(db, build, stmt, id)
	return build, err
}

func (db *datastore) GetBuildNumber(repo *model.Repo, num int) (*model.Build, error) {
	stmt := sql.Lookup(db.driver, "build-find-number")
	build := new(model.Build)
	err := meddler.QueryRow(db, build, stmt, repo.ID, num)
	return build, err
}

func (db *datastore) incrementRepoRetry(id int64) (int, error) {
	repo, err := db.RepoLoad(id)
	if err != nil {
		return 0, fmt.Errorf("database: cannot fetch repository: %s", err)
	}
	for i := 0; i < 10; i++ {
		seq, err := db.incrementRepo(repo.ID, repo.Counter+i, repo.Counter+i+1)
		if err != nil {
			return 0, err
		}
		if seq == 0 {
			continue
		}
		return seq, nil
	}
	return 0, fmt.Errorf("cannot increment next build number")
}

func (db *datastore) incrementRepo(id int64, old, new int) (int, error) {
	results, err := db.Exec(sql.Lookup(db.driver, "repo-update-counter"), new, old, id)
	if err != nil {
		return 0, fmt.Errorf("database: update repository counter: %s", err)
	}
	updated, err := results.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("database: update repository counter: %s", err)
	}
	if updated == 0 {
		return 0, nil
	}
	return new, nil
}

func (db *datastore) BuildUpdate(build *model.Build) error {
	return meddler.Update(db, buildTable, build)
}

const buildTable = "builds"
