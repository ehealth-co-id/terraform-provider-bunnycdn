package bunnycdn_api

import "fmt"
import "context"
// import "encoding/json"
import "github.com/go-resty/resty/v2"
import "terraform-provider-bunnycdn/internal/model"
import "github.com/hashicorp/terraform-plugin-log/tflog"
import "github.com/hashicorp/terraform-plugin-framework/types"

type Hostname struct {
	Id       int64  `json:"Id"`
	Hostname string `json:"Hostname"`
}

func HostnameToHostnameResourceModel(pullzoneId int64, resource *Hostname) (model.HostnameResourceModel) {
    return model.HostnameResourceModel{
		PullzoneId: types.Int64Value(pullzoneId),
        Id: types.Int64Value(resource.Id),
        Hostname: types.StringValue(resource.Hostname),
    }
}

func HostnameResourceModelToHostname(resource model.HostnameResourceModel) (Hostname) {
    return Hostname{
        Id: resource.Id.ValueInt64(),
        Hostname: resource.Hostname.ValueString(),
    }
}

func (api *BunnycdnApi) HostnameGet(ctx context.Context, pullzoneId int64, hostname string) (*Hostname, error) {
    tflog.Info(ctx, "hostname get")

    pullzone, err := api.PullzoneGet(ctx, pullzoneId)
    if (err != nil) {
        return nil, err
    }

    for _, item := range pullzone.Hostnames {
        if item.Value == hostname {
            return &Hostname {
                Id: item.Id,
                Hostname: item.Value,
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


    if (err != nil) {
        return err
    }

    if (response.StatusCode() == 204) {
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

    if (err != nil) {
        return err
    }

    if (response.StatusCode() == 204) {
        return nil
    }

    return model.NewHostnameError(response.StatusCode(), resource.Hostname)
}