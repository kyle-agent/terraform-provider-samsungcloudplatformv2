package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcPrivateNatResource{}
	_ resource.ResourceWithConfigure   = &vpcPrivateNatResource{}
	_ resource.ResourceWithImportState = &vpcPrivateNatResource{}
)

// NewVpcPrivateNatResource is a helper function to simplify the provider implementation.
func NewVpcPrivateNatResource() resource.Resource {
	return &vpcPrivateNatResource{}
}

// vpcPrivateNatResource is the resource implementation.
type vpcPrivateNatResource struct {
	config    *scpsdk.Configuration
	client1d2 *vpcv1d2.Client
	clients   *client.SCPClient
}

// Metadata returns the data source type name.
func (d *vpcPrivateNatResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_private_nat"
}

// Schema defines the schema for the data source.
func (d *vpcPrivateNatResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Private NAT.",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "Private NAT IP range \n" +
					"  - example : 192.167.0.0/24",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : PrivateNat Description",
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Computed: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Private NAT Name \n" +
					"  - example : PrivateNatName",
				Required: true,
			},
			common.ToSnakeCase("ServiceResourceId"): schema.StringAttribute{
				Description: "Private NAT connected Service Resource ID \n" +
					"  - example : 3f342bf9a557405b997c2cf48c89cbc2",
				Required: true,
			},
			common.ToSnakeCase("ServiceType"): schema.StringAttribute{
				Description: "Private NAT connected Service Type \n" +
					"  - example : DIRECT_CONNECT",
				Required: true,
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),

			// Output
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Private NAT ID \n" +
					"  - example : 12f56e27070248a6a240a497e43fbe18",
				Computed: true,
			},
			common.ToSnakeCase("PrivateNat"): schema.SingleNestedAttribute{
				Description: "Private NAT details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID \n" +
							"  - example : f1e6c81a2b054582878cb9724dc2ce9f",
						Computed: true,
					},
					common.ToSnakeCase("Cidr"): schema.StringAttribute{
						Description: "Private NAT IP range \n" +
							"  - example : 192.167.0.0/24",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At \n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description \n" +
							"  - example : PrivateNat Description",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Private NAT ID \n" +
							"  - example : 12f56e27070248a6a240a497e43fbe18",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At \n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Private NAT Name \n" +
							"  - example : PrivateNatName",
						Computed: true,
					},
					common.ToSnakeCase("ServiceResourceId"): schema.StringAttribute{
						Description: "Private NAT connected Service Resource ID \n" +
							"  - example : 3f342bf9a557405b997c2cf48c89cbc2",
						Computed: true,
					},
					common.ToSnakeCase("ServiceResourceName"): schema.StringAttribute{
						Description: "Private NAT connected Service Resource Name \n" +
							"  - example : PrivateNatName",
						Computed: true,
					},
					common.ToSnakeCase("ServiceType"): schema.StringAttribute{
						Description: "Private NAT connected Service Type \n" +
							"  - example : DIRECT_CONNECT",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "Private NAT State \n" +
							"  - example : ACTIVE",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *vpcPrivateNatResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client1d2 = inst.Client.VpcV1Dot2
	d.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPrivateNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpcv1d2.PrivateNatResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client1d2.CreatePrivateNat(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Private NAT",
			"Could not create Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(privateNat.Id)

	privateNatModel := vpcv1d2.PrivateNat{
		Id:                  types.StringValue(privateNat.Id),
		Name:                types.StringValue(privateNat.Name),
		AccountId:           types.StringValue(privateNat.AccountId),
		Cidr:                types.StringValue(privateNat.Cidr),
		State:               types.StringValue(string(privateNat.State)),
		Description:         types.StringPointerValue(privateNat.Description.Get()),
		ServiceResourceId:   types.StringValue(privateNat.ServiceResourceId),
		ServiceResourceName: types.StringValue(privateNat.ServiceResourceName),
		ServiceType:         types.StringValue(string(privateNat.ServiceType)),
		CreatedAt:           types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(privateNat.CreatedBy),
		ModifiedAt:          types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, diags := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	plan.PrivateNat = privateNatObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForPrivateNatStatus(ctx, r.client1d2, privateNat.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Private NAT",
			"Error waiting for Private NAT to become active: "+err.Error(),
		)
		return
	}

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
func (r *vpcPrivateNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1d2.PrivateNatResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Private NAT
	data, err := r.client1d2.GetPrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Private NAT",
			"Could not read Private NAT ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat

	privateNatModel := vpcv1d2.PrivateNat{
		Id:                  types.StringValue(privateNat.Id),
		Name:                types.StringValue(privateNat.Name),
		AccountId:           types.StringValue(privateNat.AccountId),
		Cidr:                types.StringValue(privateNat.Cidr),
		State:               types.StringValue(string(privateNat.State)),
		Description:         types.StringPointerValue(privateNat.Description.Get()),
		ServiceResourceId:   types.StringValue(privateNat.ServiceResourceId),
		ServiceResourceName: types.StringValue(privateNat.ServiceResourceName),
		ServiceType:         types.StringValue(string(privateNat.ServiceType)),
		CreatedAt:           types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(privateNat.CreatedBy),
		ModifiedAt:          types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, _ := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	state.PrivateNat = privateNatObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPrivateNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan vpcv1d2.PrivateNatResource  // Changed Data
	var state vpcv1d2.PrivateNatResource // Stored data
	req.Plan.Get(ctx, &plan)
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Private NAT
	_, err := r.client1d2.UpdatePrivateNat(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Private NAT",
			"Could not update Private NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetPrivateNat as UpdatePrivateNat items are not populated.
	data, err := r.client1d2.GetPrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Private NAT",
			"Could not read Private NAT ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateNat := data.PrivateNat
	plan.Id = types.StringValue(privateNat.Id)

	privateNatModel := vpcv1d2.PrivateNat{
		Id:                  types.StringValue(privateNat.Id),
		Name:                types.StringValue(privateNat.Name),
		AccountId:           types.StringValue(privateNat.AccountId),
		Cidr:                types.StringValue(privateNat.Cidr),
		State:               types.StringValue(string(privateNat.State)),
		Description:         types.StringPointerValue(privateNat.Description.Get()),
		ServiceResourceId:   types.StringValue(privateNat.ServiceResourceId),
		ServiceResourceName: types.StringValue(privateNat.ServiceResourceName),
		ServiceType:         types.StringValue(string(privateNat.ServiceType)),
		CreatedAt:           types.StringValue(privateNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(privateNat.CreatedBy),
		ModifiedAt:          types.StringValue(privateNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(privateNat.ModifiedBy),
	}
	privateNatObjectValue, _ := types.ObjectValueFrom(ctx, privateNatModel.AttributeTypes(), privateNatModel)
	plan.PrivateNat = privateNatObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPrivateNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1d2.PrivateNatResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Private NAT
	err := r.client1d2.DeletePrivateNat(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Private NAT",
			"Could not delete Private NAT unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForPrivateNatStatus(ctx, r.client1d2, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting Private NAT",
			"Error waiting for Private NAT to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForPrivateNatStatus(ctx context.Context, vpcClient *vpcv1d2.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetPrivateNat(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.PrivateNat.State), nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *vpcPrivateNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
