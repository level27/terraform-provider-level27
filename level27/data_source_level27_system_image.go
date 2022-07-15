package level27

/* import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLevel27SystemImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27SystemImageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ubuntu_1604lts_server", "ubuntu_1804lts_server", "ubuntu_2004lts_server"}, true),
				Description: `The name of the system image.
Possible values:
  - Level27
	- ubuntu_1604lts_server
	- ubuntu_1804lts_server
	- ubuntu_2004lts_server`,
			},
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID for the system provider.`,
			},
		},
	}
}

func dataSourceLevel27SystemImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var systemProviders SystemProviders
	apiClient := meta.(*Client)
	systemImages := make([]int, 0)

	systemProviders = apiClient.SystemProviders()

	log.Printf("data_source_level27_system_image.go: Read routine called.")

	for _, systemProvider := range systemProviders.Providers {
		if strconv.Itoa(systemProvider.ID) == d.Get("provider_id").(string) {
			for _, systemImage := range systemProvider.Images {
				if systemImage.Name == d.Get("name").(string) {
					systemImages = append(systemImages, systemImage.ID)
				}
			}
		}
	}

	if len(systemImages) != 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error SystemImage",
			Detail:   fmt.Sprintf("The query returned %d system images, we cannot continue", len(systemImages)),
		})
		return diags
	}

	d.SetId(strconv.Itoa(systemImages[0]))

	return nil
}
*/
