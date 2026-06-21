package servicewatch

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serviceWatchLogGroupResource{}
	_ resource.ResourceWithConfigure   = &serviceWatchLogGroupResource{}
	_ resource.ResourceWithImportState = &serviceWatchLogGroupResource{}
)

// NewServiceWatchLogGroupResource is a helper function to simplify the provider implementation.
func NewServiceWatchLogGroupResource() resource.Resource {
	return &serviceWatchLogGroupResource{}
}

// serviceWatchLogGroupResource is the data source implementation.
type serviceWatchLogGroupResource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *serviceWatchLogGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_log_group"
}

// Schema defines the schema for the resource.
func (r *serviceWatchLogGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Group Resource",
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Log group ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Log group name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("RetentionPeriod"): schema.Int32Attribute{
				Description: "Log group retention period",
				Required:    true,
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("LogGroup"): schema.SingleNestedAttribute{
				Description: "Log group",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Log group ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Log group name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("RetentionPeriod"): schema.Int32Attribute{
						Description: "Log group retention period",
						Computed:    true,
					},
					common.ToSnakeCase("RetentionPeriodName"): schema.StringAttribute{
						Description: "Log group retention period name",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Log group status\n" +
							"Allowed values: ACTIVE, DELETING, DELETED",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created date time",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Creator ID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified date time",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modifier ID",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *serviceWatchLogGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.ServiceWatch
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *serviceWatchLogGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicewatch.LogGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new LogGroup
	data, err := r.client.CreateLogGroup(ctx, plan) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating LogGroup",
			"Could not create LogGroup, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	logGroup := convertLogGroup(&data.LogGroup)
	logGroupObjectValue, diags := types.ObjectValueFrom(ctx, logGroup.AttributeTypes(), logGroup)

	plan.Id = types.StringValue(logGroup.Id.ValueString())
	plan.LogGroup = logGroupObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceWatchLogGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicewatch.LogGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Resource Group
	data, err := r.client.GetLogGroup(ctx, state.Id.ValueString())
	if err != nil && data == nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Log Group",
			"Could not read Log Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	// 존재하지 않는 log group 조회 시 null 로 return
	if data.LogGroup.GetId() == "" {
		state.LogGroup = types.ObjectNull(state.LogGroup.AttributeTypes(ctx))
		resp.State.Set(ctx, &state)
		return
	}

	// Map response body to schema and populate Computed attribute values
	logGroup := convertLogGroup(&data.LogGroup)
	logGroupObjectValue, diags := types.ObjectValueFrom(ctx, logGroup.AttributeTypes(), logGroup)
	state.LogGroup = logGroupObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceWatchLogGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var state servicewatch.LogGroupResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Resource Group
	_, err := r.client.UpdateLogGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Log group",
			"Could not update dashboard, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetLogGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading resourceGroup",
			"Could not read resourceGroup ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	logGroup := convertLogGroup(&data.LogGroup)
	logGroupObjectValue, diags := types.ObjectValueFrom(ctx, logGroup.AttributeTypes(), logGroup)
	state.LogGroup = logGroupObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceWatchLogGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicewatch.LogGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	logGroupId := state.Id.ValueString()

	// Delete existing Resource Group
	_, err := r.client.DeleteLogGroup(ctx, []string{logGroupId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Log Group",
			"Could not delete Log Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertLogGroup(logGroupResp *servicewatch2.LogGroupDTO) servicewatch.LogGroup {
	return servicewatch.LogGroup{
		Id:                  types.StringValue(logGroupResp.Id),
		Name:                types.StringValue(logGroupResp.Name),
		AccountId:           types.StringValue(logGroupResp.AccountId),
		RetentionPeriod:     types.Int32Value(logGroupResp.RetentionPeriod),
		RetentionPeriodName: types.StringValue(logGroupResp.RetentionPeriodName),
		Status:              types.StringValue(string(logGroupResp.Status)),
		CreatedAt:           types.StringValue(logGroupResp.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(logGroupResp.CreatedBy),
		ModifiedAt:          types.StringValue(logGroupResp.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(logGroupResp.ModifiedBy),
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *serviceWatchLogGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
