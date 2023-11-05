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

var _ resource.Resource = &PullzoneResource{}
var _ resource.ResourceWithImportState = &PullzoneResource{}

func NewPullzoneResource() resource.Resource {
	return &PullzoneResource{}
}

type PullzoneResource struct {
	api bunnycdn_api.BunnycdnApi
}

type PullzoneResourceModel struct {
	Name                  types.String `tfsdk:"name"`
	OriginUrl             types.String `tfsdk:"origin_url"`
	Id                    types.Int64  `tfsdk:"id"`
}

func (r *PullzoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Pullzone"
}

func (r *PullzoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *PullzoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func PullzoneToPullzoneResourceModel(resource *bunnycdn_api.Pullzone) (PullzoneResourceModel) {
    return PullzoneResourceModel{
        Id: types.Int64Value(resource.Id),
        Name: types.StringValue(resource.Name),
        OriginUrl: types.StringValue(resource.OriginUrl),
    }
}

func PullzoneResourceModelToPullzone(resource PullzoneResourceModel) (bunnycdn_api.Pullzone) {
    return bunnycdn_api.Pullzone{
        Id: resource.Id.ValueInt64(),
        Name: resource.Name.ValueString(),
        OriginUrl: resource.OriginUrl.ValueString(),
    }
}

func (r *PullzoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PullzoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdResource, err := r.api.PullzoneCreate(bunnycdn_api.Pullzone{
		Name: data.Name.ValueString(),
		OriginUrl: data.OriginUrl.ValueString(),
	})
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pull zone, got error: %s", err))
	    return
	}

	data = PullzoneToPullzoneResourceModel(createdResource)
	tflog.Trace(ctx, "created a pull zone")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PullzoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullzoneGet(data.Id.ValueInt64())
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pull zone, got error: %s", err))
	    return
	}

	data = PullzoneToPullzoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PullzoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullzoneUpdate(PullzoneResourceModelToPullzone(data))
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pull zone, got error: %s", err))
	    return
	}

	data = PullzoneToPullzoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PullzoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.PullzoneDelete(PullzoneResourceModelToPullzone(data))
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pull zone, got error: %s", err))
	    return
	}
}

func (r *PullzoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
