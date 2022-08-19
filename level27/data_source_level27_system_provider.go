package level27

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/l27-go"
)

type dataSourceSystemProviderType struct{}

// GetSchema implements tfsdk.DataSourceType
func (dataSourceSystemProviderType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The ID of the located system image",
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name for the specified provider. You usually want 'Level27'",
			},
		},
	}, nil
}

// NewDataSource implements tfsdk.DataSourceType
func (dataSourceSystemProviderType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceSystemProvider{
		p: p.(*provider),
	}, nil
}

type dataSourceSystemProvider struct {
	p *provider
}

type dataSourceSystemProviderModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Read implements tfsdk.DataSource
func (d dataSourceSystemProvider) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config dataSourceSystemProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerName := config.Name.Value

	providers, err := d.p.Client.GetSystemProviders()
	if err != nil {
		resp.Diagnostics.AddError("Error reading providers", "API request failed:\n"+err.Error())
		return
	}

	provider := findSystemProviderByName(providers, providerName)
	if provider == nil {
		resp.Diagnostics.AddError("Unable to find provider", fmt.Sprintf("No system provider with name '%s' could be found!", providerName))
		return
	}

	config.ID = tfStringId(provider.ID)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func findSystemProviderByName(providers []l27.SystemProvider, name string) *l27.SystemProvider {
	for _, provider := range providers {
		if provider.Name == name {
			return &provider
		}
	}
	return nil
}
