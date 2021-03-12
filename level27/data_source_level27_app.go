package level27

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLevel27App() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLevel27AppRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLevel27AppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var app App
	apiClient := meta.(*Client)

	app = apiClient.App("GET", d.Get("id"), nil)

	log.Printf("data_source_level27_app.go: Read routine called. Object built: %d", app.App.ID)

	d.SetId(strconv.Itoa(app.App.ID))
	d.Set("owner", app.App.Organisation.Name)

	return nil
}
