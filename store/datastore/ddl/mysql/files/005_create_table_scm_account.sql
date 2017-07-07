-- name: create-table-scm-account

CREATE TABLE IF NOT EXISTS scm_account (
 scm_id            INTEGER PRIMARY KEY AUTO_INCREMENT
,scm_host          VARCHAR(500)
,scm_login         VARCHAR(250)
,scm_password      VARCHAR(250)
,scm_type          VARCHAR(50)
);
