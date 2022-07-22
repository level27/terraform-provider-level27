package level27

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceUserType struct {
}

func (r dataSourceUserType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Required:   true,
				Type:       types.StringType,
				Validators: []tfsdk.AttributeValidator{validateTwoID()},
			},
			"email": {
				Type:     types.StringType,
				Computed: true,
			},
			"user_name": {
				Type:     types.StringType,
				Computed: true,
			},
			"first_name": {
				Type:        types.StringType,
				Description: "",
				Computed:    true,
			},
			"last_name": {
				Type:        types.StringType,
				Description: "",
				Computed:    true,
			},
			"full_name": {
				Type:        types.StringType,
				Description: "",
				Computed:    true,
			},
			"organisation": {
				Type:        types.StringType,
				Description: "",
				Computed:    true,
			},
			"roles": {
				Computed: true,
				Type:     types.ListType{}.WithElementType(types.StringType),
			},
		},
	}, nil
}

func (r dataSourceUserType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUser{
		p: p.(*provider),
	}, nil
}

type dataSourceUser struct {
	p *provider
}

type dataSourceUserModel struct {
	ID           string       `tfsdk:"id"`
	Email        types.String `tfsdk:"email"`
	UserName     types.String `tfsdk:"user_name"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	FullName     types.String `tfsdk:"full_name"`
	Organisation types.String `tfsdk:"organisation"`
	Roles        types.List   `tfsdk:"roles"`
}

func (r dataSourceUser) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config dataSourceUserModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idOrg, idUser, _ := appComponentParseID(config.ID)

	user, err := r.p.Client.OrganisationUserGetSingle(idOrg, idUser)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", "API request failed:\n"+err.Error())
		return
	}

	result := dataSourceUserModel{
		ID:           config.ID,
		Email:        types.String{Value: user.Email},
		FirstName:    types.String{Value: user.FirstName},
		LastName:     types.String{Value: user.LastName},
		UserName:     types.String{Value: user.Username},
		Organisation: types.String{Value: strconv.Itoa(idOrg)},
		Roles:        types.List{ElemType: types.StringType},
	}

	for _, v := range user.Roles {
		result.Roles.Elems = append(result.Roles.Elems, types.String{Value: v})
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}
