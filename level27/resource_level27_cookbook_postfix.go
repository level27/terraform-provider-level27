package level27

/*
import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27CookbookPostfix() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceLevel27CookbookPostfixCreate,
		ReadContext:   resourceLevel27CookbookPostfixRead,
		DeleteContext: resourceLevel27CookbookPostfixDelete,
		UpdateContext: resourceLevel27CookbookPostfixUpdate,
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
			"listening_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1",
				Description: `Listening address of the mail server`,
			},
			"contentfilter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Content filter`,
			},
		},
	}
}

func resourceLevel27CookbookPostfixCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_postfix.go: Create cookbook")

	request := cookbookCreateParameters(d, "postfix")
	systemID := d.Get("system_id").(string)

	log.Printf("resource_level27_cookbook_postfix.go: %s", request.String())

	err := cookbookCreate(ctx, d, meta, systemID, request)
	if err != nil {
		return err
	}

	return resourceLevel27CookbookPostfixRead(ctx, d, meta)
}

func resourceLevel27CookbookPostfixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_postfix.go: Read cookbook")

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
	d.Set("listening_address", cookbook.Cookbook.Cookbookparameters["listening_address"].Value)
	d.Set("contentfilter", cookbook.Cookbook.Cookbookparameters["contentfilter"].Value)

	return nil
}

func resourceLevel27CookbookPostfixUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_postfix.go: Update cookbook")

	request := cookbookCreateParameters(d, "postfix")

	log.Printf("resource_level27_cookbook_postfix.go: %s", request.String())
	diags := cookbookUpdate(ctx, d, meta, request)
	if diags != nil {
		return diags
	}

	return resourceLevel27CookbookPostfixRead(ctx, d, meta)
}

func resourceLevel27CookbookPostfixDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource_level27_cookbook_postfix.go: Delete resource")
	return cookbookDelete(d.Id(), meta)
}
*/
