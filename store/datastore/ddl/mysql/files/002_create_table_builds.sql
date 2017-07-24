-- name: create-table-builds

CREATE TABLE IF NOT EXISTS builds (
 build_id        INTEGER PRIMARY KEY AUTO_INCREMENT
,build_config_id INTEGER
,build_repo_id   INTEGER
,build_number    INTEGER
,build_event     VARCHAR(500)
,build_status    VARCHAR(500)
,build_error     VARCHAR(500)
,build_enqueued  INTEGER
,build_created   INTEGER
,build_started   INTEGER
,build_finished  INTEGER
,build_link      VARCHAR(1000)
,build_commit    VARCHAR(500)
,build_branch    VARCHAR(500)
,build_ref       VARCHAR(500)
,build_refspec   VARCHAR(1000)
,build_remote    VARCHAR(500)

,UNIQUE(build_number, build_repo_id)
);

-- name: create-index-builds-repo

CREATE INDEX ix_build_repo ON builds (build_repo_id);
