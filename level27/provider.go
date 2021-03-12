package level27

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions map[string]string

// Provider defines the TF provider for Level27
func Provider() *schema.Provider {
	return &schema.Provider{
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
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := NewAPIClient(d.Get("uri").(string), d.Get("key").(string))
	return client, diags
}
