// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/level27/terraform-provider-level27/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &Level27Provider{}
var _ provider.ProviderWithFunctions = &Level27Provider{}

// Level27Provider is the provider implementation.
type Level27Provider struct {
	version string
}

type providerData struct {
	client         *client.Client
	organisationID int64
}

// Level27ProviderModel maps provider schema data to a Go type.
type Level27ProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
	APIURL types.String `tfsdk:"api_url"`
}

// New returns a provider.Provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Level27Provider{version: version}
	}
}

func (p *Level27Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "level27"
	resp.Version = p.version
}

func (p *Level27Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The **Level27** Terraform provider lets you manage your Level27 hosting resources (apps, components, URLs, SSL certificates) via the Level27 CP4 API.\n\n" +
			"Authenticate with your Level27 API key. Set `LEVEL27_API_KEY` in the environment to avoid storing credentials in your Terraform code.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Level27 API key. Can also be set via the `LEVEL27_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "Level27 API base URL. Defaults to `https://api.level27.eu/v1`. Can also be set via `LEVEL27_API_URL`.",
				Optional:            true,
			},
		},
	}
}

func (p *Level27Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Level27ProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API key: config > env var
	apiKey := os.Getenv("LEVEL27_API_KEY")
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing Level27 API Key",
			"Set the api_key provider attribute or the LEVEL27_API_KEY environment variable.",
		)
		return
	}

	// API URL: config > env var > default
	apiURL := os.Getenv("LEVEL27_API_URL")
	if !config.APIURL.IsNull() && !config.APIURL.IsUnknown() {
		apiURL = config.APIURL.ValueString()
	}

	c := client.New(apiKey, apiURL)

	me, err := c.GetWhoAmIUser(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to resolve authenticated user", fmt.Sprintf("Calling /whoami failed: %s", err))
		return
	}
	if me.Organisation == nil {
		resp.Diagnostics.AddError("Missing organisation in /whoami response", "The authenticated user did not include an organisation in /whoami response.")
		return
	}

	pd := &providerData{
		client:         c,
		organisationID: int64(me.Organisation.ID),
	}

	resp.DataSourceData = pd
	resp.ResourceData = pd
}

func (p *Level27Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAppResource,
		NewAppComponentResource,
		NewAppComponentURLResource,
		NewSSLCertificateResource,
		NewSystemResource,
	}
}

func (p *Level27Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAppDataSource,
	}
}

func (p *Level27Provider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
