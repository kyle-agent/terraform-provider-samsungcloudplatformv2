package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/securitygroup" // securitygroup client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
)

// NewSecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

// securityGroupResource is the data source implementation.
type securityGroupResource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *securityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group"
}

// Schema defines the schema for the data source.
func (r *securityGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Security group",
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
				Description: "Name \n" +
					"  - example : sg_0911",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : sg_description",
				Optional: true,
			},
			common.ToSnakeCase("Loggable"): schema.BoolAttribute{
				Description: "loggable \n" +
					"  - example : True",
				Optional: true,
			},
			common.ToSnakeCase("SecurityGroup"): schema.SingleNestedAttribute{
				Description: "Security group",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Loggable"): schema.BoolAttribute{
						Description: "loggable",
						Computed:    true,
					},
					common.ToSnakeCase("RuleCount"): schema.Int32Attribute{
						Description: "RuleCount",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
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
func (r *securityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.SecurityGroup
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securitygroup.SecurityGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new SecurityGroup
	data, err := r.client.CreateSecurityGroup(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating security group",
			"Could not create security group, unexpected error: "+err.Error(),
		)
		return
	}

	securityGroup := data.SecurityGroup

	plan.Id = types.StringValue(securityGroup.Id)

	sgModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		RuleCount:   types.Int32PointerValue(securityGroup.RuleCount),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, diags := types.ObjectValueFrom(ctx, sgModel.AttributeTypes(), sgModel)
	plan.SecurityGroup = sgObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securitygroup.SecurityGroupResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from vpc
	data, err := r.client.GetSecurityGroup(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroup := data.SecurityGroup
	sgrModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	state.SecurityGroup = sgObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state securitygroup.SecurityGroupResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.UpdateSecurityGroup(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetSecurityGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading security group",
			"Could not read security group ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	securityGroup := data.SecurityGroup
	sgrModel := securitygroup.SecurityGroup{
		Id:          types.StringValue(securityGroup.Id),
		AccountId:   types.StringValue(securityGroup.AccountId),
		Name:        types.StringValue(securityGroup.Name),
		Description: types.StringPointerValue(securityGroup.Description.Get()),
		Loggable:    types.BoolValue(securityGroup.Loggable),
		State:       types.StringValue(securityGroup.State),
		CreatedAt:   types.StringValue(securityGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(securityGroup.CreatedBy),
		ModifiedAt:  types.StringValue(securityGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(securityGroup.ModifiedBy),
	}

	sgObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	state.SecurityGroup = sgObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securitygroup.SecurityGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing network logging storage
	err := r.client.DeleteSecurityGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting security group",
			"Could not delete security group, unexpected error: "+err.Error(),
		)
		return
	}
}
