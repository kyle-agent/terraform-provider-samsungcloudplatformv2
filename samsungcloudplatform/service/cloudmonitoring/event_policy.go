package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpcloudmonitoring "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/cloudmonitoring"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &cloudMonitoringEventPolicyResource{}
	_ resource.ResourceWithConfigure = &cloudMonitoringEventPolicyResource{}
)

func NewCloudMonitoringEventPolicyResource() resource.Resource {
	return &cloudMonitoringEventPolicyResource{}
}

type cloudMonitoringEventPolicyResource struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

func (r *cloudMonitoringEventPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_policy" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (r *cloudMonitoringEventPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.CloudMonitoring
	r.clients = inst.Client
}

func (r *cloudMonitoringEventPolicyResource) Schema(ctx context.Context, request resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Event Policy",
		Attributes: map[string]schema.Attribute{
			"event_policy_id": schema.Int64Attribute{
				Description:   "Identifier of the resource.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{},
				//PlanModifiers: []planmodifier.String{
				//	stringplanmodifier.UseStateForUnknown(),
				//},
			},
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Optional:    true,
			},
			common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
				Description: "ProductResourceId",
				Optional:    true,
			},
			common.ToSnakeCase("EventLevel"): schema.StringAttribute{
				Description: "EventLevel",
				Optional:    true,
			},
			common.ToSnakeCase("DisableYn"): schema.StringAttribute{
				Description: "DisableYn",
				Optional:    true,
			},
			common.ToSnakeCase("EventMessagePrefix"): schema.StringAttribute{
				Description: "EventMessagePrefix",
				Optional:    true,
			},
			common.ToSnakeCase("FtCount"): schema.Int64Attribute{
				Description: "FtCount",
				Optional:    true,
			},
			common.ToSnakeCase("IsLogMetric"): schema.StringAttribute{
				Description: "IsLogMetric",
				Optional:    true,
			},
			common.ToSnakeCase("MetricKey"): schema.StringAttribute{
				Description: "MetricKey",
				Optional:    true,
			},
			common.ToSnakeCase("ObjectName"): schema.StringAttribute{
				Description: "ObjectName",
				Optional:    true,
			}, common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
				Description: "ObjectDisplayName",
				Optional:    true,
			},
			common.ToSnakeCase("EventOccurTimeZone"): schema.StringAttribute{
				Description: "EventOccurTimeZone",
				Optional:    true,
			},
			common.ToSnakeCase("ObjectType"): schema.StringAttribute{
				Description: "ObjectType",
				Optional:    true,
			},
			common.ToSnakeCase("ObjectTypeName"): schema.StringAttribute{
				Description: "ObjectTypeName",
				Optional:    true,
			},
			//common.ToSnakeCase("EventThreshold"): schema.SingleNestedAttribute{
			//	Description: "EventThreshold",
			//	Optional:    true,
			//	Attributes: map[string]schema.Attribute{
			//		common.ToSnakeCase("ThresholdType"): schema.StringAttribute{
			//			Description: "ThresholdType",
			//			Optional:    true,
			//		},
			//		common.ToSnakeCase("MetricFunction"): schema.StringAttribute{
			//			Description: "MetricFunction",
			//			Optional:    true,
			//		},
			//		common.ToSnakeCase("SingleThreshold"): schema.SingleNestedAttribute{
			//			Description: "SingleThreshold",
			//			Optional:    true,
			//			Attributes: map[string]schema.Attribute{
			//				common.ToSnakeCase("ComparisonOperator"): schema.StringAttribute{
			//					Description: "ComparisonOperator",
			//					Optional:    true,
			//				},
			//				common.ToSnakeCase("Value"): schema.Float64Attribute{
			//					Description: "Value",
			//					Optional:    true,
			//				},
			//			},
			//		},
			//	},
			//},
			common.ToSnakeCase("EventPolicy"): schema.SingleNestedAttribute{
				Description: "EventPolicy",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("ModifiedDt"): schema.StringAttribute{
						Description: "ModifiedDt",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedByName"): schema.StringAttribute{
						Description: "CreatedByName",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedByName"): schema.StringAttribute{
						Description: "ModifiedByName",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedById"): schema.StringAttribute{
						Description: "CreatedById",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedDt"): schema.StringAttribute{
						Description: "CreatedDt",
						Optional:    true,
					},
					common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
						Description: "EventPolicyId",
						Optional:    true,
					},
					common.ToSnakeCase("MetricName"): schema.StringAttribute{
						Description: "MetricName",
						Optional:    true,
					},
					common.ToSnakeCase("MetricType"): schema.StringAttribute{
						Description: "MetricType",
						Optional:    true,
					},
					common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
						Description: "MetricType",
						Optional:    true,
					},
					common.ToSnakeCase("ProductSq"): schema.Int64Attribute{
						Description: "ProductSq",
						Optional:    true,
					},
					common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
						Description: "ProductResourceId",
						Optional:    true,
					},
					common.ToSnakeCase("MetricKey"): schema.StringAttribute{
						Description: "MetricKey",
						Optional:    true,
					},
					common.ToSnakeCase("MetricDescription"): schema.StringAttribute{
						Description: "MetricDescription",
						Optional:    true,
					},
					common.ToSnakeCase("MetricDescriptionEn"): schema.StringAttribute{
						Description: "MetricDescriptionEn",
						Optional:    true,
					},
					common.ToSnakeCase("ProductTargetType"): schema.StringAttribute{
						Description: "ProductTargetType",
						Optional:    true,
					},
					common.ToSnakeCase("ProductTargetTypeEn"): schema.StringAttribute{
						Description: "ProductTargetTypeEn",
						Optional:    true,
					},
					common.ToSnakeCase("IsLogMetric"): schema.StringAttribute{
						Description: "IsLogMetric",
						Optional:    true,
					},
					common.ToSnakeCase("ObjectName"): schema.StringAttribute{
						Description: "ObjectName",
						Optional:    true,
					}, common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
						Description: "ObjectDisplayName",
						Optional:    true,
					},
					common.ToSnakeCase("EventLevel"): schema.StringAttribute{
						Description: "EventLevel",
						Optional:    true,
					}, common.ToSnakeCase("FtCount"): schema.Int64Attribute{
						Description: "FtCount",
						Optional:    true,
					}, common.ToSnakeCase("EventMessagePrefix"): schema.StringAttribute{
						Description: "EventMessagePrefix",
						Optional:    true,
					}, common.ToSnakeCase("ObjectType"): schema.StringAttribute{
						Description: "ObjectType",
						Optional:    true,
					},
					common.ToSnakeCase("ObjectTypeName"): schema.StringAttribute{
						Description: "ObjectTypeName",
						Optional:    true,
					},
					common.ToSnakeCase("ProductInfoAttrs"): schema.StringAttribute{
						Description: "ProductInfoAttrs",
						Optional:    true,
					},
					common.ToSnakeCase("DisableObject"): schema.StringAttribute{
						Description: "DisableObject",
						Optional:    true,
					},
					common.ToSnakeCase("UserNames"): schema.StringAttribute{
						Description: "UserNames",
						Optional:    true,
					},
					common.ToSnakeCase("UserNameStr"): schema.StringAttribute{
						Description: "UserNameStr",
						Optional:    true,
					},
					common.ToSnakeCase("DisableYn"): schema.StringAttribute{
						Description: "DisableYn",
						Optional:    true,
					},
					common.ToSnakeCase("AttrListStr"): schema.StringAttribute{
						Description: "AttrListStr",
						Optional:    true,
					},
					common.ToSnakeCase("AsgYn"): schema.StringAttribute{
						Description: "AsgYn",
						Optional:    true,
					},
					common.ToSnakeCase("StartDt"): schema.StringAttribute{
						Description: "StartDt",
						Optional:    true,
					},
					common.ToSnakeCase("DisplayEventRule"): schema.StringAttribute{
						Description: "DisplayEventRule",
						Optional:    true,
					},
					common.ToSnakeCase("CheckAsg"): schema.BoolAttribute{
						Description: "CheckAsg",
						Optional:    true,
					},
					common.ToSnakeCase("EventOccurTimeZone"): schema.StringAttribute{
						Description: "EventOccurTimeZone",
						Optional:    true,
					},
					common.ToSnakeCase("EventThreshold"): schema.SingleNestedAttribute{
						Description: "EventThreshold",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ThresholdType"): schema.StringAttribute{
								Description: "ThresholdType",
								Optional:    true,
							},
							common.ToSnakeCase("MetricFunction"): schema.StringAttribute{
								Description: "MetricFunction",
								Optional:    true,
							},
							common.ToSnakeCase("SingleThreshold"): schema.SingleNestedAttribute{
								Description: "SingleThreshold",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("ComparisonOperator"): schema.StringAttribute{
										Description: "ComparisonOperator",
										Optional:    true,
									},
									common.ToSnakeCase("Value"): schema.Float64Attribute{
										Description: "Value",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			common.ToSnakeCase("NotificationRecipient"): schema.ListNestedAttribute{
				Description: "NotificationRecipient",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("RecipientType"): schema.StringAttribute{
							Description: "RecipientType",
							Optional:    true,
						},
						common.ToSnakeCase("RecipientKey"): schema.StringAttribute{
							Description: "RecipientKey",
							Optional:    true,
						},
						common.ToSnakeCase("RecipientTargetTable"): schema.StringAttribute{
							Description: "RecipientTargetTable",
							Optional:    true,
						},
						common.ToSnakeCase("RcvUserName"): schema.StringAttribute{
							Description: "RcvUserName",
							Optional:    true,
						},
						common.ToSnakeCase("NotificationMethod"): schema.ListNestedAttribute{
							Description: "NotificationMethod",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("EventLevel"): schema.StringAttribute{
										Description: "EventLevel",
										Optional:    true,
									},
									common.ToSnakeCase("SendMethod"): schema.ListAttribute{
										ElementType: types.StringType,
										Description: "SendMethod",
										Optional:    true,
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

func (r *cloudMonitoringEventPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudmonitoring.EventPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateEventPolicy(plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	newState, err := r.MapGetResponseToState(ctx, data, plan.ResourceType.ValueString(), plan.ProductResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *cloudMonitoringEventPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudmonitoring.EventPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetEventPolicy(state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	newState, err := r.MapGetResponseToState(ctx, data, state.ResourceType.ValueString(), state.ProductResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cloudMonitoringEventPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	//var plan cloudmonitoring.EventPolicyResource
	var state cloudmonitoring.EventPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateEventPolicy(state, state.EventPolicyId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating event policy",
			"Could not update event policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetEventPolicy(state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading EventPolicy",
			"Could not read EventPolicy ID "+state.EventPolicyId.String()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	newState, err := r.MapGetResponseToState(ctx, data, state.ResourceType.ValueString(), state.ProductResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *cloudMonitoringEventPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudmonitoring.EventPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.DeleteEventPolicy(state.EventPolicyId, state.ResourceType)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting event policy",
			"Could not delete event policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

}

func (r *cloudMonitoringEventPolicyResource) MapGetResponseToState(ctx context.Context,
	resp *scpcloudmonitoring.EventPolicyDetailResponse, xResourceType string, productResourceId string) (cloudmonitoring.EventPolicyResource, error) {

	singleThreshold := cloudmonitoring.SingleThreshold{
		ComparisonOperator: types.StringValue(resp.EventThreshold.GetSingleThreshold().ComparisonOperator),
		Value:              types.Float64Value(resp.EventThreshold.GetSingleThreshold().Value),
	}

	//}
	//rangeThreshold := cloudmonitoring.RangeThreshold{
	//	MaxComparisonOperator: types.StringValue(data.EventThreshold.GetRangeThreshold().MaxComparisonOperator),
	//	MinComparisonOperator: types.StringValue(data.EventThreshold.GetRangeThreshold().MinComparisonOperator),
	//	MaxValue:              types.Float64Value(data.EventThreshold.RangeThreshold.MaxValue),
	//	MinValue:              types.Float64Value(data.EventThreshold.RangeThreshold.MinValue),
	//}

	//	EventThreshold
	threshold := resp.GetEventThreshold()
	eventThreshold := cloudmonitoring.EventThreshold{
		ThresholdType:   types.StringValue(threshold.ThresholdType),
		MetricFunction:  types.StringValue(threshold.GetMetricFunction()),
		SingleThreshold: singleThreshold,
		//RangeThreshold:  rangeThreshold,
	}

	eventPolicyElement := cloudmonitoring.EventPolicy{
		//CreatedByName: types.StringValue(data.GetCreatedByName()),
		//CreateBy:      createby,
		//UpdateById: types.StringValue(data.GetUpdateById()),
		//UpdateBy:            createby,
		//EventPolicyId: types.Int64Value(resp.EventPolicyId),
		//ProductSummary:      productSummary,ne()),
		EventThreshold: eventThreshold,
		//MetricSummary:       metricSummary,
		//ProductSq:           types.Int64Value(resp.GetProductSq()),
		//ProductResourceId:   types.StringValue(resp.GetProductResourceId()),
		MetricKey:           types.StringValue(resp.MetricKey),
		MetricDescription:   types.StringValue(resp.GetMetricDescription()),
		MetricDescriptionEn: types.StringValue(resp.GetMetricDescriptionEn()),
		//ProductTargetType:   types.StringValue(resp.GetProductTargetType()),
		//ProductTargetTypeEn: types.StringValue(resp.GetProductTargetTypeEn()),
		IsLogMetric: types.StringValue(resp.GetIsLogMetric()),
		//ObjectName:         types.StringValue(resp.GetObjectName()),
		ObjectDisplayName:  types.StringValue(resp.GetObjectDisplayName()),
		EventLevel:         types.StringValue(resp.EventLevel),
		FtCount:            types.Int64Value(resp.FtCount),
		EventMessagePrefix: types.StringValue(resp.GetEventMessagePrefix()),
		ObjectType:         types.StringValue(resp.GetObjectType()),
		ObjectTypeName:     types.StringValue(resp.GetObjectTypeName()),
		//ProductInfoAttrs:    types.StringValue(resp.GetProductInfoAttrs(),
		DisableObject:    types.StringValue(resp.GetDisableObject()),
		UserNames:        types.StringValue(resp.GetUserNames()),
		UserNameStr:      types.StringValue(resp.GetUserNameStr()),
		DisableYn:        types.StringValue("N"),
		AttrListStr:      types.StringValue(resp.GetAttrListStr()),
		AsgYn:            types.StringValue(resp.GetAsgYn()),
		StartDt:          types.StringValue(resp.GetStartDt().String()),
		DisplayEventRule: types.StringValue(resp.GetDisplayEventRule()),
		//CheckAsg:           types.BoolValue(resp.GetCheckAsg()),
		EventOccurTimeZone: types.StringValue(resp.GetEventOccurTimeZone()),
		//EventPolicyStatistics: eventPolicyStatistics,
	}

	return cloudmonitoring.EventPolicyResource{
		EventPolicyId:     types.Int64Value(resp.EventPolicyId),
		EventPolicy:       eventPolicyElement,
		ResourceType:      types.StringValue(xResourceType),
		ProductResourceId: types.StringValue(productResourceId),
		//EventThreshold:        eventThreshold,
		NotificationRecipient: types.ListNull(types.ObjectType{AttrTypes: cloudmonitoring.NotificationRecipient{}.AttributeTypes()}),
		EventLevel:            types.StringValue(resp.GetEventLevel()),
		DisableYn:             types.StringValue("N"),
		EventMessagePrefix:    types.StringValue(resp.GetEventMessagePrefix()),
		ObjectType:            types.StringValue(resp.GetObjectType()),
		MetricKey:             types.StringValue(resp.GetMetricKey()),
		FtCount:               types.Int64Value(resp.FtCount),
		ObjectName:            types.StringValue(resp.GetObjectName()),
		IsLogMetric:           types.StringValue(resp.GetIsLogMetric()),
		ObjectTypeName:        types.StringValue(resp.GetObjectTypeName()),
		EventOccurTimeZone:    types.StringValue(resp.GetEventOccurTimeZone()),
	}, nil

}
