package level27

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

//SystemCookbook gets/sets cookbooks for a system
func (c *Client) SystemCookbook(method string, systemid interface{}, cookbookid interface{}, data interface{}) (Cookbook, diag.Diagnostics) {
	var cookbook Cookbook
	var diags diag.Diagnostics

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemid, cookbookid)
		diags = c.invokeAPI("GET", endpoint, nil, &cookbook)
		if diags != nil {
			return cookbook, diags
		}
	case "CREATE":
		endpoint := fmt.Sprintf("systems/%v/cookbooks", systemid)
		diags = c.invokeAPI("POST", endpoint, data, &cookbook)
		if diags != nil {
			return cookbook, diags
		}
		c.updateSystemCookbooks(systemid)
	case "UPDATE":
		endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemid, cookbookid)
		diags = c.invokeAPI("PUT", endpoint, data, nil)
		if diags != nil {
			return cookbook, diags
		}
		c.updateSystemCookbooks(systemid)
	case "DELETE":
		endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemid, cookbookid)
		diags = c.invokeAPI("DELETE", endpoint, nil, nil)
		if diags != nil {
			return cookbook, diags
		}
		c.updateSystemCookbooks(systemid)
	}

	return cookbook, nil
}

func (c *Client) updateSystemCookbooks(systemID interface{}) {
	endpoint := fmt.Sprintf("systems/%v/actions", systemID)
	data := "{\"type\": \"update_cookbooks\"}"
	c.invokeAPI("POST", endpoint, data, nil)
}
