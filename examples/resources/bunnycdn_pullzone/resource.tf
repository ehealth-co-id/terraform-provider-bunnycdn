resource "bunnycdn_pullzone" "test" {
  name = "test-ehealth-co-id"
  origin_type = 0
  origin_url = "https://lb.a.ehealth.id"
  origin_host_header = "test.ehealth.co.id"
  # origin_type = 2
  # storage_zone_id = 999999
  enable_smart_cache = true
  disable_cookie = false
}