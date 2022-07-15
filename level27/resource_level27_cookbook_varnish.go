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
)

func resourceLevel27CookbookVarnish() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookVarnishCreate,
		ReadContext:   resourceLevel27CookbookVarnishRead,
		DeleteContext: resourceLevel27CookbookVarnishDelete,
		UpdateContext: resourceLevel27CookbookVarnishUpdate,
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
			"backend_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Backend ip address`,
			},
			"backend_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "8080",
				Description: `Backend port`,
			},
			"frontend_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.0.0",
				Description: `Varnish listen address`,
			},
			"frontend_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "80",
				Description: `Varnish port`,
			},
			"storage_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "256",
				Description: `Cache size in MB`,
			},
			"admin_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Admin listen address`,
			},
			"admin_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "6082",
				Description: `Admin port`,
			},
		},
	}
}

func resourceLevel27CookbookVarnishCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_varnish.go: Create cookbook")

	request := cookbookCreateParameters(d, "web_varnish")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_varnish.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookVarnishRead(ctx, d, meta)
}

func resourceLevel27CookbookVarnishRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_varnish.go: Read cookbook")

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
	d.Set("backend_ip", cookbook.Cookbook.Cookbookparameters["backend_ip"].Value)
	d.Set("backend_port", cookbook.Cookbook.Cookbookparameters["backend_port"].Value)
	d.Set("frontend_ip", cookbook.Cookbook.Cookbookparameters["frontend_ip"].Value)
	d.Set("frontend_port", cookbook.Cookbook.Cookbookparameters["frontend_port"].Value)
	d.Set("storage_size", cookbook.Cookbook.Cookbookparameters["storage_size"].Value)
	d.Set("admin_ip", cookbook.Cookbook.Cookbookparameters["admin_ip"].Value)
	d.Set("admin_port", cookbook.Cookbook.Cookbookparameters["admin_port"].Value)

	return nil
}

func resourceLevel27CookbookVarnishUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_varnish.go: Update cookbook")

	request := cookbookCreateParameters(d, "web_varnish")

	log.Printf("resource_level27_cookbook_varnish.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookVarnishRead(ctx, d, meta)
}

func resourceLevel27CookbookVarnishDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_varnish.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
*/
