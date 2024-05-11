// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-bunnycdn/internal/bunnycdn_api"
	"terraform-provider-bunnycdn/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.Resource = &HostnameResource{}
var _ resource.ResourceWithImportState = &HostnameResource{}

func NewHostnameResource() resource.Resource {
	return &HostnameResource{}
}

type HostnameResource struct {
	api *bunnycdn_api.BunnycdnApi
}

type Certificate struct {
	Certificate    *string `json:"Certificate"`
	CertificateKey *string `json:"CertificateKey"`
}

func (r *HostnameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hostname"
}

func (r *HostnameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hostname resource",

		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The name of the hostname.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pullzone_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the pull zone",
				PlanModifiers:       []planmodifier.Int64{},
			},
			"enable_ssl": schema.BoolAttribute{
				MarkdownDescription: "Sets enable SSL",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{},
			},
			"force_ssl": schema.BoolAttribute{
				MarkdownDescription: "Sets force SSL",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the pull zone",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "Hostname custom certificate",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{},
			},
			"certificate_key": schema.StringAttribute{
				MarkdownDescription: "Hostname custom certificate key",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{},
			},
		},
	}
}

func (r *HostnameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostnameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.HostnameResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.HostnameCreate(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create hostname, got error: %s", err))
		return
	}

	if data.EnableSsl.ValueBool() {
		if data.Certificate.ValueStringPointer() == nil {
			err := r.api.HostnameLoadFreeCertificate(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
			if err != nil {
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to load free certificate, got error: %s", err))
			}
		} else {
			err := r.api.HostnameAddCertificate(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
			if err != nil {
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to add certificate, got error: %s", err))
			}
			certificateEncoded, err := json.Marshal(Certificate{
				Certificate:    data.Certificate.ValueStringPointer(),
				CertificateKey: data.CertificateKey.ValueStringPointer(),
			})
			if err != nil {
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to encode certificate to json: %s", err))
			}
			addCertificate := resp.Private.SetKey(ctx, "certificate", certificateEncoded)
			resp.Diagnostics.Append(addCertificate...)
		}

		err = r.api.HostnameUpdateForceSsl(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
		if err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to update force_ssl, got error: %s", err))
		}
	}

	// get and update state after updates enable_ssl and force_ssl
	remoteResource, err := r.api.HostnameGet(ctx, data.PullzoneId.ValueInt64(), data.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hostname, got error: %s", err))
		return
	}

	if data.EnableSsl.ValueBool() && data.Certificate.ValueStringPointer() != nil {
		remoteResource.Certificate = data.Certificate.ValueStringPointer()
		remoteResource.CertificateKey = data.CertificateKey.ValueStringPointer()
	}

	data = bunnycdn_api.HostnameToHostnameResourceModel(data.PullzoneId.ValueInt64(), remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostnameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.HostnameResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	remoteResource, err := r.api.HostnameGet(ctx, data.PullzoneId.ValueInt64(), data.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hostname, got error: %s", err))
		return
	}

	certificateEncoded, getCertificate := req.Private.GetKey(ctx, "certificate")
	resp.Diagnostics.Append(getCertificate...)
	if certificateEncoded != nil {
		var certificate Certificate
		err := json.Unmarshal(certificateEncoded, &certificate)
		if err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to decode certificate from json: %s", err))
		}
		remoteResource.Certificate = certificate.Certificate
		remoteResource.CertificateKey = certificate.CertificateKey
	}

	data = bunnycdn_api.HostnameToHostnameResourceModel(data.PullzoneId.ValueInt64(), remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostnameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state model.HostnameResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.EnableSsl.ValueBool() {
		if !state.EnableSsl.ValueBool() {
			if data.Certificate.ValueStringPointer() == nil {
				err := r.api.HostnameLoadFreeCertificate(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
				if err != nil {
					resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to load free certificate, got error: %s", err))
				}
			} else {
				err := r.api.HostnameAddCertificate(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
				if err != nil {
					resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to add certificate, got error: %s", err))
				}
				certificateEncoded, err := json.Marshal(Certificate{
					Certificate:    data.Certificate.ValueStringPointer(),
					CertificateKey: data.CertificateKey.ValueStringPointer(),
				})
				if err != nil {
					resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to encode certificate to json: %s", err))
				}
				addCertificate := resp.Private.SetKey(ctx, "certificate", certificateEncoded)
				resp.Diagnostics.Append(addCertificate...)
			}
		}

		err := r.api.HostnameUpdateForceSsl(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
		if err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to update force_ssl, got error: %s", err))
		}

		// get and update state after updates enable_ssl and force_ssl
		remoteResource, err := r.api.HostnameGet(ctx, data.PullzoneId.ValueInt64(), data.Hostname.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hostname, got error: %s", err))
			return
		}

		if data.EnableSsl.ValueBool() && data.Certificate.ValueStringPointer() != nil {
			remoteResource.Certificate = data.Certificate.ValueStringPointer()
			remoteResource.CertificateKey = data.CertificateKey.ValueStringPointer()
		}

		data = bunnycdn_api.HostnameToHostnameResourceModel(data.PullzoneId.ValueInt64(), remoteResource)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *HostnameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.HostnameResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.HostnameDelete(ctx, data.PullzoneId.ValueInt64(), bunnycdn_api.HostnameResourceModelToHostname(data))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete hostname, got error: %s", err))
		return
	}
}

func (r *HostnameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
