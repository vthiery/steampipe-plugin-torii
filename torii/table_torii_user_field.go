package torii

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// UserField represents a custom field available on user records from connected integrations.
type UserField struct {
	ID          int    `json:"id"`
	IDOrg       int    `json:"idOrg"`
	SourceIDApp int    `json:"sourceIdApp"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	IsDeleted   bool   `json:"isDeleted"`
}

// userFieldsResponse is the envelope returned by GET /v1.0/users/fields.
type userFieldsResponse struct {
	Fields []UserField `json:"fields"`
}

//// TABLE DEFINITION

func tableToriiUserField() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_user_field",
		Description: "List custom user fields from connected integrations in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listUserFields,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier of the field."},
			{Name: "id_org", Type: proto.ColumnType_INT, Transform: transform.FromField("IDOrg"), Description: "Unique identifier of the organization."},
			{Name: "source_id_app", Type: proto.ColumnType_INT, Transform: transform.FromField("SourceIDApp"), Description: "Unique identifier of the source application providing this field."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Display name of the field."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the field (e.g. singleLine, multiLine, date)."},
			{Name: "key", Type: proto.ColumnType_STRING, Description: "Internal key used to reference the field."},
			{Name: "is_deleted", Type: proto.ColumnType_BOOL, Description: "True if the field has been deleted."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUserFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	if q := d.EqualsQualString("name"); q != "" {
		params["q"] = q
	}

	var result userFieldsResponse
	if err := client.get(ctx, "/v1.0/users/fields", params, &result); err != nil {
		return nil, err
	}

	for _, f := range result.Fields {
		d.StreamListItem(ctx, f)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
