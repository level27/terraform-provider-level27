package main

import (
	"context"
	"fmt"
	"os"

	"bitbucket.org/level27/terraform-provider-level27/level27"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), level27.Provider, providerserver.ServeOpts{
		Address: "registry.terraform.io/level27/level27",
	})

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
