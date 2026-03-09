package torii

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// AppOwner represents a Torii app owner reference.
type AppOwner struct {
	ID                         int    `json:"id"`
	FirstName                  string `json:"firstName"`
	LastName                   string `json:"lastName"`
	FullName                   string `json:"fullName"`
	Email                      string `json:"email"`
	PhotoURL                   string `json:"photoUrl"`
	LifecycleStatus            string `json:"lifecycleStatus"`
	IsDeletedInIdentitySources bool   `json:"isDeletedInIdentitySources"`
	Status                     string `json:"status"`
}

// App represents a Torii application.
type App struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	State            string   `json:"state"`
	Category         string   `json:"category"`
	URL              string   `json:"url"`
	ImageURL         string   `json:"imageUrl"`
	Description      string   `json:"description"`
	Tags             []string `json:"tags"`
	IsCustom         bool     `json:"isCustom"`
	IsHidden         bool     `json:"isHidden"`
	CreationTime     string   `json:"creationTime"`
	LastVisitTime    FlexTime `json:"lastVisitTime"`
	ActiveUsersCount int      `json:"activeUsersCount"`
	PrimaryOwner     AppOwner `json:"primaryOwner"`
}

// appsResponse is the envelope returned by GET /v1.0/apps.
type appsResponse struct {
	Apps       []App  `json:"apps"`
	Count      int    `json:"count"`
	Total      int    `json:"total"`
	NextCursor string `json:"nextCursor"`
}

//// TABLE DEFINITION

func tableToriiApp() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_app",
		Description: "List applications discovered or managed in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listApps,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "state", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique application identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the application."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the application (e.g. Discovered, Managed, Ignored)."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Category of the application (e.g. Developer Tools, HR)."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL of the application."},
			{Name: "image_url", Type: proto.ColumnType_STRING, Description: "URL of the application's logo image."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the application."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags associated with the application."},
			{Name: "is_custom", Type: proto.ColumnType_BOOL, Description: "True if this is a custom (manually added) application."},
			{Name: "is_hidden", Type: proto.ColumnType_BOOL, Description: "True if the application is hidden."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the application was discovered or created."},
			{Name: "last_visit_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastVisitTime").Transform(flexTimeTransform), Description: "Time of the last recorded user visit."},
			{Name: "active_users_count", Type: proto.ColumnType_INT, Description: "Number of active users for this application."},
			{Name: "primary_owner", Type: proto.ColumnType_JSON, Description: "Primary owner of the application."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"size":   "1000",
		"fields": "id,name,primaryOwner,state,category,url,imageUrl,description,tags,isCustom,isHidden,creationTime,lastVisitTime,activeUsersCount",
	}
	if v := d.EqualsQualString("state"); v != "" {
		params["state"] = v
	}

	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		var result appsResponse
		if err := client.get(ctx, "/v1.0/apps", params, &result); err != nil {
			return nil, fmt.Errorf("listing apps: %w", err)
		}

		for _, a := range result.Apps {
			d.StreamListItem(ctx, a)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextCursor == "" || len(result.Apps) == 0 {
			break
		}
		cursor = result.NextCursor
	}

	return nil, nil
}
