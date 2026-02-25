// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/level27/terraform-provider-level27/internal/provider"
)

// testAccProtoV6ProviderFactories is used in every acceptance test to wire up
// the provider under test against a real Level27 API endpoint.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"level27": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// testAccPreCheck verifies that the minimum required environment variables are
// present before any acceptance test runs.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("LEVEL27_API_KEY"); v == "" {
		t.Fatal("LEVEL27_API_KEY must be set for acceptance tests")
	}
}

// testAccEnvOrFatal returns the value of the given environment variable or
// marks the test as failed if the variable is empty.
func testAccEnvOrFatal(t *testing.T, env string) string {
	t.Helper()
	v := os.Getenv(env)
	if v == "" {
		t.Fatalf("%s must be set for this acceptance test", env)
	}
	return v
}
