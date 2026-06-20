package vpc

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcNatGatewayResource{}
	_ resource.ResourceWithConfigure   = &vpcNatGatewayResource{}
	_ resource.ResourceWithImportState = &vpcNatGatewayResource{}
)

// NewVpcNatGatewayResource is a helper function to simplify the provider implementation.
func NewVpcNatGatewayResource() resource.Resource {
	return &vpcNatGatewayResource{}
}

// vpcNatGatewayResource is the data source implementation.
type vpcNatGatewayResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcNatGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_nat_gateway"
}

// Schema defines the schema for the data source.
func (r *vpcNatGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "natgateway",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 607e0938521643b5b4b266f343fae693",
				Required: true,
			},
			common.ToSnakeCase("PublicipId"): schema.StringAttribute{
				Description: "Public IP ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : NAT Gateway description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("NatGateway"): schema.SingleNestedAttribute{
				Description: "NatGateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("NatGatewayIpAddress"): schema.StringAttribute{
						Description: "NatGatewayIpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Computed:    true,
					},
					common.ToSnakeCase("SubnetName"): schema.StringAttribute{
						Description: "SubnetName",
						Computed:    true,
					},
					common.ToSnakeCase("SubnetCidr"): schema.StringAttribute{
						Description: "SubnetCidr",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcNatGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcNatGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.NatGatewayResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateNatGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating nat gateway",
			"Could not create nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	natgateway := data.NatGateway
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(natgateway.Id)

	natGatewayModel := vpc.NatGateway{
		Id:                  types.StringValue(natgateway.Id),
		Name:                types.StringValue(natgateway.Name),
		NatGatewayIpAddress: types.StringValue(natgateway.NatGatewayIpAddress),
		VpcId:               types.StringValue(natgateway.VpcId),
		VpcName:             types.StringValue(natgateway.VpcName),
		SubnetId:            types.StringValue(natgateway.SubnetId),
		SubnetName:          types.StringValue(natgateway.SubnetName),
		SubnetCidr:          types.StringValue(natgateway.SubnetCidr),
		AccountId:           types.StringValue(natgateway.AccountId),
		State:               types.StringValue(natgateway.State),
		Description:         types.StringPointerValue(natgateway.Description.Get()),
		CreatedAt:           types.StringValue(natgateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(natgateway.CreatedBy),
		ModifiedAt:          types.StringValue(natgateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(natgateway.ModifiedBy),
	}
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	plan.NatGateway = natGatewayObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForNatGatewayStatus(ctx, r.client, natgateway.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating nat gateway",
			"Error waiting for nat gateway to become active: "+err.Error(),
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
// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *vpcNatGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpcNatGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.NatGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetNatGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading nat gateway",
			"Could not read nat gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	natgateway := data.NatGateway

	natGatewayModel := vpc.NatGateway{
		Id:                  types.StringValue(natgateway.Id),
		Name:                types.StringValue(natgateway.Name),
		NatGatewayIpAddress: types.StringValue(natgateway.NatGatewayIpAddress),
		VpcId:               types.StringValue(natgateway.VpcId),
		VpcName:             types.StringValue(natgateway.VpcName),
		SubnetId:            types.StringValue(natgateway.SubnetId),
		SubnetName:          types.StringValue(natgateway.SubnetName),
		SubnetCidr:          types.StringValue(natgateway.SubnetCidr),
		AccountId:           types.StringValue(natgateway.AccountId),
		State:               types.StringValue(natgateway.State),
		Description:         types.StringPointerValue(natgateway.Description.Get()),
		CreatedAt:           types.StringValue(natgateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(natgateway.CreatedBy),
		ModifiedAt:          types.StringValue(natgateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(natgateway.ModifiedBy),
	}
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	state.NatGateway = natGatewayObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcNatGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.NatGatewayResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateNatGateway(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating nat gateway",
			"Could not update nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetNatGateway as UpdateVpc items are not populated.
	data, err := r.client.GetNatGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading nat gateway",
			"Could not read nat gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	natgateway := data.NatGateway

	natGatewayModel := vpc.NatGateway{
		Id:                  types.StringValue(natgateway.Id),
		Name:                types.StringValue(natgateway.Name),
		NatGatewayIpAddress: types.StringValue(natgateway.NatGatewayIpAddress),
		VpcId:               types.StringValue(natgateway.VpcId),
		VpcName:             types.StringValue(natgateway.VpcName),
		SubnetId:            types.StringValue(natgateway.SubnetId),
		SubnetName:          types.StringValue(natgateway.SubnetName),
		SubnetCidr:          types.StringValue(natgateway.SubnetCidr),
		AccountId:           types.StringValue(natgateway.AccountId),
		State:               types.StringValue(natgateway.State),
		Description:         types.StringPointerValue(natgateway.Description.Get()),
		CreatedAt:           types.StringValue(natgateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(natgateway.CreatedBy),
		ModifiedAt:          types.StringValue(natgateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(natgateway.ModifiedBy),
	}
	natGatewayObjectValue, diags := types.ObjectValueFrom(ctx, natGatewayModel.AttributeTypes(), natGatewayModel)
	state.NatGateway = natGatewayObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcNatGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.NatGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing NatGateway
	err := r.client.DeleteNatGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting nat gateway",
			"Could not delete nat gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForNatGatewayStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting nat gateway",
			"Error waiting for nat gateway to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForNatGatewayStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetNatGateway(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, info.NatGateway.State, nil
	})
}
