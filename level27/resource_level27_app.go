package level27

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27App() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27AppCreate,
		Read:   resourceLevel27AppRead,
		Delete: resourceLevel27AppDelete,
		Update: resourceLevel27AppUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"organisation_id": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceLevel27AppRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_app.go: Read resource")

	var app App
	apiClient := meta.(*Client)

	app = apiClient.App("GET", d.Id(), nil)

	log.Printf("resource_level27_app.go: Read routine called. Object built: %d", app.App.ID)

	d.SetId(strconv.Itoa(app.App.ID))
	d.Set("name", app.App.Name)
	d.Set("organisation_id", strconv.Itoa(app.App.Organisation.ID))
	return nil
}

func resourceLevel27AppCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_app.go: Create resource")

	var app App
	apiClient := meta.(*Client)
	request := AppRequest{
		Name:         d.Get("name").(string),
		Organisation: d.Get("organisation_id").(string),
	}

	app = apiClient.App("CREATE", nil, request.String())

	d.SetId(strconv.Itoa(app.App.ID))
	return resourceLevel27AppRead(d, meta)
}

func resourceLevel27AppUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_app.go: Update resource")

	apiClient := meta.(*Client)
	request := AppRequest{
		Name:         d.Get("name").(string),
		Organisation: d.Get("organisation_id").(string),
	}

	apiClient.App("UPDATE", d.Id(), request.String())

	return resourceLevel27AppRead(d, meta)
}

func resourceLevel27AppDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_app.go: Delete resource")

	apiClient := meta.(*Client)
	apiClient.App("DELETE", d.Id(), nil)

	return nil
}
