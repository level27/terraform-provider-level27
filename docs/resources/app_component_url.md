# level27_app_component_url

Manages a **URL** linked to a Level27 App Component (PHP, ASP, etc.).

## Example Usage

```hcl
resource "level27_app_component_url" "www" {
  app_id       = level27_app.my_project.id
  component_id = level27_app_component.php.id
  content      = "www.example.com"
  ssl_force    = true
  # Set handle_dns = true only if Level27 manages DNS for this domain.
  handle_dns   = false
}
```

## Notes

- **`handle_dns`** — Only set to `true` if Level27 manages the DNS zone for the domain. Setting it to `true` for externally managed domains will cause a 400 validation error from the API.
- **`auto_ssl_certificate`** — Requires `handle_dns = true`. For manually provisioned certificates, use `level27_ssl_certificate` with `depends_on` instead.
- `app_id`, `component_id`, and `content` are immutable — changing them forces a replacement.

## Schema

### Required

- `app_id` (Number) — ID of the parent app. Forces replacement when changed.
- `component_id` (Number) — ID of the parent component. Forces replacement when changed.
- `content` (String) — The hostname, e.g. `www.example.com`. Forces replacement when changed.

### Optional

- `ssl_force` (Boolean) — Force HTTPS redirect. Default: `true`.
- `handle_dns` (Boolean) — Automatically create DNS records. Only valid when Level27 manages the DNS zone. Default: `false`.
- `auto_ssl_certificate` (Boolean) — Automatically provision a Let's Encrypt certificate. Requires `handle_dns = true`. Default: `false`.
- `caching` (Boolean) — Enable caching. Default: `false`.
- `authentication` (Boolean) — Enable HTTP basic authentication. Default: `false`.
- `ssl_certificate_id` (Number) — ID of an existing `level27_ssl_certificate` to attach.

### Read-Only (Computed)

- `id` (Number) — Unique identifier of the URL.
- `status` (String) — Current status.
- `https` (Boolean) — Whether HTTPS is currently active.

## Import

Import using `<app_id>/<component_id>/<url_id>`:

```sh
terraform import level27_app_component_url.www 20129/68695/12345
```
