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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

var _ resource.Resource = &SSLCertificateResource{}
var _ resource.ResourceWithImportState = &SSLCertificateResource{}

type SSLCertificateResource struct {
	client *client.Client
}

type SSLCertificateResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	AppID                  types.Int64  `tfsdk:"app_id"`
	Name                   types.String `tfsdk:"name"`
	SSLType                types.String `tfsdk:"ssl_type"`
	AutoSSLCertificateURLs types.String `tfsdk:"auto_ssl_certificate_urls"`
	SSLKey                 types.String `tfsdk:"ssl_key"`
	SSLCrt                 types.String `tfsdk:"ssl_crt"`
	SSLCabundle            types.String `tfsdk:"ssl_cabundle"`
	AutoURLLink            types.Bool   `tfsdk:"auto_url_link"`
	SSLForce               types.Bool   `tfsdk:"ssl_force"`
	// Computed
	SSLStatus types.String `tfsdk:"ssl_status"`
	Status    types.String `tfsdk:"status"`
}

func NewSSLCertificateResource() resource.Resource {
	return &SSLCertificateResource{}
}

func (r *SSLCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssl_certificate"
}

func (r *SSLCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an **SSL Certificate** for a Level27 App. Supports Let's Encrypt (`letsencrypt`), Xolphin (`xolphin`), and custom (`own`) certificates.",
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
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the certificate.",
			},
			"ssl_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Certificate type: `letsencrypt`, `xolphin`, or `own`.",
				Validators: []validator.String{
					stringvalidator.OneOf("letsencrypt", "xolphin", "own"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"auto_ssl_certificate_urls": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comma-separated list of URLs for Let's Encrypt certificate (required when `ssl_type = letsencrypt`).",
			},
			"ssl_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "PEM-encoded private key (required when `ssl_type = own`).",
			},
			"ssl_crt": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "PEM-encoded certificate (required when `ssl_type = own`).",
			},
			"ssl_cabundle": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "PEM-encoded CA bundle (required when `ssl_type = own`).",
			},
			"auto_url_link": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Automatically link matching app URLs to this certificate (default: `false`).",
			},
			"ssl_force": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Force HTTPS on all linked URLs (default: `true`).",
			},
			// Computed
			"ssl_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the SSL certificate itself (e.g. `ok`, `pending`).",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Platform status of the SSL certificate resource.",
			},
		},
	}
}

func (r *SSLCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *providerData, got %T", req.ProviderData))
		return
	}
	r.client = pd.client
}

func (r *SSLCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSLCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.CreateSSLCertificate(ctx, int(plan.AppID.ValueInt64()), client.CreateSSLCertificateRequest{
		Name:                   plan.Name.ValueString(),
		SSLType:                plan.SSLType.ValueString(),
		AutoSSLCertificateURLs: plan.AutoSSLCertificateURLs.ValueString(),
		SSLKey:                 plan.SSLKey.ValueString(),
		SSLCrt:                 plan.SSLCrt.ValueString(),
		SSLCabundle:            plan.SSLCabundle.ValueString(),
		AutoURLLink:            plan.AutoURLLink.ValueBool(),
		SSLForce:               plan.SSLForce.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSL certificate", err.Error())
		return
	}

	flattenSSLCert(cert, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSLCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSLCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.GetSSLCertificate(ctx, int(state.AppID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SSL certificate", err.Error())
		return
	}

	flattenSSLCert(cert, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SSLCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SSLCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateSSLCertificate(ctx, int(plan.AppID.ValueInt64()), int(plan.ID.ValueInt64()),
		map[string]string{"name": plan.Name.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error updating SSL certificate", err.Error())
		return
	}

	cert, err := r.client.GetSSLCertificate(ctx, int(plan.AppID.ValueInt64()), int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated SSL certificate", err.Error())
		return
	}
	flattenSSLCert(cert, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSLCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSLCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSSLCertificate(ctx, int(state.AppID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting SSL certificate", err.Error())
	}
}

func (r *SSLCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Format: "<app_id>/<cert_id>"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: <app_id>/<cert_id>")
		return
	}
	appID, e1 := strconv.ParseInt(parts[0], 10, 64)
	certID, e2 := strconv.ParseInt(parts[1], 10, 64)
	if e1 != nil || e2 != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Both app_id and cert_id must be integers")
		return
	}

	cert, err := r.client.GetSSLCertificate(ctx, int(appID), int(certID))
	if err != nil {
		resp.Diagnostics.AddError("Error importing SSL certificate", err.Error())
		return
	}

	var state SSLCertificateResourceModel
	state.AppID = types.Int64Value(appID)
	flattenSSLCert(cert, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenSSLCert(cert *client.SSLCertificate, model *SSLCertificateResourceModel) {
	model.ID = types.Int64Value(int64(cert.ID))
	model.Name = types.StringValue(cert.Name)
	model.SSLType = types.StringValue(cert.SSLType)
	model.SSLStatus = types.StringValue(cert.SSLStatus)
	model.Status = types.StringValue(cert.Status)
}
