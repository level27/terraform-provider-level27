package level27

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLevel27AppComponentMysql() *schema.Resource {

	return &schema.Resource{
		Create: resourceLevel27AppComponentMysqlCreate,
		Read:   resourceLevel27AppComponentMysqlRead,
		Delete: resourceLevel27AppComponentMysqlDelete,
		Update: resourceLevel27AppComponentMysqlUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"system_group_id"},
			},
			"system_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"system_id"},
			},
			"host": {
				Optional: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceLevel27AppComponentMysqlRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_appcomponent_Mysql.go: Read resource")

	var ac AppComponent
	apiClient := meta.(*Client)

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return err
	}
	ac = apiClient.AppComponent("GET", appID, appComponentID, nil)

	log.Printf("resource_level27_appcomponent_Mysql.go: Read routine called. Object built: %d", ac.Component.ID)

	d.SetId(fmt.Sprintf("%d:%d", ac.Component.App.ID, ac.Component.ID))
	d.Set("name", ac.Component.Name)
	d.Set("username", ac.Component.Appcomponentparameters.User)
	systemID := strconv.Itoa(ac.Component.Systems[0].ID)
	if systemID == "0" {
		systemID = ""
	}
	d.Set("system_id", systemID)

	systemGroupID := strconv.Itoa(ac.Component.Systemgroup)
	if systemGroupID == "0" {
		systemGroupID = ""
	}
	d.Set("system_group_id", systemGroupID)

	d.Set("host", ac.Component.Appcomponentparameters.Host)

	return nil
}

func resourceLevel27AppComponentMysqlCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_appcomponent_Mysql.go: Create resource")

	log.Printf("resource_level27_appcomponent_Mysql.go: Create routine called. Object to be built: %v", d)

	//var ac AppComponent
	request := AppComponentRequest{
		Name:        d.Get("name").(string),
		Type:        "mysql",
		Category:    "config",
		System:      d.Get("system_id").(string),
		SystemGroup: d.Get("system_group_id").(string),
		Host:        d.Get("host").(string),
	}

	apiClient := meta.(*Client)

	ac := apiClient.AppComponent("CREATE", d.Get("app_id"), nil, request.String())

	d.SetId(fmt.Sprintf("%d:%d", ac.Component.App.ID, ac.Component.ID))
	return resourceLevel27AppComponentMysqlRead(d, meta)
}

func resourceLevel27AppComponentMysqlUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_appcomponent_Mysql.go: Update resource")
	log.Printf("resource_level27_appcomponent_Mysql.go: Update routine called. Object to be built: %v", d)

	appID, appComponentID, err := appComponentParseID(d.Id())
	if err != nil {
		return err
	}

	//var ac AppComponent
	request := AppComponentRequest{
		Name:        d.Get("name").(string),
		Type:        "mysql",
		Category:    "config",
		System:      d.Get("system_id").(string),
		SystemGroup: d.Get("system_group_id").(string),
		Host:        d.Get("host").(string),
	}

	apiClient := meta.(*Client)

	apiClient.AppComponent("UPDATE", appID, appComponentID, request.String())

	return resourceLevel27AppComponentMysqlRead(d, meta)
}

func resourceLevel27AppComponentMysqlDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("resource_level27_appcomponent_Mysql.go: Delete resource")
	return nil
}
