package level27

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27CookbookHaproxy() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookHaproxyCreate,
		ReadContext:   resourceLevel27CookbookHaproxyRead,
		DeleteContext: resourceLevel27CookbookHaproxyDelete,
		UpdateContext: resourceLevel27CookbookHaproxyUpdate,
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
			"haip_ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `High available ipv4 address`,
			},
			"haip_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `High available ipv6 address`,
			},
			"haip_routerid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Router id of the high available ip addresses`,
			},
			"backend_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "8080",
				Description: `Port of the backend the loadbalancer should proxy to`,
			},
			"backend_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Ip of the backend the loadbalancer should proxy to`,
			},
			"frontend_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "80",
				Description: `Port the loadbalancer should bind to`,
			},
			"primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Is the loadbalancer the primary one`,
			},
			"varnish": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Is varnish present on the webservers`,
			},
			"expected_traffic": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: `The amount of traffic that should be expected`,
			},
		},
	}
}

func resourceLevel27CookbookHaproxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_haproxy.go: Create cookbook")

	request := cookbookCreateParameters(d, "haproxy")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_haproxy.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookHaproxyRead(ctx, d, meta)
}

func resourceLevel27CookbookHaproxyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_haproxy.go: Read cookbook")

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
	if cookbook.Cookbook.Cookbookparameters["haip:ipv4"].Value != "nil" {
		d.Set("haip_ipv4", cookbook.Cookbook.Cookbookparameters["haip:ipv4"].Value)
	}
	if cookbook.Cookbook.Cookbookparameters["haip:ipv6"].Value != "nil" {
		d.Set("haip_ipv6", cookbook.Cookbook.Cookbookparameters["haip:ipv6"].Value)
	}
	if cookbook.Cookbook.Cookbookparameters["haip:routerid"].Value != "nil" {
		d.Set("haip_routerid", cookbook.Cookbook.Cookbookparameters["haip:routerid"].Value)
	}
	d.Set("backend_port", cookbook.Cookbook.Cookbookparameters["backend:port"].Value)
	d.Set("backend_ip", cookbook.Cookbook.Cookbookparameters["backend:ip"].Value)
	d.Set("frontend_port", cookbook.Cookbook.Cookbookparameters["port"].Value)
	d.Set("primary", cookbook.Cookbook.Cookbookparameters["primary"].Value)
	d.Set("varnish", cookbook.Cookbook.Cookbookparameters["varnish"].Value)
	d.Set("expected_traffic", cookbook.Cookbook.Cookbookparameters["expected_traffic"].Value)

	return nil
}

func resourceLevel27CookbookHaproxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_haproxy.go: Update cookbook")

	request := cookbookCreateParameters(d, "haproxy")

	log.Printf("resource_level27_cookbook_haproxy.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookHaproxyRead(ctx, d, meta)
}

func resourceLevel27CookbookHaproxyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_haproxy.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
