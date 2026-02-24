// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

var _ datasource.DataSource = &AppDataSource{}

// AppDataSource allows looking up an existing app by its numeric ID.
type AppDataSource struct {
	client *client.Client
}

type AppDataSourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Status            types.String `tfsdk:"status"`
	StatusCategory    types.String `tfsdk:"status_category"`
	HostingType       types.String `tfsdk:"hosting_type"`
	CustomPackageName types.String `tfsdk:"custom_package_name"`
	BillingStatus     types.String `tfsdk:"billing_status"`
	OrganisationID    types.Int64  `tfsdk:"organisation_id"`
}

func NewAppDataSource() datasource.DataSource {
	return &AppDataSource{}
}

func (d *AppDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *AppDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an existing Level27 **App** by its numeric ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Numeric ID of the app to look up.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the app.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current status of the app.",
			},
			"status_category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status category: `green`, `yellow`, `red`, or `grey`.",
			},
			"hosting_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hosting type: `Agency` or `Classic`.",
			},
			"custom_package_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the custom package.",
			},
			"billing_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Billing status of the app.",
			},
			"organisation_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the owning organisation.",
			},
		},
	}
}

func (d *AppDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *AppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := d.client.GetApp(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading app", err.Error())
		return
	}

	state.Name = types.StringValue(app.Name)
	state.Status = types.StringValue(app.Status)
	state.StatusCategory = types.StringValue(app.StatusCategory)
	state.HostingType = types.StringValue(app.HostingType)
	state.CustomPackageName = types.StringValue(app.CustomPackageName)
	state.BillingStatus = types.StringValue(app.BillingStatus)
	if app.Organisation != nil {
		state.OrganisationID = types.Int64Value(int64(app.Organisation.ID))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
