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
	"build-find-id":       buildFindId,
	"build-find-number":   buildFindNumber,
	"config-find-id":      configFindId,
	"config-find-repo":    configFindRepo,
	"proc-find-build":     procsFindBuild,
	"proc-find-id":        procFindId,
	"procs-delete-build":  procsDeleteBuild,
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
,repo_link
,repo_branch
,repo_full_name
,repo_owner
,repo_name
,repo_timeout
,repo_private
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_allow_manual
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
,repo_link
,repo_full_name
,repo_owner
,repo_name
,repo_timeout
,repo_private
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_allow_manual
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
,repo_link
,repo_full_name
,repo_owner
,repo_name
,repo_timeout
,repo_private
,repo_allow_pr
,repo_allow_push
,repo_allow_tags
,repo_allow_manual
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

var buildFindId = `
SELECT
 build_id
,build_config_id
,build_repo_id
,build_number
,build_event
,build_status
,build_error
,build_enqueued
,build_created
,build_started
,build_finished
,build_link
,build_commit
,build_branch
,build_ref
,build_refspec
,build_remote
FROM builds
WHERE build_id = ?
`

var buildFindNumber = `
SELECT
 build_id
,build_config_id
,build_repo_id
,build_number
,build_event
,build_status
,build_error
,build_enqueued
,build_created
,build_started
,build_finished
,build_link
,build_commit
,build_branch
,build_ref
,build_refspec
,build_remote
FROM builds
WHERE build_repo_id = ?
  AND build_number = ?
LIMIT 1;
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

var procsFindBuild = `
SELECT
 proc_id
,proc_build_id
,proc_pid
,proc_ppid
,proc_name
,proc_state
,proc_error
,proc_exit_code
,proc_started
,proc_stopped
,proc_machine
,proc_platform
,proc_environ
FROM procs
WHERE proc_build_id = ?
ORDER BY proc_id ASC
`

var procFindId = `
SELECT
 proc_id
,proc_build_id
,proc_pid
,proc_ppid
,proc_name
,proc_state
,proc_error
,proc_exit_code
,proc_started
,proc_stopped
,proc_machine
,proc_platform
,proc_environ
FROM procs
WHERE proc_id = ?
`

var procsDeleteBuild = `
DELETE FROM procs WHERE proc_build_id = ?
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
