package mysql

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"config-find-id":      configFindId,
	"scm-account-list":    scmAccountList,
	"scm-account-find-id": scmAccountFindId,
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

var scmAccountList = `
SELECT
 scm_id
,scm_host
,scm_login
,scm_password
,scm_type
FROM scm_account
`

var scmAccountFindId = `
SELECT
 scm_id
,scm_host
,scm_login
,scm_password
,scm_type
FROM scm_account
WHERE scm_id = ?
`
