package level27

import (
	"fmt"
)

// Networks gets all networks from the API
func (c *Client) Networks(method string, systemProvider string) Networks {
	var networks Networks

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("networks?limit=10000&excludefull=1&status=ok&provider=%s", systemProvider)
		c.invokeAPI("GET", endpoint, nil, &networks)
	}
	return networks
}
