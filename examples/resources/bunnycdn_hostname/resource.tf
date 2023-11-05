resource "bunnycdn_hostname" "test" {
  pullzone_id = resource.bunnycdn_pullzone.test.id
  hostname = "test.ehealth.co.id"
  enable_ssl = true
  force_ssl = false
}