package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1"
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
	_ resource.Resource              = &vpcPeeringResource{}
	_ resource.ResourceWithConfigure = &vpcPeeringResource{}
)

// NewVpcVpcPeeringResource is a helper function to simplify the provider implementation.
func NewVpcPeeringResource() resource.Resource {
	return &vpcPeeringResource{}
}

// vpcVpcPeeringResource is the data source implementation.
type vpcPeeringResource struct {
	config  *scpsdk.Configuration
	client  *vpcv1.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcPeeringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc_peering"
}

// Schema defines the schema for the data source.
func (r *vpcPeeringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Vpc peering",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
				Description: "Approver VPC Account ID",
				Required:    true,
			},
			common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
				Description: "Approver VPC ID",
				Required:    true,
			},
			common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
				Description: "Approver VPC Name. Computed: the create API derives it " +
					"from approver_vpc_id and returns it in the response (it is not a " +
					"create request field).",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
				Description: "Requester VPC ID",
				Required:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Peering Name\n" +
					"  - Minimum length: 3\n" +
					"  - Maximum length: 20\n" +
					"  - Pattern: ^[a-zA-Z0-9-]*$",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			common.ToSnakeCase("VpcPeering"): schema.SingleNestedAttribute{
				Description: "VpcPeering",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountType"): schema.StringAttribute{
						Description: "Account Type\n" +
							"  - Enum: SAME | DIFFERENT",
						Computed: true,
					},
					common.ToSnakeCase("ApproverVpcAccountId"): schema.StringAttribute{
						Description: "Approver VPC Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("ApproverVpcId"): schema.StringAttribute{
						Description: "Approver VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("ApproverVpcName"): schema.StringAttribute{
						Description: "Approver VPC Name",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "VPC Peering Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "VPC Peering ID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "VPC Peering Name\n" +
							"  - Minimum length: 3\n" +
							"  - Maximum length: 20\n" +
							"  - Pattern: ^[a-zA-Z0-9-]*$",
						Computed: true,
					},
					common.ToSnakeCase("RequesterVpcAccountId"): schema.StringAttribute{
						Description: "Requester VPC Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("RequesterVpcId"): schema.StringAttribute{
						Description: "Requester VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("RequesterVpcName"): schema.StringAttribute{
						Description: "Requester VPC Name",
						Computed:    true,
					},
					common.ToSnakeCase("DeleteRequesterAccountId"): schema.StringAttribute{
						Description: "Requester VPC Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State\n" +
							"  - Enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING | CREATING_REQUESTING | REJECTED | CANCELED | DELETING_REQUESTING",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcPeeringResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcPeeringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpcv1.VpcPeeringResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Issue #61: the documented VpcPeeringCreateRequest schema (api_docs.json,
	// /v1/vpc-peerings, supported version vpc 1.2) has exactly six fields and
	// approver_vpc_name is NOT one of them. The required server-side fields are
	// approver_vpc_account_id, approver_vpc_id, name and requester_vpc_id;
	// description and tags are optional. approver_vpc_name appears ONLY in the
	// VpcPeeringShowResponse. The api-test-automation suite gets a 202 sending
	// exactly {requester_vpc_id, approver_vpc_id, approver_vpc_account_id, name,
	// description, tags:[]} with NO approver_vpc_name.
	//
	// The previous auto-resolve-from-approver_vpc_id logic forced an
	// approver_vpc_name into the request body. A synthetic/derived name that did
	// not match the real approver VPC produced the misleading 400
	// "no value given for required property approver_vpc_name". The correct fix
	// is to NOT send approver_vpc_name at all and let the API derive it from
	// approver_vpc_id (the client wrapper already omits the field when it is
	// unset). We therefore clear any plan-supplied approver_vpc_name before the
	// create call; the real value is read back from the create response below.
	plan.ApproverVpcName = types.StringNull()

	// Create new vpc
	data, err := r.client.CreateVpcPeering(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc peering",
			"Could not create vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	plan.Id = types.StringValue(vpcPeering.Id)
	plan.ApproverVpcName = types.StringValue(vpcPeering.ApproverVpcName)
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)
	plan.VpcPeering = vpcObjectValue

	diags = resp.State.Set(ctx, plan)

	err = waitForVpcPeeringStatus(ctx, r.client, vpcPeering.Id, []string{}, []string{"ACTIVE", "CREATING_REQUESTING"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc peering",
			"Error waiting for vpc peering to become active: "+err.Error(),
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
	//diags = resp.State.Set(ctx, plan)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
}

// Read refreshes the Terraform state with the latest data.
func (r *vpcPeeringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcv1.VpcPeeringResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc peering",
			"Could not read vpc peering ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)

	state.ApproverVpcName = types.StringValue(vpcPeering.ApproverVpcName)
	state.VpcPeering = vpcObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPeeringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpcv1.VpcPeeringResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpcPeering(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc peering",
			"Could not update vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpcPeering as UpdateVpc items are not populated.
	data, err := r.client.GetVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc peering",
			"Could not read vpc peering ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcPeering := data.VpcPeering

	vpcPeeringModel := vpcv1.VpcPeering{
		Id:                       types.StringValue(vpcPeering.Id),
		Name:                     types.StringValue(vpcPeering.Name),
		AccountType:              types.StringValue(string(vpcPeering.AccountType)),
		ApproverVpcAccountId:     types.StringValue(vpcPeering.ApproverVpcAccountId),
		ApproverVpcId:            types.StringValue(vpcPeering.ApproverVpcId),
		ApproverVpcName:          types.StringValue(vpcPeering.ApproverVpcName),
		Description:              types.StringPointerValue(vpcPeering.Description.Get()),
		RequesterVpcAccountId:    types.StringValue(vpcPeering.RequesterVpcAccountId),
		RequesterVpcId:           types.StringValue(vpcPeering.RequesterVpcId),
		RequesterVpcName:         types.StringValue(vpcPeering.RequesterVpcName),
		DeleteRequesterAccountId: stringFromNullable(vpcPeering.DeleteRequesterAccountId.Get()),
		CreatedAt:                types.StringValue(vpcPeering.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(vpcPeering.CreatedBy),
		ModifiedAt:               types.StringValue(vpcPeering.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(vpcPeering.ModifiedBy),
		State:                    types.StringValue(string(vpcPeering.State)),
	}
	vpcObjectValue, diags := types.ObjectValueFrom(ctx, vpcPeeringModel.AttributeTypes(), vpcPeeringModel)

	state.ApproverVpcName = types.StringValue(vpcPeering.ApproverVpcName)
	state.VpcPeering = vpcObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPeeringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1.VpcPeeringResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing VpcPeering
	err := r.client.DeleteVpcPeering(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc peering",
			"Could not delete vpc peering, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcPeeringDeleted(ctx, r.client, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpc peering",
			"Error waiting for vpc peering to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForVpcPeeringStatus(ctx context.Context, vpcClient *vpcv1.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetVpcPeering(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpcPeering.State), nil
	})
}

// waitForVpcPeeringDeleted waits until the peering is gone. A real 404 (the
// peering no longer exists) is treated as the terminal DELETED state rather than
// a hard error, and DELETED/DELETING_REQUESTING are accepted as terminal too.
func waitForVpcPeeringDeleted(ctx context.Context, vpcClient *vpcv1.Client, id string) error {
	return client.WaitForStatus(ctx, nil, []string{"ACTIVE", "DELETING", "EDITING"}, []string{"DELETED", "DELETING_REQUESTING"}, func() (interface{}, string, error) {
		info, statusCode, err := vpcClient.GetVpcPeeringWithStatus(ctx, id)
		if err != nil {
			if statusCode == 404 {
				return &scpvpcShowSentinel{}, "DELETED", nil
			}
			return nil, "", err
		}
		return info, string(info.VpcPeering.State), nil
	})
}

// scpvpcShowSentinel is a non-nil placeholder so the state-change machinery does
// not interpret a successful 404 (resource gone) as an error.
type scpvpcShowSentinel struct{}

func stringFromNullable(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}
	return types.StringValue(*value)
}
