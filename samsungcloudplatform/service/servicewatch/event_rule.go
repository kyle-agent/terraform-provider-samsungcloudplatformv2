package servicewatch

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &serviceWatchEventRuleResource{}
	_ resource.ResourceWithConfigure   = &serviceWatchEventRuleResource{}
	_ resource.ResourceWithImportState = &serviceWatchEventRuleResource{}
)

// NewServiceWatchEventRuleResource is a helper function to simplify the provider implementation.
func NewServiceWatchEventRuleResource() resource.Resource {
	return &serviceWatchEventRuleResource{}
}

// serviceWatchEventRuleResource is the data source implementation.
type serviceWatchEventRuleResource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *serviceWatchEventRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_event_rule"
}

func (r *serviceWatchEventRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Event Rule Resource",
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Event Rule ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Optional:    true,
				Description: "Event rule description",
			},
			common.ToSnakeCase("EventIds"): schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of Event IDs",
			},
			common.ToSnakeCase("EventRuleId"): schema.StringAttribute{
				Optional:    true,
				Description: "Event rule ID",
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Optional:    true,
				Description: "Event rule name",
			},
			common.ToSnakeCase("RecipientIds"): schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Notification recipient IDs",
			},
			common.ToSnakeCase("ResourceTypeId"): schema.StringAttribute{
				Optional:    true,
				Description: "Resource type ID",
			},
			common.ToSnakeCase("ServiceId"): schema.StringAttribute{
				Required:    true,
				Description: "Service ID",
			},
			common.ToSnakeCase("SrnList"): schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of SDS cloud Resource Names",
			},
			common.ToSnakeCase("ActiveYn"): schema.StringAttribute{
				Optional:    true,
				Description: "ActiveYn",
				Validators: []validator.String{
					stringvalidator.OneOf("Y", "N"),
				},
			},
			common.ToSnakeCase("NoneAttributes"): schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of attributes to assign to None",
			},
			"tags": tag.ResourceSchema(),
			"event_rule": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:    true,
						Description: "Account ID",
					},
					"active_yn": schema.StringAttribute{
						Computed:    true,
						Description: "Whether the Event rule is active",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Created date time",
					},
					"created_by": schema.StringAttribute{
						Computed:    true,
						Description: "Creator ID",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Event rule description",
					},
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "Event rule ID",
					},
					"modified_at": schema.StringAttribute{
						Computed:    true,
						Description: "Modified date time",
					},
					"modified_by": schema.StringAttribute{
						Computed:    true,
						Description: "Modifier ID",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Event rule name",
					},
					"resource_type_id": schema.StringAttribute{
						Computed:    true,
						Description: "Resource type ID",
					},
					"service_id": schema.StringAttribute{
						Computed:    true,
						Description: "Service ID",
					},
				},
				Computed:    true,
				Description: "Event rule",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *serviceWatchEventRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *serviceWatchEventRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicewatch.EventRuleResource
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new EventRule
	data, err := r.client.CreateEventRule(ctx, plan) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Event Rule",
			"Could not create Event Rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	eventRule := convertEventRule(&data.EventRule)
	eventRuleObjectValue, diags := types.ObjectValueFrom(ctx, eventRule.AttributeTypes(), eventRule)

	plan.Id = types.StringValue(eventRule.Id.ValueString())
	plan.EventRule = eventRuleObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceWatchEventRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicewatch.EventRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Resource Group
	data, err := r.client.GetEventRule(ctx, state.Id.ValueString())
	if err != nil && data == nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Event Rule",
			"Could not read Event Rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	// 존재하지 않는 event rule 조회 시 null 로 return
	if data.EventRule.GetId() == "" {
		state.EventRule = types.ObjectNull(state.EventRule.AttributeTypes(ctx))
		resp.State.Set(ctx, &state)
		return
	}

	// Map response body to schema and populate Computed attribute values
	eventRule := convertEventRule(&data.EventRule)
	eventRuleObjectValue, diags := types.ObjectValueFrom(ctx, eventRule.AttributeTypes(), eventRule)
	state.EventRule = eventRuleObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceWatchEventRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var state servicewatch.EventRuleResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Resource Group
	_, err := r.client.UpdateEventRule(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Event Rule",
			"Could not update Event Rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetEventRule(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading resourceGroup",
			"Could not read resourceGroup ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	eventRule := convertEventRule(&data.EventRule)
	eventRuleObjectValue, diags := types.ObjectValueFrom(ctx, eventRule.AttributeTypes(), eventRule)
	state.EventRule = eventRuleObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceWatchEventRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicewatch.EventRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventRuleId := state.Id.ValueString()

	// Delete existing Resource Group
	_, err := r.client.DeleteEventRule(ctx, []string{eventRuleId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Event Rule",
			"Could not delete Event Rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertEventRule(eventRuleResp *servicewatch2.EventRuleDTO) servicewatch.EventRule {
	return servicewatch.EventRule{
		AccountId:      types.StringValue(eventRuleResp.AccountId),
		ActiveYn:       types.StringValue(string(eventRuleResp.ActiveYn)),
		CreatedAt:      types.StringValue(eventRuleResp.CreatedAt.Format("2006-01-02 15:04:05")),
		CreatedBy:      types.StringValue(eventRuleResp.CreatedBy),
		Description:    types.StringPointerValue(eventRuleResp.Description.Get()),
		Id:             types.StringValue(eventRuleResp.Id),
		ModifiedAt:     types.StringValue(eventRuleResp.ModifiedAt.Format("2006-01-02 15:04:05")),
		ModifiedBy:     types.StringValue(eventRuleResp.ModifiedBy),
		Name:           types.StringValue(eventRuleResp.Name),
		ResourceTypeId: types.StringPointerValue(eventRuleResp.ResourceTypeId.Get()),
		ServiceId:      types.StringValue(eventRuleResp.ServiceId),
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *serviceWatchEventRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
