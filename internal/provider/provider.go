package provider

import (
	"context"

    "terraform-provider-bunnycdn/internal/provider/bunnycdn_api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &BunnyCdnProvider{}

type BunnyCdnProvider struct {
	version string
}

type BunnyCdnProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *BunnyCdnProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bunnycdn"
	resp.Version = p.version
}

func (p *BunnyCdnProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Your API Key is sent in the request header.",
				Optional:            true,
			},
		},
	}
}

func (p *BunnyCdnProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BunnyCdnProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	api := bunnycdn_api.NewBunnycdnApi(data.ApiKey.String())
	resp.DataSourceData = api
	resp.ResourceData = api
}

func (p *BunnyCdnProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPullZoneResource,
	}
}

func (p *BunnyCdnProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BunnyCdnProvider{
			version: version,
		}
	}
}
