# level27_ssl_certificate

Manages an **SSL Certificate** for a Level27 App. Supports Let's Encrypt (`letsencrypt`), Xolphin (`xolphin`), and custom (`own`) certificates.

## Example Usage

### Let's Encrypt

```hcl
resource "level27_ssl_certificate" "letsencrypt" {
  app_id                    = level27_app.my_project.id
  name                      = "letsencrypt-www"
  ssl_type                  = "letsencrypt"
  auto_ssl_certificate_urls = "www.example.com"
  auto_url_link             = true
  ssl_force                 = true

  # The URL must exist before the certificate is created.
  depends_on = [level27_app_component_url.www]
}
```

### Custom / own certificate

```hcl
resource "level27_ssl_certificate" "own" {
  app_id        = level27_app.my_project.id
  name          = "my-custom-cert"
  ssl_type      = "own"
  ssl_crt       = file("cert.pem")
  ssl_key       = file("key.pem")
  ssl_cabundle  = file("chain.pem")
  auto_url_link = false
  ssl_force     = true
}
```

## Notes

- **`auto_ssl_certificate_urls`** — Comma-separated list of hostnames. Each hostname must already exist as a `level27_app_component_url` on the app. Use `depends_on` to ensure ordering.
- **`ssl_type`** is immutable — changing it forces a replacement.

## Schema

### Required

- `app_id` (Number) — ID of the parent app. Forces replacement when changed.
- `name` (String) — Name of the certificate.

### Optional

- `ssl_type` (String) — Certificate type: `letsencrypt`, `xolphin`, or `own`. Forces replacement when changed.
- `auto_ssl_certificate_urls` (String) — Comma-separated hostnames to link the certificate to. Each URL must already exist on the app.
- `auto_url_link` (Boolean) — Automatically link the certificate to matching URLs. Default: `false`.
- `ssl_force` (Boolean) — Force HTTPS on linked URLs. Default: `false`.
- `ssl_crt` (String, Sensitive) — PEM-encoded certificate (for `ssl_type = "own"`).
- `ssl_key` (String, Sensitive) — PEM-encoded private key (for `ssl_type = "own"`).
- `ssl_cabundle` (String, Sensitive) — PEM-encoded CA bundle (for `ssl_type = "own"`).

### Read-Only (Computed)

- `id` (Number) — Unique identifier of the certificate.
- `ssl_status` (String) — Status of the certificate itself, e.g. `ok`, `pending`.
- `status` (String) — Platform status of the SSL certificate resource.

## Import

Import using `<app_id>/<certificate_id>`:

```sh
terraform import level27_ssl_certificate.letsencrypt 20129/789
```
