package certificatemanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/certificatemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpcertificatemanager "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/certificatemanager/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &certificateManagerResource{}
	_ resource.ResourceWithConfigure   = &certificateManagerResource{}
	_ resource.ResourceWithImportState = &certificateManagerResource{}
)

func NewCertificateManagerResource() resource.Resource {
	return &certificateManagerResource{}
}

type certificateManagerResource struct {
	config *scpsdk.Configuration
	client *certificatemanager.Client
}

func (r *certificateManagerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

// Metadata returns the data source type name.
func (r *certificateManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_manager"
}

// Schema defines the schema for the data source.
func (r *certificateManagerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "certificate manager",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("CertBody"): schema.StringAttribute{
				Description: "Certificate body\n" +
					"  - Example: encoded certificate body data",
				Required: true,
			},
			common.ToSnakeCase("CertChain"): schema.StringAttribute{
				Description: "Certificate chain\n" +
					"  - Example: encoded certificate chain data",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Certificate Name\n" +
					"  - Example: test-certificate",
				Required: true,
			},
			common.ToSnakeCase("PrivateKey"): schema.StringAttribute{
				Description: "Private key\n" +
					"  - Example: encoded private key data",
				Required: true,
			},
			common.ToSnakeCase("Recipients"): schema.ListAttribute{
				Description: "Recipients\n" +
					"  - Example: [{\"region\":\"\",\"user_id\":\"sdaFDQSDADZ2488e195c0e97d9b9eb\",\"user_name\":\"kildong.hong\"}]",
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			common.ToSnakeCase("region"): schema.StringAttribute{
				Description: "Name of region\n" +
					"  - Example: west1",
				Required: true,
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description: "Timezone\n" +
					"  - Example: Asia/Seoul",
				Required: true,
			},
			common.ToSnakeCase("Certificate"): schema.SingleNestedAttribute{
				Description: "Certificate",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CertKind"): schema.StringAttribute{
						Description: "Certificate type\n" +
							"  - Example: PRD",
						Computed: true,
					},
					common.ToSnakeCase("Cn"): schema.StringAttribute{
						Description: "Certificate Common Name\n" +
							"  - Example: test.go.kr",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "ID\n" +
							"  - Example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Certificate Name\n" +
							"  - Example: test-certificate",
						Computed: true,
					},
					common.ToSnakeCase("NotAfterDt"): schema.StringAttribute{
						Description: "Certificate Expire Date\n" +
							"  - Example: 2026-02-07T18:07:59",
						Computed: true,
					},
					common.ToSnakeCase("NotBeforeDt"): schema.StringAttribute{
						Description: "Certificate Start Date\n" +
							"  - Example: 2025-02-08T18:07:00",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Certificate State\n" +
							"  - Example: VALID\n",
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *certificateManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.CertificateManager
}

func (r *certificateManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificatemanager.CertificateManagerResource
	fmt.Printf("-----------------------------------------------Start Create------------------------------------\n")

	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new certificate manager
	data, err := r.client.CreateCertificateManager(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating certificate manager",
			"Could not create certificate manager, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	plan.Id = types.StringValue(data.Certificate.Id)
	vgModel := certificatemanager.Certificate{
		Id:          types.StringValue(data.Certificate.Id),
		Name:        types.StringValue(data.Certificate.Name),
		CertKind:    types.StringValue(*data.Certificate.CertKind),
		Cn:          types.StringValue(data.Certificate.Cn),
		NotBeforeDt: types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		NotAfterDt:  types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		State:       types.StringValue(data.Certificate.State),
	}

	certificateObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	plan.Certificate = certificateObjectValue

	diags = resp.State.Set(ctx, plan)

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *certificateManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state certificatemanager.CertificateManagerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetCertificateManager(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading certificate manager",
			"Could not read certificate manager ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vgModel := createCertificateManagerModel(data)

	vgObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	state.Certificate = vgObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificateManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state certificatemanager.CertificateManagerResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificateManager(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting certificate manager",
			"Could not delete certificate manager, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForCertificateManagerStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting certificate manager",
			"Error waiting for certificate manager to become deleted: "+err.Error(),
		)
		return
	}
}

func createCertificateManagerModel(data *scpcertificatemanager.CertificateDetailResponse) certificatemanager.Certificate {
	return certificatemanager.Certificate{
		Id:          types.StringValue(data.Certificate.Id),
		Name:        types.StringValue(data.Certificate.Name),
		CertKind:    types.StringValue(*data.Certificate.CertKind),
		Cn:          types.StringValue(data.Certificate.Cn),
		NotBeforeDt: types.StringValue(data.Certificate.NotBeforeDt.Format(time.RFC3339)),
		NotAfterDt:  types.StringValue(data.Certificate.NotAfterDt.Format(time.RFC3339)),
		State:       types.StringValue(data.Certificate.State),
	}
}

func waitForCertificateManagerStatus(ctx context.Context, certificateManagerClient *certificatemanager.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := certificateManagerClient.GetCertificateManager(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Certificate.State), nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *certificateManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
