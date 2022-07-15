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

func resourceLevel27CookbookMemcached() *schema.Resource {

	return &schema.Resource{
		Description:   `Installs and configures MongoDB cookbook on a system`,
		CreateContext: resourceLevel27CookbookMemcachedCreate,
		ReadContext:   resourceLevel27CookbookMemcachedRead,
		DeleteContext: resourceLevel27CookbookMemcachedDelete,
		UpdateContext: resourceLevel27CookbookMemcachedUpdate,
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
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "11211",
				Description: `Port`,
			},
			"listen_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Listen adress`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     64,
				Description: `Memory size`,
			},
		},
	}
}

func resourceLevel27CookbookMemcachedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_memcached.go: Create cookbook")

	request := cookbookCreateParameters(d, "memcached")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_memcached.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookMemcachedRead(ctx, d, meta)
}

func resourceLevel27CookbookMemcachedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_memcached.go: Read cookbook")

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
	d.Set("port", cookbook.Cookbook.Cookbookparameters["port"].Value)
	d.Set("listen_address", cookbook.Cookbook.Cookbookparameters["listen_address"].Value)
	d.Set("memory", cookbook.Cookbook.Cookbookparameters["memory"].Value)

	return nil
}

func resourceLevel27CookbookMemcachedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_memcached.go: Update cookbook")

	request := cookbookCreateParameters(d, "memcached")

	log.Printf("resource_level27_cookbook_memcached.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookMemcachedRead(ctx, d, meta)
}

func resourceLevel27CookbookMemcachedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_memcached.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
*/
