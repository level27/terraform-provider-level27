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

func dataSourceLevel27SystemProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27SystemProviderRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Level27"}, true),
				Description: `The name for the specified provider.
Possible values:
  - Level27`,
			},
		},
	}
}

func dataSourceLevel27SystemProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var systemProviders SystemProviders
	apiClient := meta.(*Client)
	systemProviderIds := make([]int, 0)

	systemProviders = apiClient.SystemProviders()

	log.Printf("data_source_level27_system_provider.go: Read routine called.")

	for _, systemProvider := range systemProviders.Providers {
		if systemProvider.Name == d.Get("name").(string) {
			systemProviderIds = append(systemProviderIds, systemProvider.ID)
		}
	}

	log.Printf("data_source_level27_system_provider.go: The query returned %d system providers; %v", len(systemProviderIds), systemProviderIds)

	if len(systemProviderIds) != 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error SystemProvider",
			Detail:   fmt.Sprintf("The query returned %d system providers, we cannot continue", len(systemProviderIds)),
		})
		return diags
	}

	d.SetId(strconv.Itoa(systemProviderIds[0]))

	return nil
}
