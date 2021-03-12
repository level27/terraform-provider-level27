package level27

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// App gets an application from the API
func (c *Client) App(method string, id interface{}, data interface{}) App {
	var app App

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("GET", endpoint, nil, &app)
	case "CREATE":
		endpoint := "apps"
		c.invokeAPI("POST", endpoint, data, &app)
	case "UPDATE":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("PUT", endpoint, data, &app)
	case "DELETE":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return app
}

// AppComponent gets an appcomponent from the API
func (c *Client) AppComponent(method string, appid interface{}, id interface{}, data interface{}) AppComponent {
	var ac AppComponent

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("apps/%s/components/%s", appid, id)
		c.invokeAPI("GET", endpoint, nil, &ac)
	case "CREATE":
		endpoint := fmt.Sprintf("apps/%s/components", appid)
		c.invokeAPI("POST", endpoint, data, &ac)
	case "UPDATE":
		endpoint := fmt.Sprintf("apps/%s/components/%s", appid, id)
		c.invokeAPI("PUT", endpoint, data, &ac)
	case "DELETE":
		endpoint := fmt.Sprintf("apps/%s/components/%s", appid, id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return ac
}

// AppComponentURL gets an appcomponent URL from the API
func (c *Client) AppComponentURL(method string, appid interface{}, componentid interface{}, urlid interface{}, data interface{}) (URL, diag.Diagnostics) {
	var url URL
	var err diag.Diagnostics

	switch method {
	case "GETALL":
		endpoint := fmt.Sprintf("apps/%s/components/%s/urls", appid, componentid)
		err = c.invokeAPI("GET", endpoint, nil, &url)
	case "GET":
		endpoint := fmt.Sprintf("apps/%s/components/%s/urls/%s", appid, componentid, urlid)
		err = c.invokeAPI("GET", endpoint, nil, &url)
	case "CREATE":
		endpoint := fmt.Sprintf("apps/%s/components/%s/urls", appid, componentid)
		err = c.invokeAPI("POST", endpoint, data, &url)
	case "UPDATE":
		endpoint := fmt.Sprintf("apps/%s/components/%s/urls/%s", appid, componentid, urlid)
		err = c.invokeAPI("PUT", endpoint, data, &url)
	case "DELETE":
		endpoint := fmt.Sprintf("apps/%s/components/%s/urls/%s", appid, componentid, urlid)
		err = c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return url, err
}
