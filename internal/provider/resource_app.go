// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &AppResource{}
var _ resource.ResourceWithImportState = &AppResource{}

// AppResource defines the resource implementation.
type AppResource struct {
	client         *client.Client
	organisationID int64
}

// AppResourceModel maps the resource schema data.
type AppResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	OrganisationID    types.Int64  `tfsdk:"organisation_id"`
	CustomPackageID   types.Int64  `tfsdk:"custom_package_id"`
	CustomPackageName types.String `tfsdk:"custom_package_name"`
	AutoTeams         types.String `tfsdk:"auto_teams"`
	AutoUpgrades      types.String `tfsdk:"auto_upgrades"`
	ExternalInfo      types.String `tfsdk:"external_info"`
	// Computed
	Status         types.String `tfsdk:"status"`
	StatusCategory types.String `tfsdk:"status_category"`
	HostingType    types.String `tfsdk:"hosting_type"`
	BillingStatus  types.String `tfsdk:"billing_status"`
}

func NewAppResource() resource.Resource {
	return &AppResource{}
}

func (r *AppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *AppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Level27 **App** (project). An app groups components like PHP web-apps, databases, mail, etc.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the app.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the app/project.",
			},
			"organisation_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the organisation that owns this app.",
			},
			"custom_package_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "ID of the custom package (billing product).",
			},
			"custom_package_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of the custom package, e.g. `drupal_run`, `wordpress_walk`.",
			},
			"auto_teams": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comma-separated list of team IDs to auto-assign to this app.",
			},
			"auto_upgrades": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comma-separated list of upgrade names to apply automatically.",
			},
			"external_info": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "External reference (required when billableitemInfo entities exist for the organisation).",
			},
			// Computed
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current status of the app (e.g. `ok`, `to_create`, `creating`).",
			},
			"status_category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status category: `green`, `yellow`, `red`, or `grey`.",
			},
			"hosting_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hosting type: `Agency` or `Classic`.",
			},
			"billing_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Billing status of the app.",
			},
		},
	}
}

func (r *AppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *providerData, got %T", req.ProviderData))
		return
	}
	r.client = pd.client
	r.organisationID = pd.organisationID
}

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creatReq := client.CreateAppRequest{
		Name:         plan.Name.ValueString(),
		Organisation: int(r.organisationID),
		AutoTeams:    plan.AutoTeams.ValueString(),
		AutoUpgrades: plan.AutoUpgrades.ValueString(),
		ExternalInfo: plan.ExternalInfo.ValueString(),
	}
	if !plan.CustomPackageID.IsNull() && !plan.CustomPackageID.IsUnknown() {
		v := int(plan.CustomPackageID.ValueInt64())
		creatReq.CustomPackage = &v
	}
	if !plan.CustomPackageName.IsNull() && !plan.CustomPackageName.IsUnknown() {
		v := plan.CustomPackageName.ValueString()
		creatReq.CustomPackageName = &v
	}

	app, err := r.client.CreateApp(ctx, creatReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating app", err.Error())
		return
	}

	flattenApp(app, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetApp(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading app", err.Error())
		return
	}

	flattenApp(app, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateAppRequest{
		Name:         plan.Name.ValueString(),
		Organisation: int(r.organisationID),
		AutoTeams:    plan.AutoTeams.ValueString(),
		AutoUpgrades: plan.AutoUpgrades.ValueString(),
		ExternalInfo: plan.ExternalInfo.ValueString(),
	}
	if !plan.CustomPackageID.IsNull() && !plan.CustomPackageID.IsUnknown() {
		v := int(plan.CustomPackageID.ValueInt64())
		updateReq.CustomPackage = &v
	}
	if !plan.CustomPackageName.IsNull() && !plan.CustomPackageName.IsUnknown() {
		v := plan.CustomPackageName.ValueString()
		updateReq.CustomPackageName = &v
	}

	err := r.client.UpdateApp(ctx, int(plan.ID.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating app", err.Error())
		return
	}

	// Re-read to pick up any server-side changes.
	app, err := r.client.GetApp(ctx, int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated app", err.Error())
		return
	}
	flattenApp(app, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApp(ctx, int(state.ID.ValueInt64()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting app", err.Error())
	}
}

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected a numeric app ID, got: %s", req.ID))
		return
	}

	app, err := r.client.GetApp(ctx, int(id))
	if err != nil {
		resp.Diagnostics.AddError("Error importing app", err.Error())
		return
	}

	var state AppResourceModel
	flattenApp(app, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// flattenApp copies API response fields into the Terraform model.
func flattenApp(app *client.App, model *AppResourceModel) {
	model.ID = types.Int64Value(int64(app.ID))
	model.Name = types.StringValue(app.Name)
	model.Status = types.StringValue(app.Status)
	model.StatusCategory = types.StringValue(app.StatusCategory)
	model.HostingType = types.StringValue(app.HostingType)
	model.BillingStatus = types.StringValue(app.BillingStatus)
	if app.CustomPackageName != "" {
		model.CustomPackageName = types.StringValue(app.CustomPackageName)
	} else {
		model.CustomPackageName = types.StringNull()
	}
	if app.Organisation != nil {
		model.OrganisationID = types.Int64Value(int64(app.Organisation.ID))
	}
}
