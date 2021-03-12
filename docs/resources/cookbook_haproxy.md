---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "level27_cookbook_haproxy Resource - terraform-provider-level27"
subcategory: ""
description: |-
  
---

# level27_cookbook_haproxy (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **system_id** (String) The ID of the system where to install this cookbook.

### Optional

- **backend_ip** (String) Ip of the backend the loadbalancer should proxy to
- **backend_port** (String) Port of the backend the loadbalancer should proxy to
- **expected_traffic** (String) The amount of traffic that should be expected
- **frontend_port** (String) Port the loadbalancer should bind to
- **haip_ipv4** (String) High available ipv4 address
- **haip_ipv6** (String) High available ipv6 address
- **haip_routerid** (String) Router id of the high available ip addresses
- **id** (String) The ID of this resource.
- **primay** (Boolean) Is the loadbalancer the primary one
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **varnish** (Boolean) Is varnish present on the webservers

### Read-Only

- **cookbook_id** (String)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)

