package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/budget" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &budgetBudgetResource{}
	_ resource.ResourceWithConfigure   = &budgetBudgetResource{}
	_ resource.ResourceWithImportState = &budgetBudgetResource{}
)

// NewBudgetBudgetResource is a helper function to simplify the provider implementation.
func NewBudgetBudgetResource() resource.Resource {
	return &budgetBudgetResource{}
}

// budgetDataSource is the data source implementation.
type budgetBudgetResource struct {
	config  *scpsdk.Configuration
	client  *budget.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *budgetBudgetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_budget_budget"
}

// Schema defines the schema for the data source.
func (r *budgetBudgetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = BudgetDataSourceSchema()
}

// Configure adds the provider configured client to the data source.
func (r *budgetBudgetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Budget
	r.clients = inst.Client
}

// Create creates the resource and sets Terraform state.
func (r *budgetBudgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan budget.BudgetResource
	diags := req.Plan.Get(ctx, &plan)

	if len(diags) > 0 {
		for i, diag := range diags {
			fmt.Printf("  [%d] Severity: %s, Summary: %s, Detail: %s\n",
				i,
				diag.Severity(),
				diag.Summary(),
				diag.Detail())
		}
	} else {
		fmt.Println("No diagnostics found.")
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateAccountBudget(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating budget",
			"Could not create budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	budgetData := data.Budget
	plan.Id = types.StringValue(budgetData.Id)
	budgetModel := budget.Budget{
		Amount:     types.Int32Value(budgetData.Amount),
		CreatedAt:  types.StringValue(budgetData.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		BudgetId:   types.StringValue(budgetData.Id),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
		Name:       types.StringValue(budgetData.Name),
		StartMonth: types.StringValue(budgetData.StartMonth),
		Unit:       types.StringValue(budgetData.Unit),
	}
	budgetObjectValue, diags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
	plan.Budget = budgetObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state budget.BudgetResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.SetAccountBudget(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Budget",
			"Could not update Budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Budget",
			"Could not read Budget ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	budgetData := data.Budget
	budgetModel := budget.Budget{
		Amount:     types.Int32Value(budgetData.Amount),
		CreatedAt:  types.StringValue(budgetData.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetData.CreatedBy),
		BudgetId:   types.StringValue(budgetData.Id),
		ModifiedAt: types.StringValue(budgetData.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetData.ModifiedBy),
		Name:       types.StringValue(budgetData.Name),
		StartMonth: types.StringValue(budgetData.StartMonth),
		Unit:       types.StringValue(budgetData.Unit),
	}
	budgetObjectValue, diags := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
	state.Budget = budgetObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state budget.BudgetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Budget",
			"Could not delete Budget, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// Read read the resource and sets the updated Terraform state on success.
func (r *budgetBudgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state budget.BudgetResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetAccountBudget(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account Budget",
			err.Error(),
		)
		return
	}

	budgetElement := data.Budget

	budgetModel := budget.Budget{
		BudgetId:   types.StringValue(budgetElement.Id),
		Name:       types.StringValue(budgetElement.Name),
		Amount:     types.Int32Value(budgetElement.Amount),
		BudgetType: types.StringValue(budgetElement.Type),
		Unit:       types.StringValue(budgetElement.Unit),
		CreatedAt:  types.StringValue(budgetElement.ModifiedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringPointerValue(budgetElement.CreatedBy),
		ModifiedAt: types.StringValue(budgetElement.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringPointerValue(budgetElement.ModifiedBy),
	}

	budgetObjectValue, _ := types.ObjectValueFrom(ctx, budgetModel.AttributeTypes(), budgetModel)
	state.Budget = budgetObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func BudgetDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LastUpdated"): schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("Budget"): schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Amount"): schema.Int32Attribute{
						Computed:            true,
						Description:         "Budget amount",
						MarkdownDescription: "Budget amount",
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Computed:            true,
						Description:         "Created datetime",
						MarkdownDescription: "Created datetime",
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Computed:            true,
						Description:         "Created user",
						MarkdownDescription: "Created user",
					},
					common.ToSnakeCase("BudgetId"): schema.StringAttribute{
						Computed:            true,
						Description:         "Budget id",
						MarkdownDescription: "Budget id",
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Computed:            true,
						Description:         "Modified datetime",
						MarkdownDescription: "Modified datetime",
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Computed:            true,
						Description:         "Modified user",
						MarkdownDescription: "Modified user",
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Computed:            true,
						Description:         "Budget name",
						MarkdownDescription: "Budget name",
					},
					common.ToSnakeCase("StartMonth"): schema.StringAttribute{
						Computed:            true,
						Description:         "Budget start month",
						MarkdownDescription: "Budget start month",
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Computed:            true,
						Description:         "Budget type",
						MarkdownDescription: "Budget type",
					},
					common.ToSnakeCase("Unit"): schema.StringAttribute{
						Computed:            true,
						Description:         "Budget management unit",
						MarkdownDescription: "Budget management unit",
					},
				},
				Computed: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name (between 1 and 64 characters)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			common.ToSnakeCase("Amount"): schema.Int32Attribute{
				Description: "Amount",
				Optional:    true,
			},
			common.ToSnakeCase("StartMonth"): schema.StringAttribute{
				Description: "StartMonth",
				Optional:    true,
			},
			common.ToSnakeCase("Unit"): schema.StringAttribute{
				Description: "Unit",
				Optional:    true,
			},
			common.ToSnakeCase("Notifications"): schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUseNotification"): schema.BoolAttribute{
						Optional:            true,
						Description:         "Notification use state",
						MarkdownDescription: "Notification use state",
					},
					common.ToSnakeCase("NotificationSendPeriod"): schema.StringAttribute{
						Optional:            true,
						Description:         "Notification send period first/daily/none",
						MarkdownDescription: "Notification send period first/daily/none",
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "List of notification recipient email addresses",
						MarkdownDescription: "List of notification recipient email addresses",
					},
					common.ToSnakeCase("Thresholds"): schema.ListAttribute{
						ElementType:         types.Int32Type,
						Optional:            true,
						Description:         "List of threshold percentages for notifications",
						MarkdownDescription: "List of threshold percentages for notifications",
					},
				},
				Optional:            true,
				Description:         "Notification settings for the budget",
				MarkdownDescription: "Notification settings for the budget",
			},
			common.ToSnakeCase("Prevention"): schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("IsUsePrevention"): schema.BoolAttribute{
						Optional:            true,
						Description:         "Auto Generation prevent use state",
						MarkdownDescription: "Auto Generation prevent use state",
					},
					common.ToSnakeCase("Receivers"): schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "List of notification recipient email addresses",
						MarkdownDescription: "List of notification recipient email addresses",
					},
					common.ToSnakeCase("Threshold"): schema.Int32Attribute{
						Optional:            true,
						Description:         "New Creation prevention thresholds value 70/80/90/100",
						MarkdownDescription: "New Creation prevention thresholds value 70/80/90/100",
					},
				},
				Optional:            true,
				Description:         "Auto generation prevention settings for the budget",
				MarkdownDescription: "Auto generation prevention settings for the budget",
			},
		},
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *budgetBudgetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
