package level27

/*
import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLevel27Network() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27NetworkRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"provider_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func dataSourceLevel27NetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var networks Networks
	apiClient := meta.(*Client)
	networkIds := make([]int, 0)
	zoneID, err := strconv.Atoi(d.Get("zone_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	networks = apiClient.Networks("GET", d.Get("provider_id").(string))

	for _, network := range networks.Networks {
		if network.Zone.ID == zoneID && (network.Name == d.Get("name").(string)) {
			networkIds = append(networkIds, network.ID)
		}
	}

	d.SetId(strconv.Itoa(networkIds[0]))

	return nil
}
*/
