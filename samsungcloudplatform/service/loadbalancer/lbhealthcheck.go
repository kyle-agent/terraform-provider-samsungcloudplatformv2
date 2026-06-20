package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	loadbalancerutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/loadbalancer/1.3"
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
	_ resource.Resource                = &loadbalancerLbHealthCheckResource{}
	_ resource.ResourceWithConfigure   = &loadbalancerLbHealthCheckResource{}
	_ resource.ResourceWithImportState = &loadbalancerLbHealthCheckResource{}
)

// NewLoadBalancerLbHealthCheckResource is a helper function to simplify the provider implementation.
func NewLoadBalancerLbHealthCheckResource() resource.Resource {
	return &loadbalancerLbHealthCheckResource{}
}

// loadbalancerLbHealthCheckResource is the data source implementation.
type loadbalancerLbHealthCheckResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLbHealthCheckResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_lb_health_check"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLbHealthCheckResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lb Health Check.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LbHealthCheck"): schema.SingleNestedAttribute{
				Description: "A detail of Lb Health Check.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
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
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "HealthCheckPort",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "HealthCheckInterval",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "HealthCheckTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "HealthCheckCount",
						Optional:    true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "HttpMethod",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "HealthCheckUrl",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "ResponseCode",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckType"): schema.StringAttribute{
						Description: "HealthCheckType",
						Optional:    true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "RequestData",
						Optional:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("LbHealthCheckCreate"): schema.SingleNestedAttribute{
				Description: "Create Lb Health Check.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("Protocol"): schema.StringAttribute{
						Description: "Protocol",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckPort"): schema.Int32Attribute{
						Description: "HealthCheckPort",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
						Description: "HealthCheckInterval",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckTimeout"): schema.Int32Attribute{
						Description: "HealthCheckTimeout",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckCount"): schema.Int32Attribute{
						Description: "HealthCheckCount",
						Optional:    true,
					},
					common.ToSnakeCase("HttpMethod"): schema.StringAttribute{
						Description: "HttpMethod",
						Optional:    true,
					},
					common.ToSnakeCase("HealthCheckUrl"): schema.StringAttribute{
						Description: "HealthCheckUrl",
						Optional:    true,
					},
					common.ToSnakeCase("ResponseCode"): schema.StringAttribute{
						Description: "ResponseCode",
						Optional:    true,
					},
					common.ToSnakeCase("RequestData"): schema.StringAttribute{
						Description: "RequestData",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLbHealthCheckResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoadBalancer
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLbHealthCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LbHealthCheckResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Health Check
	data, err := r.client.CreateLbHealthCheck(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Could not create Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.LbHealthCheck.Id)

	// Map response body to schema and populate Computed attribute values
	lbHealthCheckModel := createLbHealthCheckModel(data)
	lbHealthCheckOjbectValue, diags := types.ObjectValueFrom(ctx, lbHealthCheckModel.AttributeTypes(), lbHealthCheckModel)
	plan.LbHealthCheck = lbHealthCheckOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLbHealthCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state loadbalancer.LbHealthCheckResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from LB Health Check
	data, err := r.client.GetLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Could not create Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	lbHealthCheckModel := createLbHealthCheckModel(data)

	lbHealthCheckObjectValue, diags := types.ObjectValueFrom(ctx, lbHealthCheckModel.AttributeTypes(), lbHealthCheckModel)
	state.LbHealthCheck = lbHealthCheckObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLbHealthCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state loadbalancer.LbHealthCheckResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateLbHealthCheck(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Could not create Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetFirewallRule as UpdateFirewallRule items are not populated.
	data, err := r.client.GetLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Lb Health Check",
			"Could not create Lb Health Check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	lbHealthCheckModel := createLbHealthCheckModel(data)

	lbHealthCheckObjectValue, diags := types.ObjectValueFrom(ctx, lbHealthCheckModel.AttributeTypes(), lbHealthCheckModel)
	state.LbHealthCheck = lbHealthCheckObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLbHealthCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LbHealthCheckResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing LB Health Check
	err := r.client.DeleteLbHealthCheck(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Health Check",
			"Could not delete lb health check, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func ToNullableInt32Value(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}

func createLbHealthCheckModel(data *scploadbalancer.LbHealthCheckShowResponse) loadbalancer.LbHealthCheckDetail {
	lbHealthCheck := data.LbHealthCheck

	return loadbalancer.LbHealthCheckDetail{
		Name:                types.StringValue(lbHealthCheck.Name),
		VpcId:               loadbalancerutil.ToNullableStringValue(lbHealthCheck.VpcId.Get()),
		SubnetId:            loadbalancerutil.ToNullableStringValue(lbHealthCheck.SubnetId.Get()),
		Protocol:            loadbalancerutil.ToNullableStringValue((*string)(lbHealthCheck.Protocol)),
		HealthCheckPort:     ToNullableInt32Value(lbHealthCheck.HealthCheckPort.Get()),
		HealthCheckInterval: types.Int32Value(*lbHealthCheck.HealthCheckInterval),
		HealthCheckTimeout:  types.Int32Value(*lbHealthCheck.HealthCheckTimeout),
		HealthCheckCount:    types.Int32Value(*lbHealthCheck.HealthCheckCount),
		HealthCheckUrl:      loadbalancerutil.ToNullableStringValue(lbHealthCheck.HealthCheckUrl.Get()),
		HttpMethod:          loadbalancerutil.ToNullableStringValue(lbHealthCheck.HttpMethod.Get()),
		ResponseCode:        loadbalancerutil.ToNullableStringValue(lbHealthCheck.ResponseCode.Get()),
		RequestData:         loadbalancerutil.ToNullableStringValue(lbHealthCheck.RequestData.Get()),
		HealthCheckType:     types.StringValue(string(lbHealthCheck.HealthCheckType)),
		State:               types.StringValue(lbHealthCheck.State),
		AccountId:           loadbalancerutil.ToNullableStringValue(lbHealthCheck.AccountId.Get()),
		Description:         loadbalancerutil.ToNullableStringValue(lbHealthCheck.Description.Get()),
		ModifiedBy:          types.StringValue(lbHealthCheck.ModifiedBy),
		ModifiedAt:          types.StringValue(lbHealthCheck.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(lbHealthCheck.CreatedBy),
		CreatedAt:           types.StringValue(lbHealthCheck.CreatedAt.Format(time.RFC3339)),
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *loadbalancerLbHealthCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
