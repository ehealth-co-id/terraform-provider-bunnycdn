resource "bunnycdn_pullzone" "test" {
  name = "test-ehealth-co-id"
  origin_url = "https://lb.a.ehealth.id"
  origin_host_header = "test.ehealth.co.id"
  enable_smart_cache = true
  disable_cookie = false
}