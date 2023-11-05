package bunnycdn_api

import (
	"context"
	"fmt"
	"terraform-provider-bunnycdn/internal/model"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types" // import "encoding/json"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Hostname struct {
	Id        int64  `json:"Id"`
	Hostname  string `json:"Hostname"`
	ForceSsl  bool   `json:"ForceSSL"`
	EnableSsl bool
}

func HostnameToHostnameResourceModel(pullzoneId int64, resource *Hostname) model.HostnameResourceModel {
	return model.HostnameResourceModel{
		PullzoneId: types.Int64Value(pullzoneId),
		Id:         types.Int64Value(resource.Id),
		Hostname:   types.StringValue(resource.Hostname),
		EnableSsl:  types.BoolValue(resource.EnableSsl),
		ForceSsl:   types.BoolValue(resource.ForceSsl),
	}
}

func HostnameResourceModelToHostname(resource model.HostnameResourceModel) Hostname {
	return Hostname{
		Id:        resource.Id.ValueInt64(),
		Hostname:  resource.Hostname.ValueString(),
		EnableSsl: resource.EnableSsl.ValueBool(),
		ForceSsl:  resource.ForceSsl.ValueBool(),
	}
}

func (api *BunnycdnApi) HostnameGet(ctx context.Context, pullzoneId int64, hostname string) (*Hostname, error) {
	tflog.Info(ctx, "hostname get")

	pullzone, err := api.PullzoneGet(ctx, pullzoneId)
	if err != nil {
		return nil, err
	}

	for _, item := range pullzone.Hostnames {
		if item.Value == hostname {
			return &Hostname{
				Id:        item.Id,
				Hostname:  item.Value,
				EnableSsl: item.HasCertificate,
				ForceSsl:  item.ForceSsl,
			}, nil
		}
	}
	return nil, model.NewHostnameError(404, hostname)
}

func (api *BunnycdnApi) HostnameCreate(ctx context.Context, pullzoneId int64, resource Hostname) error {
	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetBody(&resource).
		Post(fmt.Sprintf("https://api.bunny.net/pullzone/%d/addHostname", pullzoneId))

	if err != nil {
		return err
	}

	if response.StatusCode() == 204 {
		return nil
	}

	return model.NewHostnameError(response.StatusCode(), resource.Hostname)
}

func (api *BunnycdnApi) HostnameDelete(ctx context.Context, pullzoneId int64, resource Hostname) error {
	tflog.Info(ctx, "hostname delete")

	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetBody(&resource).
		Delete(fmt.Sprintf("https://api.bunny.net/pullzone/%d/removeHostname", pullzoneId))

	if err != nil {
		return err
	}

	if response.StatusCode() == 204 {
		return nil
	}

	return model.NewHostnameError(response.StatusCode(), resource.Hostname)
}

func (api *BunnycdnApi) HostnameEnableSsl(ctx context.Context, pullzoneId int64, resource Hostname) error {
	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetQueryParams(map[string]string{
			"hostname": resource.Hostname,
		}).
		Get(fmt.Sprintf("https://api.bunny.net/pullzone/loadFreeCertificate"))

	if err != nil {
		return err
	}

	if response.StatusCode() == 200 {
		return nil
	}

	return model.NewEnableSslError(response.StatusCode(), resource.Hostname, string(response.Body()))
}

func (api *BunnycdnApi) HostnameUpdateForceSsl(ctx context.Context, pullzoneId int64, resource Hostname) error {
	response, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("AccessKey", api.ApiKey).
		SetBody(map[string]interface{}{
			"Hostname": resource.Hostname,
			"ForceSSL": resource.ForceSsl,
		}).
		Post(fmt.Sprintf("https://api.bunny.net/pullzone/%d/setForceSSL", pullzoneId))

	if err != nil {
		return err
	}

	if response.StatusCode() == 204 {
		return nil
	}

	return model.NewEnableSslError(response.StatusCode(), resource.Hostname, string(response.Body()))
}
