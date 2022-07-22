package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"bitbucket.org/level27/terraform-provider-level27/level27"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "run with debugger support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/level27/level27",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), level27.Provider, opts)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
