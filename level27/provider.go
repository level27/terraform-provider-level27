package level27

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/level27/l27-go"
)

// Provider defines the TF provider for Level27
func Provider() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	Client     *l27.Client
	loginInfo  *l27.Login
	loginError error
	loginMutex sync.Mutex
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
		"level27_app":            resourceAppType{},
		"level27_organisation":   resourceOrganisationType{},
		"level27_domain":         resourceDomainType{},
		"level27_domain_contact": resourceDomainContactType{},
		"level27_system":         resourceSystemType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"level27_user":            dataSourceUserType{},
		"level27_app":             dataSourceResourceType{ResourceType: resourceAppType{}},
		"level27_organisation":    dataSourceResourceType{ResourceType: resourceOrganisationType{}},
		"level27_system_image":    dataSourceSystemImageType{},
		"level27_system_provider": dataSourceSystemProviderType{},
		"level27_system_zone":     dataSourceSystemZoneType{},
	}, nil
}

// Fetches the login info for the configured API key. Is cached, so a request only happens once.
func (p *provider) GetLoginInfo() (l27.Login, error) {
	// Mutex to avoid
	p.loginMutex.Lock()
	defer p.loginMutex.Unlock()

	if p.loginError != nil {
		return l27.Login{}, p.loginError
	}

	if p.loginInfo != nil {
		return *p.loginInfo, nil
	}

	login, err := p.Client.LoginInfo()
	if err != nil {
		p.loginError = err
		return l27.Login{}, err
	} else {
		p.loginInfo = &login
		return login, nil
	}
}

type tfTracer struct {
	Context context.Context
}

func (t tfTracer) TraceRequest(method string, url string, reqData []byte) {
	tflog.Debug(t.Context, "Request", map[string]interface{}{"method": method, "url": url, "body": string(reqData)})
}

func (t tfTracer) TraceResponse(response *http.Response) {
	tflog.Debug(t.Context, "Request Response", map[string]interface{}{"status": response.StatusCode})
}

func (t tfTracer) TraceResponseBody(response *http.Response, respData []byte) {
	tflog.Debug(t.Context, "Request Response", map[string]interface{}{"body": string(respData)})
}
