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

func resourceLevel27CookbookPhp() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookPhpCreate,
		ReadContext:   resourceLevel27CookbookPhpRead,
		DeleteContext: resourceLevel27CookbookPhpDelete,
		UpdateContext: resourceLevel27CookbookPhpUpdate,
		// https://www.terraform.io/docs/extend/resources/import.html
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
			"versions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"5.3", "5.6", "7.0", "7.1", "7.2", "7.3", "7.4", "8.0"}, true),
				},
				Required:    true,
				Description: `PHP versions`,
			},
			"apache_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "8080",
				Description: `apache listening port`,
			},
			"apache_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `apache listening address`,
			},
			"apache_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `ssl`,
			},
			"memory_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     256,
				Description: `Memory limit`,
			},
			"max_execution_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: `Max execution time`,
			},
			"upload_max_filesize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: `Max upload filesize`,
			},
			"disable_functions": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "passthru,escapeshellarg,escapeshellcmd,popen,pcntl_exec,dl,proc_open,proc_close,proc_get_status,proc_nice,proc_terminate",
				Description: `Disabled PHP functions`,
			},
			"expected_traffic": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "low",
				ValidateFunc: validation.StringInSlice([]string{"low", "medium", "high"}, true),
				Description:  `Traffic`,
			},
			"max_session_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1440,
				Description: `Max php session lifetime`,
			},
			"phpmyadmin_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Enable phpmyadmin`,
			},
			"process_manager": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ondemand",
				ValidateFunc: validation.StringInSlice([]string{"ondemand", "dynamic"}, true),
				Description:  `PHP process manager`,
			},
			"display_errors": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "no",
				ValidateFunc: validation.StringInSlice([]string{"no", "yes"}, true),
				Description:  `show display errors`,
			},
		},
	}
}

func resourceLevel27CookbookPhpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_php.go: Create cookbook")

	request := cookbookCreateParameters(d, "php")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_php.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookPhpRead(ctx, d, meta)
}

func resourceLevel27CookbookPhpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_php.go: Read cookbook")

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
	d.Set("versions", cookbook.Cookbook.Cookbookparameters["versions"].Value)
	d.Set("apache_port", cookbook.Cookbook.Cookbookparameters["apache:port"].Value)
	d.Set("apache_ip", cookbook.Cookbook.Cookbookparameters["apache:bind"].Value)
	d.Set("apache_ssl", cookbook.Cookbook.Cookbookparameters["apache:ssl"].Value)
	d.Set("memory_limit", cookbook.Cookbook.Cookbookparameters["memory_limit"].Value)
	d.Set("max_execution_time", cookbook.Cookbook.Cookbookparameters["max_execution_time"].Value)
	d.Set("upload_max_filesize", cookbook.Cookbook.Cookbookparameters["upload_max_filesize"].Value)
	d.Set("disable_functions", cookbook.Cookbook.Cookbookparameters["disable_functions"].Value)
	d.Set("expected_traffic", cookbook.Cookbook.Cookbookparameters["expected_traffic"].Value)
	d.Set("max_session_lifetime", cookbook.Cookbook.Cookbookparameters["max_session_lifetime"].Value)
	d.Set("phpmyadmin_enable", cookbook.Cookbook.Cookbookparameters["phpmyadmin:enable"].Value)
	d.Set("process_manager", cookbook.Cookbook.Cookbookparameters["process_manager"].Value)
	d.Set("display_errors", cookbook.Cookbook.Cookbookparameters["display_errors"].Value)

	return nil
}

func resourceLevel27CookbookPhpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_php.go: Update cookbook")

	request := cookbookCreateParameters(d, "php")

	log.Printf("resource_level27_cookbook_php.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookPhpRead(ctx, d, meta)
}

func resourceLevel27CookbookPhpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_php.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
