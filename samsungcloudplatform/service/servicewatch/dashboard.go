package servicewatch

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serviceWatchDashboardResource{}
	_ resource.ResourceWithConfigure   = &serviceWatchDashboardResource{}
	_ resource.ResourceWithImportState = &serviceWatchDashboardResource{}
)

// NewServiceWatchDashboardResource is a helper function to simplify the provider implementation.
func NewServiceWatchDashboardResource() resource.Resource {
	return &serviceWatchDashboardResource{}
}

// serviceWatchDashboardResource is the resource implementation.
type serviceWatchDashboardResource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *serviceWatchDashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_dashboard"
}

// Schema defines the schema for the resource.
func (r *serviceWatchDashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Dashboard Resource",
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Dashboard",
				Computed:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Dashboard ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Dashboard name",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Dashboard type",
				Computed:    true,
			},
			common.ToSnakeCase("Srn"): schema.StringAttribute{
				Description: "Service resource name",
				Computed:    true,
			},
			common.ToSnakeCase("ShareType"): schema.StringAttribute{
				Description: "Sharing type",
				Computed:    true,
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
			common.ToSnakeCase("Widgets"): schema.ListNestedAttribute{
				Description: "List of widgets",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Widget ID",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Widget type",
							Required:    true,
						},
						common.ToSnakeCase("Width"): schema.Int32Attribute{
							Description: "Widget width",
							Required:    true,
						},
						common.ToSnakeCase("Height"): schema.Int32Attribute{
							Description: "Widget height",
							Required:    true,
						},
						common.ToSnakeCase("Order"): schema.Int32Attribute{
							Description: "Widget's order in the dashboard",
							Required:    true,
						},
						common.ToSnakeCase("Properties"): schema.SingleNestedAttribute{
							Description: "Widget's detailed properties",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Title"): schema.StringAttribute{
									Description: "Widget title",
									Required:    true,
								},
								common.ToSnakeCase("Period"): schema.Int32Attribute{
									Description: "Query period (seconds)",
									Optional:    true,
								},
								common.ToSnakeCase("Stacked"): schema.BoolAttribute{
									Description: "Whether the graph is stacked",
									Required:    true,
								},
								common.ToSnakeCase("StatisticType"): schema.StringAttribute{
									Description: "Statistical function",
									Optional:    true,
								},
								common.ToSnakeCase("View"): schema.StringAttribute{
									Description: "View type",
									Required:    true,
								},
								common.ToSnakeCase("Metrics"): schema.ListNestedAttribute{
									Description: "List of metrics included in the widget",
									Required:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description: "Metric name",
												Required:    true,
											},
											common.ToSnakeCase("NamespaceName"): schema.StringAttribute{
												Description: "Namespace name",
												Required:    true,
											},
											common.ToSnakeCase("DisplayName"): schema.StringAttribute{
												Description: "Display name (label) of the metric",
												Required:    true,
											},
											common.ToSnakeCase("Color"): schema.StringAttribute{
												Description: "Metric line color",
												Required:    true,
											},
											common.ToSnakeCase("Dimensions"): schema.ListNestedAttribute{
												Description: "List of dimensions",
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														common.ToSnakeCase("Key"): schema.StringAttribute{
															Description: "Dimension key",
															Required:    true,
														},
														common.ToSnakeCase("Value"): schema.StringAttribute{
															Description: "Dimension value",
															Required:    true,
														},
													},
												},
											},
											common.ToSnakeCase("Period"): schema.Int32Attribute{
												Description: "Query period (seconds)",
												Optional:    true,
											},
											common.ToSnakeCase("StatisticType"): schema.StringAttribute{
												Description: "Statistical function",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *serviceWatchDashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			ErrUnexpectedConfigure,
			fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData),
		)

		return
	}

	r.client = inst.Client.ServiceWatch
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *serviceWatchDashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicewatch.DashboardResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Dashboard
	dashboard, err := r.client.CreateDashboard(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrCreateDashboard,
			fmt.Sprintf(ErrCreateDashboardFmt, err.Error(), detail),
		)
		return
	}

	// convert widget list response
	widgetResponses, ok := dashboard.GetWidgetsOk()
	var widgets types.List
	if ok && widgetResponses != nil {
		widgets, diags = convertWidget(ctx, widgetResponses)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(dashboard.Id)
	plan.Name = types.StringValue(dashboard.Name)
	plan.Type = types.StringValue(dashboard.Type)
	plan.Srn = types.StringValue(dashboard.Srn)
	plan.ShareType = types.StringValue(dashboard.ShareType)
	plan.CreatedAt = types.StringValue(dashboard.GetCreatedAt().Format(TimeFormatDisplay))
	plan.ModifiedAt = types.StringValue(dashboard.GetModifiedAt().Format(TimeFormatDisplay))
	plan.CreatedBy = types.StringValue(dashboard.GetCreatedBy())
	plan.ModifiedBy = types.StringValue(dashboard.GetModifiedBy())
	plan.Widgets = widgets

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceWatchDashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicewatch.DashboardResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Dashboard
	dashboard, err := r.client.GetDashboard(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadDashboard,
			fmt.Sprintf(ErrReadDashboardFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}

	// convert widget list response
	widgetResponses, ok := dashboard.GetWidgetsOk()
	var widgets types.List
	if ok && widgetResponses != nil {
		widgets, diags = convertWidget(ctx, widgetResponses)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// update state
	state.Id = types.StringValue(dashboard.Id)
	state.Name = types.StringValue(dashboard.Name)
	state.Type = types.StringValue(dashboard.Type)
	state.Srn = types.StringValue(dashboard.Srn)
	state.ShareType = types.StringValue(dashboard.ShareType)
	state.CreatedAt = types.StringValue(dashboard.GetCreatedAt().Format(TimeFormatDisplay))
	state.ModifiedAt = types.StringValue(dashboard.GetModifiedAt().Format(TimeFormatDisplay))
	state.CreatedBy = types.StringValue(dashboard.GetCreatedBy())
	state.ModifiedBy = types.StringValue(dashboard.GetModifiedBy())
	state.Widgets = widgets

	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceWatchDashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state servicewatch.DashboardResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Dashboard
	_, err := r.client.UpdateDashboard(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrUpdateDashboard,
			fmt.Sprintf(ErrUpdateDashboardFmt, err.Error(), detail),
		)
		return
	}

	// Fetch updated items from GetDashboard as UpdateDashboard items are not populated.
	dashboard, err := r.client.GetDashboard(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadDashboard,
			fmt.Sprintf(ErrReadDashboardFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}

	// convert widget list response
	widgetResponses, ok := dashboard.GetWidgetsOk()
	var widgets types.List
	if ok && widgetResponses != nil {
		widgets, diags = convertWidget(ctx, widgetResponses)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// update state
	state.Id = types.StringValue(dashboard.Id)
	state.Name = types.StringValue(dashboard.Name)
	state.Type = types.StringValue(dashboard.Type)
	state.Srn = types.StringValue(dashboard.Srn)
	state.ShareType = types.StringValue(dashboard.ShareType)
	state.CreatedAt = types.StringValue(dashboard.GetCreatedAt().Format(TimeFormatDisplay))
	state.ModifiedAt = types.StringValue(dashboard.GetModifiedAt().Format(TimeFormatDisplay))
	state.CreatedBy = types.StringValue(dashboard.GetCreatedBy())
	state.ModifiedBy = types.StringValue(dashboard.GetModifiedBy())
	state.Widgets = widgets
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceWatchDashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicewatch.DashboardResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboardId := state.Id.ValueString()

	// Delete existing Dashboard
	_, err := r.client.DeleteDashboard(ctx, []string{dashboardId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrDeleteDashboard,
			fmt.Sprintf(ErrDeleteDashboardFmt, err.Error(), detail),
		)
		return
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *serviceWatchDashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
