# Minimal example — PHP component with auto-detected version.
resource "level27_app_component" "php" {
  app_id           = level27_app.my_project.id
  name             = "web"
  appcomponenttype = "php"
  system           = 943    # replace with your system ID
  path             = "/public_html"
  # version is auto-detected as the highest available version on the system.
  # Pin it explicitly if you need a specific version:
  # version = "8.2"
}

# MySQL component — version is auto-detected from the system's cookbooks.
resource "level27_app_component" "db" {
  app_id           = level27_app.my_project.id
  name             = "db"
  appcomponenttype = "mysql"
  system           = 943    # replace with your system ID
  pass             = "YourSecurePassword123!"
  # version = "8.4"  # pin if needed
}
