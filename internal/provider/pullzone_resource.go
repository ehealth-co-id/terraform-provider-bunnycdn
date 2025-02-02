// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-bunnycdn/internal/bunnycdn_api"
	"terraform-provider-bunnycdn/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PullzoneResource{}
var _ resource.ResourceWithImportState = &PullzoneResource{}

func NewPullzoneResource() resource.Resource {
	return &PullzoneResource{}
}

type PullzoneResource struct {
	api *bunnycdn_api.BunnycdnApi
}

func (r *PullzoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pullzone"
}

func (r *PullzoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pull zone resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the pull zone.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"origin_type": schema.Int64Attribute{
				MarkdownDescription: "Sets the origin type of the pull zone (0 = OriginUrl, 2 = StorageZone)",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				PlanModifiers:       []planmodifier.Int64{},
			},
			"storage_zone_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the storage zone that will be used as the origin",
				Optional:            true,
				PlanModifiers:       []planmodifier.Int64{},
			},
			"origin_url": schema.StringAttribute{
				MarkdownDescription: "Sets the origin URL of the pull zone",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"origin_host_header": schema.StringAttribute{
				MarkdownDescription: "Sets the host header that will be sent to the origin",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"enable_smart_cache": schema.BoolAttribute{
				MarkdownDescription: "Sets the smart cache",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{},
			},
			"disable_cookie": schema.BoolAttribute{
				MarkdownDescription: "Sets disable cookie",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers:       []planmodifier.Bool{},
			},
			"error_page_enable_custom_code": schema.BoolAttribute{
				MarkdownDescription: "Sets enable custom error page",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers:       []planmodifier.Bool{},
			},
			"error_page_custom_code": schema.StringAttribute{
				MarkdownDescription: "Sets template custom error page",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the pull zone",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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

	api, ok := req.ProviderData.(*bunnycdn_api.BunnycdnApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected bunnycdn_api.BunnycdnApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func (r *PullzoneResource) Validate(ctx context.Context, data model.PullzoneResourceModel, diagnostics *diag.Diagnostics) {
	if data.OriginType.ValueInt64() != 0 && data.OriginType.ValueInt64() != 2 {
		diagnostics.AddError("Validation Error", "origin_type is not implemented")
	}
	if data.OriginType.ValueInt64() == 0 && data.OriginUrl.ValueStringPointer() == nil {
		diagnostics.AddError("Validation Error", "when origin_type = 0 origin_type must not be null")
	}
	if data.OriginType.ValueInt64() == 2 && data.StorageZoneId.ValueInt64Pointer() == nil {
		diagnostics.AddError("Validation Error", "when origin_type = 2 storage_zone_id must not be null")
	}
}

func (r *PullzoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.PullzoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	r.Validate(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	createdResource, err := r.api.PullzoneCreate(ctx, bunnycdn_api.PullzoneResourceModelToPullzone(data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pull zone, got error: %s", err))
		return
	}

	data = bunnycdn_api.PullzoneToPullzoneResourceModel(createdResource)
	tflog.Trace(ctx, "created a pull zone")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.PullzoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullzoneGet(ctx, data.Id.ValueInt64())
	if err != nil {
		pullzoneError, ok := err.(*model.PullzoneError)
		if ok && pullzoneError.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pull zone, got error: %s", err))
		return
	}

	data = bunnycdn_api.PullzoneToPullzoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.PullzoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	r.Validate(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.PullzoneUpdate(ctx, bunnycdn_api.PullzoneResourceModelToPullzone(data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pull zone, got error: %s", err))
		return
	}

	data = bunnycdn_api.PullzoneToPullzoneResourceModel(remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PullzoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.PullzoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.PullzoneDelete(ctx, bunnycdn_api.PullzoneResourceModelToPullzone(data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pull zone, got error: %s", err))
		return
	}
}

func (r *PullzoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
