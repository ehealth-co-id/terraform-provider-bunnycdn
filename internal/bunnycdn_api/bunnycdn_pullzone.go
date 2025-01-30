package bunnycdn_api

import (
	"context"
	"fmt"
	"terraform-provider-bunnycdn/internal/model"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PullzoneHostname struct {
	Id             int64  `json:"Id"`
	Value          string `json:"Value"`
	HasCertificate bool   `json:"HasCertificate"`
	ForceSsl       bool   `json:"ForceSSL"`
}

type Pullzone struct {
	Id               int64              `json:"Id"`
	Name             string             `json:"Name"`
	OriginType       int64              `json:"OriginType"`
	StorageZoneId    *int64             `json:"StorageZoneId"`
	OriginUrl        *string            `json:"OriginUrl"`
	OriginHostHeader *string            `json:"OriginHostHeader"`
	EnableTLS1       bool               `json:"EnableTLS1"`
	EnableTLS1_1     bool               `json:"EnableTLS1_1"`
	EnableSmartCache bool               `json:"EnableSmartCache"`
	DisableCookies   bool               `json:"DisableCookies"`
	Hostnames        []PullzoneHostname `json:"Hostnames"`
	ErrorPageEnableCustomCode	bool	`json:"ErrorPageEnableCustomCode"`
	ErrorPageCustomCode			string	`json:"ErrorPageCustomCode"`
}

func ifEmptyThenNil(value *string) *string {
	if *value == "" {
		return nil
	}
	return value
}

func ifZeroThenNil(value *int64) *int64 {
	if *value == 0 {
		return nil
	}
	return value
}

func PullzoneToPullzoneResourceModel(resource *Pullzone) model.PullzoneResourceModel {

	return model.PullzoneResourceModel{
		Id:               types.Int64Value(resource.Id),
		Name:             types.StringValue(resource.Name),
		OriginType:       types.Int64Value(resource.OriginType),
		StorageZoneId:    types.Int64PointerValue(ifZeroThenNil(resource.StorageZoneId)),
		OriginUrl:        types.StringPointerValue(ifEmptyThenNil(resource.OriginUrl)),
		EnableSmartCache: types.BoolValue(resource.EnableSmartCache),
		DisableCookies:   types.BoolValue(resource.DisableCookies),
		OriginHostHeader: types.StringPointerValue(ifEmptyThenNil(resource.OriginHostHeader)),
		ErrorPageEnableCustomCode: types.BoolValue(resource.ErrorPageEnableCustomCode),
		ErrorPageCustomCode: types.StringValue(ifEmptyThenNil(resource.ErrorPageCustomCode)),
	}
}

func PullzoneResourceModelToPullzone(resource model.PullzoneResourceModel) Pullzone {
	return Pullzone{
		Id:               resource.Id.ValueInt64(),
		Name:             resource.Name.ValueString(),
		OriginType:       resource.OriginType.ValueInt64(),
		StorageZoneId:    resource.StorageZoneId.ValueInt64Pointer(),
		OriginUrl:        resource.OriginUrl.ValueStringPointer(),
		EnableSmartCache: resource.EnableSmartCache.ValueBool(),
		DisableCookies:   resource.DisableCookies.ValueBool(),
		OriginHostHeader: resource.OriginHostHeader.ValueStringPointer(),
		ErrorPageEnableCustomCode:	resource.ErrorPageEnableCustomCode.ValueBool(),
		ErrorPageCustomCode:		resource.ErrorPageCustomCode.ValueString(),
		EnableTLS1:       false,
		EnableTLS1_1:     false,
	}
}

func (api *BunnycdnApi) PullzoneGet(ctx context.Context, id int64) (*Pullzone, error) {
	var resource Pullzone

	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("AccessKey", api.ApiKey).
		SetResult(&resource).
		Get(fmt.Sprintf("https://api.bunny.net/pullzone/%d", id))

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 200 {
		return &resource, nil
	}

	return nil, model.NewPullzoneError(response.StatusCode(), id)
}

func (api *BunnycdnApi) PullzoneCreate(ctx context.Context, resource Pullzone) (*Pullzone, error) {
	var createdResource Pullzone

	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetBody(&resource).
		SetResult(&createdResource).
		Post("https://api.bunny.net/pullzone/")

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 201 {
		return &createdResource, nil
	}

	return nil, model.NewPullzoneError(response.StatusCode(), resource.Id)
}

func (api *BunnycdnApi) PullzoneUpdate(ctx context.Context, resource Pullzone) (*Pullzone, error) {
	var updatedResource Pullzone

	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetBody(&resource).
		SetResult(&updatedResource).
		Post(fmt.Sprintf("https://api.bunny.net/pullzone/%d", resource.Id))

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 200 {
		return &updatedResource, nil
	}

	return nil, model.NewPullzoneError(response.StatusCode(), resource.Id)
}

func (api *BunnycdnApi) PullzoneDelete(ctx context.Context, resource Pullzone) error {
	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("AccessKey", api.ApiKey).
		Delete(fmt.Sprintf("https://api.bunny.net/pullzone/%d", resource.Id))

	if err != nil {
		return err
	}

	if response.StatusCode() == 204 {
		return nil
	}

	return model.NewPullzoneError(response.StatusCode(), resource.Id)
}
