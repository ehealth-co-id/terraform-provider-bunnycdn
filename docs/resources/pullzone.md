---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bunnycdn_pullzone Resource - terraform-provider-bunnycdn"
subcategory: ""
description: |-
  Pull zone resource
---

# bunnycdn_pullzone (Resource)

Pull zone resource

## Example Usage

```terraform
resource "bunnycdn_pullzone" "test" {
  name = "test-ehealth-co-id"
  origin_url = "https://lb.a.ehealth.id"
  origin_host_header = "test.ehealth.co.id"
  enable_smart_cache = true
  disable_cookie = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the pull zone.
- `origin_host_header` (String) Sets the host header that will be sent to the origin
- `origin_url` (String) Sets the origin URL of the pull zone

### Optional

- `disable_cookie` (Boolean) Sets disable cookie
- `enable_smart_cache` (Boolean) Sets the smart cache

### Read-Only

- `id` (Number) The ID of the pull zone

## Import

Import is supported using the following syntax:

```shell
terraform import bunnycdn_pullzone.test 1
```
