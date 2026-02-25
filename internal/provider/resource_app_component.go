// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

var _ resource.Resource = &AppComponentResource{}
var _ resource.ResourceWithImportState = &AppComponentResource{}

// AppComponentResource defines the resource implementation.
type AppComponentResource struct {
	client *client.Client
}

// AppComponentResourceModel maps the resource schema data.
type AppComponentResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	AppID            types.Int64  `tfsdk:"app_id"`
	Name             types.String `tfsdk:"name"`
	AppComponentType types.String `tfsdk:"appcomponenttype"`
	System           types.Int64  `tfsdk:"system"`
	Systemgroup      types.Int64  `tfsdk:"systemgroup"`
	LimitGroup       types.String `tfsdk:"limit_group"`
	// Type-specific
	Version     types.String `tfsdk:"version"`
	Path        types.String `tfsdk:"path"`
	Pass        types.String `tfsdk:"pass"`
	SSHKeys     types.List   `tfsdk:"sshkeys"`
	ExtraConfig types.String `tfsdk:"extraconfig"`
	// Computed
	Status         types.String `tfsdk:"status"`
	StatusCategory types.String `tfsdk:"status_category"`
	Category       types.String `tfsdk:"category"`
}

func NewAppComponentResource() resource.Resource {
	return &AppComponentResource{}
}

func (r *AppComponentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_component"
}

func (r *AppComponentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	supportedTypes := []string{
		"mysql", "php", "keyvalue_redis", "sftp", "asp", "mssql", "solr",
		"memcached", "elasticsearch", "redis", "ruby", "url", "mongodb",
		"nodejs", "python", "postgresql", "dotnet", "rabbitmq", "varnish",
		"java", "mail", "wasp", "wmssql", "mailpit", "meilisearch",
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Level27 **App Component**. A component is one of the buildable blocks of an app: PHP, MySQL, SFTP, mail, etc.\n\n" +
			"> **Note:** `appcomponenttype`, `system`, and `systemgroup` are immutable after creation — changing them forces a replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the component.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"app_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the parent app.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the component.",
			},
			"appcomponenttype": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: fmt.Sprintf("Component type. One of: `%s`.", strings.Join(supportedTypes, "`, `")),
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf(supportedTypes...),
				},
			},
			"system": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "ID of the system to deploy on. Mutually exclusive with `systemgroup`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"systemgroup": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "ID of the system group to deploy on. Mutually exclusive with `system`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"limit_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Limit group for this component, e.g. `database` or `application`.",
			},
			"version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Version of the component runtime (e.g. PHP version `8.2`, MySQL version `8.4`). If omitted, the version is auto-detected from the system's installed cookbooks.",
			},
			"path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Relative web root path. Defaults to `/public_html` for PHP/ASP components.",
			},
			"pass": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for the component. Auto-generated if omitted.",
			},
			"sshkeys": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				MarkdownDescription: "List of SSH key IDs to deploy onto this component.",
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"extraconfig": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Extra configuration for PHP/ASP components (custom php.ini directives, etc.).",
			},
			// Computed
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current status of the component.",
			},
			"status_category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status category: `green`, `yellow`, `red`, or `grey`.",
			},
			"category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Component category assigned by Level27 (e.g. `Web-apps`, `Databases`).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}

	_ = booldefault.StaticBool(false) // imported, suppress lint
}

func (r *AppComponentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *AppComponentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppComponentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKeys := expandInt64List(ctx, plan.SSHKeys, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	version := plan.Version.ValueString()
	if version == "" && !plan.System.IsNull() && !plan.System.IsUnknown() {
		sys, sysErr := r.client.GetSystem(ctx, int(plan.System.ValueInt64()))
		if sysErr == nil {
			if v := sys.FindCookbookVersion(plan.AppComponentType.ValueString()); v != "" {
				version = v
			}
		}
	}

	createReq := client.CreateAppComponentRequest{
		Name:             plan.Name.ValueString(),
		AppComponentType: plan.AppComponentType.ValueString(),
		System:           int(plan.System.ValueInt64()),
		SystemGroup:      int(plan.Systemgroup.ValueInt64()),
		LimitGroup:       plan.LimitGroup.ValueString(),
		Version:          version,
		Path:             plan.Path.ValueString(),
		Pass:             plan.Pass.ValueString(),
		SSHKeys:          sshKeys,
		ExtraConfig:      plan.ExtraConfig.ValueString(),
	}

	comp, err := r.client.CreateAppComponent(ctx, int(plan.AppID.ValueInt64()), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating app component", err.Error())
		return
	}

	// Poll until the component leaves transitional status.
	comp, err = r.client.WaitForComponentStatus(ctx, int(plan.AppID.ValueInt64()), comp.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for app component", err.Error())
		return
	}

	flattenComponent(comp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppComponentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	comp, err := r.client.GetAppComponent(ctx, int(state.AppID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading app component", err.Error())
		return
	}

	flattenComponent(comp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppComponentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppComponentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateVersion := plan.Version.ValueString()
	if updateVersion == "" && !plan.System.IsNull() && !plan.System.IsUnknown() {
		sys, sysErr := r.client.GetSystem(ctx, int(plan.System.ValueInt64()))
		if sysErr == nil {
			if v := sys.FindCookbookVersion(plan.AppComponentType.ValueString()); v != "" {
				updateVersion = v
			}
		}
	}

	err := r.client.UpdateAppComponent(ctx, int(plan.AppID.ValueInt64()), int(plan.ID.ValueInt64()), client.UpdateAppComponentRequest{
		Name:        plan.Name.ValueString(),
		LimitGroup:  plan.LimitGroup.ValueString(),
		Version:     updateVersion,
		Path:        plan.Path.ValueString(),
		ExtraConfig: plan.ExtraConfig.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating app component", err.Error())
		return
	}

	// Poll until the component leaves transitional status.
	comp, err := r.client.WaitForComponentStatus(ctx, int(plan.AppID.ValueInt64()), int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for app component update", err.Error())
		return
	}
	flattenComponent(comp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppComponentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAppComponent(ctx, int(state.AppID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting app component", err.Error())
	}
}

func (r *AppComponentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "<app_id>/<component_id>"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: <app_id>/<component_id>")
		return
	}
	appID, err1 := strconv.ParseInt(parts[0], 10, 64)
	compID, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Both app_id and component_id must be integers")
		return
	}

	comp, err := r.client.GetAppComponent(ctx, int(appID), int(compID))
	if err != nil {
		resp.Diagnostics.AddError("Error importing app component", err.Error())
		return
	}

	var state AppComponentResourceModel
	state.AppID = types.Int64Value(appID)
	flattenComponent(comp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenComponent(comp *client.AppComponent, model *AppComponentResourceModel) {
	model.ID = types.Int64Value(int64(comp.ID))
	model.Name = types.StringValue(comp.Name)
	model.AppComponentType = types.StringValue(comp.AppComponentType)
	model.Status = types.StringValue(comp.Status)
	model.StatusCategory = types.StringValue(comp.StatusCategory)
	model.Category = types.StringValue(comp.Category)
}
