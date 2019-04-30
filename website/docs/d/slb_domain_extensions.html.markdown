---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_domaine_extensions"
sidebar_current: "docs-alicloud-datasource-slb-domain-extensions"
description: |-
    Provides a list of server load balancer domain extensions to the user.
---

# alicloud\_slb_domain_extensions

This data source provides the domain extensions associated with a server load balancer listener.

## Example Usage

```
data "alicloud_slb_domain_extensions" "slb_domain_extensions" {
	ids = ["${alicloud_slb_domain_extension.example.id}"]
	load_balancer_id      = "${alicloud_slb.instance.id}"
	frontend_port         = "${alicloud_slb_listener.https.frontend_port}"
}
output "first_slb_domain_extension_id" {
  value = "${data.alicloud_slb_domain_extensions.slb_domain_extensions.slb_domain_extensions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - ID of the SLB with listener rules.
* `frontend_port` - SLB listener port.
* `ids` - (Optional) A list of rules IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `slb_domain_extensions` - A list of domain name extensions:
  * `id` - The ID of the domain name extension.
  * `domain` - The domain name.
  * `server_certificate_id` - The ID of the certificate used by the domain name
  * `url` - Path in the HTTP request where the rule applies (e.g. "/image").
  * `server_group_id` - ID of the linked VServer group.
