package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/securitygroup" // securitygroup client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &securityGroupRuleResource{}
	_ resource.ResourceWithConfigure = &securityGroupRuleResource{}
)

// NewSecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupRuleResource() resource.Resource {
	return &securityGroupRuleResource{}
}

// securityGroupResource is the data source implementation.
type securityGroupRuleResource struct {
	config  *scpsdk.Configuration
	client  *securitygroup.Client
	clients *client.SCPClient
}

func (r *securityGroupRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

// Metadata returns the data source type name.
func (r *securityGroupRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_security_group_rule"
}

// Schema defines the schema for the data source.
func (r *securityGroupRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Security group rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
				Description: "SecurityGroupId \n" +
					"  - example : cff990e6d5ed43d3ab239e4aba0b4c3e",
				Required: true,
			},
			common.ToSnakeCase("ethertype"): schema.StringAttribute{
				Description: "ethertype \n" +
					"  - example : IPV4",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4"),
				},
			},
			common.ToSnakeCase("protocol"): schema.StringAttribute{
				Description: "protocol \n" +
					"  - example : TCP",
				Optional: true,
			},
			common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
				Description: "portRangeMin \n" +
					"  - example : 22",
				Optional: true,
			},
			common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
				Description: "portRangeMax \n" +
					"  - example : 22",
				Optional: true,
			},
			common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
				Description: "RemoteIpPrefix \n" +
					"  - example : 1.1.1.1/32",
				Optional: true,
			},
			common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
				Description: "RemoteGroupId \n" +
					"  - example : 8a8048af06b048329867e57284347066",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : securityGroupRuleDescription",
				Optional: true,
			},
			common.ToSnakeCase("Direction"): schema.StringAttribute{
				Description: "Direction \n" +
					"  - example : ingress",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ingress", "egress"),
				},
			},
			common.ToSnakeCase("SecurityGroupRule"): schema.SingleNestedAttribute{
				Description: "Security group rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("SecurityGroupId"): schema.StringAttribute{
						Description: "SecurityGroupId",
						Computed:    true,
					},
					common.ToSnakeCase("ethertype"): schema.StringAttribute{
						Description: "ethertype",
						Computed:    true,
					},
					common.ToSnakeCase("protocol"): schema.StringAttribute{
						Description: "protocol",
						Computed:    true,
					},
					common.ToSnakeCase("portRangeMin"): schema.Int32Attribute{
						Description: "portRangeMin",
						Computed:    true,
					},
					common.ToSnakeCase("portRangeMax"): schema.Int32Attribute{
						Description: "portRangeMax",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteIpPrefix"): schema.StringAttribute{
						Description: "RemoteIpPrefix",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteGroupId"): schema.StringAttribute{
						Description: "RemoteGroupId",
						Computed:    true,
					},
					common.ToSnakeCase("RemoteGroupName"): schema.StringAttribute{
						Description: "RemoteGroupName",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "Direction",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *securityGroupRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securitygroup.SecurityGroupRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new internet gateway
	data, err := r.client.CreateSecurityGroupRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating security group rule",
			"Could not create security group rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroupRule := data.SecurityGroupRule
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(data.SecurityGroupRule.Id)

	sgrModel := securitygroup.SecurityGroupRule{
		Id:              types.StringValue(securityGroupRule.Id),
		SecurityGroupId: types.StringValue(securityGroupRule.SecurityGroupId),
		Ethertype:       types.StringPointerValue(securityGroupRule.Ethertype.Get()),
		Protocol:        types.StringPointerValue(securityGroupRule.Protocol.Get()),
		PortRangeMin:    types.Int32PointerValue(securityGroupRule.PortRangeMin.Get()),
		PortRangeMax:    types.Int32PointerValue(securityGroupRule.PortRangeMax.Get()),
		RemoteIpPrefix:  types.StringPointerValue(securityGroupRule.RemoteIpPrefix.Get()),
		RemoteGroupId:   types.StringPointerValue(securityGroupRule.RemoteGroupId.Get()),
		RemoteGroupName: types.StringPointerValue(securityGroupRule.RemoteGroupName.Get()),
		Description:     types.StringPointerValue(securityGroupRule.Description.Get()),
		Direction:       types.StringValue(string(securityGroupRule.Direction)),
		CreatedAt:       types.StringValue(securityGroupRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(securityGroupRule.CreatedBy),
		ModifiedAt:      types.StringValue(securityGroupRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:      types.StringValue(securityGroupRule.ModifiedBy),
	}

	sgrObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	plan.SecurityGroupRule = sgrObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securitygroup.SecurityGroupRuleResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from security group rule
	data, err := r.client.GetSecurityGroupRule(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading security group rule",
			"Could not read security group rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	securityGroupRule := data.SecurityGroupRule

	sgrModel := securitygroup.SecurityGroupRule{
		Id:              types.StringValue(securityGroupRule.Id),
		SecurityGroupId: types.StringValue(securityGroupRule.SecurityGroupId),
		Ethertype:       types.StringPointerValue(securityGroupRule.Ethertype.Get()),
		Protocol:        types.StringPointerValue(securityGroupRule.Protocol.Get()),
		PortRangeMin:    types.Int32PointerValue(securityGroupRule.PortRangeMin.Get()),
		PortRangeMax:    types.Int32PointerValue(securityGroupRule.PortRangeMax.Get()),
		RemoteIpPrefix:  types.StringPointerValue(securityGroupRule.RemoteIpPrefix.Get()),
		RemoteGroupId:   types.StringPointerValue(securityGroupRule.RemoteGroupId.Get()),
		RemoteGroupName: types.StringPointerValue(securityGroupRule.RemoteGroupName.Get()),
		Description:     types.StringPointerValue(securityGroupRule.Description.Get()),
		Direction:       types.StringValue(string(securityGroupRule.Direction)),
		CreatedAt:       types.StringValue(securityGroupRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:       types.StringValue(securityGroupRule.CreatedBy),
		ModifiedAt:      types.StringValue(securityGroupRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:      types.StringValue(securityGroupRule.ModifiedBy),
	}
	sgrObjectValue, diags := types.ObjectValueFrom(ctx, sgrModel.AttributeTypes(), sgrModel)
	state.SecurityGroupRule = sgrObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securitygroup.SecurityGroupRuleResource
	diags := req.State.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing security group rule
	err := r.client.DeleteSecurityGroupRule(ctx, state.Id.ValueString()) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting security group rule",
			"Could not delete security group rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
