package level27

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceSystemType struct {
}

// GetSchema implements tfsdk.ResourceType
func (resourceSystemType) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:          types.StringType,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: `The name of the system.`,
			},
			"hostname": {
				Type:       types.StringType,
				Optional:   true,
				Validators: []tfsdk.AttributeValidator{validateRegex(`(\b(?:(?:[^.-/]{0,1})[\w-]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63})\b)`)},
			},
			"fqdn": {
				Type:          types.StringType,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
				Description:   `The FQDN of the system.`,
			},
			"organisation": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
				Description:   `The ID of the organisation that owns the system.`,
			},
			"management_type": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{DefaultValue(tfString("basic"))},
				Validators:    []tfsdk.AttributeValidator{validateStringInSlice([]string{"basic", "professional", "enterprise", "professional_level27"})},
			},
			"system_image": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
				Description:   `The ID of the system image to use for creation.`,
			},
			"system_provider_configuration": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
				Description:   `The ID of the system provider to use for creation.`,
			},
			"zone": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{validateID()},
				Description:   `The ID of the zone where to create the system.`,
			},
			"networks": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:          types.StringType,
						Computed:      true,
						PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.UseStateForUnknown()},
						Validators:    []tfsdk.AttributeValidator{validateID()},
					},
					"ip": {
						Type:     types.StringType,
						Computed: true,
					},
					"ipv6": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
			"cpu": {
				Type:        types.Int64Type,
				Required:    true,
				Description: `The number of CPUs.`,
			},
			"memory": {
				Type:        types.Int64Type,
				Required:    true,
				Description: `The amount of memory in GB.`,
			},
			"disk": {
				Type:        types.Int64Type,
				Required:    true,
				Description: `The amount of diskspace in GB.`,
			},
			"public_ipv4": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
			"public_ipv6": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
			"internal_ipv4": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
			"internal_ipv6": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
			"customer_ipv4": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
			"customer_ipv6": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
		},
	}, nil
}

// NewResource implements tfsdk.ResourceType
func (r resourceSystemType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceSystem{
		p: p.(*provider),
	}, nil
}

type resourceSystem struct {
	p *provider
}

// Create implements tfsdk.Resource
func (resourceSystem) Create(context.Context, tfsdk.CreateResourceRequest, *tfsdk.CreateResourceResponse) {
	panic("unimplemented")
}

// Delete implements tfsdk.Resource
func (resourceSystem) Delete(context.Context, tfsdk.DeleteResourceRequest, *tfsdk.DeleteResourceResponse) {
	panic("unimplemented")
}

// Read implements tfsdk.Resource
func (resourceSystem) Read(context.Context, tfsdk.ReadResourceRequest, *tfsdk.ReadResourceResponse) {
	panic("unimplemented")
}

// Update implements tfsdk.Resource
func (resourceSystem) Update(context.Context, tfsdk.UpdateResourceRequest, *tfsdk.UpdateResourceResponse) {
	panic("unimplemented")
}

/*
import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLevel27System() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27SystemCreate,
		Read:   resourceLevel27SystemRead,
		Delete: resourceLevel27SystemDelete,
		Update: resourceLevel27SystemUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the system.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the system.`,
			},
			"customer_fqdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`(\b(?:(?:[^.-/]{0,1})[\w-]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63})\b)`), "Valid FQDN required"),
				Description:  `The customer fqdn of the system.`,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fqdn of the system.`,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the organisation that owns the system.`,
			},
			"management_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "professional", "enterprise", "professional_level27"}, false),
				Description: `The management type of the system.
Possible values:
  - basic
  - professional
  - enterprise
  - professional_level27`,
			},
			"systemimage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the system image to use for creation.`,
			},
			"systemprovider_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the system provider to use for creation.`,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the zone where to create the system.`,
			},
			"auto_networks": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the zone where to create the system.`,
			},
			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of CPU's.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The amount of memory in GB.`,
			},
			"disk": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The amount of diskspace in GB.`,
			},
			"public_ipv4": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ipv6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internal_ipv4": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internal_ipv6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"customer_ipv4": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"customer_ipv6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceLevel27SystemCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_system.go: Create resource")

	apiClient := meta.(*Client)

	request := SystemRequest{
		Name:                        d.Get("name").(string),
		Organisation:                d.Get("organisation_id").(string),
		Systemimage:                 d.Get("systemimage_id").(string),
		SystemproviderConfiguration: d.Get("systemprovider_configuration_id").(string),
		Zone:                        d.Get("zone_id").(string),
		CPU:                         d.Get("cpu").(int),
		Memory:                      d.Get("memory").(int),
		Disk:                        d.Get("disk").(int),
		AutoNetworks:                d.Get("auto_networks").([]interface{}),
		CustomerFqdn:                d.Get("customer_fqdn").(string),
	}

	system := apiClient.System("CREATE", nil, request.String())
	d.SetId(strconv.Itoa(system.System.ID))

	createStateConf := &resource.StateChangeConf{
		Pending:                   []string{"stopped", "to_create", "shutting_down", "starting", "rebooting", "busy"},
		Target:                    []string{"running"},
		Refresh:                   serverCreatingStateRefresh(apiClient, system.System.ID),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	_, err := createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Instance not created. %s", err)
	}

	d.SetId(strconv.Itoa(system.System.ID))

	return resourceLevel27SystemRead(d, meta)
}

func resourceLevel27SystemRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_system.go: Read resource")

	apiClient := meta.(*Client)
	system := apiClient.System("GET", d.Id(), nil)

	d.SetId(strconv.Itoa(system.System.ID))
	d.Set("customer_fqdn", system.System.Hostname)
	d.Set("fqdn", system.System.Fqdn)
	d.Set("organisation_id", strconv.Itoa(system.System.Organisation.ID))
	d.Set("uid", system.System.UID)
	d.Set("systemimage_id", system.System.Systemimage.ID)
	d.Set("systemprovider_configuration_id", system.System.SystemproviderConfiguration.ID)
	d.Set("zone_id", system.System.Zone.ID)
	d.Set("name", system.System.Name)
	d.Set("cpu", system.System.CPU)
	d.Set("disk", system.System.Disk)
	d.Set("memory", system.System.Memory)
	// autonetworks

	networks := flattenNetworks(&system)
	d.Set("public_ipv4", networks["public_ipv4"])
	d.Set("public_ipv6", networks["public_ipv6"])
	d.Set("internal_ipv4", networks["internal_ipv4"])
	d.Set("internal_ipv6", networks["internal_ipv6"])
	d.Set("customer_ipv4", networks["customer_ipv4"])
	d.Set("customer_ipv6", networks["customer_ipv6"])

	return nil
}

func resourceLevel27SystemUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_system.go: Update system")

	apiClient := meta.(*Client)

	request := SystemRequest{
		Name:                        d.Get("name").(string),
		Organisation:                d.Get("organisation_id").(string),
		Systemimage:                 d.Get("systemimage_id").(string),
		SystemproviderConfiguration: d.Get("systemprovider_configuration_id").(string),
		Zone:                        d.Get("zone_id").(string),
		CPU:                         d.Get("cpu").(int),
		Memory:                      d.Get("memory").(int),
		Disk:                        d.Get("disk").(int),
		CustomerFqdn:                d.Get("customer_fqdn").(string),
	}

	_ = apiClient.System("UPDATE", d.Id(), request.UpdateString())
	systemID, _ := strconv.Atoi(d.Id())

	createStateConf := &resource.StateChangeConf{
		Pending:    []string{"starting", "rebooting", "busy"},
		Target:     []string{"ok"},
		Refresh:    serverUpdatingStateRefresh(apiClient, systemID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err := createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Instance not created. %s", err)
	}

	return resourceLevel27SystemRead(d, meta)
}

func resourceLevel27SystemDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_system.go: Delete resource")

	apiClient := meta.(*Client)
	_ = apiClient.System("DELETE", d.Id(), nil)

	return nil
}

func serverCreatingStateRefresh(client *Client, systemID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		system := client.System("GET", systemID, nil)
		return system, system.System.RunningStatus, nil
	}
}

func serverUpdatingStateRefresh(client *Client, systemID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		system := client.System("GET", systemID, nil)
		return system, system.System.Status, nil
	}
}

func flattenNetworks(system *System) map[string][]string {
	networks := make(map[string][]string)

	for _, network := range system.System.Networks {
		networks["ids"] = append(networks["ids"], strconv.Itoa(network.ID))
		if network.NetPublic {
			for _, ip := range network.Ips {
				if ip.PublicIpv4 != "" {
					networks["public_ipv4"] = append(networks["public_ipv4"], ip.PublicIpv4)
				}
				if ip.PublicIpv6 != "" {
					networks["public_ipv6"] = append(networks["public_ipv6"], ip.PublicIpv6)
				}
			}
		}

		if network.NetInternal {
			for _, ip := range network.Ips {
				if ip.Ipv4 != "" {
					networks["internal_ipv4"] = append(networks["internal_ipv4"], ip.Ipv4)
				}
				if ip.Ipv6 != "" {
					networks["internal_ipv6"] = append(networks["internal_ipv6"], ip.Ipv6)
				}
			}
		}

		if network.NetCustomer {
			for _, ip := range network.Ips {
				if ip.Ipv4 != "" {
					networks["customer_ipv4"] = append(networks["customer_ipv4"], ip.Ipv4)
				}
				if ip.Ipv6 != "" {
					networks["customer_ipv6"] = append(networks["customer_ipv6"], ip.Ipv6)
				}
			}
		}
	}
	sort.Strings(networks["ids"])
	return networks
}
*/
