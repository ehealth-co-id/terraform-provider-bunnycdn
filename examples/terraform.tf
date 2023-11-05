terraform {
  required_providers {
    bunnycdn = {
      source = "registry.terraform.io/ehealth-co-id/bunnycdn"
    }
  }
}

provider "bunnycdn" {
  api_key = "some-api-key"
}

resource "bunnycdn_pullzone" "test" {
  name = "test-ehealth-co-id"
  origin_url = "https://lb.a.ehealth.id"
  origin_host_header = "test.ehealth.co.id"
  enable_smart_cache = true
  disable_cookie = false
}

resource "bunnycdn_hostname" "test" {
  pullzone_id = resource.bunnycdn_pullzone.test.id
  hostname = "test.ehealth.co.id"
  enable_ssl = true
  force_ssl = false
}