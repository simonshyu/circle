package mysql

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"scm-account-list":    scmAccountList,
	"scm-account-find-id": scmAccountFindId,
	"repo-list":           repoList,
	"repo-find-id":        repoFindId,
	"repo-find-scm-id":    repoFindScmId,
	"repo-update-counter": repoUpdateCounter,
	"config-find-id":      configFindId,
	"config-find-repo":    configFindRepo,
	"task-list":           taskList,
	"task-delete":         taskDelete,
}

var scmAccountList = `
SELECT
 scm_id
,scm_host
,scm_login
,scm_password
,scm_private_token
,scm_type
FROM scm_account
`

var scmAccountFindId = `
SELECT
 scm_id
,scm_host
,scm_login
,scm_password
,scm_private_token
,scm_type
FROM scm_account
WHERE scm_id = ?
`

var repoList = `
SELECT
 repo_id
,repo_scm_id
,repo_clone
,repo_branch
,repo_full_name
,repo_owner
,repo_name
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_counter
,repo_hash
FROM repos
`

var repoFindId = `
SELECT
 repo_id
,repo_scm_id
,repo_clone
,repo_branch
,repo_full_name
,repo_owner
,repo_name
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_counter
,repo_hash
FROM repos
WHERE repo_id = ?
`

var repoFindScmId = `
SELECT
 repo_id
,repo_scm_id
,repo_clone
,repo_branch
,repo_full_name
,repo_owner
,repo_name
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_counter
,repo_hash
FROM repos
WHERE repo_scm_id = ?
ORDER BY repo_id ASC
`

var repoUpdateCounter = `
UPDATE repos SET repo_counter = ?
WHERE repo_counter = ?
  AND repo_id = ?
`

var configFindId = `
SELECT
 config_id
,config_repo_id
,config_hash
,config_data
FROM config
WHERE config_id = ?
`

var configFindRepo = `
SELECT
 config_id
,config_repo_id
,config_hash
,config_data
FROM config
WHERE config_repo_id = ?
LIMIT 1
`

var taskList = `
SELECT
 task_id
,task_data
,task_labels
FROM tasks
`

var taskDelete = `
DELETE FROM tasks WHERE task_id = ?
`
