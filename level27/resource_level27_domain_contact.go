package level27

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/l27-go"
)

type resourceDomainContactType struct{}

func (r resourceDomainContactType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
			"type": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			},
			"first_name": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			},
			"last_name": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			},
			"organisation_name": {
				Required:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			},
			"street": {
				Required: true,
				Type:     types.StringType,
			},
			"house_number": {
				Required: true,
				Type:     types.StringType,
			},
			"zip": {
				Required: true,
				Type:     types.StringType,
			},
			"city": {
				Required: true,
				Type:     types.StringType,
			},
			"state": {
				Optional: true,
				Type:     types.StringType,
			},
			"phone": {
				Required: true,
				Type:     types.StringType,
			},
			"fax": {
				Optional: true,
				Type:     types.StringType,
			},
			"email": {
				Required:   true,
				Type:       types.StringType,
				Validators: []tfsdk.AttributeValidator{validateEmail()},
			},
			"tax_number": {
				Required: true,
				Type:     types.StringType,
			},
			"country": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (r resourceDomainContactType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDomainContact{
		p: p.(*provider),
	}, nil
}

type resourceDomainContact struct {
	p *provider
}

type resourceDomainContactModel struct {
	ID               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	FirstName        types.String `tfsdk:"first_name"`
	LastName         types.String `tfsdk:"last_name"`
	OrganisationName types.String `tfsdk:"organisation_name"`
	Street           types.String `tfsdk:"street"`
	HouseNumber      types.String `tfsdk:"house_number"`
	Zip              types.String `tfsdk:"zip"`
	City             types.String `tfsdk:"city"`
	State            types.String `tfsdk:"state"`
	Phone            types.String `tfsdk:"phone"`
	Fax              types.String `tfsdk:"fax"`
	Email            types.String `tfsdk:"email"`
	TaxNumber        types.String `tfsdk:"tax_number"`
	Country          types.String `tfsdk:"country"`
}

// Create implements tfsdk.Resource
func (r resourceDomainContact) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan resourceDomainContactModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := domainContactRequestFrom(plan)

	contact, err := r.p.Client.DomainContactCreate(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating domain contact", "API request failed:\n"+err.Error())
		return
	}

	result := domainContactModelFrom(contact)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

// Delete implements tfsdk.Resource
func (r resourceDomainContact) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state resourceDomainContactModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(state.ID.Value)

	err := r.p.Client.DomainContactDelete(orgID)
	if err != nil {
		resp.Diagnostics.AddError("Error destroying domain contact", "API request failed:\n"+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// Read implements tfsdk.Resource
func (r resourceDomainContact) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state resourceDomainContactModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainContactID, _ := strconv.Atoi(state.ID.Value)

	contact, err := r.p.Client.DomainContactGetSingle(domainContactID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain contact", "API request failed:\n"+err.Error())
		return
	}

	result := domainContactModelFrom(contact)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

// Update implements tfsdk.Resource
func (r resourceDomainContact) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Retrieve values from plan
	var plan resourceDomainContactModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceDomainContactModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := domainContactRequestFrom(plan)

	err := r.p.Client.DomainContactUpdate(tfDStringId(state.ID), request)

	if err != nil {
		resp.Diagnostics.AddError("Error updating domain contact", "API request failed:\n"+err.Error())
		return
	}
}

func (r resourceDomainContact) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func domainContactModelFrom(contact l27.DomainContact) resourceDomainContactModel {
	result := resourceDomainContactModel{
		ID:               tfStringId(contact.ID),
		Type:             tfString(contact.Type),
		FirstName:        tfString(contact.FirstName),
		LastName:         tfString(contact.LastName),
		OrganisationName: tfString(contact.OrganisationName),
		Street:           tfString(contact.Street),
		HouseNumber:      tfString(contact.HouseNumber),
		Zip:              tfString(contact.Zip),
		City:             tfString(contact.City),
		State:            tfStringP(contact.State),
		Phone:            tfString(contact.Phone),
		Email:            tfString(contact.Email),
		TaxNumber:        tfString(contact.TaxNumber),
		Country:          tfString(contact.Country.ID),
		Fax:              tfStringP(contact.Fax),
	}

	return result
}

func domainContactRequestFrom(plan resourceDomainContactModel) l27.DomainContactRequest {
	request := l27.DomainContactRequest{
		Type:             plan.Type.Value,
		FirstName:        plan.FirstName.Value,
		LastName:         plan.LastName.Value,
		OrganisationName: plan.OrganisationName.Value,
		Street:           plan.Street.Value,
		HouseNumber:      plan.HouseNumber.Value,
		Zip:              plan.Zip.Value,
		City:             plan.City.Value,
		State:            tfDStringP(plan.State),
		Phone:            plan.Phone.Value,
		Fax:              tfDStringP(plan.Fax),
		Email:            plan.Email.Value,
		TaxNumber:        plan.TaxNumber.Value,
		Country:          plan.Country.Value,
	}

	return request
}
