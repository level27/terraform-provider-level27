package level27

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27Mailbox() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27MailboxCreate,
		Read:   resourceLevel27MailboxRead,
		Delete: resourceLevel27MailboxDelete,
		Update: resourceLevel27MailboxUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"password": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ooo_enabled": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"ooo_subject": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ooo_text": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"mailgroup_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"addresses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceLevel27MailboxCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailbox.go: Create resource")

	var mailbox Mailbox
	request := MailboxRequest{
		Name:       d.Get("name").(string),
		Password:   d.Get("password").(string),
		OooEnabled: d.Get("ooo_enabled").(bool),
		OooSubject: d.Get("ooo_subject").(string),
		OooText:    d.Get("ooo_text").(string),
	}

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%s/mailboxes", d.Get("mailgroup_id"))

	body, err := apiClient.sendRequest("POST", endpoint, request.String())

	if err != nil {
		log.Printf("resource_level27_mailbox.go: No Mailbox created")
		d.SetId("")
		return err
	}

	err = json.Unmarshal(body, &mailbox)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mailbox.Mailbox.ID))
	log.Printf("resource_level27_mailbox.go: Mailbox created: %s", d.Id())
	log.Printf("resource_level27_mailbox.go: Mailbox addresses: %v", d.Get("addresses"))
	addresses := d.Get("addresses").(*schema.Set)

	updateAddresses(addresses, d, meta, strconv.Itoa(mailbox.Mailbox.Mailgroup.ID))
	d.Set("addresses", addresses)

	return resourceLevel27MailboxRead(d, meta)
}

func resourceLevel27MailboxRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailbox.go: Read resource")

	var mailbox Mailbox
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%s/mailboxes/%s", d.Get("mailgroup_id"), d.Id())

	body, err := apiClient.sendRequest("GET", endpoint, nil)

	if err != nil {
		log.Printf("resource_level27_mailbox.go: No Mailbox created")
		d.SetId("")
		return err
	}

	err = json.Unmarshal(body, &mailbox)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mailbox.Mailbox.ID))
	d.Set("mailgroup_id", strconv.Itoa(mailbox.Mailbox.Mailgroup.ID))
	d.Set("name", mailbox.Mailbox.Name)
	d.Set("ooo_enabled", mailbox.Mailbox.OooEnabled)
	d.Set("ooo_subject", mailbox.Mailbox.OooSubject)
	d.Set("ooo_text", mailbox.Mailbox.OooText)
	d.Set("addresses", splatAddresses(originalAddresses(strconv.Itoa(mailbox.Mailbox.Mailgroup.ID), strconv.Itoa(mailbox.Mailbox.ID), meta)))
	return nil
}

func resourceLevel27MailboxUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailbox.go: Update resource id: %s", d.Id())

	addresses := d.Get("addresses").(*schema.Set)
	updateAddresses(addresses, d, meta, d.Get("mailgroup_id").(string))
	d.Set("addresses", addresses)

	return resourceLevel27MailboxRead(d, meta)
}

func resourceLevel27MailboxDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_mailbox.go: Delete resource id: %s", d.Id())

	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%s/mailboxes/%s", d.Get("mailgroup_id"), d.Id())

	_, err := apiClient.sendRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_mailbox.go: No Mailbox deleted: %s", d.Id())
		d.SetId("")
	}

	return nil
}

// Custom functions //
func updateAddresses(addresses *schema.Set, d *schema.ResourceData, meta interface{}, mailgroup string) {
	log.Printf("resource_level27_mailbox.go: updateAddresses")

	addressesOriginal := originalAddresses(mailgroup, d.Id(), meta)

	apiClient := meta.(*Client)
	//linkedAddresses := []string{}

	// loop over the new list to see if we need to add some domains
	for _, address := range addresses.List() {
		if !(containsAddress(addressesOriginal, address.(string))) {
			//add the domain
			endpoint := fmt.Sprintf("mailgroups/%s/mailboxes/%s/addresses", mailgroup, d.Id())
			log.Printf("resource_level27_mailbox.go: Link new address; %s", address.(string))
			s := fmt.Sprintf("{\"address\": \"%s\"}", address)
			_, err := apiClient.sendRequest("POST", endpoint, s)
			if err != nil {
				log.Printf("resource_level27_mailbox.go: No Mailbox Address linked: %s", err)
			}
		}
		log.Printf("resource_level27_mailbox.go: Mailbox Address linked: %s", address.(string))
	}
	// loop over the old list to see if we need to delete some domains
	for _, address := range addressesOriginal.MailboxAddresses {
		if !addresses.Contains(address.Address) {
			log.Printf("resource_level27_mailbox.go: Delete linked address: %s", address.Address)
			// delete the domain
			endpoint := fmt.Sprintf("mailgroups/%s/mailboxes/%s/addresses/%d", mailgroup, d.Id(), address.ID)
			_, err := apiClient.sendRequest("DELETE", endpoint, "")
			if err != nil {
				log.Printf("resource_level27_mailbox.go: No Mailgroup Domain Link created: %s", err)
			}
			// remove it from linkedDomains
			addresses.Remove(address.Address)
		}
	}
}

func originalAddresses(mailgroup string, mailbox string, meta interface{}) MailboxAddresses {
	var mailboxaddresses MailboxAddresses
	apiClient := meta.(*Client)
	endpoint := fmt.Sprintf("mailgroups/%s/mailboxes/%s/addresses", mailgroup, mailbox)

	body, err := apiClient.sendRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("resource_level27_mailbox.go: No MailboxAddresses found.")
	}
	err = json.Unmarshal(body, &mailboxaddresses)

	return mailboxaddresses
}

func splatAddresses(addresses MailboxAddresses) []string {
	splat := []string{}
	for _, address := range addresses.MailboxAddresses {
		splat = append(splat, address.Address)
	}
	return splat
}

func containsAddress(a MailboxAddresses, x string) bool {
	for _, n := range a.MailboxAddresses {
		if x == n.Address {
			return true
		}
	}
	return false
}
*/
