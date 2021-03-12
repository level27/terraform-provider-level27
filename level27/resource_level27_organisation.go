package level27

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27Organisation() *schema.Resource {
	return &schema.Resource{
		Create: resourceLevel27OrganisationCreate,
		Read:   resourceLevel27OrganisationRead,
		Delete: resourceLevel27OrganisationDelete,
		Update: resourceLevel27OrganisationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the organisation.`,
			},
			"tax_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Valid tax number for the organisation.`,
				ForceNew:    false,
			},
			"reseller_organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the reseller organisation.`,
				ForceNew:    false,
			},
			"parent_organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the parent organisation.`,
				ForceNew:    false,
			},
		},
	}
}

func resourceLevel27OrganisationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Starting Create")
	var organisation Organisation

	request := OrganisationRequest{
		Name:                 d.Get("name").(string),
		TaxNumber:            d.Get("tax_number").(string),
		ResellerOrganisation: d.Get("reseller_organisation_id").(string),
	}
	apiClient := meta.(*Client)
	endpoint := "organisations"

	body, err := apiClient.sendRequest("POST", endpoint, request)
	if err != nil {
		log.Printf("resource_level27_organisation.go: No Organisation created: %s", err)
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &organisation)
	if err != nil {
		return err
	}

	// Create is OK, we now set the properties of the object
	d.SetId(strconv.Itoa(organisation.Organisation.ID))
	return resourceLevel27OrganisationRead(d, meta)
}

func resourceLevel27OrganisationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Starting Read")
	log.Printf("resource_level27_organisation.go: Resource id: %s", d.Id())

	var organisation Organisation
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("organisations/%s", d.Id())

	body, err := apiClient.sendRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_organisation.go: No Organisation found: %s", d.Id())
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &organisation)
	if err != nil {
		return err
	}

	log.Printf("resource_level27_organisation.go: Read routine called. Object built: %d", organisation.Organisation.ID)

	d.SetId(strconv.Itoa(organisation.Organisation.ID))
	d.Set("name", organisation.Organisation.Name)
	d.Set("tax_number", organisation.Organisation.TaxNumber)
	if organisation.Organisation.ResellerOrganisation.ID != 0 {
		d.Set("reseller_organisation_id", strconv.Itoa(organisation.Organisation.ResellerOrganisation.ID))
	}

	return nil
}

func resourceLevel27OrganisationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Starting Update")

	request := OrganisationRequest{
		Name:                 d.Get("name").(string),
		TaxNumber:            d.Get("tax_number").(string),
		ResellerOrganisation: d.Get("reseller_organisation_id").(string),
	}
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("organisations/%s", d.Id())

	_, err := apiClient.sendRequest("PUT", endpoint, request)
	if err != nil {
		log.Printf("resource_level27_organisation.go: No Organisation updated: %s", err)
		d.SetId("")
		return nil
	}

	// Update is OK, we now set the properties of the object
	return resourceLevel27OrganisationRead(d, meta)
}

func resourceLevel27OrganisationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("organisations/%s", d.Id())

	_, err := apiClient.sendRequest("DELETE", endpoint, "")
	if err != nil {
		log.Printf("resource_level27_organisation.go: No Organisation deleted: %s", err)
		return err
	}

	d.SetId("")
	return nil
}
