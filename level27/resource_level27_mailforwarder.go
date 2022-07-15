package level27

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27Mailforwarder() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27MailforwarderCreate,
		Read:   resourceLevel27MailforwarderRead,
		Delete: resourceLevel27MailforwarderDelete,
		Update: resourceLevel27MailforwarderUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mailgroup_id": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"address": {
				Required: true,
				Type:     schema.TypeString,
			},
			"destination": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceLevel27MailforwarderCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailforwarder.go: Create resource.")

	var mailforwarder Mailforwarder
	request := MailforwarderRequest{
		Address:     d.Get("address").(string),
		Destination: d.Get("destination").(string),
	}
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", d.Get("mailgroup_id"))

	body, err := apiClient.sendRequest("POST", endpoint, request.String())

	if err != nil {
		log.Printf("resource_level27_mailforwarder.go: No Mailforwarder created: %s", err)
		d.SetId("")
		return nil
	}

	err = json.Unmarshal(body, &mailforwarder)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mailforwarder.Mailforwarder.ID))
	d.Set("mailgroup_id", mailforwarder.Mailforwarder.Mailgroup.ID)
	return resourceLevel27MailforwarderRead(d, meta)
}

func resourceLevel27MailforwarderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailforwarder.go: Read resource id: %s", d.Id())

	var mailforwarder Mailforwarder
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%s", d.Get("mailgroup_id"), d.Id())

	body, err := apiClient.sendRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_mailforwarder.go: No Mailforwarder found: %s", d.Id())
		d.SetId("")
	}
	err = json.Unmarshal(body, &mailforwarder)
	if err != nil {
		return err
	}
	log.Printf("resource_level27_mailforwarder.go: Read routine called. Object built: %d", mailforwarder.Mailforwarder.ID)

	d.SetId(strconv.Itoa(mailforwarder.Mailforwarder.ID))
	d.Set("mailgroup_id", mailforwarder.Mailforwarder.Mailgroup.ID)
	d.Set("address", mailforwarder.Mailforwarder.Address)
	d.Set("destination", strings.Join(mailforwarder.Mailforwarder.Destination, ","))

	return nil
}

func resourceLevel27MailforwarderUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailforwarder.go: Update resource id: %s", d.Id())

	apiClient := meta.(*Client)
	request := MailforwarderRequest{
		Address:     d.Get("address").(string),
		Destination: d.Get("destination").(string),
	}
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%s", d.Get("mailgroup_id"), d.Id())

	_, err := apiClient.sendRequest("PUT", endpoint, request.String())
	if err != nil {
		log.Printf("resource_level27_mailforwarder.go: No Mailforwarder updated: %s", d.Id())
		d.SetId("")
	}

	return resourceLevel27MailforwarderRead(d, meta)
}

func resourceLevel27MailforwarderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailforwarder.go: Delete resource id: %s", d.Id())

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%s", d.Get("mailgroup_id"), d.Id())

	_, err := apiClient.sendRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_mailforwarder.go: No Mailforwarder deleted: %s", d.Id())
		d.SetId("")
	}

	return nil
}
*/
