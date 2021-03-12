package level27

import (
	"fmt"
)

//Domain gets a system from the API
func (c *Client) Domain(method string, id interface{}, data interface{}) System {
	var system System

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("systems/%s", id)
		c.invokeAPI("GET", endpoint, nil, &system)
	case "CREATE":
		endpoint := "systems"
		c.invokeAPI("POST", endpoint, data, &system)
	case "UPDATE":
		endpoint := fmt.Sprintf("systems/%s", id)
		c.invokeAPI("PUT", endpoint, data, &system)
	case "DELETE":
		endpoint := fmt.Sprintf("systems/%s", id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return system
}
