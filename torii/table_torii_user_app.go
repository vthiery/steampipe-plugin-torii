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

// UserApp represents an application used by a specific Torii user.
type UserApp struct {
	ID                   int      `json:"id"`
	IDUser               int      `json:"-"` // injected from path param
	Name                 string   `json:"name"`
	PrimaryOwner         int      `json:"primaryOwner"`
	State                string   `json:"state"`
	Category             string   `json:"category"`
	URL                  string   `json:"url"`
	ImageURL             string   `json:"imageUrl"`
	Description          string   `json:"description"`
	Tags                 []string `json:"tags"`
	Score                float64  `json:"score"`
	IsCustom             bool     `json:"isCustom"`
	AddedBy              int      `json:"addedBy"`
	CreationTime         string   `json:"creationTime"`
	IsHidden             bool     `json:"isHidden"`
	Sources              []string `json:"sources"`
	LastUsageTime        string   `json:"lastUsageTime"`
	IsUserRemovedFromApp bool     `json:"isUserRemovedFromApp"`
	AnnualCost           float64  `json:"annualCost"`
	Currency             string   `json:"currency"`
}

// userAppsResponse is the envelope returned by GET /v1.0/users/{idUser}/apps.
type userAppsResponse struct {
	Apps []UserApp `json:"apps"`
}

// userAppFields lists all fields requested from the API.
const userAppFields = "id,name,primaryOwner,state,category,url,imageUrl,description,tags,score,isCustom,addedBy,creationTime,isHidden,sources,lastUsageTime,isUserRemovedFromApp,annualCost,currency"

//// TABLE DEFINITION

func tableToriiUserApp() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_user_app",
		Description: "List applications used by a specific user in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listUserApps,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_id", Require: plugin.Required},
				{Name: "state", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "user_id", Type: proto.ColumnType_INT, Transform: transform.FromField("IDUser"), Description: "Unique identifier of the user."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier of the application."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the application."},
			{Name: "primary_owner", Type: proto.ColumnType_INT, Description: "User ID of the primary application owner."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the application (Managed, Discovered, etc.)."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Category of the application."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "URL of the application."},
			{Name: "image_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ImageURL"), Description: "URL of the application image."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the application."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags associated with the application."},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "Usage score of the application for the user."},
			{Name: "is_custom", Type: proto.ColumnType_BOOL, Description: "True if the application was manually added."},
			{Name: "added_by", Type: proto.ColumnType_INT, Description: "User ID of the person who added the application."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the application was created in Torii."},
			{Name: "is_hidden", Type: proto.ColumnType_BOOL, Description: "True if the application is hidden."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "Identity sources reporting this application."},
			{Name: "last_usage_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time of the user's last usage of the application."},
			{Name: "is_user_removed_from_app", Type: proto.ColumnType_BOOL, Description: "True if the user was removed from the application."},
			{Name: "annual_cost", Type: proto.ColumnType_DOUBLE, Description: "Annual cost of the license in the license currency."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "Currency of the annual cost."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUserApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	userID := d.EqualsQuals["user_id"].GetInt64Value()
	if userID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"fields": userAppFields,
	}
	if v := d.EqualsQualString("state"); v != "" {
		params["state"] = v
	}

	var result userAppsResponse
	path := fmt.Sprintf("/v1.0/users/%s/apps", strconv.FormatInt(userID, 10))
	if err := client.get(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("listing apps for user %d: %w", userID, err)
	}

	for _, a := range result.Apps {
		a.IDUser = int(userID)
		d.StreamListItem(ctx, a)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
