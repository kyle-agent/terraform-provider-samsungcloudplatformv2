package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpcTgwRuleResource{}
	_ resource.ResourceWithConfigure = &vpcTgwRuleResource{}
)

// NewVpcTgwRuleResource is a helper function to simplify the provider implementation.
func NewVpcTgwRuleResource() resource.Resource {
	return &vpcTgwRuleResource{}
}

type vpcTgwRuleResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (v vpcTgwRuleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_transit_gateway_rule"
}

// Schema defines the schema for the data source.
func (v *vpcTgwRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit gateway rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "Transit Gateway Id ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Routing Rule description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
				Description: "Destination CIDR \n" +
					"  - example : 10.10.10.0/24",
				Required: true,
			},
			common.ToSnakeCase("DestinationType"): schema.StringAttribute{
				Description: "Destination Type \n" +
					"  - example : VPC | TGW",
				Required: true,
			},
			common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
				Description: "Tgw Connection Vpc ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("RoutingRule"): schema.SingleNestedAttribute{
				Description: "Routing rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationCidr"): schema.StringAttribute{
						Description: "DestinationCidr",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationResourceId"): schema.StringAttribute{
						Description: "DestinationResourceId",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationResourceName"): schema.StringAttribute{
						Description: "DestinationResourceName",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationType"): schema.StringAttribute{
						Description: "DestinationType",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
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
					common.ToSnakeCase("SourceResourceId"): schema.StringAttribute{
						Description: "SourceResourceId",
						Computed:    true,
					},
					common.ToSnakeCase("SourceResourceName"): schema.StringAttribute{
						Description: "SourceResourceName",
						Computed:    true,
					},
					common.ToSnakeCase("SourceType"): schema.StringAttribute{
						Description: "SourceType",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("TgwConnectionVpcId"): schema.StringAttribute{
						Description: "TgwConnectionVpcId",
						Computed:    true,
					},
					common.ToSnakeCase("TgwConnectionVpcName"): schema.StringAttribute{
						Description: "TgwConnectionVpcName",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcTgwRuleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = inst.Client.Vpc
	r.clients = inst.Client
}

// Create the resource and sets the initial Terraform state.
func (r *vpcTgwRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan vpc.RoutingRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new routing rule
	data, err := r.client.CreateTgwRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Could not create routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRule := data.TransitGatewayRule
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(routingRule.Id)
	diags = resp.State.Set(ctx, plan)

	routingRuleModel := createRoutingRuleModel(&routingRule)

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	plan.RoutingRule = routingRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	// Non-empty Pending lets StateChangeConf short-circuit on a parked/ERROR state
	// instead of polling for the full timeout (issue #76).
	err = waitForRoutingRuleStatus(ctx, r.client, plan.TransitGatewayId.ValueString(), data.TransitGatewayRule.Id,
		[]string{common.CreatingState},
		[]string{common.ActiveState})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become active: "+err.Error(),
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
func (r *vpcTgwRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from routing rule
	data, err := r.client.GetRoutingRule(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading transit gateway routing rule",
			"Could not read routing rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	routingRuleModel := createRoutingRuleModel(&data.TransitGatewayRules[0])

	routingRuleObjectValue, diags := types.ObjectValueFrom(ctx, routingRuleModel.AttributeTypes(), routingRuleModel)
	state.RoutingRule = routingRuleObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (v vpcTgwRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r vpcTgwRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.RoutingRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing routing rule
	err := r.client.DeleteRoutingRule(ctx, state.TransitGatewayId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting transit gateway routing rule",
			"Could not delete transit gateway routing rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForRoutingRuleStatus(ctx, r.client, state.TransitGatewayId.ValueString(), state.Id.ValueString(),
		[]string{common.DeletingState, common.ActiveState},
		[]string{common.DeletedState})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become deleted: "+err.Error(),
		)
		return
	}
}

func createRoutingRuleModel(data *scpvpc.TransitGatewayRule) vpc.RoutingRule {
	return vpc.RoutingRule{
		AccountId:               types.StringValue(data.AccountId),
		CreatedAt:               types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:               types.StringValue(data.CreatedBy),
		Description:             types.StringValue(data.Description),
		DestinationCidr:         types.StringValue(data.DestinationCidr),
		DestinationResourceId:   types.StringPointerValue(data.DestinationResourceId.Get()),
		DestinationResourceName: types.StringPointerValue(data.DestinationResourceName.Get()),
		DestinationType:         types.StringValue(string(data.DestinationType)),
		Id:                      types.StringValue(data.Id),
		ModifiedAt:              types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:              types.StringValue(data.ModifiedBy),
		SourceResourceId:        types.StringPointerValue(data.SourceResourceId.Get()),
		SourceResourceName:      types.StringPointerValue(data.SourceResourceName.Get()),
		SourceType:              types.StringValue(string(data.SourceType)),
		State:                   types.StringValue(string(data.State)),
		TgwConnectionVpcId:      types.StringPointerValue(data.TgwConnectionVpcId.Get()),
		TgwConnectionVpcName:    types.StringPointerValue(data.TgwConnectionVpcName.Get()),
	}
}

func waitForRoutingRuleStatus(ctx context.Context, routingRuleClient *vpc.Client, transitGatewayId string, routingRuleId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := routingRuleClient.GetRoutingRule(ctx, transitGatewayId, routingRuleId)
		if err != nil {
			return nil, "", err
		}
		if len(info.TransitGatewayRules) == 0 {
			return info, "DELETED", nil
		}
		return info, string(info.TransitGatewayRules[0].State), nil
	})
}
