package model

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HostnameResourceModel struct {
	PullzoneId     types.Int64  `tfsdk:"pullzone_id"`
	Id             types.Int64  `tfsdk:"id"`
	Hostname       types.String `tfsdk:"hostname"`
	EnableSsl      types.Bool   `tfsdk:"enable_ssl"`
	ForceSsl       types.Bool   `tfsdk:"force_ssl"`
	Certificate    types.String `tfsdk:"certificate"`
	CertificateKey types.String `tfsdk:"certificate_key"`
}

type HostnameError struct {
	StatusCode int
	Hostname   string
}

func NewHostnameError(statusCode int, hostname string) *HostnameError {
	return &HostnameError{
		StatusCode: statusCode,
		Hostname:   hostname,
	}
}

func (e *HostnameError) Error() string {
	if e.StatusCode == 400 {
		return "Invalid request"
	}
	if e.StatusCode == 401 {
		return "Request authorization failed"
	}
	if e.StatusCode == 404 {
		return fmt.Sprintf("Hostname %s does not exist", e.Hostname)
	}
	if e.StatusCode >= 500 {
		return fmt.Sprintf("Bunnycdn server error. status code: %d", e.StatusCode)
	}
	return fmt.Sprintf("Unexpected status code %d", e.StatusCode)
}

type EnableSslError struct {
	StatusCode int
	Hostname   string
	Body       string
}

func NewEnableSslError(statusCode int, hostname string, body string) *EnableSslError {
	return &EnableSslError{
		StatusCode: statusCode,
		Hostname:   hostname,
		Body:       body,
	}
}

func (e *EnableSslError) Error() string {
	if e.StatusCode == 400 {
		return fmt.Sprintf("Failed configuring free certificate. response: %s", e.Body)
	}
	if e.StatusCode == 401 {
		return "Request authorization failed"
	}
	if e.StatusCode == 404 {
		return fmt.Sprintf("Hostname %s does not exist", e.Hostname)
	}
	if e.StatusCode >= 500 {
		return fmt.Sprintf("Bunnycdn server error. status code: %d", e.StatusCode)
	}
	return fmt.Sprintf("Unexpected status code %d", e.StatusCode)
}

type DeleteCertificateError struct {
	StatusCode int
	Hostname   string
	Body       string
}

func NewDeleteCertificateError(statusCode int, hostname string, body string) *DeleteCertificateError {
	return &DeleteCertificateError{
		StatusCode: statusCode,
		Hostname:   hostname,
		Body:       body,
	}
}

func (e *DeleteCertificateError) Error() string {
	if e.StatusCode == 400 {
		return fmt.Sprintf("Failed removing certificate. response: %s", e.Body)
	}
	if e.StatusCode == 401 {
		return "Request authorization failed"
	}
	if e.StatusCode == 404 {
		return fmt.Sprintf("Hostname %s does not exist", e.Hostname)
	}
	if e.StatusCode >= 500 {
		return fmt.Sprintf("Bunnycdn server error. status code: %d", e.StatusCode)
	}
	return fmt.Sprintf("Unexpected status code %d", e.StatusCode)
}
