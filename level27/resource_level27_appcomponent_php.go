package level27

/*
import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var appSchema = map[string]*schema.Schema{
	"app_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"username": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"system_id": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"system_group_id"},
	},
	"system_group_id": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"system_id"},
	},
	"version": {
		Optional: true,
		Type:     schema.TypeString,
	},
	"path": {
		Optional: true,
		Type:     schema.TypeString,
		Default:  "/public_html",
	},
	"sshkeys": {
		Optional: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"url": {
		Optional: true,
		Type:     schema.TypeSet,

		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"content": {
					Type:     schema.TypeString,
					Required: true,
				},
				"ssl_force": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"ssl_certificate": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"handle_dns": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	},
}

func resourceLevel27AppComponentPhp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLevel27AppComponentPhpCreate,
		ReadContext:   resourceLevel27AppComponentPhpRead,
		DeleteContext: resourceLevel27AppComponentPhpDelete,
		UpdateContext: resourceLevel27AppComponentPhpUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: appSchema,
	}
}

func resourceLevel27AppComponentPhpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Printf("resource_level27_appcomponent_php.go: Read resource")

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*Client)
	ac := apiClient.AppComponent("GET", appID, appComponentID, nil)

	log.Printf("resource_level27_appcomponent_php.go: Read routine called. Object built: %d", ac.Component.ID)

	d.SetId(fmt.Sprintf("%d:%d", ac.Component.App.ID, ac.Component.ID))
	d.Set("name", ac.Component.Name)
	d.Set("username", ac.Component.Appcomponentparameters.User)
	systemID := strconv.Itoa(ac.Component.Systems[0].ID)
	if systemID == "0" {
		systemID = ""
	}
	d.Set("system_id", systemID)

	systemGroupID := strconv.Itoa(ac.Component.Systemgroup)
	if systemGroupID == "0" {
		systemGroupID = ""
	}
	d.Set("system_group_id", systemGroupID)

	d.Set("version", ac.Component.Appcomponentparameters.Version)
	d.Set("path", ac.Component.Appcomponentparameters.Path)

	sshKeysRead := ac.Component.Appcomponentparameters.Sshkeys
	sshKeysNew := make([]string, len(sshKeysRead))
	for i, s := range sshKeysRead {
		sshKeysNew[i] = strconv.Itoa(s.ID)
	}
	d.Set("sshkeys", sshKeysNew)

	log.Printf("resource_level27_appcomponent_php.go: Read resource: %v", d.State())

	urls, err1 := apiClient.AppComponentURL("GETALL", appID, appComponentID, nil, nil)
	if err1 != nil {
		return err1
	}
	d.Set("url", flattenURLs(urls, d.Id()))

	log.Printf("resource_level27_appcomponent_php.go: Read resource: %v", d.State())

	return diags
}

func resourceLevel27AppComponentPhpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_appcomponent_php.go: Create resource")

	request := AppComponentRequest{
		Name:        d.Get("name").(string),
		Type:        "php",
		Category:    "config",
		System:      d.Get("system_id").(string),
		SystemGroup: d.Get("system_group_id").(string),
		Version:     d.Get("version").(string),
		Path:        d.Get("path").(string),
		SSHKeys:     d.Get("sshkeys").([]interface{}),
	}

	apiClient := meta.(*Client)
	ac := apiClient.AppComponent("CREATE", d.Get("app_id"), nil, request.String())

	d.SetId(fmt.Sprintf("%d:%d", ac.Component.App.ID, ac.Component.ID))
	return resourceLevel27AppComponentPhpRead(ctx, d, meta)
}

func resourceLevel27AppComponentPhpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_appcomponent_php.go: Update resource")

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") ||
		d.HasChange("system_id") ||
		d.HasChange("system_group_id") ||
		d.HasChange("version") ||
		d.HasChange("path") ||
		d.HasChange("sshkeys") {
		request := AppComponentRequest{
			Name:        d.Get("name").(string),
			Type:        "php",
			Category:    "config",
			System:      d.Get("system_id").(string),
			SystemGroup: d.Get("system_group_id").(string),
			Version:     d.Get("version").(string),
			Path:        d.Get("path").(string),
			SSHKeys:     d.Get("sshkeys").([]interface{}),
		}

		apiClient := meta.(*Client)
		apiClient.AppComponent("UPDATE", appID, appComponentID, request.String())
	}

	if d.HasChange("url") {
		old, new := d.GetChange("url")
		err := updateURLs(old, new, d, meta)
		if err != nil {
			return err
		}
	}

	return resourceLevel27AppComponentPhpRead(ctx, d, meta)
}

func resourceLevel27AppComponentPhpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_appcomponent_php.go: Delete resource")

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*Client)
	apiClient.AppComponent("DELETE", appID, appComponentID, nil)

	return nil
}

func updateURLs(old interface{}, new interface{}, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)
	log.Printf("resource_level27_appcomponent_php.go: URLs CHANGES OLD %v", oldSet.List())
	log.Printf("resource_level27_appcomponent_php.go: URLs CHANGES NEW %v", newSet.List())

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*Client)
	add := newSet.Difference(oldSet).List()
	delete := oldSet.Difference(newSet).List()

	log.Printf("resource_level27_appcomponent_php.go: URLs ADD    %v", add)
	log.Printf("resource_level27_appcomponent_php.go: URLs DELETE %v", delete)

	for _, d := range delete {
		url := d.(map[string]interface{})
		log.Printf("resource_level27_appcomponent_php.go: URLs DELETE    %s:%s:%s", appID, appComponentID, url["id"].(string))
		_, err := apiClient.AppComponentURL("DELETE", appID, appComponentID, url["id"].(string), nil)
		if err != nil {
			return err
		}
	}

	for _, a := range add {
		url := a.(map[string]interface{})
		certificateID, err := strconv.Atoi(url["ssl_certificate"].(string))
		if err != nil {
			certificateID = 0
		}
		request := URLRequest{
			Content:        url["content"].(string),
			SSLForce:       url["ssl_force"].(bool),
			SSLCertificate: certificateID,
			HandleDNS:      url["handle_dns"].(bool),
		}

		log.Printf("resource_level27_appcomponent_php.go: URLs ADD    %s", request.String())
		_, err1 := apiClient.AppComponentURL("CREATE", appID, appComponentID, nil, request.String())
		if err1 != nil {
			return err1
		}
	}

	return nil
}

func flattenURLs(in URL, componentID string) []interface{} {
	var out = make([]interface{}, len(in.URLs), len(in.URLs))
	for i, url := range in.URLs {
		u := make(map[string]interface{})
		u["content"] = url.Content
		u["id"] = fmt.Sprintf("%d", url.ID)
		u["ssl_force"] = url.SslForce
		u["ssl_certificate"] = strconv.Itoa(url.SslCertificate.ID)
		if u["ssl_certificate"] == "0" {
			u["ssl_certificate"] = ""
		}
		u["handle_dns"] = false
		if url.Type == "linked" {
			u["handle_dns"] = true
		}
		out[i] = u
	}

	return out
}
*/
