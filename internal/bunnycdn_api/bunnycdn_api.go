package bunnycdn_api

type BunnycdnApi struct {
	ApiKey string
}

func NewBunnycdnApi(apiKey string) *BunnycdnApi {
	return &BunnycdnApi{
		ApiKey: apiKey,
	}
}
