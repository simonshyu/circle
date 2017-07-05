-- name: create-table-secrets

CREATE TABLE IF NOT EXISTS secrets (
 secret_id          INTEGER PRIMARY KEY AUTO_INCREMENT
,secret_repo_id     INTEGER
,secret_name        VARCHAR(250)
,secret_value       MEDIUMBLOB

,UNIQUE(secret_name, secret_repo_id)
);

-- name: create-index-secrets-repo

CREATE INDEX ix_secrets_repo  ON secrets  (secret_repo_id);
