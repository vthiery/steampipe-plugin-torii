// Package torii provides a Steampipe plugin for querying Torii SaaS management resources using SQL.
package torii

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the Torii Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-torii",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"torii_user":           tableToriiUser(),
			"torii_app":            tableToriiApp(),
			"torii_app_field":      tableToriiAppField(),
			"torii_app_user":       tableToriiAppUser(),
			"torii_contract":       tableToriiContract(),
			"torii_contract_field": tableToriiContractField(),
			"torii_role":           tableToriiRole(),
			"torii_audit_log":      tableToriiAuditLog(),
			"torii_user_app":       tableToriiUserApp(),
			"torii_user_field":     tableToriiUserField(),
			"torii_organization":   tableToriiOrganization(),
		},
	}
	return p
}
