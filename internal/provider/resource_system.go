// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

var _ resource.Resource = &SystemResource{}
var _ resource.ResourceWithImportState = &SystemResource{}

// SystemResource manages a Level27 system (server).
type SystemResource struct {
	client *client.Client
}

// SystemResourceModel maps the resource schema data.
type SystemResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	CustomerFqdn           types.String `tfsdk:"customer_fqdn"`
	Type                   types.String `tfsdk:"type"`
	OrganisationID         types.Int64  `tfsdk:"organisation_id"`
	SystemimageID          types.Int64  `tfsdk:"systemimage_id"`
	SystemproviderConfigID types.Int64  `tfsdk:"systemprovider_configuration_id"`
	ZoneID                 types.Int64  `tfsdk:"zone_id"`
	CPU                    types.Int64  `tfsdk:"cpu"`
	Disk                   types.Int64  `tfsdk:"disk"`
	Memory                 types.Int64  `tfsdk:"memory"`
	ManagementType         types.String `tfsdk:"management_type"`
	Period                 types.Int64  `tfsdk:"period"`
	ExternalInfo           types.String `tfsdk:"external_info"`
	ParentsystemID         types.Int64  `tfsdk:"parentsystem_id"`
	Networks               types.Map    `tfsdk:"networks"`
	AutoInstall            types.Bool   `tfsdk:"auto_install"`
	// Computed
	Status         types.String `tfsdk:"status"`
	StatusCategory types.String `tfsdk:"status_category"`
}

func NewSystemResource() resource.Resource {
	return &SystemResource{}
}

func (r *SystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (r *SystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Level27 **System** (virtual server). " +
			"The `type`, `organisation_id`, `systemimage_id`, `systemprovider_configuration_id`, `zone_id`, and `parentsystem_id` are immutable after creation — changing them forces a replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the system.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Hostname (FQDN) of the system, e.g. `myserver.example.com`.",
			},
			"customer_fqdn": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Customer-facing FQDN. Defaults to `name` when omitted.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "System type, e.g. `kvmguest`, `vmware`, `baremetal`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organisation_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the organisation that owns this system.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"systemimage_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the system image to use for provisioning.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"systemprovider_configuration_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the system provider configuration (hypervisor / datacenter profile).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"zone_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the zone (network/datacenter zone) to deploy the system in.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"cpu": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Number of virtual CPUs.",
			},
			"disk": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Disk size in GB.",
			},
			"memory": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Memory in GB.",
			},
			"management_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Management level. Common values: `infra_plus`, `basic`, `enterprise`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"period": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Billing period in months (default: `1`).",
			},
			"external_info": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "External reference / billing info.",
			},
			"parentsystem_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "ID of the parent (host) system for nested virtualisation. Assigned automatically by the API.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplaceIfConfigured()},
			},
			// Computed
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current provisioning status of the system (e.g. `ok`, `creating`).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status_category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status category: `green`, `yellow`, `red`, or `grey`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"networks": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				MarkdownDescription: "Map of network-name → IPv4 address to assign. " +
					"Use `\"auto\"` to have the provider pick a free address from `GET /networks/{id}/locate`. " +
					"Use an empty string `\"\"` for no IP on that interface. " +
					"Once `\"auto\"` is applied, the state preserves `\"auto\"` so subsequent plans are stable.",
			},
			"auto_install": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				MarkdownDescription: "When `true` (default), triggers OS installation (`autoInstall`) immediately after " +
					"the system is allocated and networks are configured, by calling `POST /systems/{id}/actions` " +
					"with `autoInstall`. The provider waits until installation completes. " +
					"Set to `false` to skip automatic installation.",
			},
		},
	}
}

func (r *SystemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SystemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := client.CreateSystemRequest{
		Name:                        plan.Name.ValueString(),
		Type:                        plan.Type.ValueString(),
		Organisation:                int(plan.OrganisationID.ValueInt64()),
		Systemimage:                 int(plan.SystemimageID.ValueInt64()),
		SystemproviderConfiguration: int(plan.SystemproviderConfigID.ValueInt64()),
		Zone:                        int(plan.ZoneID.ValueInt64()),
		CPU:                         int(plan.CPU.ValueInt64()),
		Disk:                        int(plan.Disk.ValueInt64()),
		Memory:                      int(plan.Memory.ValueInt64()),
	}

	if !plan.CustomerFqdn.IsNull() && !plan.CustomerFqdn.IsUnknown() {
		body.CustomerFqdn = plan.CustomerFqdn.ValueString()
	}
	if !plan.ManagementType.IsNull() && !plan.ManagementType.IsUnknown() {
		body.ManagementType = plan.ManagementType.ValueString()
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		body.Period = int(plan.Period.ValueInt64())
	}
	if !plan.ExternalInfo.IsNull() && !plan.ExternalInfo.IsUnknown() {
		v := plan.ExternalInfo.ValueString()
		body.ExternalInfo = &v
	}
	if !plan.ParentsystemID.IsNull() && !plan.ParentsystemID.IsUnknown() {
		v := int(plan.ParentsystemID.ValueInt64())
		body.Parentsystem = &v
	}

	sys, err := r.client.CreateSystem(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating system", err.Error())
		return
	}

	// Persist the ID immediately so Terraform can track it even if the wait fails.
	plan.ID = types.Int64Value(int64(sys.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	// Wait for provisioning to complete.
	sys, err = r.client.WaitForSystemStatus(ctx, sys.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for system creation", err.Error())
		return
	}

	flattenSystem(sys, &plan)

	// Attach requested networks.
	if !plan.Networks.IsNull() && !plan.Networks.IsUnknown() {
		if diags := r.syncNetworks(ctx, int(plan.ID.ValueInt64()), plan.Networks); diags != nil {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// Read back actual networks into state.
	if diags := r.readNetworks(ctx, int(plan.ID.ValueInt64()), &plan, plan.Networks); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Trigger OS installation if requested and the system is in "stopped" state.
	if (plan.AutoInstall.IsNull() || plan.AutoInstall.ValueBool()) && sys.Status == "stopped" {
		if err := r.client.PerformSystemAction(ctx, int(plan.ID.ValueInt64()), "autoInstall"); err != nil {
			resp.Diagnostics.AddError("Error triggering autoInstall", err.Error())
			return
		}
		sys, err = r.client.WaitForSystemStatus(ctx, int(plan.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for system installation", err.Error())
			return
		}
		flattenSystem(sys, &plan)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SystemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sys, err := r.client.GetSystem(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading system", err.Error())
		return
	}

	flattenSystem(sys, &state)

	if diags := r.readNetworks(ctx, int(state.ID.ValueInt64()), &state, state.Networks); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SystemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := client.UpdateSystemRequest{
		Name:                        plan.Name.ValueString(),
		Type:                        plan.Type.ValueString(),
		Organisation:                int(plan.OrganisationID.ValueInt64()),
		Systemimage:                 int(plan.SystemimageID.ValueInt64()),
		SystemproviderConfiguration: int(plan.SystemproviderConfigID.ValueInt64()),
		Zone:                        int(plan.ZoneID.ValueInt64()),
		CPU:                         int(plan.CPU.ValueInt64()),
		Disk:                        int(plan.Disk.ValueInt64()),
		Memory:                      int(plan.Memory.ValueInt64()),
	}
	if !plan.CustomerFqdn.IsNull() && !plan.CustomerFqdn.IsUnknown() {
		body.CustomerFqdn = plan.CustomerFqdn.ValueString()
	}
	if !plan.ManagementType.IsNull() && !plan.ManagementType.IsUnknown() {
		body.ManagementType = plan.ManagementType.ValueString()
	}
	if !plan.ExternalInfo.IsNull() {
		v := plan.ExternalInfo.ValueString()
		body.ExternalInfo = &v
	}

	err := r.client.UpdateSystem(ctx, int(plan.ID.ValueInt64()), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating system", err.Error())
		return
	}

	// Wait for the update to complete.
	sys, err := r.client.WaitForSystemStatus(ctx, int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for system update", err.Error())
		return
	}

	flattenSystem(sys, &plan)
	// period is not returned by the API; plan.Period already has the correct value from the plan/default.

	// Sync networks: diff vs live API.
	if diags := r.syncNetworks(ctx, int(plan.ID.ValueInt64()), plan.Networks); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := r.readNetworks(ctx, int(plan.ID.ValueInt64()), &plan, plan.Networks); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SystemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystem(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting system", err.Error())
		return
	}

	// Wait until the system is fully removed.
	waitCtx := ctx
	for {
		_, err := r.client.GetSystem(waitCtx, int(state.ID.ValueInt64()))
		if err != nil {
			if client.IsNotFound(err) {
				return
			}
			resp.Diagnostics.AddError("Error waiting for system deletion", err.Error())
			return
		}
		select {
		case <-waitCtx.Done():
			resp.Diagnostics.AddError("Timeout waiting for system deletion", waitCtx.Err().Error())
			return
		default:
		}
	}
}

func (r *SystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected a numeric system ID, got: %s", req.ID))
		return
	}

	sys, err := r.client.GetSystem(ctx, int(id))
	if err != nil {
		resp.Diagnostics.AddError("Error importing system", err.Error())
		return
	}

	var state SystemResourceModel
	flattenSystem(sys, &state)
	if diags := r.readNetworks(ctx, int(id), &state, types.MapNull(types.StringType)); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// flattenSystem copies API response fields into the Terraform model.
func flattenSystem(sys *client.System, m *SystemResourceModel) {
	m.ID = types.Int64Value(int64(sys.ID))
	m.Name = types.StringValue(sys.Name)
	m.Type = types.StringValue(sys.Type)
	m.Status = types.StringValue(sys.Status)
	m.StatusCategory = types.StringValue(sys.StatusCategory)
	m.CPU = types.Int64Value(int64(sys.CPU))
	m.Disk = types.Int64Value(int64(sys.Disk))
	m.Memory = types.Int64Value(int64(sys.Memory))

	if sys.CustomerFqdn != "" {
		m.CustomerFqdn = types.StringValue(sys.CustomerFqdn)
	} else {
		m.CustomerFqdn = types.StringNull()
	}
	if sys.ManagementType != "" {
		m.ManagementType = types.StringValue(sys.ManagementType)
	} else {
		m.ManagementType = types.StringNull()
	}
	if sys.ExternalInfo != "" {
		m.ExternalInfo = types.StringValue(sys.ExternalInfo)
	} else {
		m.ExternalInfo = types.StringNull()
	}
	if sys.Organisation != nil {
		m.OrganisationID = types.Int64Value(int64(sys.Organisation.ID))
	}
	if sys.Zone != nil {
		m.ZoneID = types.Int64Value(int64(sys.Zone.ID))
	}
	if sys.Systemimage != nil {
		m.SystemimageID = types.Int64Value(int64(sys.Systemimage.ID))
	}
	if sys.SystemproviderConfiguration != nil {
		m.SystemproviderConfigID = types.Int64Value(int64(sys.SystemproviderConfiguration.ID))
	}
	if sys.Parentsystem != nil {
		m.ParentsystemID = types.Int64Value(int64(sys.Parentsystem.ID))
	} else {
		m.ParentsystemID = types.Int64Null()
	}
}

// readNetworks fetches current system networks (with IPs) and stores them in m.Networks
// as a map of network-name → IPv4 (empty string if no IP assigned).
// If desired specifies "auto" for a network that already has an IP, "auto" is preserved
// in state so subsequent plans don't show a spurious diff.
func (r *SystemResource) readNetworks(ctx context.Context, systemID int, m *SystemResourceModel, desired types.Map) diag.Diagnostics {
	nets, err := r.client.ListSystemNetworks(ctx, systemID)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Error reading system networks", err.Error())}
	}
	desiredElems := desired.Elements()
	result := make(map[string]attr.Value, len(nets))
	for _, shn := range nets {
		ips, err := r.client.ListSystemNetworkIPs(ctx, systemID, shn.ID)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error reading network IPs",
				fmt.Sprintf("network %q: %s", shn.Network.Name, err),
			)}
		}
		ipv4 := ""
		if len(ips) > 0 {
			ipv4 = ips[0].IPv4
		}
		// If the user configured "auto" for this network and an IP is already assigned,
		// preserve "auto" in state to prevent perpetual plan diffs.
		if ipv4 != "" {
			if dv, ok := desiredElems[shn.Network.Name]; ok {
				if ds, ok := dv.(types.String); ok && !ds.IsNull() && !ds.IsUnknown() && ds.ValueString() == "auto" {
					ipv4 = "auto"
				}
			}
		}
		result[shn.Network.Name] = types.StringValue(ipv4)
	}
	netMap, diags := types.MapValue(types.StringType, result)
	if diags.HasError() {
		return diags
	}
	m.Networks = netMap
	return nil
}

// syncNetworks reconciles the desired networks map (name → ipv4) against live API state.
// "auto" as an IPv4 value triggers automatic address selection via the locate endpoint.
// An empty string means "attach network but assign no IP".
func (r *SystemResource) syncNetworks(ctx context.Context, systemID int, desired types.Map) diag.Diagnostics {
	// Parse desired entries.
	type wantNet struct {
		ipv4 string // "" = no IP, "auto" = locate, else specific
	}
	wants := map[string]wantNet{}
	for name, v := range desired.Elements() {
		s, ok := v.(types.String)
		ipv4 := ""
		if ok && !s.IsNull() && !s.IsUnknown() {
			ipv4 = s.ValueString()
		}
		wants[name] = wantNet{ipv4: ipv4}
	}

	// Get current attachments from the API.
	currentNets, err := r.client.ListSystemNetworks(ctx, systemID)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Error listing system networks", err.Error())}
	}
	currentByName := map[string]client.SystemHasNetwork{}
	for _, shn := range currentNets {
		currentByName[shn.Network.Name] = shn
	}

	// Remove networks no longer wanted.
	for name, shn := range currentByName {
		if _, wanted := wants[name]; !wanted {
			if err := r.client.RemoveSystemNetwork(ctx, systemID, shn.ID); err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error removing network from system",
					fmt.Sprintf("network %q: %s", name, err),
				)}
			}
			delete(currentByName, name)
		}
	}

	// Add or update desired networks.
	for name, want := range wants {
		shn, exists := currentByName[name]

		// Attach network if not yet present.
		if !exists {
			net, err := r.client.GetNetworkByName(ctx, name)
			if err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error resolving network name",
					fmt.Sprintf("%q: %s", name, err),
				)}
			}
			added, err := r.client.AddSystemNetwork(ctx, systemID, net.ID)
			if err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error adding network to system",
					fmt.Sprintf("%q: %s", name, err),
				)}
			}
			shn = *added
		}

		if want.ipv4 == "" {
			// Enforce no IP: remove any that exist on this interface.
			ips, err := r.client.ListSystemNetworkIPs(ctx, systemID, shn.ID)
			if err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error reading network IPs",
					fmt.Sprintf("network %q: %s", name, err),
				)}
			}
			for _, ip := range ips {
				if err := r.client.RemoveSystemNetworkIP(ctx, systemID, shn.ID, ip.ID); err != nil {
					return diag.Diagnostics{diag.NewErrorDiagnostic(
						"Error removing IP",
						fmt.Sprintf("%q from network %q: %s", ip.IPv4, name, err),
					)}
				}
			}
			continue
		}

		// Read current IPs for this interface.
		ips, err := r.client.ListSystemNetworkIPs(ctx, systemID, shn.ID)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error reading network IPs",
				fmt.Sprintf("network %q: %s", name, err),
			)}
		}

		// For "auto": only assign an IP if none is assigned yet.
		// Once an IP has been located and assigned, never replace it automatically.
		if want.ipv4 == "auto" {
			if len(ips) > 0 {
				continue // already has an IP — leave it alone
			}
			located, err := r.client.LocateNetworkIP(ctx, shn.Network.ID)
			if err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error locating free IP",
					fmt.Sprintf("network %q: %s", name, err),
				)}
			}
			if err := r.client.AddSystemNetworkIP(ctx, systemID, shn.ID, located); err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error assigning auto IP",
					fmt.Sprintf("%q to network %q: %s", located, name, err),
				)}
			}
			continue
		}

		// Explicit IP: remove any existing IP that differs, then assign.
		targetIPv4 := want.ipv4
		alreadyAssigned := false
		for _, ip := range ips {
			if ip.IPv4 == targetIPv4 {
				alreadyAssigned = true
			} else {
				if err := r.client.RemoveSystemNetworkIP(ctx, systemID, shn.ID, ip.ID); err != nil {
					return diag.Diagnostics{diag.NewErrorDiagnostic(
						"Error removing outdated IP",
						fmt.Sprintf("%q from network %q: %s", ip.IPv4, name, err),
					)}
				}
			}
		}
		if !alreadyAssigned {
			if err := r.client.AddSystemNetworkIP(ctx, systemID, shn.ID, targetIPv4); err != nil {
				return diag.Diagnostics{diag.NewErrorDiagnostic(
					"Error assigning IP",
					fmt.Sprintf("%q to network %q: %s", targetIPv4, name, err),
				)}
			}
		}
	}
	return nil
}
