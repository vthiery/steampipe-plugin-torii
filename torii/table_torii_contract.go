package torii

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// Contract represents a Torii contract.
type Contract struct {
	ID     int    `json:"id"`
	IDApp  int    `json:"idApp"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Status string `json:"status"`
}

// contractsResponse is the envelope returned by GET /v1.0/contracts.
type contractsResponse struct {
	Contracts  []Contract `json:"contracts"`
	Count      int        `json:"count"`
	Total      int        `json:"total"`
	NextCursor string     `json:"nextCursor"`
}

//// TABLE DEFINITION

func tableToriiContract() *plugin.Table {
	return &plugin.Table{
		Name:        "torii_contract",
		Description: "List software contracts managed in your Torii organization.",
		List: &plugin.ListConfig{
			Hydrate: listContracts,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContract,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique contract identifier."},
			{Name: "id_app", Type: proto.ColumnType_INT, Transform: transform.FromField("IDApp"), Description: "Unique identifier of the application this contract belongs to."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the contract."},
			{Name: "owner", Type: proto.ColumnType_STRING, Description: "Name of the contract owner."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the contract (e.g. active, expired)."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listContracts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"fields": "id,name,idApp,owner,status",
	}
	if v := d.EqualsQualString("status"); v != "" {
		params["status"] = v
	}

	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}

		var result contractsResponse
		if err := client.get(ctx, "/v1.0/contracts", params, &result); err != nil {
			return nil, fmt.Errorf("listing contracts: %w", err)
		}

		for _, c := range result.Contracts {
			d.StreamListItem(ctx, c)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextCursor == "" || len(result.Contracts) == 0 {
			break
		}
		cursor = result.NextCursor
	}

	return nil, nil
}

func getContract(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result Contract
	if err := client.get(ctx, fmt.Sprintf("/v1.0/contracts/%d", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting contract %d: %w", id, err)
	}

	return result, nil
}
