package bunnycdn_api

import "fmt"
import "context"
import "github.com/go-resty/resty/v2"
import "terraform-provider-bunnycdn/internal/model"
import "github.com/hashicorp/terraform-plugin-framework/types"

type PullzoneHostname struct {
    Id       int64  `json:"Id"`
    Value    string `json:"Value"`
}

type Pullzone struct {
    Id int64 `json:"Id"`
    Name string `json:"Name"`
    OriginUrl string `json:"OriginUrl"`
    OriginHostHeader string `json:"OriginHostHeader"`
    EnableSmartCache bool `json:"EnableSmartCache"`
    DisableCookies bool `json:"DisableCookies"`
    Hostnames []PullzoneHostname `json:"Hostnames"`
}

func PullzoneToPullzoneResourceModel(resource *Pullzone) (model.PullzoneResourceModel) {
    return model.PullzoneResourceModel{
        Id: types.Int64Value(resource.Id),
        Name: types.StringValue(resource.Name),
        OriginUrl: types.StringValue(resource.OriginUrl),
        EnableSmartCache: types.BoolValue(resource.EnableSmartCache),
        DisableCookies: types.BoolValue(resource.DisableCookies),
		OriginHostHeader: types.StringValue(resource.OriginHostHeader),
    }
}

func PullzoneResourceModelToPullzone(resource model.PullzoneResourceModel) (Pullzone) {
    return Pullzone{
        Id: resource.Id.ValueInt64(),
        Name: resource.Name.ValueString(),
        OriginUrl: resource.OriginUrl.ValueString(),
		EnableSmartCache: resource.EnableSmartCache.ValueBool(),
		DisableCookies: resource.DisableCookies.ValueBool(),
        OriginHostHeader: resource.OriginHostHeader.ValueString(),
    }
}

func (api *BunnycdnApi) PullzoneGet(ctx context.Context, id int64) (*Pullzone, error) {
	var resource Pullzone

    response, err := resty.New().R().
        SetContext(ctx).
        SetHeader("AccessKey", api.ApiKey).
        SetResult(&resource).
        Get(fmt.Sprintf("https://api.bunny.net/pullzone/%d", id))

    if (err != nil) {
        return nil, err
    }

    if (response.StatusCode() == 200) {
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
    
    if (err != nil) {
        return nil, err
    }

    if (response.StatusCode() == 201) {
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

    if (err != nil) {
        return nil, err
    }

    if (response.StatusCode() == 200) {
        return &updatedResource, nil
    }

    return nil, model.NewPullzoneError(response.StatusCode(), resource.Id)
}

func (api *BunnycdnApi) PullzoneDelete(ctx context.Context, resource Pullzone) error {
    response, err := resty.New().R().
        SetContext(ctx).
        SetHeader("AccessKey", api.ApiKey).
        Delete(fmt.Sprintf("https://api.bunny.net/pullzone/%d", resource.Id))    

    if (err != nil) {
        return err
    }

    if (response.StatusCode() == 204) {
        return nil
    }

    return model.NewPullzoneError(response.StatusCode(), resource.Id)
}