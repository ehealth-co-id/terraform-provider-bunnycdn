// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

    "terraform-provider-bunnycdn/internal/provider/bunnycdn_api"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PullZoneResource{}
var _ resource.ResourceWithImportState = &PullZoneResource{}

func NewPullZoneResource() resource.Resource {
	return &PullZoneResource{}
}

type PullZoneResource struct {
	api bunnycdn_api.BunnycdnApi
}

type PullZoneResourceModel struct {
	Name                  types.String `tfsdk:"name"`
	OriginUrl             types.String `tfsdk:"origin_url"`
	Id                    types.Int64  `tfsdk:"id"`
}

func (r *PullZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *PullZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pull zone resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the pull zone.",
				Optional:            false,
			},
			"origin_url": schema.StringAttribute{
				MarkdownDescription: "Sets the origin URL of the pull zone",
				Optional:            false,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the pull zone",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PullZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(bunnycdn_api.BunnycdnApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *bunnycdn_api.BunnycdnApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func PullZoneToPullZoneResourceModel(resource *bunnycdn_api.PullZone) (PullZoneResourceModel) {
    return PullZoneResourceModel{
        Id: types.Int64Value(resource.Id),
        Name: types.StringValue(resource.Name),
        OriginUrl: types.StringValue(resource.OriginUrl),
    }
}

func PullZoneResourceModelToPullZone(resource PullZoneResourceModel) (bunnycdn_api.PullZone) {
    return bunnycdn_api.PullZone{
        Id: resource.Id.ValueInt64(),
        Name: resource.Name.ValueString(),
        OriginUrl: resource.OriginUrl.ValueString(),
    }
}

func (r *PullZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PullZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdResource, err := r.api.PullZoneCreate(bunnycdn_api.PullZone{
		Name: data.Name.ValueString(),
		OriginUrl: data.OriginUrl.ValueString(),
	})
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pull zone, got error: %s", err))
	    return
	}

	data = PullZoneToPullZoneResourceModel(createdResource)
	tflog.Trace(ctx, "created a pull zone")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PullZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullZoneGet(data.Id.ValueInt64())
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pull zone, got error: %s", err))
	    return
	}

	data = PullZoneToPullZoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PullZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullZoneUpdate(PullZoneResourceModelToPullZone(data))
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pull zone, got error: %s", err))
	    return
	}

	data = PullZoneToPullZoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PullZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.PullZoneDelete(PullZoneResourceModelToPullZone(data))
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pull zone, got error: %s", err))
	    return
	}
}

func (r *PullZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
