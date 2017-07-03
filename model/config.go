package model

// ConfigStore persists pipeline configuration to storage.
type ConfigStore interface {
	ConfigCreate(*Config) error
}

// Config represents a pipeline configuration.
type Config struct {
	ID   int64  `json:"-"    meddler:"config_id,pk"`
	Data string `json:"data" meddler:"config_data"`
	Hash string `json:"hash" meddler:"config_hash"`
}
