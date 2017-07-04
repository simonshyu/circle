package sqlite

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"config-find-id": configFindId,
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
