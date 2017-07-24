package model

// ConfigStore persists pipeline configuration to storage.
// Often access by handler.Config.Storage.Config
type ConfigStore interface {
	ConfigCreate(*Config) error
	ConfigLoad(int64) (*Config, error)
}

// Config represents a pipeline configuration.
type Config struct {
	ID     int64  `json:"-"    meddler:"config_id,pk"`
	RepoID int64  `json:"-"    meddler:"config_repo_id"`
	Data   string `json:"data" meddler:"config_data"`
	Hash   string `json:"hash" meddler:"config_hash"`
}

/*
Field Explanation:
*/
