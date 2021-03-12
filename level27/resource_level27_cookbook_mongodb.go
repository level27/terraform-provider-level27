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

func resourceLevel27CookbookMongodb() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookMongodbCreate,
		ReadContext:   resourceLevel27CookbookMongodbRead,
		DeleteContext: resourceLevel27CookbookMongodbDelete,
		UpdateContext: resourceLevel27CookbookMongodbUpdate,
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
				Default:     "27017",
				Description: `Port`,
			},
			"listen": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Listen address`,
			},
		},
	}
}

func resourceLevel27CookbookMongodbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mongodb.go: Create cookbook")

	request := cookbookCreateParameters(d, "mongodb")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_mongodb.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookMongodbRead(ctx, d, meta)
}

func resourceLevel27CookbookMongodbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mongodb.go: Read cookbook")

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
	d.Set("listen", cookbook.Cookbook.Cookbookparameters["listen"].Value)

	return nil
}

func resourceLevel27CookbookMongodbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mongodb.go: Update cookbook")

	request := cookbookCreateParameters(d, "mongodb")

	log.Printf("resource_level27_cookbook_mongodb.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookMongodbRead(ctx, d, meta)
}

func resourceLevel27CookbookMongodbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mongodb.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
