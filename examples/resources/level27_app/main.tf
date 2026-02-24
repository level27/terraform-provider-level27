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
  # api_key = "..." — or set LEVEL27_API_KEY env var
}

variable "organisation_id" {
  description = "Your Level27 organisation ID."
  type        = number
}

# -------------------------------------------------------------------
# App
# -------------------------------------------------------------------
resource "level27_app" "my_project" {
  name            = "my-terraform-project"
  organisation_id = var.organisation_id

  # custom_package_id   = 456  # optional billing product ID
  # custom_package_name = "drupal_run"  # optional billing product name
}

# -------------------------------------------------------------------
# PHP web-app component
# The 'version' is auto-detected as the highest available PHP version
# on the system when omitted.
# -------------------------------------------------------------------
resource "level27_app_component" "php" {
  app_id           = level27_app.my_project.id
  name             = "web"
  appcomponenttype = "php"
  system           = 943 # replace with your system or systemgroup ID
  path             = "/public_html"
  # version        = "8.4"  # pin to a specific version if needed
}

# -------------------------------------------------------------------
# MySQL database component
# The 'version' is auto-detected from the system's installed cookbooks
# when omitted. A 'pass' is required for MySQL.
# -------------------------------------------------------------------
resource "level27_app_component" "db" {
  app_id           = level27_app.my_project.id
  name             = "db"
  appcomponenttype = "mysql"
  system           = 943 # replace with your system or systemgroup ID
  pass             = "YourSecurePassword123!"
  # version        = "8.4"  # pin to a specific version if needed
}

# -------------------------------------------------------------------
# URL linked to the PHP component
# Set handle_dns = true only if Level27 manages the DNS for this domain.
# -------------------------------------------------------------------
resource "level27_app_component_url" "www" {
  app_id       = level27_app.my_project.id
  component_id = level27_app_component.php.id
  content      = "www.example.com"
  ssl_force    = true
  handle_dns   = false # set to true if Level27 manages DNS for this domain
}

# -------------------------------------------------------------------
# Let's Encrypt SSL certificate
# auto_ssl_certificate_urls must match an existing URL on the component.
# -------------------------------------------------------------------
resource "level27_ssl_certificate" "letsencrypt" {
  app_id                    = level27_app.my_project.id
  name                      = "letsencrypt-www"
  ssl_type                  = "letsencrypt"
  auto_ssl_certificate_urls = "www.example.com"
  auto_url_link             = true
  ssl_force                 = true

  depends_on = [level27_app_component_url.www]
}
