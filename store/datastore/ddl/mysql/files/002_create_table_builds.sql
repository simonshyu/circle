-- name: create-table-builds

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

-- name: create-index-builds-repo

CREATE INDEX ix_build_repo ON builds (build_repo_id);
