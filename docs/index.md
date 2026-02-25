# Level27 Provider

The **Level27** Terraform provider lets you manage apps, components, URLs, and SSL certificates on the [Level27](https://level27.eu) hosting platform.

## Authentication

You can get your API key at (https://app.level27.eu/account/profile/security)

Set your API key via the `LEVEL27_API_KEY` environment variable (recommended):

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

## Schema

### Optional

- `api_key` (String, Sensitive) — Level27 API key. Can also be set via the `LEVEL27_API_KEY` environment variable.
- `api_url` (String) — Base URL of the Level27 API. Defaults to `https://api.level27.eu/v1`.
