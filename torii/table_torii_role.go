package torii

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Role represents a Torii role.
type Role struct {
	ID          int    `json:"id"`
	SystemKey   string `json:"systemKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsAdmin     bool   `json:"isAdmin"`
	UsersCount  int    `json:"usersCount"`
}

// rolesResponse is the envelope returned by GET /v1.0/roles.
type rolesResponse struct {
	Roles []Role `json:"roles"`
}

//// TABLE DEFINITION

func tableToriiRole() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_role",
		Description: "List roles defined in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listRoles,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique role identifier."},
			{Name: "system_key", Type: proto.ColumnType_STRING, Description: "System key for the role (e.g. admin, member)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Display name of the role."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the role and its permissions."},
			{Name: "is_admin", Type: proto.ColumnType_BOOL, Description: "True if this role grants administrative access."},
			{Name: "users_count", Type: proto.ColumnType_INT, Description: "Number of users assigned to this role."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result rolesResponse
	if err := client.get(ctx, "/v1.0/roles", nil, &result); err != nil {
		return nil, fmt.Errorf("listing roles: %w", err)
	}

	for _, r := range result.Roles {
		d.StreamListItem(ctx, r)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
