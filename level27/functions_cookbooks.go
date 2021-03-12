package level27

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func cookbookParseID(id string) (string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Level27 API error",
			Detail:   fmt.Sprintf("unexpected format of ID (%s), expected attribute1:attribute2", id),
		})
	}

	return parts[0], parts[1], diags
}

func cookbookStateRefresh(client *Client, systemID interface{}, cookbookID interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cookbook, _ := client.SystemCookbook("GET", systemID, cookbookID, nil)
		return cookbook.Cookbook, cookbook.Cookbook.Status, nil
	}
}

func cookbookImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	systemID, cookbookID, err := cookbookParseID(d.Id())

	if err != nil {
		return nil, fmt.Errorf("Error parsing cookbook ID")
	}

	d.Set("system_id", systemID)
	d.Set("cookbook_id", cookbookID)
	d.SetId(fmt.Sprintf("%s:%s", systemID, cookbookID))

	return []*schema.ResourceData{d}, nil
}

func cookbookDelete(id string, meta interface{}) diag.Diagnostics {
	systemID, cookbookID, err := cookbookParseID(id)
	if err != nil {
		return err
	}

	apiClient := meta.(*Client)
	var diags diag.Diagnostics
	_, diags = apiClient.SystemCookbook("DELETE", systemID, cookbookID, nil)
	if diags != nil {
		return diags
	}

	return nil
}

func cookbookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}, systemID string, request CookbookRequest) diag.Diagnostics {
	apiClient := meta.(*Client)
	var diags diag.Diagnostics
	cookbook, diags := apiClient.SystemCookbook("CREATE", systemID, nil, request.String())
	if diags != nil {
		return diags
	}

	createStateConf := cookbookStateConf(d, apiClient, systemID, cookbook.Cookbook.ID)
	_, createError := createStateConf.WaitForStateContext(ctx)
	if createError != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Level27 %s cookbook error", request.Type),
			Detail:   fmt.Sprintf("Error returned from Level27 API\n%v", createError.Error()),
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%s:%d", systemID, cookbook.Cookbook.ID))
	return nil
}

func cookbookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, request CookbookRequest) diag.Diagnostics {
	apiClient := meta.(*Client)
	var diags diag.Diagnostics
	systemID, cookbookID, diags := cookbookParseID(d.Id())
	if diags != nil {
		return diags
	}

	_, diags = apiClient.SystemCookbook("UPDATE", systemID, cookbookID, request.String())
	if diags != nil {
		return diags
	}
	updateStateConf := cookbookStateConf(d, apiClient, systemID, cookbookID)
	_, updateError := updateStateConf.WaitForStateContext(ctx)
	if updateError != nil {
		return diag.FromErr(updateError)
	}
	return nil
}

func cookbookCreateParameters(d *schema.ResourceData, cookbookType string) CookbookRequest {
	cookbookParameters := make(map[string]interface{})
	switch cookbookType {
	case "docker":
		// no cookbook parameters here :-)
	case "haproxy":
		cookbookParameters["haip:ipv4"] = d.Get("haip_ipv4").(string)
		cookbookParameters["haip:ipv6"] = d.Get("haip_ipv6").(string)
		cookbookParameters["haip:routerid"] = d.Get("haip_routerid").(string)
		cookbookParameters["backend:port"] = d.Get("backend_port").(string)
		cookbookParameters["backend:ip"] = d.Get("backend_ip").(string)
		cookbookParameters["port"] = d.Get("frontend_port").(string)
		cookbookParameters["primary"] = d.Get("primary").(bool)
		cookbookParameters["varnish"] = d.Get("varnish").(bool)
		cookbookParameters["expected_traffic"] = d.Get("expected_traffic").(string)
	case "memcached":
		cookbookParameters["port"] = d.Get("port").(string)
		cookbookParameters["listen_address"] = d.Get("listen_address").(string)
		cookbookParameters["memory"] = d.Get("memory").(int)
	case "mongodb":
		cookbookParameters["port"] = d.Get("port").(string)
		cookbookParameters["listen"] = d.Get("listen").(string)
	case "mysql":
		cookbookParameters["query_cache_size"] = d.Get("query_cache_size").(int)
		cookbookParameters["query_cache_limit"] = d.Get("query_cache_limit").(int)
		cookbookParameters["table_cache"] = d.Get("table_cache").(int)
		cookbookParameters["port"] = d.Get("port").(string)
		cookbookParameters["bind_address"] = d.Get("bind_address").(string)
		cookbookParameters["version"] = d.Get("version").([]interface{})
	case "php":
		cookbookParameters["versions"] = d.Get("versions").([]interface{})
		cookbookParameters["apache:port"] = d.Get("apache_port").(string)
		cookbookParameters["apache:bind"] = d.Get("apache_ip").(string)
		cookbookParameters["apache:ssl"] = d.Get("apache_ssl").(bool)
		cookbookParameters["memory_limit"] = d.Get("memory_limit").(int)
		cookbookParameters["max_execution_time"] = d.Get("max_execution_time").(int)
		cookbookParameters["upload_max_filesize"] = d.Get("upload_max_filesize").(int)
		cookbookParameters["disable_functions"] = d.Get("disable_functions").(string)
		cookbookParameters["expected_traffic"] = d.Get("expected_traffic").(string)
		cookbookParameters["max_session_lifetime"] = d.Get("max_session_lifetime").(int)
		cookbookParameters["phpmyadmin:enable"] = d.Get("phpmyadmin_enable").(bool)
		cookbookParameters["process_manager"] = d.Get("process_manager").(string)
		cookbookParameters["display_errors"] = d.Get("display_errors").(string)
	case "postfix":
		cookbookParameters["listening_address"] = d.Get("listening_address").(string)
		cookbookParameters["contentfilter"] = d.Get("contentfilter").(bool)
	case "url":
		cookbookParameters["proxy_service"] = d.Get("proxy_service").(string)
		cookbookParameters["expected_traffic"] = d.Get("expected_traffic").(string)
		cookbookParameters["backend:port"] = d.Get("backend_port").(string)
		cookbookParameters["backend:ip"] = d.Get("backend_ip").(string)
		cookbookParameters["default_vhost"] = d.Get("default_vhost").(bool)
		cookbookParameters["max_upload_filesize"] = d.Get("max_upload_filesize").(int)
	case "web_varnish":
		cookbookParameters["backend_ip"] = d.Get("backend_ip").(string)
		cookbookParameters["backend_port"] = d.Get("backend_port").(string)
		cookbookParameters["frontend_ip"] = d.Get("frontend_ip").(string)
		cookbookParameters["frontend_port"] = d.Get("frontend_port").(string)
		cookbookParameters["storage_size"] = d.Get("storage_size").(string)
		cookbookParameters["admin_ip"] = d.Get("admin_ip").(string)
		cookbookParameters["admin_port"] = d.Get("admin_port").(string)
	}

	r := CookbookRequest{
		Type:               cookbookType,
		Cookbookparameters: cookbookParameters,
	}
	return r
}

func cookbookStateConf(d *schema.ResourceData, apiClient *Client, systemID interface{}, cookbookID interface{}) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"to_create", "creating", "updating"},
		Target:     []string{"ok"},
		Refresh:    cookbookStateRefresh(apiClient, systemID, cookbookID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
