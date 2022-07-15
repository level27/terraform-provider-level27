package level27

/*
import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/level27/l27-go"
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
	var diags *diag.Diagnostics

	var app l27.App
	apiClient := meta.(*l27.Client)

	app, err := apiClient.App(d.Get("id").(int))

	log.Printf("data_source_level27_app.go: Read routine called. Object built: %d", app.App.ID)

	d.SetId(strconv.Itoa(app.ID))
	d.Set("owner", app.Organisation.Name)

	return nil
}
*/
