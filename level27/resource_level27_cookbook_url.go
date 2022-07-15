package level27

/*
import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLevel27CookbookURL() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookURLCreate,
		ReadContext:   resourceLevel27CookbookURLRead,
		DeleteContext: resourceLevel27CookbookURLDelete,
		UpdateContext: resourceLevel27CookbookURLUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: cookbookImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cookbook_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the system where to install this cookbook.`,
			},
			"proxy_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "nginx",
				Description: `Proxy service`,
			},
			"expected_traffic": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "low",
				ValidateFunc: validation.StringInSlice([]string{"low", "medium", "high"}, true),
				Description:  `Traffic`,
			},
			"backend_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "8080",
				Description: `Backend port`,
			},
			"backend_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Backend ip`,
			},
			"default_vhost": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Create default nginx vhost`,
			},
			"max_upload_filesize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: `Max upload filesize`,
			},
		},
	}
}

func resourceLevel27CookbookURLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_url.go: Create cookbook")

	request := cookbookCreateParameters(d, "url")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_url.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookURLRead(ctx, d, meta)
}

func resourceLevel27CookbookURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_url.go: Read cookbook")

	apiClient := meta.(*Client)
	systemID, cookbookID, err := cookbookParseID(d.Id())
	if err != nil {
		return err
	}

	cookbook, err := apiClient.SystemCookbook("GET", systemID, cookbookID, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d:%d", cookbook.Cookbook.System.ID, cookbook.Cookbook.ID))
	d.Set("cookbook_id", strconv.Itoa(cookbook.Cookbook.ID))
	d.Set("system_id", strconv.Itoa(cookbook.Cookbook.System.ID))
	d.Set("proxy_service", cookbook.Cookbook.Cookbookparameters["proxy_service"].Value)
	d.Set("expected_traffic", cookbook.Cookbook.Cookbookparameters["expected_traffic"].Value)
	d.Set("backend_port", cookbook.Cookbook.Cookbookparameters["backend:port"].Value)
	d.Set("backend_ip", cookbook.Cookbook.Cookbookparameters["backend:ip"].Value)
	d.Set("default_vhost", cookbook.Cookbook.Cookbookparameters["default_vhost"].Value)
	d.Set("max_upload_filesize", cookbook.Cookbook.Cookbookparameters["max_upload_filesize"].Value)

	return nil
}

func resourceLevel27CookbookURLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_url.go: Update cookbook")

	request := cookbookCreateParameters(d, "url")

	log.Printf("resource_level27_cookbook_url.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookURLRead(ctx, d, meta)
}

func resourceLevel27CookbookURLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_url.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
*/
