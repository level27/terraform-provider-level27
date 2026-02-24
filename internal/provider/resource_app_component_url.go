// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

var _ resource.Resource = &AppComponentURLResource{}
var _ resource.ResourceWithImportState = &AppComponentURLResource{}

type AppComponentURLResource struct {
	client *client.Client
}

type AppComponentURLResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	AppID          types.Int64  `tfsdk:"app_id"`
	ComponentID    types.Int64  `tfsdk:"component_id"`
	Content        types.String `tfsdk:"content"`
	SSLForce       types.Bool   `tfsdk:"ssl_force"`
	HandleDNS      types.Bool   `tfsdk:"handle_dns"`
	AutoSSLCert    types.Bool   `tfsdk:"auto_ssl_certificate"`
	Caching        types.Bool   `tfsdk:"caching"`
	Authentication types.Bool   `tfsdk:"authentication"`
	SSLCertificate types.Int64  `tfsdk:"ssl_certificate_id"`
	// Computed
	Status types.String `tfsdk:"status"`
	HTTPS  types.Bool   `tfsdk:"https"`
}

func NewAppComponentURLResource() resource.Resource {
	return &AppComponentURLResource{}
}

func (r *AppComponentURLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_component_url"
}

func (r *AppComponentURLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a **URL** linked to a Level27 App Component (PHP, ASP, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"app_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the parent app.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"component_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the parent component.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The hostname/URL, e.g. `www.example.com`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssl_force": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Force HTTPS redirect (default: `true`).",
			},
			"handle_dns": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Automatically create DNS records for this URL (default: `false`).",
			},
			"auto_ssl_certificate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Automatically create a Let's Encrypt certificate. Requires `handle_dns = true` (default: `false`).",
			},
			"caching": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable caching for this URL (default: `false`).",
			},
			"authentication": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable HTTP basic authentication (default: `false`).",
			},
			"ssl_certificate_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "ID of an existing SSL certificate to attach.",
			},
			// Computed
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current status of the URL.",
			},
			"https": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether HTTPS is active on this URL.",
			},
		},
	}
}

func (r *AppComponentURLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppComponentURLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppComponentURLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := r.client.CreateAppComponentURL(ctx, int(plan.AppID.ValueInt64()), int(plan.ComponentID.ValueInt64()),
		client.CreateAppComponentURLRequest{
			Content:            plan.Content.ValueString(),
			SSLForce:           plan.SSLForce.ValueBool(),
			HandleDNS:          plan.HandleDNS.ValueBool(),
			AutoSSLCertificate: plan.AutoSSLCert.ValueBool(),
			Caching:            plan.Caching.ValueBool(),
			Authentication:     plan.Authentication.ValueBool(),
			SSLCertificate:     int(plan.SSLCertificate.ValueInt64()),
		})
	if err != nil {
		resp.Diagnostics.AddError("Error creating URL", err.Error())
		return
	}

	flattenURL(u, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppComponentURLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppComponentURLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := r.client.GetAppComponentURL(ctx, int(state.AppID.ValueInt64()), int(state.ComponentID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading URL", err.Error())
		return
	}

	flattenURL(u, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppComponentURLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppComponentURLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAppComponentURL(ctx, int(plan.AppID.ValueInt64()), int(plan.ComponentID.ValueInt64()), int(plan.ID.ValueInt64()),
		client.UpdateAppComponentURLRequest{
			Content:        plan.Content.ValueString(),
			SSLForce:       plan.SSLForce.ValueBool(),
			HandleDNS:      plan.HandleDNS.ValueBool(),
			Caching:        plan.Caching.ValueBool(),
			Authentication: plan.Authentication.ValueBool(),
			SSLCertificate: int(plan.SSLCertificate.ValueInt64()),
		})
	if err != nil {
		resp.Diagnostics.AddError("Error updating URL", err.Error())
		return
	}

	u, err := r.client.GetAppComponentURL(ctx, int(plan.AppID.ValueInt64()), int(plan.ComponentID.ValueInt64()), int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated URL", err.Error())
		return
	}
	flattenURL(u, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppComponentURLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppComponentURLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAppComponentURL(ctx, int(state.AppID.ValueInt64()), int(state.ComponentID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting URL", err.Error())
	}
}

func (r *AppComponentURLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Format: "<app_id>/<component_id>/<url_id>"
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: <app_id>/<component_id>/<url_id>")
		return
	}
	appID, e1 := strconv.ParseInt(parts[0], 10, 64)
	compID, e2 := strconv.ParseInt(parts[1], 10, 64)
	urlID, e3 := strconv.ParseInt(parts[2], 10, 64)
	if e1 != nil || e2 != nil || e3 != nil {
		resp.Diagnostics.AddError("Invalid import ID", "All IDs must be integers")
		return
	}

	u, err := r.client.GetAppComponentURL(ctx, int(appID), int(compID), int(urlID))
	if err != nil {
		resp.Diagnostics.AddError("Error importing URL", err.Error())
		return
	}

	var state AppComponentURLResourceModel
	state.AppID = types.Int64Value(appID)
	state.ComponentID = types.Int64Value(compID)
	flattenURL(u, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenURL(u *client.AppComponentURL, model *AppComponentURLResourceModel) {
	model.ID = types.Int64Value(int64(u.ID))
	model.Content = types.StringValue(u.Content)
	model.SSLForce = types.BoolValue(u.SSLForce)
	model.Caching = types.BoolValue(u.Caching)
	model.Authentication = types.BoolValue(u.Authentication)
	model.Status = types.StringValue(u.Status)
	model.HTTPS = types.BoolValue(u.HTTPS)
}
