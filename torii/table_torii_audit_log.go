package torii

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// AuditRequestDetails holds HTTP request metadata for an audit log entry.
type AuditRequestDetails struct {
	Path          string `json:"path"`
	Method        string `json:"method"`
	RemoteAddress string `json:"remoteAddress"`
}

// AuditLog represents a Torii admin audit log entry.
type AuditLog struct {
	PerformedBy          int                    `json:"performedBy"`
	PerformedByFirstName string                 `json:"performedByFirstName"`
	PerformedByLastName  string                 `json:"performedByLastName"`
	PerformedByEmail     string                 `json:"performedByEmail"`
	IDTargetOrg          int                    `json:"idTargetOrg"`
	CreationTime         string                 `json:"creationTime"`
	Type                 string                 `json:"type"`
	RequestDetails       AuditRequestDetails    `json:"requestDetails"`
	Properties           map[string]interface{} `json:"properties"`
}

// auditLogsResponse is the envelope returned by GET /v1.0/audit.
type auditLogsResponse struct {
	Audit      []AuditLog `json:"audit"`
	NextCursor string     `json:"nextCursor"`
	Count      int        `json:"count"`
}

//// TABLE DEFINITION

func tableToriiAuditLog() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_audit_log",
		Description: "List admin audit log events in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listAuditLogs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "performed_by", Type: proto.ColumnType_INT, Description: "User ID of the user who performed the action."},
			{Name: "performed_by_first_name", Type: proto.ColumnType_STRING, Description: "First name of the user who performed the action."},
			{Name: "performed_by_last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user who performed the action."},
			{Name: "performed_by_email", Type: proto.ColumnType_STRING, Description: "Email address of the user who performed the action."},
			{Name: "id_target_org", Type: proto.ColumnType_INT, Transform: transform.FromField("IDTargetOrg"), Description: "Organization ID the action was performed in."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the audit log entry was created."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type/category of the audited action (e.g. create_workflow, update_user)."},
			{Name: "request_details", Type: proto.ColumnType_JSON, Description: "HTTP request details (path, method, remote address) of the action."},
			{Name: "properties", Type: proto.ColumnType_JSON, Description: "Additional properties and context for the audited action."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"size": "1000",
		"sort": "desc",
	}
	if v := d.EqualsQualString("type"); v != "" {
		params["entity"] = v
	}

	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		var result auditLogsResponse
		if err := client.get(ctx, "/v1.0/audit", params, &result); err != nil {
			return nil, fmt.Errorf("listing audit logs: %w", err)
		}

		for _, a := range result.Audit {
			d.StreamListItem(ctx, a)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextCursor == "" || len(result.Audit) == 0 {
			break
		}
		cursor = result.NextCursor
	}

	return nil, nil
}
