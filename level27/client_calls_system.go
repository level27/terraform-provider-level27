package level27

import (
	"fmt"
)

//System gets a system from the API
func (c *Client) System(method string, id interface{}, data interface{}) System {
	var system System

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("systems/%v", id)
		c.invokeAPI("GET", endpoint, nil, &system)
	case "CREATE":
		endpoint := "systems"
		c.invokeAPI("POST", endpoint, data, &system)
	case "UPDATE":
		endpoint := fmt.Sprintf("systems/%v", id)
		c.invokeAPI("PUT", endpoint, data, &system)
	case "DELETE":
		endpoint := fmt.Sprintf("systems/%v", id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return system
}

//SystemProviders gets all system providers and their system images
func (c *Client) SystemProviders() SystemProviders {
	var providers SystemProviders

	endpoint := fmt.Sprintf("systems/providers")
	c.invokeAPI("GET", endpoint, nil, &providers)

	return providers
}

//SystemProviderConfigurations gets all system provider configurations
func (c *Client) SystemProviderConfigurations() ProviderConfigurations {
	var configurations ProviderConfigurations

	endpoint := fmt.Sprintf("systems/provider/configurations")
	c.invokeAPI("GET", endpoint, nil, &configurations)

	return configurations
}

//SystemRegions gets all system regions
func (c *Client) SystemRegions() Regions {
	var regions Regions

	endpoint := fmt.Sprintf("systems/regions")
	c.invokeAPI("GET", endpoint, nil, &regions)

	return regions
}
