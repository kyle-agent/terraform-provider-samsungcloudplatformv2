package servicewatch

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &serviceWatchAlertResource{}
	_ resource.ResourceWithConfigure = &serviceWatchAlertResource{}
)

// NewServiceWatchAlertResource is a helper function to simplify the provider implementation.
func NewServiceWatchAlertResource() resource.Resource {
	return &serviceWatchAlertResource{}
}

// serviceWatchAlertResource is the data source implementation.
type serviceWatchAlertResource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *serviceWatchAlertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_alert"
}

// Schema defines the schema for the resource.
func (r *serviceWatchAlertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Alert Resource",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("LastUpdated"): schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Alert ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Alert name",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 100),
				},
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Alert type",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(AlertTypeMetric, AlertTypeService, AlertTypeComposite),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Alert description",
				Optional:    true,
			},
			common.ToSnakeCase("ActivatedYn"): schema.StringAttribute{
				Description: "Whether the Alert is activated or not",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(YnYes, YnNo),
				},
			},
			common.ToSnakeCase("Level"): schema.StringAttribute{
				Description: "Alert level - HIGH, MIDDLE, LOW",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(AlertLevelHigh, AlertLevelMiddle, AlertLevelLow),
				},
			},
			common.ToSnakeCase("NamespaceId"): schema.StringAttribute{
				Description: "Namespace ID",
				Computed:    true,
			},
			common.ToSnakeCase("NamespaceName"): schema.StringAttribute{
				Description: "Namespace name",
				Required:    true,
			},
			common.ToSnakeCase("MetricId"): schema.StringAttribute{
				Description: "Sharing type",
				Computed:    true,
			},
			common.ToSnakeCase("MetricName"): schema.StringAttribute{
				Description: "Metric name",
				Required:    true,
			},
			common.ToSnakeCase("Dimensions"): schema.ListNestedAttribute{
				Description: "List of dimension",
				Optional:    true,
				Computed:    true,
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
				Description: "Period (seconds)",
				Required:    true,
			},
			common.ToSnakeCase("Statistic"): schema.StringAttribute{
				Description: "Statistic - SUM, AVG, MAX, MIN",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(StatSum, StatAvg, StatMax, StatMin),
				},
			},
			common.ToSnakeCase("EvaluationCount"): schema.Int32Attribute{
				Description: "Evaluation count for the Alert condition",
				Optional:    true,
				Computed:    true,
			},
			common.ToSnakeCase("Threshold"): schema.Float32Attribute{
				Description: "Threshold for the Alert condition (except from RANGE operator)",
				Optional:    true,
				Computed:    true,
			},
			common.ToSnakeCase("UpperBound"): schema.Float32Attribute{
				Description: "Upper bound for the Alert range operator",
				Optional:    true,
				Computed:    true,
			},
			common.ToSnakeCase("LowerBound"): schema.Float32Attribute{
				Description: "Lower bound for the Alert range operator",
				Optional:    true,
				Computed:    true,
			},
			common.ToSnakeCase("Operator"): schema.StringAttribute{
				Description: "Operator - EQ, NOT_EQ, GT, GTE, LT, LTE, RANGE",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(OpEQ, OpNotEQ, OpGT, OpGTE, OpLT, OpLTE, OpRange),
				},
			},
			common.ToSnakeCase("ViolationCount"): schema.Int32Attribute{
				Description: "Violation count for the Alert condition",
				Optional:    true,
				Computed:    true,
			},
			common.ToSnakeCase("MissingDataOption"): schema.StringAttribute{
				Description: "Missing data option - MISSING, BREACHING, NOT_BREACHING, IGNORE",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(MissingDataMissing, MissingDataBreaching, MissingDataNotBreaching, MissingDataIgnore),
				},
			},
			common.ToSnakeCase("RecipientIds"): schema.ListAttribute{
				Description: "List of user IDs",
				Optional:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),
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
	}
}

// Configure adds the provider configured client to the data source.
func (r *serviceWatchAlertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *serviceWatchAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicewatch.AlertResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Dimension Keys for Reading Metric List
	dimensionKeys, diags := getDimensionKeys(ctx, plan.Dimensions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Metric List By NamespaceName, MetricName, DimensionKeys
	r.readMetrics(ctx, &plan, dimensionKeys, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Alert
	data, err := r.client.CreateAlert(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrCreateAlert,
			fmt.Sprintf(ErrCreateAlertFmt, err.Error(), detail),
		)
		return
	}

	// Fetch updated items from GetAlert as UpdateAlert items are not populated.
	alertResp, err := r.client.GetAlert(ctx, data.GetId())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadAlert,
			fmt.Sprintf(ErrReadAlertFmt, data.GetId(), err.Error(), detail),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	convertFromAlertDetailResponse(ctx, &plan, alertResp)
	plan.Id = types.StringValue(data.GetId())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceWatchAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicewatch.AlertResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Resource Group
	alertResp, err := r.client.GetAlert(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadAlert,
			fmt.Sprintf(ErrReadAlertFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}
	// Map response body to schema and populate Computed attribute values
	convertFromAlertDetailResponse(ctx, &state, alertResp)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceWatchAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state servicewatch.AlertResource
	var plan servicewatch.AlertResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 변경 사항이 없으면 state 값 셋팅
	plan.ActivatedYn = useStateIfUnset(plan.ActivatedYn, state.ActivatedYn)
	plan.Description = useStateIfUnset(plan.Description, state.Description)
	plan.MetricId = useStateIfUnset(plan.MetricId, state.MetricId)
	plan.NamespaceId = useStateIfUnset(plan.NamespaceId, state.NamespaceId)
	plan.EvaluationCount = useStateIfUnset(plan.EvaluationCount, state.EvaluationCount)
	plan.ViolationCount = useStateIfUnset(plan.ViolationCount, state.ViolationCount)

	// operator 변경 시 threshold/bound 값 처리
	if !plan.Operator.Equal(state.Operator) {
		if plan.Operator.ValueString() == "RANGE" {
			plan.Threshold = types.Float32Null()
		} else {
			plan.UpperBound = types.Float32Null()
			plan.LowerBound = types.Float32Null()
		}
	} else {
		plan.Threshold = useStateIfUnset(plan.Threshold, state.Threshold)
		plan.UpperBound = useStateIfUnset(plan.UpperBound, state.UpperBound)
		plan.LowerBound = useStateIfUnset(plan.LowerBound, state.LowerBound)
	}

	// activatedYn 이 변경되면, activated Update 수행
	if !plan.ActivatedYn.Equal(state.ActivatedYn) {
		_, err := r.client.UpdateAlertActivated(ctx, plan.Id.ValueString(), plan.ActivatedYn.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				ErrUpdateActivatedAlert,
				fmt.Sprintf(ErrUpdateActivatedAlertFmt, err.Error(), detail),
			)
			return
		}
		state.ActivatedYn = plan.ActivatedYn
	}

	// description 이 변경되면, description Update 수행
	if !plan.Description.Equal(state.Description) {
		_, err := r.client.UpdateAlertDescription(ctx, plan.Id.ValueString(), plan.Description.ValueString())
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				ErrUpdateDescriptionAlert,
				fmt.Sprintf(ErrUpdateDescriptionAlertFmt, err.Error(), detail),
			)
			return
		}
		state.Description = plan.Description
	}
	// metric 정보 변경 시, GetMetric 호출하여 Id 조회
	planDimensionKeys, diags := getDimensionKeys(ctx, plan.Dimensions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateDimensionKeys, diags := getDimensionKeys(ctx, state.Dimensions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.NamespaceName.Equal(state.NamespaceName) || !plan.MetricName.Equal(state.MetricName) || !reflect.DeepEqual(planDimensionKeys, stateDimensionKeys) {
		r.readMetrics(ctx, &plan, planDimensionKeys, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	needsUpdateResult := needsUpdate(plan, state)
	if needsUpdateResult {
		_, err := r.client.UpdateAlert(ctx, plan.Id.ValueString(), plan)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				ErrUpdateAlert,
				fmt.Sprintf(ErrUpdateAlertFmt, err.Error(), detail),
			)
			return
		}
	}

	// Fetch updated items from GetAlert as UpdateAlert items are not populated.
	alertResp, err := r.client.GetAlert(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadAlert,
			fmt.Sprintf(ErrReadAlertFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	convertFromAlertDetailResponse(ctx, &state, alertResp)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceWatchAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicewatch.AlertResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alertId := state.Id.ValueString()

	// Delete existing Resource Group
	_, err := r.client.DeleteAlert(ctx, []string{alertId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrDeleteAlert,
			fmt.Sprintf(ErrDeleteAlertFmt, err.Error(), detail),
		)
		return
	}
}

func (r *serviceWatchAlertResource) readMetrics(ctx context.Context, plan *servicewatch.AlertResource, dimensionKeys [][]string, diagnostics *diag.Diagnostics) {
	// call metric info list api
	metrics, err := r.client.GetMetrics(ctx, plan.NamespaceName.ValueString(), plan.MetricName.ValueString(), dimensionKeys)
	if err != nil {
		detail := client.GetDetailFromError(err)
		diagnostics.AddError(
			ErrGetMetrics,
			fmt.Sprintf(ErrGetMetricsFmt, err.Error(), detail, plan.NamespaceName.ValueString(), plan.MetricName.ValueString(), dimensionKeys),
		)
		return
	}

	// if namespace does not exist, generate 404 error
	if len(metrics.Namespaces) == 0 {
		diagnostics.AddError(
			ErrReadMetrics,
			fmt.Sprintf(ErrReadMetricsFmt, plan.NamespaceName.ValueString(), plan.MetricName.ValueString(), dimensionKeys),
		)
		return
	}

	// if metric does not exist, generate 404 error
	namespace := metrics.Namespaces[0]
	plan.NamespaceId = types.StringValue(namespace.GetId())
	if len(namespace.Dimensions) > 0 && len(namespace.Dimensions[0].Metrics) > 0 {
		plan.MetricId = types.StringValue(namespace.Dimensions[0].Metrics[0].GetId())
	} else {
		diagnostics.AddError(
			ErrReadMetrics,
			fmt.Sprintf(ErrReadMetricsFmt, plan.NamespaceName.ValueString(), plan.MetricName.ValueString(), dimensionKeys),
		)
		return
	}
}

func getDimensionKeys(ctx context.Context, dimensions types.List) ([][]string, diag.Diagnostics) {
	var keys []string
	if !dimensions.IsNull() && !dimensions.IsUnknown() {
		var items []servicewatch.Dimension
		diags := dimensions.ElementsAs(ctx, &items, false)
		if diags.HasError() {
			return nil, diags
		}

		for _, item := range items {
			keys = append(keys, item.Key.ValueString())
		}
	}

	// When the alert has no dimensions, return an empty outer slice so the
	// metric search request omits the dimensions filter (sends []) instead of
	// wrapping an empty key list ([[]]), which would not match a dimensionless
	// metric and cause Create to fail with a 404.
	if len(keys) == 0 {
		return [][]string{}, nil
	}

	return [][]string{keys}, nil
}

func convertFromAlertDetailResponse(ctx context.Context, state *servicewatch.AlertResource, alertResp *servicewatch2.AlertDetailResponse) servicewatch.AlertResource {
	var dimensions []servicewatch.Dimension
	for _, dimension := range alertResp.Dimensions {
		dimensions = append(dimensions, servicewatch.Dimension{
			Key:   types.StringValue(dimension.GetKey()),
			Value: types.StringValue(dimension.GetValue()),
		})
	}
	state.Dimensions, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.Dimension{}.AttributeTypes()}, dimensions)

	state.Name = types.StringValue(alertResp.GetName())
	state.Description = nullableStringTypes(alertResp.GetDescriptionOk())
	state.ActivatedYn = types.StringValue(string(alertResp.GetActivatedYn()))
	state.Type = types.StringValue(string(alertResp.GetType()))
	state.Level = types.StringValue(string(alertResp.GetLevel()))
	state.NamespaceId = types.StringValue(alertResp.GetNamespaceId())
	state.NamespaceName = types.StringValue(alertResp.GetNamespaceName())
	state.MetricId = types.StringValue(alertResp.GetMetricId())
	state.MetricName = types.StringValue(alertResp.GetMetricName())
	state.Period = types.Int32Value(alertResp.GetPeriod())
	state.Statistic = types.StringValue(string(alertResp.GetStatistic()))
	state.EvaluationCount = types.Int32Value(alertResp.GetEvaluationCount())
	state.Threshold = nullableFloat32Types(alertResp.GetThresholdOk())
	state.UpperBound = nullableFloat32Types(alertResp.GetUpperBoundOk())
	state.LowerBound = nullableFloat32Types(alertResp.GetLowerBoundOk())
	state.Operator = types.StringValue(string(alertResp.GetOperator()))
	state.ViolationCount = types.Int32Value(alertResp.GetViolationCount())
	state.MissingDataOption = types.StringValue(string(alertResp.GetMissingDataOption()))
	state.CreatedAt = types.StringValue(alertResp.GetCreatedAt().Format(TimeFormatDisplay))
	state.CreatedBy = types.StringValue(alertResp.GetCreatedBy())
	state.ModifiedAt = types.StringValue(alertResp.GetModifiedAt().Format(TimeFormatDisplay))
	state.ModifiedBy = types.StringValue(alertResp.GetModifiedBy())

	return *state
}

func useStateIfUnset[T attr.Value](plan, state T) T {
	if plan.IsNull() || plan.IsUnknown() {
		return state
	}
	return plan
}

func needsUpdate(plan, state servicewatch.AlertResource) bool {
	planReq := convertUpdateModel(plan)
	stateReq := convertUpdateModel(state)

	return !reflect.DeepEqual(planReq, stateReq)
}

func convertUpdateModel(model servicewatch.AlertResource) servicewatch.AlertUpdateModel {

	return servicewatch.AlertUpdateModel{
		Level:             model.Level,
		NamespaceId:       model.NamespaceId,
		MetricId:          model.MetricId,
		Dimensions:        model.Dimensions,
		Period:            model.Period,
		Statistic:         model.Statistic,
		EvaluationCount:   model.EvaluationCount,
		Threshold:         model.Threshold,
		UpperBound:        model.UpperBound,
		LowerBound:        model.LowerBound,
		Operator:          model.Operator,
		ViolationCount:    model.ViolationCount,
		MissingDataOption: model.MissingDataOption,
	}
}
