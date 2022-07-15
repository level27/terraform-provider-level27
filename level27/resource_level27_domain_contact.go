package level27

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27DomainContact() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27DomainContactCreate,
		Read:   resourceLevel27DomainContactRead,
		Delete: resourceLevel27DomainContactDelete,
		Update: resourceLevel27DomainContactUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    true,
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    true,
			},
			"organisation_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"tax_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"street": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"zip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"city": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
				ForceNew:    false,
			},
		},
	}
}

func resourceLevel27DomainContactCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_contact.go: Starting Create")

	request := DomainContactRequest{
		Type:             d.Get("type").(string),
		FirstName:        d.Get("first_name").(string),
		LastName:         d.Get("last_name").(string),
		Organisation:     d.Get("organisation_id").(string),
		Street:           d.Get("street").(string),
		Zip:              d.Get("zip").(string),
		City:             d.Get("city").(string),
		Country:          d.Get("country").(string),
		TaxNumber:        d.Get("tax_number").(string),
		Email:            d.Get("email").(string),
		Phone:            d.Get("phone").(string),
		OrganisationName: d.Get("organisation_name").(string),
	}

	apiClient := meta.(*Client)
	endpoint := "domaincontacts"
	body, err := apiClient.sendRequest("POST", endpoint, request)

	if err != nil {
		log.Printf("resource_level27_domain_contact.go: No Domain record created: %s", err)
		d.SetId("")
		return err
	}

	var domainContact DomainContact
	err = json.Unmarshal(body, &domainContact)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(domainContact.Domaincontact.ID))

	return resourceLevel27DomainContactRead(d, meta)
}

func resourceLevel27DomainContactRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_contact.go: Starting Read")

	var domainContact DomainContact
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domaincontacts/%s", d.Id())
	body, err := apiClient.sendRequest("GET", endpoint, nil)

	if err != nil {
		log.Printf("resource_level27_domain_contact.go: No Domain contact found: %s", d.Id())
		d.SetId("")
		return err
	}

	err = json.Unmarshal(body, &domainContact)
	if err != nil {
		return err
	}

	log.Printf("resource_level27_domain_contact.go:  Domain contact found: %d", domainContact.Domaincontact.ID)

	d.SetId(strconv.Itoa(domainContact.Domaincontact.ID))
	d.Set("type", domainContact.Domaincontact.Type)
	d.Set("first_name", domainContact.Domaincontact.FirstName)
	d.Set("last_name", domainContact.Domaincontact.LastName)
	d.Set("organisation_name", domainContact.Domaincontact.OrganisationName)
	d.Set("street", domainContact.Domaincontact.Street)
	d.Set("zip", domainContact.Domaincontact.Zip)
	d.Set("city", domainContact.Domaincontact.City)
	d.Set("country", domainContact.Domaincontact.Country.ID)
	d.Set("phone", domainContact.Domaincontact.Phone)
	d.Set("email", domainContact.Domaincontact.Email)
	d.Set("tax_number", domainContact.Domaincontact.TaxNumber)
	orgid := strconv.Itoa(domainContact.Domaincontact.Organisation.ID)
	if orgid == "null" {
		orgid = "0"
	}
	d.Set("organisation_id", strconv.Itoa(domainContact.Domaincontact.Organisation.ID))

	return nil
}

func resourceLevel27DomainContactUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_contact.go: Starting Update")

	request := DomainContactRequest{
		Type:             d.Get("type").(string),
		FirstName:        d.Get("first_name").(string),
		LastName:         d.Get("last_name").(string),
		Organisation:     d.Get("organisation_id").(string),
		Street:           d.Get("street").(string),
		Zip:              d.Get("zip").(string),
		City:             d.Get("city").(string),
		Country:          d.Get("country").(string),
		TaxNumber:        d.Get("tax_number").(string),
		Email:            d.Get("email").(string),
		Phone:            d.Get("phone").(string),
		OrganisationName: d.Get("organisation_name").(string),
	}

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domaincontacts/%s", d.Id())
	body, err := apiClient.sendRequest("PUT", endpoint, request)

	if err != nil {
		log.Printf("resource_level27_domain_contact.go: No Domain contact created: %s", err)
		d.SetId("")
		return err
	}

	var domainContact DomainContact
	err = json.Unmarshal(body, &domainContact)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(domainContact.Domaincontact.ID))

	return resourceLevel27DomainContactRead(d, meta)
}

func resourceLevel27DomainContactDelete(d *schema.ResourceData, meta interface{}) error { return nil }
*/
