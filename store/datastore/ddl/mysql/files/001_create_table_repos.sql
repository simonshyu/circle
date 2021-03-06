-- name: create-table-repos

CREATE TABLE IF NOT EXISTS repos (
 repo_id            INTEGER PRIMARY KEY AUTO_INCREMENT
,repo_scm_id        INTEGER
,repo_clone         VARCHAR(1000)
,repo_link          VARCHAR(1000)
,repo_branch        VARCHAR(500)
,repo_full_name     VARCHAR(500)
,repo_owner         VARCHAR(250)
,repo_name          VARCHAR(250)
,repo_timeout       INTEGER
,repo_private       BOOLEAN
,repo_allow_pr      BOOLEAN
,repo_allow_push    BOOLEAN
,repo_allow_tags    BOOLEAN
,repo_allow_manual  BOOLEAN
,repo_counter       INTEGER
,repo_hash          VARCHAR(500)

,UNIQUE(repo_scm_id, repo_owner, repo_name)
);
