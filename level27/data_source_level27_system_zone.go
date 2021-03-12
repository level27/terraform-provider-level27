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

func dataSourceLevel27SystemZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27SystemZoneRead,

		Schema: map[string]*schema.Schema{
			"region_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "be",
				ValidateFunc: validation.StringInSlice([]string{"be", "ams", "blr", "eu-central-1", "eu-west-1", "fra", "lon", "nyc", "sfo", "sgp", "tor"}, true),
				Description: `The region for the specified zone.  
Possible values:  
  - be (Default)
  - ams
  - blr
  - eu-central-1
  - eu-west-1
  - fra
  - lon
  - nyc
  - sfo
  - sgp
  - tor`,
			},
			"zone_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BRUIX1Z1",
				ValidateFunc: validation.StringInSlice([]string{"BRUIX2", "BRUIX2-C", "BRUIX1Z1", "BRUIX1Z2", "1", "2", "3", "a", "b", "c", "be2"}, true),
				Description: `The region for the specified zone.  
Possible values:  
  - be
    - BRUIX1Z1 (Default)
    - BRUIX1Z2
	- BRUIX2
	- BRUIX2-C
  - ams
	- 2
	- 3
  - blr
    - 1
  - eu-central-1
	- a
	- b
	- c
  - eu-west-1
	- a
	- b
	- c
  - fra
    - 1
  - lon
	- 1
  - nyc
	- 1
	- 2
    - 3
  - sfo
	- 1
	- 2
  - sgp
    - 1
  - tor
    - 1`,
			},
		},
	}
}

func dataSourceLevel27SystemZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var regions Regions
	apiClient := meta.(*Client)
	zoneIds := make([]int, 0)

	regions = apiClient.SystemRegions()

	log.Printf("data_source_level27_system_zone.go: Read routine called.")

	for _, region := range regions.Regions {
		if region.Name == d.Get("region_name").(string) {
			for _, zone := range region.Zones {
				if zone.ShortName == d.Get("zone_name").(string) {
					zoneIds = append(zoneIds, zone.ID)
				}
			}
		}
	}

	log.Printf("data_source_level27_system_zone.go: The query returned %d system zones; %v", len(zoneIds), zoneIds)

	if len(zoneIds) != 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error System Zone",
			Detail:   fmt.Sprintf("The query returned %d system zones, we cannot continue", len(zoneIds)),
		})
		return diags
	}

	log.Printf("data_source_level27_system_zone.go: zones found %v", zoneIds)

	d.SetId(strconv.Itoa(zoneIds[0]))

	return nil
}
