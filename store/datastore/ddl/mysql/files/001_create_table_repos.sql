-- name: create-table-repos

CREATE TABLE IF NOT EXISTS repos (
 repo_id            INTEGER PRIMARY KEY AUTO_INCREMENT
,repo_clone         VARCHAR(1000)
,repo_branch        VARCHAR(500)
,repo_scm           VARCHAR(50)
);
