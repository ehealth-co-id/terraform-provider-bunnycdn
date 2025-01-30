package model

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PullzoneResourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	OriginType       types.Int64  `tfsdk:"origin_type"`
	StorageZoneId    types.Int64  `tfsdk:"storage_zone_id"`
	OriginUrl        types.String `tfsdk:"origin_url"`
	OriginHostHeader types.String `tfsdk:"origin_host_header"`
	EnableSmartCache types.Bool   `tfsdk:"enable_smart_cache"`
	DisableCookies   types.Bool   `tfsdk:"disable_cookie"`
	ErrorPageEnableCustomCode	types.Bool		`tfsdk:"error_page_enable_custom_code"`
	ErrorPageCustomCode			types.String	`tfsdk:"error_page_custom_code"`
}

type PullzoneError struct {
	StatusCode int
	PullzoneId int64
}

func NewPullzoneError(statusCode int, pullzoneId int64) *PullzoneError {
	return &PullzoneError{
		StatusCode: statusCode,
		PullzoneId: pullzoneId,
	}
}

func (e *PullzoneError) Error() string {
	if e.StatusCode == 400 {
		return "Invalid request"
	}
	if e.StatusCode == 401 {
		return "Request authorization failed"
	}
	if e.StatusCode == 404 {
		return fmt.Sprintf("Pull zone with ID %d does not exist", e.PullzoneId)
	}
	if e.StatusCode >= 500 {
		return fmt.Sprintf("Bunnycdn server error. status code: %d", e.StatusCode)
	}
	return fmt.Sprintf("Unexpected status code %d", e.StatusCode)
}
