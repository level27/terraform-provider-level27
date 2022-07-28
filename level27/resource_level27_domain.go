package level27

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/l27-go"
)

type resourceDomainType struct{}

func (r resourceDomainType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	attrNameserver := tfsdk.Attribute{
		Type:     types.StringType,
		Optional: true,
	}

	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
			"organisation": {
				Computed:      true,
				Optional:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
			"name": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown(), tfsdk.RequiresReplace()},
			},
			"full_name": {
				Computed:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
			},
			"extension": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			},
			"contact_licensee": {
				Required:   true,
				Type:       types.StringType,
				Validators: []tfsdk.AttributeValidator{validateID()},
			},
			"contact_onsite": {
				Type:       types.StringType,
				Optional:   true,
				Validators: []tfsdk.AttributeValidator{validateID()},
			},
			"handle_dns": {
				Type:          types.BoolType,
				Computed:      true,
				Optional:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{DefaultValue(types.Bool{Value: true})},
			},
			"name_server_1":      attrNameserver,
			"name_server_2":      attrNameserver,
			"name_server_3":      attrNameserver,
			"name_server_4":      attrNameserver,
			"name_server_1_ip":   attrNameserver,
			"name_server_2_ip":   attrNameserver,
			"name_server_3_ip":   attrNameserver,
			"name_server_4_ip":   attrNameserver,
			"name_server_1_ipv6": attrNameserver,
			"name_server_2_ipv6": attrNameserver,
			"name_server_3_ipv6": attrNameserver,
			"name_server_4_ipv6": attrNameserver,
			"ttl": {
				Type:     types.Int64Type,
				Computed: true,
				Optional: true,
				// Default is 8hrs, same as lvl
				PlanModifiers: tfsdk.AttributePlanModifiers{DefaultValue(types.Int64{Value: 28800})},
			},
			"teams": {
				Type:          types.ListType{}.WithElementType(types.StringType),
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{DefaultValue(types.List{ElemType: types.StringType})},
			},
		},
	}, nil
}

func (r resourceDomainType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDomain{
		p: p.(*provider),
	}, nil
}

type resourceDomain struct {
	p *provider
}

type resourceDomainModel struct {
	ID              types.String `tfsdk:"id"`
	Organisation    types.String `tfsdk:"organisation"`
	Name            types.String `tfsdk:"name"`
	FullName        types.String `tfsdk:"full_name"`
	Extension       types.String `tfsdk:"extension"`
	ContactLicensee types.String `tfsdk:"contact_licensee"`
	ContactOnsite   types.String `tfsdk:"contact_onsite"`
	HandleDns       types.Bool   `tfsdk:"handle_dns"`
	NameServer1     types.String `tfsdk:"name_server_1"`
	NameServer2     types.String `tfsdk:"name_server_2"`
	NameServer3     types.String `tfsdk:"name_server_3"`
	NameServer4     types.String `tfsdk:"name_server_4"`
	NameServer1Ip   types.String `tfsdk:"name_server_1_ip"`
	NameServer2Ip   types.String `tfsdk:"name_server_2_ip"`
	NameServer3Ip   types.String `tfsdk:"name_server_3_ip"`
	NameServer4Ip   types.String `tfsdk:"name_server_4_ip"`
	NameServer1Ipv6 types.String `tfsdk:"name_server_1_ipv6"`
	NameServer2Ipv6 types.String `tfsdk:"name_server_2_ipv6"`
	NameServer3Ipv6 types.String `tfsdk:"name_server_3_ipv6"`
	NameServer4Ipv6 types.String `tfsdk:"name_server_4_ipv6"`
	Teams           types.List   `tfsdk:"teams"`
	Ttl             types.Int64  `tfsdk:"ttl"`
}

func (r resourceDomain) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan resourceDomainModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := organisationOrDefault(r.p, plan.Organisation, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var request l27.DomainRequest
	domainRequestFrom(r, orgID, plan, &resp.Diagnostics, &request)

	domain, err := r.p.Client.DomainCreate(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating domain", "API request failed:\n"+err.Error())
		return
	}

	result := domainModelFrom(domain)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDomain) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state resourceDomainModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, ok := parseID(state.ID.Value, &resp.Diagnostics)
	if !ok {
		return
	}

	err := r.p.Client.DomainDelete(appID)
	if err != nil {
		resp.Diagnostics.AddError("Error destroying domain", "API request failed:\n"+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceDomain) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state resourceDomainModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID, _ := strconv.Atoi(state.ID.Value)

	domain, err := r.p.Client.Domain(domainID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain", "API request failed:\n"+err.Error())
		return
	}

	result := domainModelFrom(domain)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDomain) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan resourceDomainModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state resourceDomainModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID, _ := strconv.Atoi(state.ID.Value)

	var request l27.DomainRequest
	domainRequestFrom(r, tfDStringId(state.Organisation), plan, &resp.Diagnostics, &request)

	err := r.p.Client.DomainUpdatePut(domainID, request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating domain", "API request failed:\n"+err.Error())
		return
	}

	// NOTE: we set state = plan her. This relies on UseStateForUnknown() so that plan already contains merged state.
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (resourceDomain) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDomainExtensions(p *provider) (map[string]int, error) {
	providers, err := p.Client.Extension()
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, provider := range providers {
		for _, domainType := range provider.Domaintypes {
			m[domainType.Extension] = domainType.ID
		}
	}

	return m, nil
}

func domainModelFrom(domain l27.Domain) resourceDomainModel {
	result := resourceDomainModel{
		ID:              tfStringId(domain.ID),
		Name:            tfString(domain.Name),
		FullName:        tfString(domain.Fullname),
		Extension:       tfString(domain.Domaintype.Extension),
		ContactLicensee: tfStringId(domain.DomaincontactLicensee.ID),
		HandleDns:       types.Bool{Value: domain.DNSIsHandled},
		NameServer1:     tfStringP(domain.Nameserver1),
		NameServer2:     tfStringP(domain.Nameserver2),
		NameServer3:     tfStringP(domain.Nameserver3),
		NameServer4:     tfStringP(domain.Nameserver4),
		NameServer1Ip:   tfStringP(domain.NameserverIP1),
		NameServer2Ip:   tfStringP(domain.NameserverIP2),
		NameServer3Ip:   tfStringP(domain.NameserverIP3),
		NameServer4Ip:   tfStringP(domain.NameserverIP4),
		NameServer1Ipv6: tfStringP(domain.NameserverIpv61),
		NameServer2Ipv6: tfStringP(domain.NameserverIpv62),
		NameServer3Ipv6: tfStringP(domain.NameserverIpv63),
		NameServer4Ipv6: tfStringP(domain.NameserverIpv64),
		Ttl:             types.Int64{Value: int64(domain.TTL)},
		ContactOnsite:   types.String{Null: true},
		Organisation:    tfStringId(domain.Organisation.ID),
	}

	if domain.DomaincontactOnsite != nil {
		result.ContactOnsite = tfStringId(domain.DomaincontactOnsite.ID)
	}

	result.Teams = types.List{ElemType: types.StringType}

	for _, v := range domain.Teams {
		result.Teams.Elems = append(result.Teams.Elems, tfStringId(v.ID))
	}

	return result
}

func domainRequestFrom(r resourceDomain, orgID int, plan resourceDomainModel, diags *diag.Diagnostics, request *l27.DomainRequest) {
	contactLicenseeID, _ := strconv.Atoi(plan.ContactLicensee.Value)
	var contactOnsiteID *int
	if !plan.ContactOnsite.Null && !plan.ContactOnsite.Unknown {
		id, _ := strconv.Atoi(plan.ContactOnsite.Value)
		contactOnsiteID = &id
	}

	extensions, err := getDomainExtensions(r.p)
	if err != nil {
		diags.AddError("Error resolving domain type", "API request failed:\n"+err.Error())
		return
	}

	domainType, ok := extensions[plan.Extension.Value]
	if !ok {
		diags.AddError(
			"Invalid domain extension",
			fmt.Sprintf("Unknown domain extension: '%s'", plan.Extension.Value))
	}

	*request = l27.DomainRequest{
		Organisation:          orgID,
		Name:                  plan.Name.Value,
		NameServer1:           tfDStringP(plan.NameServer1),
		NameServer2:           tfDStringP(plan.NameServer2),
		NameServer3:           tfDStringP(plan.NameServer3),
		NameServer4:           tfDStringP(plan.NameServer4),
		NameServer1Ip:         tfDStringP(plan.NameServer1Ip),
		NameServer2Ip:         tfDStringP(plan.NameServer2Ip),
		NameServer3Ip:         tfDStringP(plan.NameServer3Ip),
		NameServer4Ip:         tfDStringP(plan.NameServer4Ip),
		NameServer1Ipv6:       tfDStringP(plan.NameServer1Ipv6),
		NameServer2Ipv6:       tfDStringP(plan.NameServer2Ipv6),
		NameServer3Ipv6:       tfDStringP(plan.NameServer3Ipv6),
		NameServer4Ipv6:       tfDStringP(plan.NameServer4Ipv6),
		Handledns:             plan.HandleDns.Value,
		Domaincontactlicensee: contactLicenseeID,
		DomainContactOnSite:   contactOnsiteID,
		TTL:                   int(plan.Ttl.Value),
		Domaintype:            domainType,
		AutoTeams:             formatTeams(plan.Teams),
	}

	request.Action = "create"
}
