# Let's Encrypt certificate — auto-linked to an existing URL.
# The URL must exist before the certificate is created (use depends_on).
resource "level27_ssl_certificate" "letsencrypt" {
  app_id                    = level27_app.my_project.id
  name                      = "letsencrypt-www"
  ssl_type                  = "letsencrypt"
  auto_ssl_certificate_urls = "www.example.com"
  auto_url_link             = true
  ssl_force                 = true

  depends_on = [level27_app_component_url.www]
}

# Custom / own certificate.
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
