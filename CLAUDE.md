# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains a Terraform provider for BunnyCDN, allowing Terraform users to manage BunnyCDN resources through infrastructure as code. The provider is built using the Terraform Plugin Framework.

## Code Architecture

- `/internal/provider/`: Contains the implementation of the Terraform provider, resource definitions, and data sources.
- `/internal/bunnycdn_api/`: Contains the BunnyCDN API client used to interact with the BunnyCDN API.
- `/internal/model/`: Contains the data models used by the provider and API client.
- `/docs/`: Contains the auto-generated documentation for the provider.
- `/examples/`: Contains example Terraform configurations for using the provider.

The provider follows the standard Terraform provider architecture:
1. `provider.go` defines the provider schema and configuration
2. Resource implementation files (`pullzone_resource.go`, `hostname_resource.go`) define the resource schemas and CRUD operations
3. API client files handle the actual HTTP communication with the BunnyCDN API

## Common Development Commands

### Building the Provider

Build the provider:

```
go build -o terraform-provider-bunnycdn
```

Install the provider locally for testing:

```
go install
```

### Testing

Run acceptance tests:

```
make testacc
```

or directly:

```
TF_ACC=1 go test ./... -v -timeout 120m
```

Run a specific test:

```
TF_ACC=1 go test ./internal/provider -v -run=TestAccPullzoneResource -timeout 120m
```

### Generating Documentation

Generate provider documentation:

```
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
```

### Debugging

Run the provider with debugging support:

```
go run main.go -debug
```

## Terraform Provider Usage

Configure the provider:

```terraform
provider "bunnycdn" {
  api_key = "your-api-key"
}
```

## Development Notes

- The provider uses the Terraform Plugin Framework (not the older SDK)
- Resource implementations follow the standard CRUD pattern
- Error handling includes special handling for 404 responses to implement graceful handling of deleted resources
- The API client uses Go Resty for HTTP requests