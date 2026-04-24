// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppResource tests the full CRUD lifecycle of a level27_app resource.
//
// Required environment variables:
//   - LEVEL27_API_KEY
func TestAccAppResource(t *testing.T) {
	resourceName := "level27_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create and verify initial state.
			{
				Config: testAccAppConfig("tf-acc-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-app"),
					resource.TestCheckResourceAttrSet(resourceName, "organisation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "status_category"),
					resource.TestCheckResourceAttrSet(resourceName, "hosting_type"),
				),
			},
			// Step 2: Import by ID and verify state is consistent.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Update name and verify the change is applied.
			{
				Config: testAccAppConfig("tf-acc-app-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-app-updated"),
				),
			},
		},
	})
}

func testAccAppConfig(name string) string {
	return fmt.Sprintf(`
resource "level27_app" "test" {
  name = %[1]q
}
`, name)
}
