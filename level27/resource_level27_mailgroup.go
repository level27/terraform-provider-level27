package level27

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27Mailgroup() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27MailgroupCreate,
		Read:   resourceLevel27MailgroupRead,
		Delete: resourceLevel27MailgroupDelete,
		Update: resourceLevel27MailgroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the owning organisation",
				ForceNew:    false,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the owning organisation",
				ForceNew:    false,
			},
			"linked_domains": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceLevel27MailgroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailgroup.go: Starting Create")
	var mailgroup Mailgroup

	request := MailgroupRequest{
		Name:         d.Get("name").(string),
		Organisation: d.Get("organisation_id").(string),
	}
	apiClient := meta.(*Client)
	endpoint := "mailgroups"

	body, err := apiClient.sendRequest("POST", endpoint, request.String())

	if err != nil {
		log.Printf("resource_level27_mailgroup.go: No Mailgroup created: %s", err)
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &mailgroup)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mailgroup.Mailgroup.ID))

	domains := d.Get("linked_domains").(*schema.Set)
	d.Set("linked_domains", updateLinkedDomains(domains, d, meta, strconv.Itoa(mailgroup.Mailgroup.ID)))

	return resourceLevel27MailgroupRead(d, meta)
}

func resourceLevel27MailgroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Starting Read")
	log.Printf("resource_level27_mailgroup.go: Resource id: %s", d.Id())

	mailgroup := getMailgroup(d, meta)

	log.Printf("resource_level27_mailgroup.go: Read routine called. Object built: %d", mailgroup.Mailgroup.ID)

	d.SetId(strconv.Itoa(mailgroup.Mailgroup.ID))
	d.Set("name", mailgroup.Mailgroup.Name)
	d.Set("organisation_id", strconv.Itoa(mailgroup.Mailgroup.Organisation.ID))
	d.Set("linked_domains", splatLinkedDomains(mailgroup.Mailgroup.Domains))

	return nil
}

func resourceLevel27MailgroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailgroup.go: Starting Update")

	domains := d.Get("linked_domains").(*schema.Set)

	d.Set("linked_domains", updateLinkedDomains(domains, d, meta, d.Id()))
	return resourceLevel27MailgroupRead(d, meta)
}

func resourceLevel27MailgroupDelete(d *schema.ResourceData, meta interface{}) error { return nil }

// Custom functions //

func updateLinkedDomains(domains *schema.Set, d *schema.ResourceData, meta interface{}, mailgroup string) []int {
	log.Printf("resource_level27_mailgroup.go: updateLinkedDomains")

	mailgroupOriginal := getMailgroup(d, meta)
	linkedDomainsOriginal := splatLinkedDomains(mailgroupOriginal.Mailgroup.Domains)

	apiClient := meta.(*Client)

	linkedDomains := []int{}
	// loop over the new list to see if we need to add some domains
	for _, domain := range domains.List() {
		if !(containsInt(linkedDomainsOriginal, domain.(int))) {
			//add the domain
			endpoint := fmt.Sprintf("mailgroups/%s/domains", d.Id())
			log.Printf("resource_level27_mailgroup.go: Link new domain; %d", domain)
			s := fmt.Sprintf("{\"domain\": %d}", domain)
			_, err := apiClient.sendRequest("POST", endpoint, s)
			if err != nil {
				log.Printf("resource_level27_mailgroup.go: No Mailgroup Domain Link created: %s", err)
			}
			linkedDomains = append(linkedDomains, domain.(int))
		} else {
			linkedDomains = append(linkedDomains, domain.(int))
		}
	}
	// loop over the old list to see if we need to delete some domains
	for _, domain := range linkedDomainsOriginal {
		if !domains.Contains(domain) {
			// delete the domain
			endpoint := fmt.Sprintf("mailgroups/%s/domains/%d", d.Id(), domain)
			log.Printf("resource_level27_mailgroup.go: Delete linked domain: %d", domain)
			_, err := apiClient.sendRequest("DELETE", endpoint, "")
			if err != nil {
				log.Printf("resource_level27_mailgroup.go: No Mailgroup Domain Link created: %s", err)
			}
			// remove it from linkedDomains
			for i, v := range linkedDomains {
				if v == domain {
					linkedDomains = append(linkedDomains[:i], linkedDomains[i+1:]...)
					break
				}
			}
		} else {
			linkedDomains = append(linkedDomains, domain)
		}
	}
	return linkedDomains
}

func splatLinkedDomains(domains []interface{}) []int {
	linkedDomains := []int{}
	for _, domain := range domains {
		id := int(domain.(map[string]interface{})["id"].(float64))
		linkedDomains = append(linkedDomains, id)
	}
	return linkedDomains
}

func getMailgroup(d *schema.ResourceData, meta interface{}) Mailgroup {
	var mailgroup Mailgroup
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%s", d.Id())

	body, err := apiClient.sendRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_mailgroup.go: No Mailgroup found: %s", d.Id())
		d.SetId("")
	}
	err = json.Unmarshal(body, &mailgroup)

	return mailgroup
}

func containsInt(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
