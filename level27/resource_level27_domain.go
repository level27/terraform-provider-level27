package level27

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27Domain() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27DomainCreate,
		Read:   resourceLevel27DomainRead,
		Delete: resourceLevel27DomainDelete,
		Update: resourceLevel27DomainUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the domain",
				Required:    true,
				ForceNew:    true,
			},
			"extension": {
				Type:        schema.TypeString,
				Description: "Extension of the domain",
				Required:    true,
				ForceNew:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "one of none, create, transfer, transferFromQuarantine",
				ForceNew:    false,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the owning organisation",
				ForceNew:    false,
			},
			"domain_contact_licensee_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the domain contact",
				ForceNew:    false,
			},
			"handle_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, use the Level27 nameservers",
				ForceNew:    false,
			},
		},
	}
}

func resourceLevel27DomainCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain.go: Starting Create")
	var domain Domain
	domainExtensions, _ := getDomainExtensions(d, meta)
	log.Printf("resource_level27_domain.go: Domain extension: %s wit ID:%d", d.Get("extension").(string), domainExtensions[d.Get("extension").(string)])

	request := DomainRequest{
		Name:                  d.Get("name").(string),
		Action:                d.Get("action").(string),
		Domaintype:            domainExtensions[d.Get("extension").(string)],
		Domaincontactlicensee: d.Get("domain_contact_licensee_id").(string),
		Organisation:          d.Get("organisation_id").(string),
		Handledns:             d.Get("handle_dns").(bool),
	}
	apiClient := meta.(*Client)
	endpoint := "domains"

	body, err := apiClient.sendRequest("POST", endpoint, request.String())

	if err != nil {
		log.Printf("resource_level27_domain.go: No Domain created: %s", err)
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &domain)
	if err != nil {
		return err
	}

	// Only check availability when registering
	//available, domainTypeId := checkDomainAvailability(d.Get("name").(string), d.Get("extension").(string))
	//log.Printf("Avalable: %t", available)
	//log.Printf("domaintype: %s", domainTypeId)

	// Create is OK, we now set the properties of the object
	d.SetId(strconv.Itoa(domain.Domain.ID))
	return resourceLevel27DomainRead(d, meta)
}

func resourceLevel27DomainRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain.go: Starting Read")
	var domain Domain
	apiClient := meta.(*Client)
	log.Printf("resource_level27_domain.go: Resource id: %s", d.Id())

	endpoint := fmt.Sprintf("domains/%s", d.Id())

	body, err := apiClient.sendRequest("GET", endpoint, "")
	if err != nil {
		log.Printf("resource_level27_domain.go: No Domain found: %s", d.Id())
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &domain)
	if err != nil {
		return err
	}

	log.Printf("resource_level27_domain.go: Read routine called. Object built: %d", domain.Domain.ID)

	d.SetId(strconv.Itoa(domain.Domain.ID))
	d.Set("name", domain.Domain.Name)
	d.Set("extension", domain.Domain.Domaintype.Extension)
	d.Set("organisation_id", strconv.Itoa(domain.Domain.Organisation.ID))
	d.Set("domain_contact_licensee_id", strconv.Itoa(domain.Domain.DomaincontactLicensee.ID))
	d.Set("handle_dns", domain.Domain.DNSIsHandled)

	return nil
}

func resourceLevel27DomainUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain.go: Starting Update")

	domainExtensions, _ := getDomainExtensions(d, meta)

	request := DomainRequest{
		Name:                  d.Get("name").(string),
		Action:                d.Get("action").(string),
		Domaintype:            domainExtensions[d.Get("extension").(string)],
		Domaincontactlicensee: d.Get("domain_contact_licensee_id").(string),
		Organisation:          d.Get("organisation_id").(string),
		Handledns:             d.Get("handle_dns").(bool),
	}
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domains/%s", d.Id())

	_, err := apiClient.sendRequest("PUT", endpoint, request.String())

	if err != nil {
		log.Printf("resource_level27_domain.go: No Domain updated: %s", err)
		d.SetId("")
		return nil
	}

	return resourceLevel27DomainRead(d, meta)
}

func resourceLevel27DomainDelete(d *schema.ResourceData, meta interface{}) error {
	endpoint := fmt.Sprintf("domains/%s", d.Id())
	apiClient := meta.(*Client)
	_, err := apiClient.sendRequest("DELETE", endpoint, "")
	if err != nil {
		log.Printf("resource_level27_domain.go: No Domain deleted: %s", d.Id())
		return nil
	}

	d.SetId("")
	return nil
}

func getDomainExtensions(d *schema.ResourceData, meta interface{}) (map[string]int, error) {
	log.Printf("resource_level27_domain.go: Start getDomainExtension")
	var domainProviders DomainProvider
	apiClient := meta.(*Client)
	endpoint := "domains/providers"

	body, err := apiClient.sendRequest("GET", endpoint, "")
	if err != nil {
		log.Printf("resource_level27_domain.go: No domainextensions found: %s", d.Id())
		d.SetId("")
		return nil, err
	}

	err = json.Unmarshal(body, &domainProviders)
	if err != nil {
		log.Printf("unmarshall gone wrong")
		log.Printf(string(body))
		return nil, err
	}

	m := make(map[string]int)
	for _, provider := range domainProviders.Providers {
		for _, domainType := range provider.Domaintypes {
			m[domainType.Extension] = domainType.ID
		}
	}

	return m, nil
}

func checkDomainAvailability(name string, extension string) (bool, string) {
	log.Printf("resource_level27_domain.go: Starting checkAvailability")
	return true, "1"
}
