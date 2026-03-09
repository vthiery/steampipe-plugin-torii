package torii

import "github.com/turbot/steampipe-plugin-sdk/v5/plugin"

// ToriiConfig holds the connection configuration for the Torii plugin.
type ToriiConfig struct {
	// API key for authenticating with the Torii API.
	APIKey *string `hcl:"api_key"`
}

// ConfigInstance returns a new, empty ToriiConfig.
func ConfigInstance() interface{} {
	return &ToriiConfig{}
}

// GetConfig retrieves the connection configuration for the given connection.
func GetConfig(connection *plugin.Connection) ToriiConfig {
	if connection == nil || connection.Config == nil {
		return ToriiConfig{}
	}
	cfg, _ := connection.Config.(ToriiConfig)
	return cfg
}
