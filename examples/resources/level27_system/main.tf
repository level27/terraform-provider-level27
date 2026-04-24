terraform {
  required_providers {
    level27 = {
      source  = "registry.terraform.io/level27/level27"
      version = "~> 0.1"
    }
  }
}

# Set LEVEL27_API_KEY env var or configure api_key in provider block.
provider "level27" {
  #api_url = "https://api.cp4dev.be/v1"  # uncomment for dev environment
  #api_key = "..." — or set LEVEL27_API_KEY env var
}

# -------------------------------------------------------------------
# Variables
# Run `l27lookup` to find the right IDs for your environment:
#
#   ./l27lookup images 1             → systemimage_id
#   ./l27lookup configs              → systemprovider_configuration_id
#   ./l27lookup zones                → zone_id
#   ./l27lookup networks public      → public network names
#   ./l27lookup networks customer    → customer network names
# -------------------------------------------------------------------

variable "systemimage_id" {
  description = "OS image ID (see: l27lookup images 1)."
  type        = number
}

variable "systemprovider_configuration_id" {
  description = "Provider hardware configuration ID (see: l27lookup configs)."
  type        = number
}

variable "zone_id" {
  description = "Datacenter zone ID (see: l27lookup zones)."
  type        = number
}

variable "name" {
  description = "Fully-qualified hostname of the server (e.g. myserver.example.com)."
  type        = string
}

variable "cpu" {
  description = "Number of virtual CPUs."
  type        = number
  default     = 1
}

variable "memory" {
  description = "Memory in GB."
  type        = number
  default     = 1
}

variable "disk" {
  description = "Disk size in GB."
  type        = number
  default     = 20
}

variable "management_type" {
  description = "Server management type."
  type        = string
  default     = "pro"
}

variable "networks" {
  description = "Map of network-name \u2192 IPv4 to assign. Use \"auto\" for automatic, \"\" for no IP."
  type        = map(string)
  default     = {}
}
variable "auto_install" {
  description = "Trigger OS installation after network attachment. Set to false to skip."
  type        = bool
  default     = true
}
# -------------------------------------------------------------------
# System (virtual server)
# -------------------------------------------------------------------
resource "level27_system" "this" {
  name                            = var.name
  type                            = "kvmguest"
  systemimage_id                  = var.systemimage_id
  systemprovider_configuration_id = var.systemprovider_configuration_id
  zone_id                         = var.zone_id

  cpu             = var.cpu
  memory          = var.memory
  disk            = var.disk
  management_type = var.management_type

  networks     = var.networks
  auto_install = var.auto_install
}

output "system_id" {
  description = "ID of the created system — use as system in level27_app_component."
  value       = level27_system.this.id
}

output "system_status" {
  value = level27_system.this.status
}
