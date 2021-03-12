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

func resourceLevel27CookbookDocker() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookDockerCreate,
		ReadContext:   resourceLevel27CookbookDockerRead,
		DeleteContext: resourceLevel27CookbookDockerDelete,
		UpdateContext: resourceLevel27CookbookDockerUpdate,
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
		},
	}
}

func resourceLevel27CookbookDockerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_docker.go: Create cookbook")

	request := cookbookCreateParameters(d, "docker")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_docker.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookDockerRead(ctx, d, meta)
}

func resourceLevel27CookbookDockerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_docker.go: Read cookbook")

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

	return nil
}

func resourceLevel27CookbookDockerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_docker.go: Update cookbook")

	request := cookbookCreateParameters(d, "docker")

	log.Printf("resource_level27_cookbook_docker.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookDockerRead(ctx, d, meta)
}

func resourceLevel27CookbookDockerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_docker.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
