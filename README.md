# Terraform Provider for Level27

A [Terraform](https://www.terraform.io) provider for managing resources on the [Level27](https://level27.be) hosting platform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build from source)

## Supported Resources

| Resource | Description |
|----------|-------------|
| `level27_app` | App (project) that groups components |
| `level27_app_component` | Component within an app (PHP, MySQL, SFTP, ...) |
| `level27_app_component_url` | URL attached to a component |
| `level27_ssl_certificate` | SSL certificate (Let's Encrypt, ...) |

## Authentication

Set your API key via the environment variable (recommended):

```sh
export LEVEL27_API_KEY="your-api-key"
```

Or configure it in the provider block (avoid for production):

```hcl
provider "level27" {
  api_key = "your-api-key"
}
```

## Usage

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

resource "level27_app" "example" {
  name = "my-project"
}

resource "level27_app_component" "php" {
  app_id           = level27_app.example.id
  name             = "web"
  appcomponenttype = "php"
  system           = 943
  path             = "/public_html"
  # version is auto-detected as the highest available version on the system
}

resource "level27_app_component" "db" {
  app_id           = level27_app.example.id
  name             = "db"
  appcomponenttype = "mysql"
  system           = 943
  pass             = "YourSecurePassword123!"
  # version is auto-detected from the system's installed cookbooks
}

resource "level27_app_component_url" "www" {
  app_id       = level27_app.example.id
  component_id = level27_app_component.php.id
  content      = "www.example.com"
  ssl_force    = true
  handle_dns   = false # set to true only if Level27 manages DNS for this domain
}

resource "level27_ssl_certificate" "letsencrypt" {
  app_id                    = level27_app.example.id
  name                      = "letsencrypt-www"
  ssl_type                  = "letsencrypt"
  auto_ssl_certificate_urls = "www.example.com"
  auto_url_link             = true
  ssl_force                 = true

  depends_on = [level27_app_component_url.www]
}
```

See the [`examples/`](examples/) directory for more complete examples.

## Development

### Building

```sh
go build -o terraform-provider-level27 .
```

### Local installation

Install the provider into your local Terraform plugin cache for testing:

```sh
VERSION=0.1.0 make install
```

Then add a `dev_overrides` block to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/level27/level27" = "/home/<you>/terraform-provider-level27"
  }
  direct {}
}
```

With this in place, `terraform plan` / `terraform apply` in any configuration using this provider will use your local build directly — no `terraform init` needed.

### Testing

Run unit tests:

```sh
make test
```

Run acceptance tests (requires a valid API key and creates real resources):

```sh
export LEVEL27_API_KEY="your-api-key"
make testacc
```

### Other targets

| Command | Description |
|---------|-------------|
| `make build` | Compile the provider binary |
| `make install VERSION=x.y.z` | Install locally for development |
| `make fmt` | Format Go source files |
| `make lint` | Run golangci-lint |
| `make tidy` | Run `go mod tidy` |
| `make docs` | Regenerate documentation from schema |

## Provider behaviour notes

### Version auto-detection

When `version` is omitted on a `level27_app_component`, the provider queries `GET /v1/systems/{id}` and selects the **highest available version** for the given component type from the system's installed cookbooks. Specify `version` explicitly to pin to a particular release.

### Status polling

After creating or updating an app component, the provider polls the API every 3 seconds until the component leaves a transitional state (`creating`, `updating`, etc.). This ensures the final `ok` status is stored in state rather than an intermediate `updating` status.

### handle_dns on URLs

Set `handle_dns = true` on a `level27_app_component_url` **only** when Level27 manages the DNS zone for that domain. Using it for externally managed domains will result in a 400 API error.

### Organisation resolution

The provider resolves your organisation automatically via `GET /whoami`. You no longer need to configure `organisation_id` in resources.

## License

[MPL-2.0](LICENSE)
