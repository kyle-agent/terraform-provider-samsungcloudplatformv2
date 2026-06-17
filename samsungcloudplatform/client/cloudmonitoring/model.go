package cloudmonitoring

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-cloudmonitoring" // 해당 서비스의 서비스 타입(keystone 에 등록된 service type)을 추가한다.

// ---------- Event(DataSource) ------------ //

type EventDataSourceId struct { //	/v2/events/{EventId} API Request
	EventId           types.String    `tfsdk:"event_id"`
	ResourceType      types.String    `tfsdk:"x_resource_type"`
	Filter            []filter.Filter `tfsdk:"filter"`
	EventLevel        types.String    `tfsdk:"event_level"`
	EventMessage      types.String    `tfsdk:"event_message"`
	EventState        types.String    `tfsdk:"event_state"`
	StartDt           types.String    `tfsdk:"start_dt"`
	EndDt             types.String    `tfsdk:"end_dt"`
	DurationSecond    types.Int64     `tfsdk:"duration_second"`
	EventPolicyId     types.Int64     `tfsdk:"event_policy_id"`
	ProductResourceId types.String    `tfsdk:"product_resource_id"`
	MetricKey         types.String    `tfsdk:"metric_key"`
	ObjectDisplayName types.String    `tfsdk:"object_display_name"`
	ObjectType        types.String    `tfsdk:"object_type"`
	MetricName        types.String    `tfsdk:"metric_name"`
	ObjectName        types.String    `tfsdk:"object_name"`
	ObjectTypeName    types.String    `tfsdk:"object_type_name"`
	ProductIpAddress  types.String    `tfsdk:"product_ip_address"`
	ProductName       types.String    `tfsdk:"product_name"`
	ProductTypeCode   types.String    `tfsdk:"product_type_code"`

	//EventDetail  types.Object    `tfsdk:"event_detail"`
	//EventLevel        types.String `tfsdk:"event_level"`
	//EventMessage      types.String `tfsdk:"event_message"`
	//ProductSummary    types.Object `tfsdk:"product_summary"`
	//MetricSummary     types.Object `tfsdk:"metric_summary"`
	////EventPolicySummary types.Object    `tfsdk:"event_policy_summary"`

}

type EventDataSourceIds struct { //	/v2/events API ERequest
	ResourceType      types.String    `tfsdk:"x_resource_type"`
	ProductResourceId types.String    `tfsdk:"product_resource_id"`
	EventState        types.String    `tfsdk:"event_state"`
	QueryStartDt      types.String    `tfsdk:"query_start_dt"`
	QueryEndDt        types.String    `tfsdk:"query_end_dt"`
	Filter            []filter.Filter `tfsdk:"filter"`
	Events            []Event         `tfsdk:"events"`
}

type EventAccountDataSourceIds struct { //	/v2/events API ERequest
	EventState   types.String    `tfsdk:"event_state"`
	QueryStartDt types.String    `tfsdk:"query_start_dt"`
	QueryEndDt   types.String    `tfsdk:"query_end_dt"`
	Filter       []filter.Filter `tfsdk:"filter"`
	Events       []Event         `tfsdk:"events"`
}

type EventNotificationStateDataSourceIds struct { //	/v2/events/{eventId}/notification-states API ERequest
	ResourceType            types.String             `tfsdk:"x_resource_type"`
	EventId                 types.String             `tfsdk:"event_id"`
	Filter                  []filter.Filter          `tfsdk:"filter"`
	EventNotificationStates []EventNotificationState `tfsdk:"event_notification_states"`
}

type EventNotificationState struct {
	//Description types.String `tfsdk:"description"`
	NotificationStatus   []NotificationStatus `tfsdk:"notification_status"`
	User                 User                 `tfsdk:"user"`
	UserEmail            types.String         `tfsdk:"user_email"`
	UserId               types.String         `tfsdk:"user_id"`
	UserMobileTelNo      types.String         `tfsdk:"user_mobile_tel_no"`
	UserNameNotification types.String         `tfsdk:"user_name_notification"`
}

// ---------- Event Policy(DataSource) ------------ //

type EventPolicyDataSource struct { //	/v2/event-policies/{eventPolicyId} API Request
	ResourceType      types.String `tfsdk:"x_resource_type"`
	EventPolicyId     types.Int64  `tfsdk:"event_policy_id"`
	EventPolicyDetail types.Object `tfsdk:"event_policy_detail"`
}

type EventPolicyDataSourceIds struct { //	/v2/event-policies API Request
	ResourceType         types.String         `tfsdk:"x_resource_type"`
	ProductResourceId    types.String         `tfsdk:"product_resource_id"`
	MetricKey            types.String         `tfsdk:"metric_key"`
	Page                 types.Int32          `tfsdk:"page"`
	Size                 types.Int32          `tfsdk:"size"`
	EventPolicySummaries []EventPolicySummary `tfsdk:"event_policy_summaries"`
}

type EventPolicyHistoryDataSourceIds struct { //	/v2/event-policies/{eventId}/histories API ERequest
	ResourceType         types.String         `tfsdk:"x_resource_type"`
	EventPolicyId        types.Int64          `tfsdk:"event_policy_id"`
	QueryStartDt         types.String         `tfsdk:"query_start_dt"`
	QueryEndDt           types.String         `tfsdk:"query_end_dt"`
	Filter               []filter.Filter      `tfsdk:"filter"`
	EventPolicyHistories []EventPolicyHistory `tfsdk:"event_policy_histories"`
}

type EventPolicyNotificationDataSourceIds struct { //	/v2/event-policies/{eventId}/notifications API ERequest
	ResourceType             types.String              `tfsdk:"x_resource_type"`
	EventPolicyId            types.Int64               `tfsdk:"event_policy_id"`
	EventPolicyNotifications []EventPolicyNotification `tfsdk:"event_policy_notifications"`
}

// ---------- Event Policy (Resource) ------------ //

type EventPolicyResource struct { //	/v2/event-policies API POST
	EventPolicyId      types.Int64  `tfsdk:"event_policy_id"`
	ResourceType       types.String `tfsdk:"x_resource_type"`
	ProductResourceId  types.String `tfsdk:"product_resource_id"`
	EventLevel         types.String `tfsdk:"event_level"`
	DisableYn          types.String `tfsdk:"disable_yn"`
	EventMessagePrefix types.String `tfsdk:"event_message_prefix"`
	FtCount            types.Int64  `tfsdk:"ft_count"`
	IsLogMetric        types.String `tfsdk:"is_log_metric"`
	MetricKey          types.String `tfsdk:"metric_key"`
	ObjectName         types.String `tfsdk:"object_name"`
	ObjectDisplayName  types.String `tfsdk:"object_display_name"`
	EventOccurTimeZone types.String `tfsdk:"event_occur_time_zone"`
	ObjectType         types.String `tfsdk:"object_type"`
	ObjectTypeName     types.String `tfsdk:"object_type_name"`
	//EventThreshold        EventThreshold          `tfsdk:"event_threshold"`
	EventPolicy           EventPolicy `tfsdk:"event_policy"`
	NotificationRecipient types.List  `tfsdk:"notification_recipient"`
}

// ---------- Event & Event Policy Response ------------ //

type EventPolicyHistory struct {
	ModifiedDt     types.String `tfsdk:"modified_dt"`      //	map[string]
	ModifiedBy     types.String `tfsdk:"modified_by"`      //	map[string]
	ModifiedByName types.String `tfsdk:"modified_by_name"` //	map[string]
	CreateById     types.String `tfsdk:"create_by_id"`     //	map[string]
	//	CreateBy              CreateBy              `tfsdk:"create_by"`      //	map[string]
	UpdateById types.String `tfsdk:"update_by_id"` //	map[string]
	//	UpdateBy              CreateBy              `tfsdk:"update_by"`      //	map[string]
	EventPolicyHistoryId   types.Int64  `tfsdk:"event_policy_history_id"`
	EventPolicyHistoryType types.String `tfsdk:"event_policy_history_type"`
	EventPolicyId          types.Int64  `tfsdk:"event_policy_id"`
	ProductResourceId      types.String `tfsdk:"product_resource_id"`
	ProductName            types.String `tfsdk:"product_name"`
	MetricKey              types.String `tfsdk:"metric_key"`
	MetricName             types.String `tfsdk:"metric_name"`
	MetricDescription      types.String `tfsdk:"metric_description"`
	MetricDescriptionEn    types.String `tfsdk:"metric_description_en"`
	MetricUnit             types.String `tfsdk:"metric_unit"`
	ProductTargetType      types.String `tfsdk:"product_target_type"`
	ProductTargetTypeEn    types.String `tfsdk:"product_target_type_en"`
	ObjectName             types.String `tfsdk:"object_name"`
	ObjectDisplayName      types.String `tfsdk:"object_display_name"`
	ObjectType             types.String `tfsdk:"object_type"`
	EventLevel             types.String `tfsdk:"event_level"`
	FtCount                types.Int64  `tfsdk:"ft_count"`
	EventMessagePrefix     types.String `tfsdk:"event_message_prefix"`
	DisableObject          types.String `tfsdk:"disable_object"`
	DisableYn              types.String `tfsdk:"disable_yn"`
	EventOccurTimeZone     types.String `tfsdk:"event_occur_time_zone"`
	//EventThreshold         EventThreshold        `tfsdk:"event_threshold"`
	EventPolicyStatistics EventPolicyStatistics `tfsdk:"event_policy_statistics"`
}

type EventPolicyStatistics struct {
	EventPolicyStatisticsType   types.String `tfsdk:"event_policy_statistics_type"`
	EventPolicyStatisticsPeriod types.Int64  `tfsdk:"event_policy_statistics_period"`
}

type EventThreshold struct {
	ThresholdType   types.String    `tfsdk:"threshold_type"`
	MetricFunction  types.String    `tfsdk:"metric_function"`
	SingleThreshold SingleThreshold `tfsdk:"single_threshold"`
	//RangeThreshold  RangeThreshold  `tfsdk:"range_threshold"`
}

type SingleThreshold struct {
	ComparisonOperator types.String  `tfsdk:"comparison_operator"`
	Value              types.Float64 `tfsdk:"value"`
}

type RangeThreshold struct {
	MaxComparisonOperator types.String  `tfsdk:"max_comparison_operator"`
	MinComparisonOperator types.String  `tfsdk:"min_comparison_operator"`
	MaxValue              types.Float64 `tfsdk:"max_value"`
	MinValue              types.Float64 `tfsdk:"min_value"`
}

type Event struct {
	DurationSecond    types.Int64  `tfsdk:"duration_second"`
	EndDt             types.String `tfsdk:"end_dt"`
	EventId           types.String `tfsdk:"event_id"`
	EventLevel        types.String `tfsdk:"event_level"`
	EventMessage      types.String `tfsdk:"event_message"`
	EventPolicyId     types.Int64  `tfsdk:"event_policy_id"`
	EventState        types.String `tfsdk:"event_state"`
	MetricKey         types.String `tfsdk:"metric_key"`
	MetricName        types.String `tfsdk:"metric_name"`
	ObjectDisplayName types.String `tfsdk:"object_display_name"`
	ObjectName        types.String `tfsdk:"object_name"`
	ObjectType        types.String `tfsdk:"object_type"`
	ObjectTypeName    types.String `tfsdk:"object_type_name"`
	ProductResourceId types.String `tfsdk:"product_resource_id"`
	ProductTypeCode   types.String `tfsdk:"product_type_code"`
	StartDt           types.String `tfsdk:"start_dt"`
}

type User struct {
	Email    types.String `tfsdk:"email"`
	Id       types.String `tfsdk:"id"`
	UserName types.String `tfsdk:"user_name"`
}
type NotificationStatus struct {
	ErrorMsg   types.String `tfsdk:"error_msg"`
	SendDt     types.String `tfsdk:"send_dt"`
	SendMethod types.String `tfsdk:"send_method"`
	SendResult types.String `tfsdk:"send_result"`
}

func (m Event) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"durationSecond":    types.Int64Type,
		"endDt":             types.StringType,
		"eventId":           types.StringType,
		"eventLevel":        types.StringType,
		"eventMessage":      types.StringType,
		"eventPolicyId":     types.StringType,
		"eventState":        types.StringType,
		"metricKey":         types.StringType,
		"metricName":        types.StringType,
		"objectDisplayName": types.StringType,
		"objectName":        types.StringType,
		"objectType":        types.StringType,
		"objectTypeName":    types.StringType,
		"productIpAddress":  types.StringType,
		"productName":       types.StringType,
		"productResourceId": types.StringType,
		"productTypeCode":   types.StringType,
		"startDt":           types.StringType,
	}
}

type EventPolicySummary struct {
	CreateById types.String `tfsdk:"create_by_id"` //	map[string]
	UpdateById types.String `tfsdk:"update_by_id"` //	map[string]
	//UpdateBy              CreateBy              `tfsdk:"update_by"`    //	map[string]
	CreatedDt      types.String `tfsdk:"created_dt"`  //	map[string]
	ModifiedDt     types.String `tfsdk:"modified_dt"` //	map[string]
	ModifiedBy     types.String `tfsdk:"modified_by"` //	map[string]
	CreatedByName  types.String `tfsdk:"created_by_name"`
	CreatedBy      types.String `tfsdk:"created_by"` //	map[string]
	ModifiedByName types.String `tfsdk:"modified_by_name"`
	CreatedById    types.String `tfsdk:"created_by_id"` //	map[string]
	//CreateBy              CreateBy              `tfsdk:"create_by"`     //	map[string]
	UpdatedById           types.String          `tfsdk:"updated_by_id"` //	map[string]
	EventPolicyId         types.Int64           `tfsdk:"event_policy_id"`
	ProductResourceId     types.String          `tfsdk:"product_resource_id"`
	ProductSq             types.Int64           `tfsdk:"product_sq"`
	ProductName           types.String          `tfsdk:"product_name"`
	MetricKey             types.String          `tfsdk:"metric_key"`
	MetricName            types.String          `tfsdk:"metric_name"`
	MetricDescription     types.String          `tfsdk:"metric_description"`
	MetricDescriptionEn   types.String          `tfsdk:"metric_description_en"`
	ProductTargetType     types.String          `tfsdk:"product_target_type"`
	ProductTargetTypeEn   types.String          `tfsdk:"product_target_type_en"`
	IsLogMetric           types.String          `tfsdk:"is_log_metric"`
	ObjectName            types.String          `tfsdk:"object_name"`
	EventLevel            types.String          `tfsdk:"event_level"`
	FtCount               types.Int64           `tfsdk:"ft_count"`
	EventMessagePrefix    types.String          `tfsdk:"event_message_prefix"`
	CheckAsg              types.Bool            `tfsdk:"check_asg"`
	EventOccurTimeZone    types.String          `tfsdk:"event_occur_time_zone"`
	EventThreshold        EventThreshold        `tfsdk:"event_threshold"`
	EventPolicyStatistics EventPolicyStatistics `tfsdk:"event_policy_statistics"`
}

type EventPolicy struct {
	ProductResourceId types.String `tfsdk:"product_resource_id"`
	CreatedByName     types.String `tfsdk:"created_by_name"`
	DisplayEventRule  types.String `tfsdk:"display_event_rule"`
	//GroupSummary        GroupSummary   `tfsdk:"group_summary"`
	ProductSq           types.Int64  `tfsdk:"product_sq"`
	ProductTargetType   types.String `tfsdk:"product_target_type"`
	ProductTargetTypeEn types.String `tfsdk:"product_target_type_en"`
	MetricDescription   types.String `tfsdk:"metric_description"`
	MetricDescriptionEn types.String `tfsdk:"metric_description_en"`
	MetricKey           types.String `tfsdk:"metric_key"`
	MetricName          types.String `tfsdk:"metric_name"`
	ModifiedByName      types.String `tfsdk:"modified_by_name"`
	CreatedDt           types.String `tfsdk:"created_dt"`  //	map[string]
	ModifiedDt          types.String `tfsdk:"modified_dt"` //	map[string]
	CreatedBy           types.String `tfsdk:"created_by"`  //	map[string]
	//CreatedName         types.String   `tfsdk:"created_name"`  //	map[string]
	ModifiedBy types.String `tfsdk:"modified_by"` //	map[string]
	//ModifiedName        types.String   `tfsdk:"modified_name"` //	map[string]
	CreateById types.String `tfsdk:"created_by_id"` //	map[string]
	//CreateBy            CreateBy       `tfsdk:"create_by"`     //	map[string]
	//UpdateById          types.String   `tfsdk:"updated_by_id"` //	map[string]
	//UpdateBy            CreateBy       `tfsdk:"update_by"`     //	map[string]
	EventPolicyId types.Int64 `tfsdk:"event_policy_id"`
	//ProductSummary types.Object `tfsdk:"product_summary"`
	//MetricSummary      MetricSummary  `tfsdk:"metric_summary"`
	MetricType         types.String   `tfsdk:"metric_type"`
	MetricUnit         types.String   `tfsdk:"metric_unit"`
	IsLogMetric        types.String   `tfsdk:"is_log_metric"`
	ObjectName         types.String   `tfsdk:"object_name"`
	ObjectDisplayName  types.String   `tfsdk:"object_display_name"`
	EventLevel         types.String   `tfsdk:"event_level"`
	FtCount            types.Int64    `tfsdk:"ft_count"`
	EventMessagePrefix types.String   `tfsdk:"event_message_prefix"`
	ObjectType         types.String   `tfsdk:"object_type"`
	ObjectTypeName     types.String   `tfsdk:"object_type_name"`
	ProductInfoAttrs   types.String   `tfsdk:"product_info_attrs"`
	DisableObject      types.String   `tfsdk:"disable_object"`
	UserNames          types.String   `tfsdk:"user_names"`
	UserNameStr        types.String   `tfsdk:"user_name_str"`
	DisableYn          types.String   `tfsdk:"disable_yn"`
	AttrListStr        types.String   `tfsdk:"attr_list_str"`
	AsgYn              types.String   `tfsdk:"asg_yn"`
	StartDt            types.String   `tfsdk:"start_dt"`
	CheckAsg           types.Bool     `tfsdk:"check_asg"`
	EventOccurTimeZone types.String   `tfsdk:"event_occur_time_zone"`
	EventThreshold     EventThreshold `tfsdk:"event_threshold"`
	//EventPolicyStatistics EventPolicyStatistics `tfsdk:"event_policy_statistics"`
}

type EventPolicyDetail struct {
	CreateById types.String `tfsdk:"create_by_id"` //	map[string]
	UpdateById types.String `tfsdk:"update_by_id"` //	map[string]
	//UpdateBy              CreateBy              `tfsdk:"update_by"`    //	map[string]
	CreatedDt      types.String `tfsdk:"created_dt"`  //	map[string]
	ModifiedDt     types.String `tfsdk:"modified_dt"` //	map[string]
	ModifiedBy     types.String `tfsdk:"modified_by"` //	map[string]
	CreatedByName  types.String `tfsdk:"created_by_name"`
	CreatedBy      types.String `tfsdk:"created_by"` //	map[string]
	ModifiedByName types.String `tfsdk:"modified_by_name"`
	CreatedById    types.String `tfsdk:"created_by_id"` //	map[string]
	//CreatedDt      time.Time `json:"createdDt"`
	//ModifiedDt     time.Time `json:"modifiedDt"`
	//CreatedBy      string    `json:"createdBy"`
	//CreatedByName types.String `tfsdk:"created_by_name"`
	//ModifiedBy     string    `json:"modifiedBy"`
	//ModifiedByName string    `json:"modifiedByName"`
	//CreateById     string    `json:"createById"`
	//CreateBy              CreateBy              `tfsdk:"create_by"`     //	map[string]
	//UpdateById types.String `tfsdk:"updated_by_id"` //	map[string]
	//UpdateBy              CreateBy              `tfsdk:"update_by"`     //	map[string]
	EventPolicyId       types.Int64    `tfsdk:"event_policy_id"`
	ProductSummary      ProductSummary `tfsdk:"product_summary"`
	MetricSummary       MetricSummary  `tfsdk:"metric_summary"`
	ProductSq           types.Int64    `tfsdk:"product_sq"`
	ProductResourceId   types.String   `tfsdk:"product_resource_id"`
	MetricKey           types.String   `tfsdk:"metric_key"`
	MetricDescription   types.String   `tfsdk:"metric_description"`
	MetricDescriptionEn types.String   `tfsdk:"metric_description_en"`
	ProductTargetType   types.String   `tfsdk:"product_target_type"`
	ProductTargetTypeEn types.String   `tfsdk:"product_target_type_en"`
	IsLogMetric         types.String   `tfsdk:"is_log_metric"`
	ObjectName          types.String   `tfsdk:"object_name"`
	ObjectDisplayName   types.String   `tfsdk:"object_display_name"`
	EventLevel          types.String   `tfsdk:"event_level"`
	FtCount             types.Int64    `tfsdk:"ft_count"`
	EventMessagePrefix  types.String   `tfsdk:"event_message_prefix"`
	ObjectType          types.String   `tfsdk:"object_type"`
	ObjectTypeName      types.String   `tfsdk:"object_type_name"`
	ProductInfoAttrs    types.String   `tfsdk:"product_info_attrs"`
	DisableObject       types.String   `tfsdk:"disable_object"`
	UserNames           types.String   `tfsdk:"user_names"`
	UserNameStr         types.String   `tfsdk:"user_name_str"`
	DisableYn           types.String   `tfsdk:"disable_yn"`
	AttrListStr         types.String   `tfsdk:"attr_list_str"`
	AsgYn               types.String   `tfsdk:"asg_yn"`
	StartDt             types.String   `tfsdk:"start_dt"`
	DisplayEventRule    types.String   `tfsdk:"display_event_rule"`
	CheckAsg            types.Bool     `tfsdk:"check_asg"`
	EventOccurTimeZone  types.String   `tfsdk:"event_occur_time_zone"`
	EventThreshold      EventThreshold `tfsdk:"event_threshold"`
	//EventPolicyStatistics EventPolicyStatistics `tfsdk:"event_policy_statistics"`
}

func (d EventPolicyDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create_by_id": types.StringType,
		"update_by_id": types.StringType,
		//"update_by": types.StringType,
		"created_dt":       types.StringType,
		"modified_dt":      types.StringType,
		"modified_by":      types.StringType,
		"created_by_name":  types.StringType,
		"created_by":       types.StringType,
		"modified_by_name": types.StringType,
		"created_by_id":    types.StringType,
		//CreateBy              CreateBy              `tfsdk:"create_by"`     //	map[string]
		"event_policy_id": types.Int64Type,
		"product_summary": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"product_sq":          types.Int64Type,
				"product_resource_id": types.StringType,
				"product_name":        types.StringType,
				"product_type_code":   types.StringType,
				"product_type_name":   types.StringType,
				"product_ip_address":  types.StringType,
				"product_state":       types.StringType,
				"agent_state":         types.StringType,
				"product_sub_name":    types.StringType,
				"product_sub_type":    types.StringType,
				"lb_name":             types.StringType,
				"vpc_name":            types.StringType,
				"lb_size":             types.StringType,
			},
		},
		"metric_summary": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_key":             types.StringType,
				"metric_name":            types.StringType,
				"metric_set_key":         types.StringType,
				"metric_set_name":        types.StringType,
				"metric_description":     types.StringType,
				"metric_description_en":  types.StringType,
				"product_target_type":    types.StringType,
				"product_target_type_en": types.StringType,
				"metric_type":            types.StringType,
				"metric_unit":            types.StringType,
				"is_object_exist":        types.StringType,
				"is_log_metric":          types.StringType,
			},
		},
		"product_sq":             types.Int64Type,
		"product_resource_id":    types.StringType,
		"metric_key":             types.StringType,
		"metric_description":     types.StringType,
		"metric_description_en":  types.StringType,
		"product_target_type":    types.StringType,
		"product_target_type_en": types.StringType,
		"is_log_metric":          types.StringType,
		"object_name":            types.StringType,
		"object_display_name":    types.StringType,
		"event_level":            types.StringType,
		"ft_count":               types.Int64Type,
		"event_message_prefix":   types.StringType,
		"object_type":            types.StringType,
		"object_type_name":       types.StringType,
		"product_info_attrs":     types.StringType,
		"disable_object":         types.StringType,
		"user_names":             types.StringType,
		"user_name_str":          types.StringType,
		"disable_yn":             types.StringType,
		"attr_list_str":          types.StringType,
		"asg_yn":                 types.StringType,
		"start_dt":               types.StringType,
		"display_event_rule":     types.StringType,
		"check_asg":              types.BoolType,
		"event_occur_time_zone":  types.StringType,
		"event_threshold": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"threshold_type":  types.StringType,
				"metric_function": types.StringType,
				"single_threshold": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"comparison_operator": types.StringType,
						"value":               types.Float64Type,
					},
				},
				//"range_threshold": types.ObjectType{
				//	AttrTypes: map[string]attr.Type{
				//		"max_comparison_operator": types.StringType,
				//		"min_comparison_operator": types.StringType,
				//		"max_value":               types.Float64Type,
				//		"min_value":               types.Float64Type,
				//	},
				//},
			},
		},
		//"event_policy_statistics": types.ObjectType{
		//	AttrTypes: map[string]attr.Type{
		//		"event_policy_statistics_type":   types.StringType,
		//		"event_policy_statistics_period": types.Int64Type,
		//	},
		//},
	}
}

func (m EventPolicy) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		//"created_by_name":    types.StringType,
		//"display_event_rule": types.StringType,
		////"group_summary":         types.ObjectType{},
		//"metric_description":    types.StringType,
		//"metric_description_en": types.StringType,
		"metric_key":  types.StringType,
		"metric_name": types.StringType,
		//"modified_by_name":      types.StringType,
		//"created_dt":            types.StringType,
		//"modified_dt":           types.StringType,
		//"created_by":            types.StringType,
		//"created_name":          types.StringType,
		//"modified_by":           types.StringType,
		//"modified_name":         types.StringType,
		//"create_by_id":          types.StringType,
		////"create_by":               types.ObjectType{},
		//"update_by_id": types.StringType,
		////"update_by":               types.ObjectType{},
		//"event_policy_id": types.Int64Type,
		////"product_summary":         types.ObjectType{},
		////"metric_summary":          types.ObjectType{},
		//"metric_type":          types.StringType,
		//"metric_unit":          types.StringType,
		//"is_log_metric":        types.StringType,
		//"object_name":          types.StringType,
		//"object_display_name":  types.StringType,
		//"event_level":          types.StringType,
		//"ft_count":             types.Int64Type,
		//"event_message_prefix": types.StringType,
		//"object_type":          types.StringType,
		//"object_type_name":     types.StringType,
		//"product_info_attrs":   types.StringType,
		//"disable_object":       types.StringType,
		//"user_names":           types.StringType,
		//"user_name_str":        types.StringType,
		//"disable_yn":           types.StringType,
		//"attr_list_str":        types.StringType,
		//"asg_yn":               types.StringType,
		//"start_dt":             types.StringType,
		////"check_asg":             types.BoolType,
		//"event_occur_time_zone": types.StringType,
		////"event_threshold":         types.ObjectType{},
		////"event_policy_statistics": types.ObjectType{},
	}
}

type GroupSummary struct {
}

type ProductSummary struct {
	ProductSq         types.Int64  `tfsdk:"product_sq"`
	ProductResourceId types.String `tfsdk:"product_resource_id"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductTypeCode   types.String `tfsdk:"product_type_code"`
	ProductTypeName   types.String `tfsdk:"product_type_name"`
	ProductIpAddress  types.String `tfsdk:"product_ip_address"`
	ProductState      types.String `tfsdk:"product_state"`
	AgentState        types.String `tfsdk:"agent_state"`
	ProductSubName    types.String `tfsdk:"product_sub_name"`
	ProductSubType    types.String `tfsdk:"product_sub_type"`
	LbName            types.String `tfsdk:"lb_name"`
	VpcName           types.String `tfsdk:"vpc_name"`
	LbSize            types.String `tfsdk:"lb_size"`
}

type MetricSummary struct {
	MetricKey           types.String `tfsdk:"metric_key"`
	MetricName          types.String `tfsdk:"metric_name"`
	MetricSetKey        types.String `tfsdk:"metric_set_key"`
	MetricSetName       types.String `tfsdk:"metric_set_name"`
	MetricDescription   types.String `tfsdk:"metric_description"`
	MetricDescriptionEn types.String `tfsdk:"metric_description_en"`
	ProductTargetType   types.String `tfsdk:"product_target_type"`
	ProductTargetTypeEn types.String `tfsdk:"product_target_type_en"`
	MetricType          types.String `tfsdk:"metric_type"`
	MetricUnit          types.String `tfsdk:"metric_unit"`
	IsObjectExist       types.String `tfsdk:"is_object_exist"`
	IsLogMetric         types.String `tfsdk:"is_log_metric"`
}

type EventPolicyInfo struct {
	MetricKey             types.String   `tfsdk:"metric_key"`
	IsLogMetric           types.String   `tfsdk:"is_log_metric"`
	ObjectName            types.String   `tfsdk:"object_name"`
	ObjectDisplayName     types.String   `tfsdk:"object_display_name"`
	PodObjectName         types.String   `tfsdk:"pod_object_name"`
	PodObjectDisplayName  types.String   `tfsdk:"pod_object_display_name"`
	EventLevel            types.String   `tfsdk:"event_level"`
	EventThreshold        EventThreshold `tfsdk:"event_threshold"`
	EventPolicyStatistics types.Object   `tfsdk:"event_policy_statistics"`
	FtCount               types.Int64    `tfsdk:"ft_count"`
	EventMessagePrefix    types.String   `tfsdk:"event_message_prefix"`
	ObjectType            types.String   `tfsdk:"object_type"`
	ObjectTypeName        types.String   `tfsdk:"object_type_name"`
	DisableYn             types.String   `tfsdk:"disable_yn"`
	EventRule             types.String   `tfsdk:"event_rule"`
	MetricName            types.String   `tfsdk:"metric_name"`
	EventOccurTimeZone    types.String   `tfsdk:"event_occur_time_zone"`
}

type EventPolicyNotification struct {
	RecipientType       types.String         `tfsdk:"recipient_type"`
	RecipientKey        types.String         `tfsdk:"recipient_key"`
	RecipientName       types.String         `tfsdk:"recipient_name"`
	NotificationMethods []NotificationMethod `tfsdk:"notification_methods"`
	UserAdditionalInfo  CreateBy             `tfsdk:"user_additional_info"`
}

type CreateBy struct { //	/v2/event-policies API POST
	CompanyName        types.String `tfsdk:"company_name"`
	Email              types.String `tfsdk:"email"`
	Id                 types.String `tfsdk:"id"`
	LastLoginAt        types.String `tfsdk:"last_login_at"`
	MemberCreatedAt    types.String `tfsdk:"member_created_at"`
	PasswordReuseCount types.Int32  `tfsdk:"password_reuse_count"`
	UserName           types.String `tfsdk:"user_name"`
	TzId               types.String `tfsdk:"tz_id"`
	Timezone           types.String `tfsdk:"timezone"`
	DstOffset          types.String `tfsdk:"dst_offset"`
	UtcOffset          types.String `tfsdk:"utc_offset"`
	AccountId          types.String `tfsdk:"account_id"`
}

type NotificationRecipient struct {
	RecipientType        types.String         `tfsdk:"recipient_type"`
	RecipientKey         types.String         `tfsdk:"recipient_key"`
	NotificationMethod   []NotificationMethod `tfsdk:"notification_method"`
	RecipientTargetTable types.String         `tfsdk:"recipient_target_table"`
	RcvUserName          types.String         `tfsdk:"rcv_user_name"`
}

type NotificationMethod struct {
	EventLevel types.String   `tfsdk:"event_level"`
	SendMethod []types.String `tfsdk:"send_method"`
}

func (m NotificationMethod) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"event_level": types.StringType,
		"send_method": types.ListType{ElemType: types.StringType},
	}
}

func (m NotificationRecipient) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"recipient_type":         types.StringType,
		"recipient_key":          types.StringType,
		"notification_method":    types.ListType{ElemType: types.ObjectType{AttrTypes: NotificationMethod{}.AttributeTypes()}},
		"recipient_target_table": types.StringType,
		"rcv_user_name":          types.StringType,
	}
}

// ListMetrics
// [GET] /v1/cloudmonitorings/product/v2/metrics
// metric 에 대한 복수형 datasource인 경우 'Metric'+'DataSourceIds' 와 같이 네이밍한다.
type MetricDataSourceIds struct {
	// request와 response 필드에 대해 정의하는 곳 -> example/datasource에서 사용하는 필드이므로 참고해서 작성

	// request -> main.tf 의 인풋필드에 해당
	ProductTypeCode types.String `tfsdk:"product_type_code"`
	ObjectType      types.String `tfsdk:"object_type"`

	// 이건 뭘까?
	Filter []filter.Filter `tfsdk:"filter"`

	// response -> main.tf의 "metrics"에 해당
	// 클라이언트 코드에서는 response가 PageResponseMetricInfoDto 이지만, datasource에서는 간단한 자원명으로 입력한다.
	Metrics []Metric `tfsdk:"metrics"`
}

// response 필드에서 사용한 클래스에 대한 생성자
// PageResponseMetricInfoDto 하위 contents에 해당하는 MetricInfoDto의 필드만 정의한다.
type Metric struct {
	DisableObject       types.String `tfsdk:"disable_object"`
	DisplayUnit         types.String `tfsdk:"display_unit"`
	FixedUnit           types.String `tfsdk:"fixed_unit"`
	IsLogMetric         types.String `tfsdk:"is_log_metric"`
	IsObjectExist       types.String `tfsdk:"is_object_exist"`
	MetricDescription   types.String `tfsdk:"metric_description"`
	MetricDescriptionEn types.String `tfsdk:"metric_description_en"`
	MetricKey           types.String `tfsdk:"metric_key"`
	MetricName          types.String `tfsdk:"metric_name"`
	MetricOrder         types.Int32  `tfsdk:"metric_order"`
	MetricSetKey        types.String `tfsdk:"metric_set_key"`
	MetricSetName       types.String `tfsdk:"metric_set_name"`
	MetricType          types.String `tfsdk:"metric_type"`
	MetricUnit          types.String `tfsdk:"metric_unit"`
	ObjectKeyName       types.String `tfsdk:"object_key_name"`
	ObjectType          types.String `tfsdk:"object_type"`
	ObjectTypeNameEng   types.String `tfsdk:"object_type_name_eng"`
	ObjectTypeNameLoc   types.String `tfsdk:"object_type_name_loc"`
	PerfTitle           types.String `tfsdk:"perf_title"`
	ProductTargetType   types.String `tfsdk:"product_target_type"`
	ProductTargetTypeEn types.String `tfsdk:"product_target_type_en"`
	ProductTypeCode     types.String `tfsdk:"product_type_code"`
	ProductTypeName     types.String `tfsdk:"product_type_name"`
}

// ListService
// [GET] /v1/cloudmonitorings/product/v1/product-types
type ProductTypeDataSourceIds struct {
	// request
	ProductCategoryCode types.String `tfsdk:"product_category_code"`

	Filter []filter.Filter `tfsdk:"filter"`

	// response
	ProductTypes []ProductType `tfsdk:"product_types"`
}
type ProductType struct {
	ParentProductTypeName types.String `tfsdk:"parent_product_type_name"`
	ProductTypeCode       types.String `tfsdk:"product_type_code"`
	ProductTypeName       types.String `tfsdk:"product_type_name"`
	StateMetricKey        types.String `tfsdk:"state_metric_key"`
}

// ListAccountResources
// [GET] /v1/cloudmonitorings/product/v2/accounts/products
type AccountProductDataSourceIds struct {
	// request
	ResourceType types.String `tfsdk:"x_resource_type"`

	Filter []filter.Filter `tfsdk:"filter"`

	// response
	AccountProducts []AccountProduct `tfsdk:"account_products"`
}
type AccountProduct struct {
	AccountId         types.String `tfsdk:"account_id"`
	EndDt             types.String `tfsdk:"end_dt"`
	LastEventLevel    types.String `tfsdk:"last_event_level"`
	PoolName          types.String `tfsdk:"pool_name"`
	ProductIpAddress  types.String `tfsdk:"product_ip_address"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductResourceId types.String `tfsdk:"product_resource_id"`
	ProductState      types.String `tfsdk:"product_state"`
	ProductTypeCode   types.String `tfsdk:"product_type_code"`
	ProductTypeName   types.String `tfsdk:"product_type_name"`
	StartDt           types.String `tfsdk:"start_dt"`
}

// ListAccountMember
// [GET] /v1/cloudmonitorings/product/v1/accounts/members
type AccountMemberDataSourceIds struct {
	Filter []filter.Filter `tfsdk:"filter"`

	// response
	AccountMembers []AccountMember `tfsdk:"account_members"`
}
type AccountMember struct {
	CompanyName     types.String `tfsdk:"company_name"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatedByEmail  types.String `tfsdk:"created_by_email"`
	CreatedByName   types.String `tfsdk:"created_by_name"`
	CreatedDt       types.String `tfsdk:"created_dt"`
	Email           types.String `tfsdk:"email"`
	LastAccessDt    types.String `tfsdk:"last_access_dt"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedByEmail types.String `tfsdk:"modified_by_email"`
	ModifiedByName  types.String `tfsdk:"modified_by_name"`
	ModifiedDt      types.String `tfsdk:"modified_dt"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	PositionName    types.String `tfsdk:"position_name"`
	UserGroupCount  types.String `tfsdk:"user_group_count"`
	UserId          types.String `tfsdk:"user_id"`
	UserName        types.String `tfsdk:"user_name"`
}

// ListAddressBooks
// [GET] /v1/cloudmonitorings/product/v2/users/addrbooks
type AddressBookDataSourceIds struct {
	Filter []filter.Filter `tfsdk:"filter"`

	// response
	AddressBooks []AddressBook `tfsdk:"address_books"`
}
type AddressBook struct {
	AddrBookName  types.String `tfsdk:"addr_book_name"`
	AddrbookId    types.String `tfsdk:"addrbook_id"`
	CreatedBy     types.String `tfsdk:"created_by"`
	CreatedByName types.String `tfsdk:"created_by_name"`
	CreatedDt     types.String `tfsdk:"created_dt"`
	MemberCount   types.Int32  `tfsdk:"member_count"`
}

// ListAddressBookMembers
// [GET] /v1/cloudmonitorings/product/v2/addrbooks/{{addrbookId}}/members
type AddressBookMemberDataSourceIds struct {
	// request
	AddrbookId int32 `tfsdk:"addrbook_id"`

	Filter []filter.Filter `tfsdk:"filter"`

	// response
	AddressBookMembers []AddressBookMember `tfsdk:"address_book_members"`
}
type AddressBookMember struct {
	UserEmail       types.String `tfsdk:"user_email"`
	UserId          types.String `tfsdk:"user_id"`
	UserLanguage    types.String `tfsdk:"user_language"`
	UserMobileTelNo types.String `tfsdk:"user_mobile_tel_no"`
	UserName        types.String `tfsdk:"user_name"`
	UserTimezone    types.String `tfsdk:"user_timezone"`
}

// ListMetricPerfData
// [POST] /v1/cloudmonitorings/product/v2/metric-data
type MetricPerfDataDataSourceIds struct {
	// request
	XResourceType types.String `tfsdk:"x_resource_type"` // header
	IgnoreInvalid types.String `tfsdk:"ignore_invalid"`
	//MetricDataConditions []MetricDataConditionOpenAPIV2 `tfsdk:"metric_data_conditions"`
	QueryStartDt types.String `tfsdk:"query_start_dt"`
	QueryEndDt   types.String `tfsdk:"query_end_dt"`
	// 계층이 아닌 구조로 변경
	MetricKey          types.String `tfsdk:"metric_key"`
	ObjectType         types.String `tfsdk:"object_type"`
	StatisticsPeriod   types.Int32  `tfsdk:"statistics_period"`
	StatisticsTypeList []string     `tfsdk:"statistics_type_list"`
	ObjectList         []string     `tfsdk:"object_list"`
	ProductResourceId  types.String `tfsdk:"product_resource_id"`

	Filter []filter.Filter `tfsdk:"filter"`

	// response
	MetricPerfDatas []MetricPerfData `tfsdk:"metric_perf_datas"`
}
type MetricDataConditionOpenAPIV2 struct {
	MetricKey            types.String          `tfsdk:"metric_key"`
	ObjectType           types.String          `tfsdk:"object_type"`
	ProductResourceInfos []ProductResourceInfo `tfsdk:"product_resource_infos"`
	StatisticsPeriod     types.Int32           `tfsdk:"statistics_period"`
	StatisticsTypeList   []string              `tfsdk:"statistics_type_list"`
}
type ProductResourceInfo struct {
	ObjectList        []string     `tfsdk:"object_list"`
	ProductResourceId types.String `tfsdk:"product_resource_id"`
}
type MetricPerfData struct {
	MetricKey         types.String         `tfsdk:"metric_key"`
	MetricName        types.String         `tfsdk:"metric_name"`
	MetricType        types.String         `tfsdk:"metric_type"`
	MetricUnit        types.String         `tfsdk:"metric_unit"`
	ObjectDisplayName types.String         `tfsdk:"object_display_name"`
	ObjectName        types.String         `tfsdk:"object_name"`
	ObjectType        types.String         `tfsdk:"object_type"`
	PerfData          []MetricPerfDataItem `tfsdk:"perf_data"`
	ProductName       types.String         `tfsdk:"product_name"`
	ProductResourceId types.String         `tfsdk:"product_resource_id"`
	StatisticsPeriod  types.Int32          `tfsdk:"statistics_period"`
	StatisticsType    types.String         `tfsdk:"statistics_type"`
}
type MetricPerfDataItem struct {
	Ts    float64 `tfsdk:"ts"`
	Value string  `tfsdk:"value"`
}
