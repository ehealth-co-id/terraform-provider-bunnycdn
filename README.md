# Terraform Provider BunnyCDN

A Terraform provider for managing BunnyCDN resources. This provider allows you to manage your BunnyCDN infrastructure as code, including pull zones, hostnames, and certificates.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Installation

### Terraform Registry

This provider is available on the [Terraform Registry](https://registry.terraform.io/providers/ehealth-co-id/bunnycdn/latest). Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    bunnycdn = {
      source  = "ehealth-co-id/bunnycdn"
      version = ">= 1.0.0"
    }
  }
}

provider "bunnycdn" {
  api_key = "your-bunnycdn-api-key"
}
```

### Local Development Installation

For local development or testing:

1. Clone this repository
2. Build the provider:
   ```
   go build -o terraform-provider-bunnycdn
   ```
3. Move the binary to the appropriate plugin directory

## Features

This provider currently supports the following BunnyCDN resources:

- **Pull Zones** - Create and manage pull zones
- **Hostnames** - Add and configure custom hostnames for your pull zones

## Usage Examples

### Pull Zone

```hcl
resource "bunnycdn_pullzone" "example" {
  name             = "example-pull-zone"
  origin_url       = "https://example.com"
  origin_type      = 0  # 0 = OriginUrl, 2 = StorageZone
  
  # Optional settings
  origin_host_header           = "example.com"
  enable_smart_cache           = true
  disable_cookie               = false
  error_page_enable_custom_code = false
  error_page_custom_code       = ""
}
```

### Hostname with Free SSL

```hcl
resource "bunnycdn_hostname" "example" {
  pull_zone_id = bunnycdn_pullzone.example.id
  hostname     = "cdn.example.com"
  
  # Using free SSL
  force_ssl    = true
}
```

### Hostname with Custom Certificate

```hcl
resource "bunnycdn_hostname" "custom_cert" {
  pull_zone_id    = bunnycdn_pullzone.example.id
  hostname        = "secure.example.com"
  force_ssl       = true
  
  # Using custom certificate
  certificate     = file("path/to/certificate.crt")
  certificate_key = file("path/to/private.key")
}
```

## Development

### Building the Provider

1. Clone the repository
2. Build the provider using the Go compiler:
   ```
   go build -o terraform-provider-bunnycdn
   ```

### Testing

To run the acceptance tests:

```
make testacc
```

or directly:

```
TF_ACC=1 go test ./... -v -timeout 120m
```

### Generating Documentation

Generate provider documentation:

```
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Mozilla Public License 2.0 - see the LICENSE file for details.

## Acknowledgements

- [BunnyCDN API Documentation](https://docs.bunny.net/reference/bunnynet-api-overview)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)