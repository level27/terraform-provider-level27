// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSystemResource tests the full CRUD lifecycle of a level27_system resource.
//
// Required environment variables:
//   - LEVEL27_API_KEY
//   - LEVEL27_TEST_SYSTEMIMAGE_ID      – System image ID to use for provisioning.
//   - LEVEL27_TEST_SYSTEMPROVIDER_ID   – System provider configuration ID.
//   - LEVEL27_TEST_ZONE_ID             – Zone ID to deploy into.
func TestAccSystemResource(t *testing.T) {
	imageID := testAccEnvOrFatal(t, "LEVEL27_TEST_SYSTEMIMAGE_ID")
	providerCfgID := testAccEnvOrFatal(t, "LEVEL27_TEST_SYSTEMPROVIDER_ID")
	zoneID := testAccEnvOrFatal(t, "LEVEL27_TEST_ZONE_ID")
	resourceName := "level27_system.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the system and verify initial state.
			{
				Config: testAccSystemConfig(imageID, providerCfgID, zoneID, "tf-acc-system.example.com", 1, 1, 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-system.example.com"),
					resource.TestCheckResourceAttr(resourceName, "type", "kvmguest"),
					resource.TestCheckResourceAttrSet(resourceName, "organisation_id"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "1"),
					resource.TestCheckResourceAttr(resourceName, "memory", "1"),
					resource.TestCheckResourceAttr(resourceName, "disk", "20"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "status_category"),
				),
			},
			// Step 2: Import by ID and verify state consistency.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// period is not returned by the API after creation.
				ImportStateVerifyIgnore: []string{"period"},
			},
			// Step 3: Update CPU and memory (scale up).
			{
				Config: testAccSystemConfig(imageID, providerCfgID, zoneID, "tf-acc-system.example.com", 2, 2, 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory", "2"),
				),
			},
		},
	})
}

func testAccSystemConfig(imageID, providerCfgID, zoneID, name string, cpu, memory, disk int) string {
	return fmt.Sprintf(`
resource "level27_system" "test" {
	name                            = %[4]q
  type                            = "kvmguest"
	systemimage_id                  = %[1]s
	systemprovider_configuration_id = %[2]s
	zone_id                         = %[3]s
	cpu                             = %[5]d
	memory                          = %[6]d
	disk                            = %[7]d
  management_type                 = "infra_plus"
}
`, imageID, providerCfgID, zoneID, name, cpu, memory, disk)
}
