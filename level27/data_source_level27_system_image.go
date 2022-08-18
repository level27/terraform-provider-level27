package level27

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/l27-go"
)

type dataSourceSystemImageType struct{}

// GetSchema implements tfsdk.DataSourceType
func (dataSourceSystemImageType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Description: "The name of the system image to use. For example ubuntu_2004lts_server",
			},
			"provider_id": {
				Type:        types.StringType,
				Required:    true,
				Description: "The ID of the system provider",
				Validators:  []tfsdk.AttributeValidator{validateID()},
			},
		},
	}, nil
}

// NewDataSource implements tfsdk.DataSourceType
func (dataSourceSystemImageType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceSystemImage{
		p: p.(*provider),
	}, nil
}

type dataSourceSystemImage struct {
	p *provider
}

type dataSourceSystemImageModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Provider types.String `tfsdk:"provider_id"`
}

// Read implements tfsdk.DataSource
func (d dataSourceSystemImage) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config dataSourceSystemImageModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerID, _ := strconv.Atoi(config.Provider.Value)
	imageName := config.Name.Value

	providers, err := d.p.Client.GetSystemProviders()
	if err != nil {
		resp.Diagnostics.AddError("Error reading providers", "API request failed:\n"+err.Error())
		return
	}

	provider := findSystemProvider(providers, providerID)
	if provider == nil {
		resp.Diagnostics.AddError("Unable to find provider", fmt.Sprintf("No system provider with ID '%d' could be found!", providerID))
		return
	}

	image := findSystemProviderImage(provider.Images, imageName)
	if image == nil {
		resp.Diagnostics.AddError("Unable to find image on provider", fmt.Sprintf("No image with name '%s' could be found on provider %s!", imageName, provider.Name))
		return
	}

	config.ID = tfStringId(image.ID)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func findSystemProvider(providers []l27.SystemProvider, id int) *l27.SystemProvider {
	for _, provider := range providers {
		if provider.ID == id {
			return &provider
		}
	}
	return nil
}

func findSystemProviderImage(images []l27.SystemProviderImage, name string) *l27.SystemProviderImage {
	for _, image := range images {
		if image.Name == name {
			return &image
		}
	}
	return nil
}
