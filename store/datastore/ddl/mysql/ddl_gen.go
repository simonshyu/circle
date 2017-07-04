package mysql

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-repos",
		stmt: createTableRepos,
	},
	{
		name: "create-table-builds",
		stmt: createTableBuilds,
	},
	{
		name: "create-index-builds-repo",
		stmt: createIndexBuildsRepo,
	},
	{
		name: "create-table-config",
		stmt: createTableConfig,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 001_create_table_repos.sql
//

var createTableRepos = `
CREATE TABLE IF NOT EXISTS repos (
 repo_id            INTEGER PRIMARY KEY AUTO_INCREMENT
,repo_host          VARCHAR(250)
,repo_login         VARCHAR(250)
,repo_password      VARCHAR(250)
,repo_clone         VARCHAR(1000)
,repo_branch        VARCHAR(500)
,repo_scm           VARCHAR(50)
);
`

//
// 002_create_table_builds.sql
//

var createTableBuilds = `
CREATE TABLE IF NOT EXISTS builds (
 build_id        INTEGER PRIMARY KEY AUTO_INCREMENT
,build_repo_id   INTEGER
,build_number    INTEGER
,build_created   INTEGER
,build_started   INTEGER
,build_finished  INTEGER
,build_config_id INTEGER

,UNIQUE(build_number, build_repo_id)
);
`

var createIndexBuildsRepo = `
CREATE INDEX ix_build_repo ON builds (build_repo_id);
`

//
// 003_create_table_config.sql
//

var createTableConfig = `
CREATE TABLE IF NOT EXISTS config (
 config_id       INTEGER PRIMARY KEY AUTO_INCREMENT
,config_repo_id  INTEGER
,config_hash     VARCHAR(250)
,config_data     MEDIUMBLOB

,UNIQUE(config_hash, config_repo_id)
);
`
