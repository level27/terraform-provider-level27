package level27

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type dataSourceResourceType struct {
	ResourceType tfsdk.ResourceType
}

func (r dataSourceResourceType) GetSchema(c context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema, diag := r.ResourceType.GetSchema(c)
	if diag.HasError() {
		return tfsdk.Schema{}, diag
	}

	newAttributes := make(map[string]tfsdk.Attribute)
	for k, v := range schema.Attributes {
		if k == "id" {
			v.Required = true
			v.Computed = false
			v.Optional = false
		} else {
			v.Required = false
			v.Computed = true
			v.Optional = true
		}

		newAttributes[k] = v
	}

	schema.Attributes = newAttributes

	return schema, nil
}

func (r dataSourceResourceType) NewDataSource(c context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	res, diagnostics := r.ResourceType.NewResource(c, p)

	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return dataSourceResource{
		Resource: res,
	}, nil
}

type dataSourceResource struct {
	Resource tfsdk.Resource
}

func (r dataSourceResource) Read(c context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	resRequest := tfsdk.ReadResourceRequest{
		State: tfsdk.State{
			Raw:    req.Config.Raw,
			Schema: req.Config.Schema,
		},
		ProviderMeta: req.ProviderMeta,
	}

	resResp := tfsdk.ReadResourceResponse{
		State: resp.State,
	}

	r.Resource.Read(c, resRequest, &resResp)

	resp.State = resResp.State
	resp.Diagnostics = resResp.Diagnostics
}
