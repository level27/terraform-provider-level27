# URL attached to a PHP component.
# Set handle_dns = true only if Level27 manages DNS for this domain.
resource "level27_app_component_url" "www" {
  app_id       = level27_app.my_project.id
  component_id = level27_app_component.php.id
  content      = "www.example.com"
  ssl_force    = true
  handle_dns   = false # set to true if Level27 manages DNS for this domain
}
