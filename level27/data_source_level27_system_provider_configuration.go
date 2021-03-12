package level27

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLevel27SystemProviderConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27SystemProviderConfigurationRead,

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID for the system provider.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large", "flexible"}, true),
				Description: `The name for the system configuration.
Possible values:
  - Level27
	- small
	- medium
	- large
	- flexible`,
			},
		},
	}
}

func dataSourceLevel27SystemProviderConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var providerConfigurations ProviderConfigurations
	apiClient := meta.(*Client)
	providerConfigurationIds := make([]int, 0)

	providerConfigurations = apiClient.SystemProviderConfigurations()

	log.Printf("data_source_level27_system_provider_configuration.go: Read routine called.")

	for _, providerConfiguration := range providerConfigurations.ProviderConfigurations {
		// log.Printf("%s - %d", providerConfiguration.ExternalID, providerConfiguration.Systemprovider.ID)
		if providerConfiguration.ExternalID == d.Get("name").(string) && strconv.Itoa(providerConfiguration.Systemprovider.ID) == d.Get("provider_id").(string) {
			providerConfigurationIds = append(providerConfigurationIds, providerConfiguration.ID)
		}
	}

	if len(providerConfigurationIds) != 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error SystemProviderConfiguration",
			Detail:   fmt.Sprintf("The query returned %d system provider configurations, we cannot continue", len(providerConfigurationIds)),
		})
		return diags
	}

	d.SetId(strconv.Itoa(providerConfigurationIds[0]))

	return nil
}
