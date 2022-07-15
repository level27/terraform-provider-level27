package level27

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27DomainRecord() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27DomainRecordCreate,
		Read:   resourceLevel27DomainRecordRead,
		Delete: resourceLevel27DomainRecordDelete,
		Update: resourceLevel27DomainRecordUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Description: "Name of the domain",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the domain record",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain record",
				ForceNew:    false,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Content of the domain record",
				ForceNew:    false,
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Priority of the domain record (only for MX records)",
				ForceNew:    false,
			},
		},
	}
}

func resourceLevel27DomainRecordCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_record.go: Starting Create")

	request := DomainRecordRequest{
		Type:     d.Get("type").(string),
		Name:     d.Get("name").(string),
		Content:  d.Get("content").(string),
		Priority: d.Get("priority").(int),
	}
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domains/%s/records", d.Get("domain_id").(string))
	body, err := apiClient.sendRequest("POST", endpoint, request)

	if err != nil {
		log.Printf("resource_level27_domain_record.go: No Domain record created: %s", err)
		d.SetId("")
		return err
	}

	var domainRecord DomainRecord
	err = json.Unmarshal(body, &domainRecord)
	if err != nil {
		return err
	}

	// Create is OK, we now set the properties of the object
	d.SetId(strconv.Itoa(domainRecord.Record.ID))
	return resourceLevel27DomainRecordRead(d, meta)
}

func resourceLevel27DomainRecordRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_record.go: Starting Read")
	log.Printf("resource_level27_domain_record.go: Resource id %s", d.Id())

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domains/%s/records/%s", d.Get("domain_id").(string), d.Id())
	body, err := apiClient.sendRequest("GET", endpoint, nil)

	if err != nil {
		log.Printf("resource_level27_domain_record.go: No Domain record found: %s", d.Id())
		d.SetId("")
		return err
	}

	var domainRecord DomainRecord
	err = json.Unmarshal(body, &domainRecord)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(domainRecord.Record.ID))
	d.Set("name", domainRecord.Record.Name)
	d.Set("type", domainRecord.Record.Type)
	d.Set("content", domainRecord.Record.Content)
	d.Set("priority", domainRecord.Record.Priority)

	return nil
}

func resourceLevel27DomainRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_record.go: Starting Update")

	request := DomainRecordRequest{
		Type:     d.Get("type").(string),
		Name:     d.Get("name").(string),
		Content:  d.Get("content").(string),
		Priority: d.Get("priority").(int),
	}
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domains/%s/records/%s", d.Get("domain_id").(string), d.Id())
	_, err := apiClient.sendRequest("PUT", endpoint, request)

	if err != nil {
		log.Printf("resource_level27_domain_record.go: No Domain record updated: %s", err)
		d.SetId("")
		return err
	}

	// Update is OK, we now set the properties of the object
	return resourceLevel27DomainRecordRead(d, meta)
}

func resourceLevel27DomainRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_domain_record.go: Starting Delete")

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("domains/%s/records/%s", d.Get("domain_id").(string), d.Id())
	_, err := apiClient.sendRequest("DELETE", endpoint, nil)

	if err != nil {
		log.Printf("resource_level27_domain_record.go: No Domain record updated: %s", err)
		return err
	}

	// Delete is OK, we now set the properties of the object
	d.SetId("")
	return nil
}
*/
