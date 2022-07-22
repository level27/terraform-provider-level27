package level27

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/level27/l27-go"
)

// Provider defines the TF provider for Level27
func Provider() tfsdk.Provider {
	return &provider{}
	/* 	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEVEL27_API_URI", nil),
				Description: "URI of the REST API endpoint. This serves as the base of all requests.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEVEL27_API_KEY", nil),
				Description: "Key used to connect to the API.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"level27_network":                       dataSourceLevel27Network(),
			"level27_system_provider":               dataSourceLevel27SystemProvider(),
			"level27_system_provider_configuration": dataSourceLevel27SystemProviderConfiguration(),
			"level27_system_image":                  dataSourceLevel27SystemImage(),
			"level27_system_zone":                   dataSourceLevel27SystemZone(),
			//"level27_app":                           dataSourceLevel27App(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"level27_organisation":       resourceLevel27Organisation(),
			"level27_system":             resourceLevel27System(),
			"level27_cookbook_php":       resourceLevel27CookbookPhp(),
			"level27_cookbook_mysql":     resourceLevel27CookbookMysql(),
			"level27_cookbook_docker":    resourceLevel27CookbookDocker(),
			"level27_cookbook_haproxy":   resourceLevel27CookbookHaproxy(),
			"level27_cookbook_memcached": resourceLevel27CookbookMemcached(),
			"level27_cookbook_mongodb":   resourceLevel27CookbookMongodb(),
			"level27_cookbook_postfix":   resourceLevel27CookbookPostfix(),
			"level27_cookbook_url":       resourceLevel27CookbookURL(),
			"level27_cookbook_varnish":   resourceLevel27CookbookVarnish(),
			//"level27_domain":         resourceLevel27Domain(),
			//"level27_domain_contact": resourceLevel27DomainContact(),
			//"level27_domain_record":  resourceLevel27DomainRecord(),
			//"level27_user":               resourceLevel27User(),
			//"level27_mailgroup":          resourceLevel27Mailgroup(),
			//"level27_mailforwarder":      resourceLevel27Mailforwarder(),
			//"level27_mailbox":            resourceLevel27Mailbox(),
			//"level27_app":                resourceLevel27App(),
			//"level27_appcomponent_php":   resourceLevel27AppComponentPhp(),
			//"level27_appcomponent_mysql": resourceLevel27AppComponentMysql(),

		},
		ConfigureContextFunc: configureProvider,
	} */
}

type provider struct {
	Client *l27.Client
}

// Provider schema struct
type providerData struct {
	ApiUrl types.String `tfsdk:"api_url"`
	Key    types.String `tfsdk:"key"`
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_url": {
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "API URL for the Level27 API.",
			},
			"key": {
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the Level27 API",
			},
		},
	}, nil
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	req.Config.Get(ctx, &config)

	var apiUrl string
	if config.ApiUrl.Unknown {
		resp.Diagnostics.AddError("Invalid API URL", "Cannot use unknown value for API URL.")
		return
	}

	if !config.ApiUrl.Null {
		apiUrl = config.ApiUrl.Value
	} else {
		apiUrl = os.Getenv("LEVEL27_API_URL")
		if apiUrl == "" {
			apiUrl = "https://api.level27.eu/v1"
		}
	}

	var apiKey string
	if config.Key.Unknown {
		resp.Diagnostics.AddError("Invalid API key", "Cannot use unknown value for API key.")
		return
	}

	if !config.Key.Null {
		apiKey = config.Key.Value
	} else {
		apiKey = os.Getenv("LEVEL27_API_KEY")
	}

	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API key", "An API Key must be provided via provider configuration or the LEVEL27_API_KEY environment variable.")
		return
	}

	p.Client = l27.NewAPIClient(apiUrl, apiKey)
	p.Client.TraceRequests(tfTracer{Context: ctx})
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"level27_app":          resourceAppType{},
		"level27_organisation": resourceOrganisationType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		//"hashicups_coffees": dataSourceCoffeesType{},
	}, nil
}

type tfTracer struct {
	Context context.Context
}

func (t tfTracer) TraceRequest(method string, url string, reqData []byte) {
	tflog.Debug(t.Context, "Request", map[string]interface{}{"method": method, "url": url, "body": string(reqData)})
}

func (t tfTracer) TraceResponse(response *http.Response) {

}

func (t tfTracer) TraceResponseBody(response *http.Response, respData []byte) {

}
