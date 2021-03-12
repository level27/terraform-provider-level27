package level27

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27User() *schema.Resource {
	return &schema.Resource{
		Create: resourceLevel27UserCreate,
		Read:   resourceLevel27UserRead,
		Delete: resourceLevel27UserDelete,
		Update: resourceLevel27UserUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"last_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				ForceNew:    false,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"roles": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceLevel27UserCreate(d *schema.ResourceData, meta interface{}) error {
	// log.Printf("Starting Create")
	return resourceLevel27UserRead(d, meta)
}

func resourceLevel27UserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Starting Read")
	log.Printf("resource_level27_organisation.go: Resource id: %s", d.Id())

	apiClient := meta.(*Client)
	user := apiClient.User("GET", nil, d.Id(), nil)

	log.Printf("resource_level27_user.go: Read routine called. Object built: %s", user.User.Email)

	return nil
}

func resourceLevel27UserUpdate(d *schema.ResourceData, meta interface{}) error {
	// log.Printf("Starting Update")
	return resourceLevel27UserRead(d, meta)
}

func resourceLevel27UserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
