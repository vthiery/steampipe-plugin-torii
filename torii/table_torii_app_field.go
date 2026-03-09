package torii

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// FieldOption represents a selectable option for a dropdown field.
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// AppField represents a custom or built-in field defined for applications.
type AppField struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	SystemKey string        `json:"systemKey"`
	Type      string        `json:"type"`
	Options   []FieldOption `json:"options"`
}

// appFieldsResponse is the envelope returned by GET /v1.0/apps/fields.
type appFieldsResponse struct {
	Fields []AppField `json:"fields"`
}

//// TABLE DEFINITION

func tableToriiAppField() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_app_field",
		Description: "List custom and built-in fields defined for applications in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listAppFields,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier of the field."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Display name of the field."},
			{Name: "system_key", Type: proto.ColumnType_STRING, Description: "Internal system key of the field."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the field (e.g. dropdown, singleLine, multiLine, date)."},
			{Name: "options", Type: proto.ColumnType_JSON, Description: "Available options for dropdown-type fields."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAppFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result appFieldsResponse
	if err := client.get(ctx, "/v1.0/apps/fields", nil, &result); err != nil {
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
