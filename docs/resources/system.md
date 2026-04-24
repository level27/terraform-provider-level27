---
page_title: "level27_system Resource - terraform-provider-level27"
subcategory: ""
description: |-
  Manages a Level27 System (virtual server).
---

# level27_system (Resource)

Manages a Level27 **System** (virtual server).

The `type`, `systemimage_id`, `systemprovider_configuration_id`, `zone_id`, and `parentsystem_id` attributes are **immutable** — changing them forces a replacement.

Use the `l27lookup` helper (bundled with this provider) to discover IDs for your environment:

```sh
./l27lookup images 1              # systemimage_id  (1 = Level27 provider)
./l27lookup configs               # systemprovider_configuration_id
./l27lookup zones                 # zone_id
./l27lookup networks public       # network names (public)
./l27lookup networks customer     # network names (customer / private)
./l27lookup mgmt <org_id>         # management types with pricing
```

## Example Usage

### Minimal

```hcl
resource "level27_system" "example" {
  name                            = "myserver.example.com"
  type                            = "kvmguest"
  systemimage_id                  = 89   # ubuntu_2404lts_server
  systemprovider_configuration_id = 17   # Level27 Flexible
  zone_id                         = 1    # Hasselt 1

  cpu    = 2
  memory = 4
  disk   = 50
}
```

### With networks and IP assignment

```hcl
resource "level27_system" "example" {
  name                            = "myserver.example.com"
  type                            = "kvmguest"
  systemimage_id                  = 89
  systemprovider_configuration_id = 17
  zone_id                         = 1

  cpu             = 2
  memory          = 4
  disk            = 50
  management_type = "infra_plus"

  networks = {
    "level27_public_13" = "auto"       # auto-assign a free IPv4
    "level27_cust_32"   = "10.0.1.5"  # specific IPv4
    "level27_internal"  = ""          # attach interface, no IP
  }
}

output "system_id" {
  value = level27_system.example.id
}
```

## Schema

### Required

- `name` (String) — Hostname (FQDN) of the system, e.g. `myserver.example.com`.
- `type` (String) — System type. Common values: `kvmguest`, `vmware`, `baremetal`. **Forces replacement.**
- `systemimage_id` (Number) — ID of the OS image to use for provisioning. **Forces replacement.**
- `systemprovider_configuration_id` (Number) — ID of the hardware/hypervisor profile. **Forces replacement.**
- `zone_id` (Number) — ID of the datacenter zone to deploy the system in. **Forces replacement.**
- `cpu` (Number) — Number of virtual CPUs.
- `disk` (Number) — Disk size in GB.
- `memory` (Number) — Memory in GB.

### Optional

- `customer_fqdn` (String) — Customer-facing FQDN. Defaults to `name` when omitted.
- `external_info` (String) — External reference or billing annotation.
- `management_type` (String) — Server management level. Common values:
  - `basic`
  - `infra_plus`
  - `enterprise`
  
  Use `./l27lookup mgmt <org_id>` to list all available types with pricing.

- `networks` (Map of String) — Network interfaces to attach to the system.  
  Each **key** is a network name (resolved to an ID by the provider; use `./l27lookup networks` to list).  
  Each **value** controls IP assignment on that interface:

  | Value | Behaviour |
  |---|---|
  | `"auto"` | Automatically assign a free IPv4 via `GET /networks/{id}/locate`. The assigned address is kept stable — subsequent plans never replace it. |
  | `"10.x.x.x"` | Assign the specified IPv4 address. Change the value to replace the IP. |
  | `""` (empty string) | Attach the interface without assigning an IP. Any existing IP is removed. |

  Networks not in the map are detached from the system.

- `auto_install` (Boolean) — When `true` (default), triggers OS installation (`autoInstall`) immediately after the system is allocated and networks are configured. The provider calls `POST /systems/{id}/actions` with `autoInstall` and waits until installation completes. Only fires when the system status is `stopped`. Set to `false` to skip automatic installation. Default: `true`.
- `parentsystem_id` (Number) — ID of the parent host system for nested virtualisation. Usually auto-assigned by the API. Setting this explicitly forces replacement when changed.
- `period` (Number) — Billing period in months. Default: `1`.

### Read-Only

- `id` (Number) — Unique identifier of the system.
- `organisation_id` (Number) — ID of the organisation that owns this system (resolved automatically via `/whoami`).
- `status` (String) — Current provisioning status (e.g. `ok`, `creating`, `deleted`).
- `status_category` (String) — Status category: `green`, `yellow`, `red`, or `grey`.

## Import

An existing system can be imported using its numeric ID:

```sh
terraform import level27_system.example 1094
```

After import, run `terraform plan` to see which optional attributes differ from their defaults and add them to your configuration as needed.

