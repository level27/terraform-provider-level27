# level27_app

Manages a Level27 **App** (project). An app groups components such as PHP web-apps, databases, mail, etc.

## Example Usage

```hcl
resource "level27_app" "my_project" {
  name            = "my-terraform-project"
  organisation_id = 414

  # Billing product — both fields are optional.
  # custom_package_id   = 456
  # custom_package_name = "drupal_run"
}
```

## Schema

### Required

- `name` (String) — Name of the app/project.
- `organisation_id` (Number) — ID of the organisation that owns this app. Forces replacement when changed.

### Optional

- `custom_package_id` (Number) — ID of the billing product (custom package).
- `custom_package_name` (String) — Name of the billing product, e.g. `drupal_run`, `wordpress_walk`.
- `auto_teams` (String) — Comma-separated list of team IDs to auto-assign to this app.
- `auto_upgrades` (String) — Comma-separated list of upgrade names to apply automatically.
- `external_info` (String) — External reference (required when billableitemInfo entities exist for the organisation).

### Read-Only (Computed)

- `id` (Number) — Unique identifier of the app.
- `status` (String) — Current status, e.g. `ok`, `to_create`, `creating`.
- `status_category` (String) — Status category: `green`, `yellow`, `red`, or `grey`.
- `hosting_type` (String) — Hosting type: `Agency` or `Classic`.
- `billing_status` (String) — Billing status of the app.

## Import

Import an existing app by its numeric ID:

```sh
terraform import level27_app.my_project 20129
```
