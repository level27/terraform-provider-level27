// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppComponentResource tests the full CRUD lifecycle of a
// level27_app_component resource (PHP component).
//
// Required environment variables:
//   - LEVEL27_API_KEY
//   - LEVEL27_TEST_ORG_ID    – Organisation to create the parent app in.
//   - LEVEL27_TEST_SYSTEM_ID – System (server) ID on which to deploy the component.
func TestAccAppComponentResource(t *testing.T) {
	orgID := testAccEnvOrFatal(t, "LEVEL27_TEST_ORG_ID")
	systemID := testAccEnvOrFatal(t, "LEVEL27_TEST_SYSTEM_ID")
	resourceName := "level27_app_component.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create and verify.
			{
				Config: testAccAppComponentConfig(orgID, systemID, "tf-acc-comp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-comp"),
					resource.TestCheckResourceAttr(resourceName, "appcomponenttype", "php"),
					resource.TestCheckResourceAttr(resourceName, "system", systemID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "category"),
				),
			},
			// Step 2: Import by "<app_id>/<component_id>" and verify.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// 'pass' is write-only; the API never returns it.
				ImportStateVerifyIgnore: []string{"pass"},
			},
			// Step 3: Update the component name.
			{
				Config: testAccAppComponentConfig(orgID, systemID, "tf-acc-comp-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-comp-updated"),
				),
			},
		},
	})
}

func testAccAppComponentConfig(orgID, systemID, componentName string) string {
	return fmt.Sprintf(`
resource "level27_app" "test" {
  name            = "tf-acc-app-comp"
  organisation_id = %[1]s
}

resource "level27_app_component" "test" {
  app_id          = level27_app.test.id
  name            = %[3]q
  appcomponenttype = "php"
  system          = %[2]s
}
`, orgID, systemID, componentName)
}
