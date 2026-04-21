# Level27 Provider

The **Level27** Terraform provider lets you manage infrastructure and applications on the [Level27](https://level27.eu) hosting platform: virtual servers, network attachments, apps, components, URLs, and SSL certificates.

## Authentication

Get your API key at <https://app.level27.eu/account/profile/security>.

Set it via the `LEVEL27_API_KEY` environment variable (recommended):

```sh
export LEVEL27_API_KEY="your-api-key"
```

Or configure it directly in the provider block (not recommended for production):

```hcl
provider "level27" {
  api_key = "your-api-key"
}
```

## Example Usage

```hcl
terraform {
  required_providers {
    level27 = {
      source  = "registry.terraform.io/level27/level27"
      version = "~> 0.1"
    }
  }
}

provider "level27" {}
```

## l27lookup helper

This provider ships a `l27lookup` CLI tool to discover resource IDs and names for your environment. Build it with:

```sh
make l27lookup
```

Commands:

| Command | Description |
|---|---|
| `./l27lookup orgs` | List organisations (reference only, `organisation_id` is resolved automatically via `/whoami`) |
| `./l27lookup zones` | List datacenter zones (→ `zone_id`) |
| `./l27lookup configs` | List provider configurations (→ `systemprovider_configuration_id`) |
| `./l27lookup images <provider_id>` | List OS images (→ `systemimage_id`) |
| `./l27lookup networks [public\|customer\|internal]` | List networks by type (→ keys for `networks` map) |
| `./l27lookup mgmt <org_id>` | List management types with EUR pricing |
| `./l27lookup apps` | List apps |

## Schema

### Optional

- `api_key` (String, Sensitive) — Level27 API key. Can also be set via the `LEVEL27_API_KEY` environment variable.
- `api_url` (String) — Base URL of the Level27 API. Defaults to `https://api.level27.eu/v1`.
