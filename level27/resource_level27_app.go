package level27

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/level27/l27-go"
)

type resourceAppType struct {
}

func (r resourceAppType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Required: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (r resourceAppType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceApp{
		p: p.(*provider),
	}, nil
}

type resourceApp struct {
	p *provider
}

type resourceAppModel struct {
	ID           types.String `tfsdk:"id"`
	Organisation types.String `tfsdk:"organisation"`
	Name         types.String `tfsdk:"name"`
}

func (r resourceApp) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan resourceAppModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var orgID int

	if plan.Organisation.Null || plan.Organisation.Value == "" {
		// Organisation defaults to the org of your login
		login, err := r.p.GetLoginInfo()
		if err != nil {
			resp.Diagnostics.AddError("Error fetching login info", "API request failed:\n"+err.Error())
			return
		}

		orgID = login.User.Organisation.ID
	} else {
		var ok bool
		orgID, ok = parseID(plan.Organisation.Value, &resp.Diagnostics)
		if !ok {
			return
		}
	}

	request := l27.AppPostRequest{
		Name:         plan.Name.Value,
		Organisation: orgID,
	}

	app, err := r.p.Client.AppCreate(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating app", "API request failed:\n"+err.Error())
		return
	}

	result := resourceAppModel{
		ID:           types.String{Value: strconv.Itoa(app.ID)},
		Name:         types.String{Value: app.Name},
		Organisation: types.String{Value: strconv.Itoa(app.Organisation.ID)},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceApp) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state resourceAppModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, ok := parseID(state.ID.Value, &resp.Diagnostics)
	if !ok {
		return
	}

	app, err := r.p.Client.App(appID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading app", "API request failed:\n"+err.Error())
		return
	}

	state.ID = types.String{Value: strconv.Itoa(app.ID)}
	state.Name = types.String{Value: app.Name}
	state.Organisation = types.String{Value: strconv.Itoa(app.Organisation.ID)}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceApp) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan resourceAppModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state resourceAppModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, ok := parseID(state.ID.Value, &resp.Diagnostics)
	if !ok {
		return
	}

	orgID, ok := parseID(plan.Organisation.Value, &resp.Diagnostics)
	if !ok {
		return
	}

	request := l27.AppPutRequest{
		Name:         plan.Name.Value,
		Organisation: orgID,
	}

	err := r.p.Client.AppUpdate(appID, request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating app", "API request failed:\n"+err.Error())
		return
	}

	state.Name = plan.Name
	state.Organisation = plan.Organisation

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceApp) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state resourceAppModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, ok := parseID(state.ID.Value, &resp.Diagnostics)
	if !ok {
		return
	}

	err := r.p.Client.AppDelete(appID)
	if err != nil {
		resp.Diagnostics.AddError("Error destroying app", "API request failed:\n"+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceApp) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
