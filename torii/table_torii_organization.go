package torii

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Organization represents the Torii organization profile.
type Organization struct {
	ID           int    `json:"id"`
	CompanyName  string `json:"companyName"`
	Domain       string `json:"domain"`
	CreationTime string `json:"creationTime"`
}

// organizationResponse is the envelope returned by GET /v1.0/orgs/my.
type organizationResponse struct {
	Org Organization `json:"org"`
}

//// TABLE DEFINITION

func tableToriiOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_organization",
		Description: "Get the organization profile for your Torii account.",
		List: &plugin.ListConfig{
			Hydrate: getOrganization,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier of the organization."},
			{Name: "company_name", Type: proto.ColumnType_STRING, Description: "Company name of the organization."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Primary domain of the organization."},
			{Name: "creation_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time the organization was created in Torii."},
		},
	}
}

//// HYDRATE FUNCTIONS

func getOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result organizationResponse
	if err := client.get(ctx, "/v1.0/orgs/my", nil, &result); err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, result.Org)
	return nil, nil
}
