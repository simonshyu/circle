package mysql

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"scm-account-list":    scmAccountList,
	"scm-account-find-id": scmAccountFindId,
	"repo-find-id":        repoFindId,
	"config-find-id":      configFindId,
	"config-find-repo":    configFindRepo,
}

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
,repo_hash
FROM repos
WHERE repo_id = ?
`
