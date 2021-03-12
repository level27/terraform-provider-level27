package level27

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

func resourceLevel27CookbookMysql() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookMysqlCreate,
		ReadContext:   resourceLevel27CookbookMysqlRead,
		DeleteContext: resourceLevel27CookbookMysqlDelete,
		UpdateContext: resourceLevel27CookbookMysqlUpdate,
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
			"query_cache_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     32,
				Description: `Query cache size`,
			},
			"query_cache_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: `the maximum size of a single resultset in the cache`,
			},
			"table_cache": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     400,
				Description: `Table cache`,
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "3306",
				Description: `Port`,
			},
			"bind_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "localhost",
				Description: `Bind address`,
			},
			"version": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"5.7", "8.0"}, true),
				},
				Required:    true,
				Description: `Mysql version`,
			},
		},
	}
}

func resourceLevel27CookbookMysqlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mysql.go: Create cookbook")

	request := cookbookCreateParameters(d, "mysql")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_mysql.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookMysqlRead(ctx, d, meta)
}

func resourceLevel27CookbookMysqlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mysql.go: Read cookbook")

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
	d.Set("query_cache_size", cookbook.Cookbook.Cookbookparameters["query_cache_size"].Value)
	d.Set("query_cache_limit", cookbook.Cookbook.Cookbookparameters["query_cache_limit"].Value)
	d.Set("table_cache", cookbook.Cookbook.Cookbookparameters["table_cache"].Value)
	d.Set("port", cookbook.Cookbook.Cookbookparameters["port"].Value)
	d.Set("bind_address", cookbook.Cookbook.Cookbookparameters["bind_address"].Value)
	if cookbook.Cookbook.Cookbookparameters["versions"].Value != "default" {
		d.Set("versions", cookbook.Cookbook.Cookbookparameters["versions"].Value)
	}

	return nil
}

func resourceLevel27CookbookMysqlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mysql.go: Update cookbook")

	request := cookbookCreateParameters(d, "mysql")

	log.Printf("resource_level27_cookbook_mysql.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookMysqlRead(ctx, d, meta)
}

func resourceLevel27CookbookMysqlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_mysql.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
