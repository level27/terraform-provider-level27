// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppComponentURLResource tests the full CRUD lifecycle of a
// level27_app_component_url resource.
//
// Required environment variables:
//   - LEVEL27_API_KEY
//   - LEVEL27_TEST_SYSTEM_ID – System (server) ID on which to deploy the component.
func TestAccAppComponentURLResource(t *testing.T) {
	systemID := testAccEnvOrFatal(t, "LEVEL27_TEST_SYSTEM_ID")
	resourceName := "level27_app_component_url.test"

	// Generate a unique hostname to avoid conflicts between runs.
	rnd := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	hostname := fmt.Sprintf("tf-acc-%s.example.com", rnd)
	hostnameUpdated := fmt.Sprintf("tf-acc-%s-upd.example.com", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create and verify.
			{
				Config: testAccAppComponentURLConfig(systemID, hostname),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", hostname),
					resource.TestCheckResourceAttr(resourceName, "ssl_force", "false"),
					resource.TestCheckResourceAttr(resourceName, "handle_dns", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			// Step 2: Import by "<app_id>/<component_id>/<url_id>" and verify.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Update URL content.
			{
				Config: testAccAppComponentURLConfig(systemID, hostnameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", hostnameUpdated),
				),
			},
		},
	})
}

func testAccAppComponentURLConfig(systemID, hostname string) string {
	return fmt.Sprintf(`
resource "level27_app" "test" {
	name = "tf-acc-app-url"
}

resource "level27_app_component" "test" {
  app_id           = level27_app.test.id
  name             = "tf-acc-comp-url"
  appcomponenttype = "php"
	system           = %[1]s
}

resource "level27_app_component_url" "test" {
  app_id       = level27_app.test.id
  component_id = level27_app_component.test.id
	content      = %[2]q
  ssl_force    = false
  handle_dns   = false
}
`, systemID, hostname)
}
