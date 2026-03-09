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

// AppUser represents a user of a specific Torii application.
type AppUser struct {
	IDUser                      int      `json:"idUser"`
	IDApp                       int      `json:"-"` // injected from path param
	Email                       string   `json:"email"`
	FirstName                   string   `json:"firstName"`
	LastName                    string   `json:"lastName"`
	FullName                    string   `json:"fullName"`
	AppName                     string   `json:"appName"`
	Status                      string   `json:"status"`
	Role                        string   `json:"role"`
	ExternalStatus              string   `json:"externalStatus"`
	Sources                     []string `json:"sources"`
	IsExternal                  bool     `json:"isExternal"`
	IsDeletedInIdentitySources  bool     `json:"isDeletedInIdentitySources"`
	IsUserRemovedFromApp        bool     `json:"isUserRemovedFromApp"`
	IsOwnerNotifiedToRemoveUser bool     `json:"isOwnerNotifiedToRemoveUser"`
	IsNewUserInApp              bool     `json:"isNewUserInApp"`
	IsOffboardingIgnored        bool     `json:"isOffboardingIgnored"`
	CreationTime                string   `json:"creationTime"`
	LastVisitTime               string   `json:"lastVisitTime"`
	CreationTimeInApp           string   `json:"creationTimeInApp"`
	LifecycleStatus             string   `json:"lifecycleStatus"`
	Score                       float64  `json:"score"`
}

// appUsersResponse is the envelope returned by GET /v1.0/apps/{idApp}/users.
type appUsersResponse struct {
	Users      []AppUser `json:"users"`
	Count      int       `json:"count"`
	Total      int       `json:"total"`
	NextCursor string    `json:"nextCursor"`
}

//// TABLE DEFINITION

func tableToriiAppUser() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_app_user",
		Description: "List users of a specific application in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listAppUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "app_id", Require: plugin.Required},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "app_id", Type: proto.ColumnType_INT, Transform: transform.FromField("IDApp"), Description: "Unique identifier of the application."},
			{Name: "id_user", Type: proto.ColumnType_INT, Transform: transform.FromField("IDUser"), Description: "Unique identifier of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name of the user."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "Full name of the user."},
			{Name: "app_name", Type: proto.ColumnType_STRING, Description: "Name of the application."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the user in the application (active, deleted, managed)."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "Role of the user in the application."},
			{Name: "external_status", Type: proto.ColumnType_STRING, Description: "Status reported by the external application source."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "List of identity sources (e.g. okta, google) that report this user's access."},
			{Name: "is_external", Type: proto.ColumnType_BOOL, Description: "True if the user is external."},
			{Name: "is_deleted_in_identity_sources", Type: proto.ColumnType_BOOL, Description: "True if the user has left the organization according to identity sources."},
			{Name: "is_user_removed_from_app", Type: proto.ColumnType_BOOL, Description: "True if the user has been removed from the application."},
			{Name: "is_owner_notified_to_remove_user", Type: proto.ColumnType_BOOL, Description: "True if the app owner has been notified to remove this user."},
			{Name: "is_new_user_in_app", Type: proto.ColumnType_BOOL, Description: "True if this user is newly detected in the application."},
			{Name: "is_offboarding_ignored", Type: proto.ColumnType_BOOL, Description: "True if the offboarding process is being ignored for this user."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the user record was created in Torii."},
			{Name: "last_visit_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time of the user's last visit to the application."},
			{Name: "creation_time_in_app", Type: proto.ColumnType_TIMESTAMP, Description: "Time the user was created in the application."},
			{Name: "lifecycle_status", Type: proto.ColumnType_STRING, Description: "Lifecycle status of the user (active, offboarding, offboarded)."},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "Usage score of the user for the application."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAppUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	appID := d.EqualsQuals["app_id"].GetInt64Value()
	if appID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"size": "1000",
	}
	if v := d.EqualsQualString("status"); v != "" {
		params["status"] = v
	}

	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		var result appUsersResponse
		path := fmt.Sprintf("/v1.0/apps/%s/users", strconv.FormatInt(appID, 10))
		if err := client.get(ctx, path, params, &result); err != nil {
			return nil, fmt.Errorf("listing app users for app %d: %w", appID, err)
		}

		for _, u := range result.Users {
			u.IDApp = int(appID)
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
