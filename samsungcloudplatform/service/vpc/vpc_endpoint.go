package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcVpcEndpointResource{}
	_ resource.ResourceWithConfigure = &vpcVpcEndpointResource{}
)

// NewVpcVpcEndpointResource is a helper function to simplify the provider implementation.
func NewVpcVpcEndpointResource() resource.Resource {
	return &vpcVpcEndpointResource{}
}

// vpcVpcEndpointResource is the data source implementation.
type vpcVpcEndpointResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcVpcEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_endpoint"
}

// Schema defines the schema for the data source.
func (r *vpcVpcEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "vpcendpoint",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Endpoint Name \n" +
					"  - example : vpcEndpointName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Required: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("ResourceType"): schema.StringAttribute{
				Description: "VPC Endpoint Resource Type \n" +
					"  - example : FS | OBS | SCR | DNS",
				Required: true,
			},
			common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
				Description: "VPC Endpoint Resource Key \n" +
					"  - example(case: SCR/DNS) : 07c5364702384471b650147321b52173 \n" +
					"  - example(case: FS/OBS) : 1.1.1.1",
				Required: true,
			},
			common.ToSnakeCase("ResourceInfo"): schema.StringAttribute{
				Description: "VPC Endpoint Resource Info \n" +
					"  - example(case: FS) : 192.168.0.1(SSD) \n" +
					"  - example(case: OBS) : https://xxx.samsungsdscloud.com \n" +
					"  - example(case: SCR) : xxx.samsungsdscloud.com(Auth) \n" +
					"  - example(case: DNS) : Private DNS Name",
				Required: true,
			},
			common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
				Description: "Endpoint IP Address \n" +
					"  - example : 10.10.10.10",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : VPC Endpoint description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("VpcEndpoint"): schema.SingleNestedAttribute{
				Description: "VpcEndpoint",
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
					common.ToSnakeCase("EndpointIpAddress"): schema.StringAttribute{
						Description: "EndpointIpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceType"): schema.StringAttribute{
						Description: "ResourceType",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceKey"): schema.StringAttribute{
						Description: "ResourceKey",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceInfo"): schema.StringAttribute{
						Description: "ResourceInfo",
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
func (r *vpcVpcEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcVpcEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.VpcEndpointResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpcEndpoint(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc endpoint",
			"Could not create vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(vpcendpoint.Id)

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		// #94: map the actual resource_key the API returns (the 1.1 show model
		// echoes a distinct resource_key field), NOT account_id. Using AccountId
		// here corrupted state (resource_key != account_id) and caused a spurious
		// diff/replace on every re-plan.
		ResourceKey:  types.StringValue(vpcendpoint.ResourceKey),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	plan.VpcEndpoint = vpcEndpointObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForVpcEndpointStatus(ctx, r.client, vpcendpoint.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc endpoint",
			"Error waiting for vpc endpoint to become active: "+err.Error(),
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
func (r *vpcVpcEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.VpcEndpointResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc endpoint",
			"Could not read vpc endpoint ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		// #94: map the real resource_key (not account_id) so refresh/update do not
		// corrupt state or force a spurious diff.
		ResourceKey:  types.StringValue(vpcendpoint.ResourceKey),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	state.VpcEndpoint = vpcEndpointObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcVpcEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.VpcEndpointResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpcEndpoint(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc endpoint",
			"Could not update vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpcEndpoint as UpdateVpc items are not populated.
	data, err := r.client.GetVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc endpoint",
			"Could not read vpc endpoint ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcendpoint := data.VpcEndpoint

	vpcEndpointModel := vpc.VpcEndpoint{
		Id:           types.StringValue(vpcendpoint.Id),
		Name:         types.StringValue(vpcendpoint.Name),
		VpcId:        types.StringValue(vpcendpoint.VpcId),
		VpcName:      types.StringValue(vpcendpoint.VpcName),
		SubnetId:     types.StringValue(vpcendpoint.SubnetId),
		SubnetName:   types.StringValue(vpcendpoint.SubnetName),
		ResourceType: types.StringValue(string(vpcendpoint.ResourceType)),
		// #94: map the real resource_key (not account_id) so refresh/update do not
		// corrupt state or force a spurious diff.
		ResourceKey:  types.StringValue(vpcendpoint.ResourceKey),
		ResourceInfo: types.StringValue(vpcendpoint.ResourceInfo),
		AccountId:    types.StringValue(vpcendpoint.AccountId),
		State:        types.StringValue(string(vpcendpoint.State)),
		Description:  types.StringPointerValue(vpcendpoint.Description.Get()),
		CreatedAt:    types.StringValue(vpcendpoint.CreatedAt.Format(time.RFC3339)),
		CreatedBy:    types.StringValue(vpcendpoint.CreatedBy),
		ModifiedAt:   types.StringValue(vpcendpoint.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:   types.StringValue(vpcendpoint.ModifiedBy),
	}
	vpcEndpointObjectValue, diags := types.ObjectValueFrom(ctx, vpcEndpointModel.AttributeTypes(), vpcEndpointModel)
	state.VpcEndpoint = vpcEndpointObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcVpcEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.VpcEndpointResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing VpcEndpoint
	err := r.client.DeleteVpcEndpoint(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc endpoint",
			"Could not delete vpc endpoint, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcEndpointStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc endpoint",
			"Error waiting for vpc endpoint to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForVpcEndpointStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetVpcEndpoint(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpcEndpoint.State), nil
	})
}
