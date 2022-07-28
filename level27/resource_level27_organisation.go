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

type resourceOrganisationType struct{}

func (r resourceOrganisationType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:      true,
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: `The name of the organisation.`,
			},
			"tax_number": {
				Type:        types.StringType,
				Required:    true,
				Description: `Valid tax number for the organisation.`,
			},
			"reseller_organisation_id": {
				Type:        types.StringType,
				Optional:    true,
				Description: `The ID of the reseller organisation.`,
				Validators:  []tfsdk.AttributeValidator{validateID()},
			},
			"parent_organisation_id": {
				Type:          types.StringType,
				Description:   `The ID of the parent organisation.`,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
		},
	}, nil
}

func (r resourceOrganisationType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceOrganisation{
		p: p.(*provider),
	}, nil
}

type resourceOrganisation struct {
	p *provider
}

type resourceOrganisationModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	TaxNumber              types.String `tfsdk:"tax_number"`
	ResellerOrganisationID types.String `tfsdk:"reseller_organisation_id"`
	ParentOrganisationID   types.String `tfsdk:"parent_organisation_id"`
}

func (r resourceOrganisation) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan resourceOrganisationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := l27.OrganisationCreate{
		Name:      plan.Name.Value,
		TaxNumber: plan.TaxNumber.Value,
	}

	if !plan.ResellerOrganisationID.Null && !plan.ResellerOrganisationID.Unknown {
		id, _ := strconv.Atoi(plan.ResellerOrganisationID.Value)
		request.ResellerOrganisation = &id
	}

	org, err := r.p.Client.OrganisationCreate(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating organisation", "API request failed:\n"+err.Error())
		return
	}

	result := organisationModelFrom(org)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceOrganisation) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state resourceOrganisationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(state.ID.Value)

	org, err := r.p.Client.Organisation(orgID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading organisation", "API request failed:\n"+err.Error())
		return
	}

	result := organisationModelFrom(org)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func organisationModelFrom(org l27.Organisation) resourceOrganisationModel {
	result := resourceOrganisationModel{
		ID:                     types.String{Value: strconv.Itoa(org.ID)},
		Name:                   types.String{Value: org.Name},
		TaxNumber:              types.String{Value: org.TaxNumber},
		ParentOrganisationID:   types.String{Value: org.ParentOrganisation},
		ResellerOrganisationID: types.String{Null: true},
	}

	if org.ResellerOrganisation != nil {
		result.ResellerOrganisationID = types.String{Value: strconv.Itoa(*org.ResellerOrganisation)}
	}

	return result
}

func (r resourceOrganisation) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan resourceOrganisationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state resourceOrganisationModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(state.ID.Value)

	request := map[string]interface{}{
		"name":      plan.Name.Value,
		"taxNumber": plan.TaxNumber.Value,
	}

	err := r.p.Client.OrganisationUpdate(orgID, request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating organisation", "API request failed:\n"+err.Error())
		return
	}

	// NOTE: we set state = plan her. This relies on UseStateForUnknown() so that plan already contains merged state.
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r resourceOrganisation) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state resourceOrganisationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(state.ID.Value)

	err := r.p.Client.OrganisationDelete(orgID)
	if err != nil {
		resp.Diagnostics.AddError("Error destroying organisation", "API request failed:\n"+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceOrganisation) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
