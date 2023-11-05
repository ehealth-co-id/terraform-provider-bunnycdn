package bunnycdn_api

import "fmt"
import "errors"
import "github.com/go-resty/resty/v2"

type PullZone struct {
    Id        int64
    Name      string
    OriginUrl string
}

func (api *BunnycdnApi) PullZoneGet(id int64) (*PullZone, error) {
	var resource PullZone

    response, err := resty.New().R().
        SetHeader("AccessKey", api.ApiKey).
        SetResult(&resource).
        Get(fmt.Sprintf("https://api.bunny.net/pullzone/%d", id))

    if (err != nil) {
        return nil, err
    }

    if (response.Status() == "200") {
        return &resource, nil
    }

    if (response.Status() == "401") {
        return nil, errors.New(fmt.Sprintf("Request authorization failed", resource.Id))
    }
    if (response.Status() == "404") {
        return nil, errors.New(fmt.Sprintf("Pull zone with ID %d does not exist", resource.Id))
    }
    return nil, errors.New(fmt.Sprintf("Unexpected status code %d on POST https://api.bunny.net/pullzone/", response.Status()))
}

func (api *BunnycdnApi) PullZoneCreate(resource PullZone) (*PullZone, error) {
	var createdResource PullZone

    response, err := resty.New().R().
        SetHeader("Content-Type", "application/json").
        SetBody(&resource).
        SetResult(&createdResource).
        Post("https://api.bunny.net/pullzone/")
    
    if (err != nil) {
        return nil, err
    }

    if (response.Status() == "200") {
        return &createdResource, nil
    }

    if (response.Status() == "401") {
        return nil, errors.New(fmt.Sprintf("Request authorization failed", resource.Id))
    }
    if (response.Status() == "404") {
        return nil, errors.New(fmt.Sprintf("Pull zone with ID %d does not exist", resource.Id))
    }
    return nil, errors.New(fmt.Sprintf("Unexpected status code %d on POST https://api.bunny.net/pullzone/", response.Status()))
}

func (api *BunnycdnApi) PullZoneUpdate(resource PullZone) (*PullZone, error) {
	var updatedResource PullZone

    response, err := resty.New().R().
        SetHeader("Content-Type", "application/json").
        SetBody(&resource).
        SetResult(&updatedResource).
        Post(fmt.Sprintf("https://api.bunny.net/pullzone/%d", resource.Id))

    if (err != nil) {
        return nil, err
    }

    if (response.Status() == "200") {
        return &updatedResource, nil
    }

    if (response.Status() == "401") {
        return nil, errors.New(fmt.Sprintf("Request authorization failed", resource.Id))
    }
    if (response.Status() == "400") {
        return nil, errors.New(fmt.Sprintf("Failed configuring the storage zone; model validation failed", resource.Id))
    }
    if (response.Status() == "404") {
        return nil, errors.New(fmt.Sprintf("Pull zone with ID %d does not exist", resource.Id))
    }
    return nil, errors.New(fmt.Sprintf("Unexpected status code %d on POST https://api.bunny.net/pullzone/%d", response.Status(), resource.Id))
}

func (api *BunnycdnApi) PullZoneDelete(resource PullZone) error {
    response, err := resty.New().R().
        Delete(fmt.Sprintf("https://api.bunny.net/pullzone/%d", resource.Id))

    if (err != nil) {
        return err
    }

    if (response.Status() == "200") {
        return nil
    }

    if (response.Status() == "401") {
        return errors.New(fmt.Sprintf("Request authorization failed", resource.Id))
    }
    if (response.Status() == "404") {
        return errors.New(fmt.Sprintf("Pull zone with ID %d does not exist", resource.Id))
    }
    return errors.New(fmt.Sprintf("Unexpected status code %d on DELETE https://api.bunny.net/pullzone/%d", response.Status(), resource.Id))
}