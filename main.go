package main

import (
	"bitbucket.org/level27/terraform-provider-level27/level27"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: level27.Provider,
	})
}
