package torii

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// User represents a Torii user.
type User struct {
	ID                         int      `json:"id"`
	IDOrg                      int      `json:"idOrg"`
	FirstName                  string   `json:"firstName"`
	LastName                   string   `json:"lastName"`
	Email                      string   `json:"email"`
	CreationTime               string   `json:"creationTime"`
	IDRole                     int      `json:"idRole"`
	Role                       string   `json:"role"`
	LifecycleStatus            string   `json:"lifecycleStatus"`
	IsDeletedInIdentitySources bool     `json:"isDeletedInIdentitySources"`
	IsExternal                 bool     `json:"isExternal"`
	ActiveAppsCount            int      `json:"activeAppsCount"`
	AdditionalEmails           []string `json:"additionalEmails"`
}

// usersResponse is the envelope returned by GET /v1.0/users.
type usersResponse struct {
	Users      []User `json:"users"`
	Count      int    `json:"count"`
	Total      int    `json:"total"`
	NextCursor string `json:"nextCursor"`
}

//// TABLE DEFINITION

func tableToriiUser() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_user",
		Description: "List users in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "email", Require: plugin.Optional},
				{Name: "lifecycle_status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique user identifier."},
			{Name: "id_org", Type: proto.ColumnType_INT, Transform: transform.FromField("IDOrg"), Description: "Unique organization identifier."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name of the user."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Primary email address of the user."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the user was created."},
			{Name: "id_role", Type: proto.ColumnType_INT, Transform: transform.FromField("IDRole"), Description: "Unique role identifier assigned to the user."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "Role name assigned to the user (e.g. Admin, Member)."},
			{Name: "lifecycle_status", Type: proto.ColumnType_STRING, Description: "Lifecycle status of the user (active, offboarding, offboarded)."},
			{Name: "is_deleted_in_identity_sources", Type: proto.ColumnType_BOOL, Description: "True if the user has left the organization according to identity sources."},
			{Name: "is_external", Type: proto.ColumnType_BOOL, Description: "True if the user is external (contractor, consultant, etc.)."},
			{Name: "active_apps_count", Type: proto.ColumnType_INT, Description: "Number of active applications the user has access to."},
			{Name: "additional_emails", Type: proto.ColumnType_JSON, Description: "List of additional email addresses associated with the user."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"size":   "1000",
		"fields": "idOrg,firstName,lastName,email,creationTime,idRole,lifecycleStatus,isDeletedInIdentitySources,isExternal,activeAppsCount,role,additionalEmails",
	}
	if v := d.EqualsQualString("email"); v != "" {
		params["email"] = v
	}
	if v := d.EqualsQualString("lifecycle_status"); v != "" {
		params["lifecycleStatus"] = v
	}

	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		var result usersResponse
		if err := client.get(ctx, "/v1.0/users", params, &result); err != nil {
			return nil, fmt.Errorf("listing users: %w", err)
		}

		for _, u := range result.Users {
			d.StreamListItem(ctx, u)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextCursor == "" || len(result.Users) == 0 {
			break
		}
		cursor = result.NextCursor
	}

	return nil, nil
}

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"fields":  "idOrg,firstName,lastName,email,creationTime,idRole,lifecycleStatus,isDeletedInIdentitySources,isExternal,activeAppsCount,role,additionalEmails",
		"idUsers": strconv.FormatInt(id, 10),
	}

	var result usersResponse
	if err := client.get(ctx, "/v1.0/users", params, &result); err != nil {
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	if len(result.Users) == 0 {
		return nil, nil
	}
	return result.Users[0], nil
}
