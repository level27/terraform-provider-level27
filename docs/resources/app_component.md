# level27_app_component

Manages a Level27 **App Component**. A component is one of the buildable blocks of an app: PHP, MySQL, SFTP, mail, etc.

> **Note:** `appcomponenttype`, `system`, and `systemgroup` are immutable after creation — changing them forces a replacement.

## Example Usage

### PHP component (version auto-detected)

```hcl
resource "level27_app_component" "php" {
  app_id           = level27_app.my_project.id
  name             = "web"
  appcomponenttype = "php"
  system           = 943
  path             = "/public_html"
  # version is auto-detected as the highest available PHP version on the system.
  # Pin it explicitly if needed:
  # version = "8.2"
}
```

### MySQL component

```hcl
resource "level27_app_component" "db" {
  app_id           = level27_app.my_project.id
  name             = "db"
  appcomponenttype = "mysql"
  system           = 943
  pass             = "YourSecurePassword123!"
  # version is auto-detected from the system's installed cookbooks.
}
```

## Version auto-detection

When `version` is omitted, the provider queries `GET /v1/systems/{system_id}` and picks the **highest available version** for the given `appcomponenttype` from the system's installed cookbooks.

- For types that expose a single version (`version` key), that version is used.
- For types that expose multiple versions (`versions` key, e.g. PHP), the highest semantic version is chosen.

Specify `version` explicitly to pin to a particular release.

## Status polling

After creating or updating a component, the provider polls the API every 3 seconds until the component leaves a transitional state (`to_create`, `creating`, `to_update`, `updating`). This prevents Terraform from storing an intermediate status like `updating` in the state.

## Schema

### Required

- `app_id` (Number) — ID of the parent app. Forces replacement when changed.
- `name` (String) — Name of the component.
- `appcomponenttype` (String) — Component type. One of: `mysql`, `php`, `keyvalue_redis`, `sftp`, `asp`, `mssql`, `solr`, `memcached`, `elasticsearch`, `redis`, `ruby`, `url`, `mongodb`, `nodejs`, `python`, `postgresql`, `dotnet`, `rabbitmq`, `varnish`, `java`, `mail`, `wasp`, `wmssql`, `mailpit`, `meilisearch`. Forces replacement when changed.

### Optional

- `system` (Number) — ID of the system to deploy on. Mutually exclusive with `systemgroup`. Forces replacement when changed.
- `systemgroup` (Number) — ID of the system group to deploy on. Mutually exclusive with `system`. Forces replacement when changed.
- `version` (String) — Version of the component runtime. Auto-detected from system cookbooks when omitted (highest available version is used).
- `path` (String) — Relative web root path (e.g. `/public_html` for PHP).
- `pass` (String, Sensitive) — Password for the component. Required for `mysql` and similar database types.
- `limit_group` (String) — Limit group, e.g. `database` or `application`.
- `sshkeys` (List of Number) — List of SSH key IDs to deploy onto this component.
- `extraconfig` (String) — Extra configuration for PHP/ASP components (custom php.ini directives, etc.).

### Read-Only (Computed)

- `id` (Number) — Unique identifier of the component.
- `status` (String) — Current status, e.g. `ok`, `creating`, `updating`.
- `status_category` (String) — Status category: `green`, `yellow`, `red`, or `grey`.
- `category` (String) — Component category assigned by Level27 (e.g. `Web-apps`, `Databases`).

## Import

Import an existing component using `<app_id>/<component_id>`:

```sh
terraform import level27_app_component.php 20129/68695
```
