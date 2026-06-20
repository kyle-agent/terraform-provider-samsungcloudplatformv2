package eventstreams

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/eventstreams"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpEventstreams "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/eventstreams/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &eventstreamsClusterResource{}
	_ resource.ResourceWithConfigure   = &eventstreamsClusterResource{}
	_ resource.ResourceWithImportState = &eventstreamsClusterResource{}
)

func NewEventstreamsClusterResource() resource.Resource {
	return &eventstreamsClusterResource{}
}

type eventstreamsClusterResource struct {
	config  *scpsdk.Configuration
	client  *eventstreams.Client
	clients *client.SCPClient
}

func (r *eventstreamsClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventstreams_cluster"
}

func (r *eventstreamsClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "eventstreams",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AkhqEnabled"): schema.BoolAttribute{
				Description: "AHKQ Enabled",
				Required:    true,
			},
			common.ToSnakeCase("AllowableIpAddresses"): schema.SetAttribute{
				Description: "Allowed IP addresses list  \n" +
					"  - example: ['192.168.10.1/32']",
				Required:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("DbaasEngineVersionId"): schema.StringAttribute{
				Description: "DBaaS engine version ID \n" +
					"  - example: '189299a34f464cac94a24f2d8d57afec' (Kafka 3.8.0)",
				Required: true,
			},
			common.ToSnakeCase("IsCombined"): schema.BoolAttribute{
				Description: "ZOOKEEPER,BROKER combined (IsCombined=true), ZOOKEEPER,BROKER seperated (IsCombined=False)",
				Required:    true,
			},
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AkhqId"): schema.StringAttribute{
						Description: "AkhqId",
						Optional:    true,
					},
					common.ToSnakeCase("AkhqPassword"): schema.StringAttribute{
						Description: "Akhq password ",
						Optional:    true,
					},
					common.ToSnakeCase("BrokerPort"): schema.Int32Attribute{
						Description: "Broker port \n" +
							"  - example: 9091 \n",
						Required: true,
					},
					common.ToSnakeCase("BrokerSaslId"): schema.StringAttribute{
						Description: "Broker Sasl ID \n" +
							"  - minLength: 2  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-z]+$ \n",
						Required: true,
					},
					common.ToSnakeCase("BrokerSaslPassword"): schema.StringAttribute{
						Description: "Broker Sasl password \n" +
							"  - minLength: 8  \n" +
							"  - maxLength: 30  \n" +
							"  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30} (\"'제외) \n",
						Required: true,
					},
					common.ToSnakeCase("ZookeeperPort"): schema.Int32Attribute{
						Description: "Zookeeper port \n" +
							"  - example: 2180 \n",
						Required: true,
					},
					common.ToSnakeCase("ZookeeperSaslId"): schema.StringAttribute{
						Description: "Zookeeper Sasl ID \n" +
							"  - minLength: 2  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-z]+$ \n",
						Required: true,
					},
					common.ToSnakeCase("ZookeeperSaslPassword"): schema.StringAttribute{
						Description: "Zookeeper Sasl password \n" +
							"  - minLength: 8  \n" +
							"  - maxLength: 30  \n" +
							"  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30} (\"'제외) \n",
						Required: true,
					},
				},
			},
			common.ToSnakeCase("InstanceGroups"): schema.ListNestedAttribute{
				Description: "Instance groups",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("BlockStorageGroups"): schema.ListNestedAttribute{
							Description: "BlockStorage groups",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("Id"): schema.StringAttribute{
										Description: "Id",
										Computed:    true,
									},
									common.ToSnakeCase("Name"): schema.StringAttribute{
										Description: "Name",
										Computed:    true,
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: "Role type \n" +
											"  - example: 'OS' \n",
										Required: true,
									},
									common.ToSnakeCase("SizeGb"): schema.Int32Attribute{
										Description: "Size in GB \n" +
											"  - example: 104 \n" +
											"  - minLength: 16  \n" +
											"  - maxLength: 5120  \n",
										Required: true,
									},
									common.ToSnakeCase("VolumeType"): schema.StringAttribute{
										Description: "Volume type \n" +
											"  - example: 'SSD' \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("SSD", "SSD_KMS", "HDD", "HDD_KMS"),
										},
									},
								},
							},
						},
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Id",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						common.ToSnakeCase("Instances"): schema.ListNestedAttribute{
							Description: "Instances",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("Name"): schema.StringAttribute{
										Description: "Name",
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("RoleType"): schema.StringAttribute{
										Description: "Role type \n" +
											"  - example: 'ZOOKEEPER_BROKER' \n" +
											"  - pattern: ZOOKEEPER_BROKER / ZOOKEEPER / BROKER / AKHQ \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("ZOOKEEPER_BROKER", "ZOOKEEPER", "BROKER", "AKHQ"),
										},
									},
									common.ToSnakeCase("ServiceIpAddress"): schema.StringAttribute{
										Description: "User subnet IP address",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									common.ToSnakeCase("PublicIpId"): schema.StringAttribute{
										Description: "Public IP ID (Required when NatEnabled=True)",
										Optional:    true,
									},
									//common.ToSnakeCase("PublicIpAddress"): schema.StringAttribute{
									//	Description: "Public IP address",
									//	Computed:    true,
									//	PlanModifiers: []planmodifier.String{
									//		stringplanmodifier.UseStateForUnknown(),
									//	},
									//},
								},
							},
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: "Role type \n" +
								"  - example: 'ZOOKEEPER_BROKER' \n" +
								"  - pattern: ZOOKEEPER_BROKER (IsCombined=True) / ZOOKEEPER, BROKER (IsCombined=False) / AKHQ (optional) \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ZOOKEEPER_BROKER", "ZOOKEEPER", "BROKER", "AKHQ"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: "Server type name \n" +
								"  - example: 'es1v2m4' \n",
							Required: true,
						},
					},
				},
			},
			common.ToSnakeCase("InstanceNamePrefix"): schema.StringAttribute{
				Description: "Instance name prefix \n" +
					"  - example: 'test'  \n" +
					"  - minLength: 3  \n" +
					"  - maxLength: 13  \n" +
					"  - pattern: ^[a-z][a-zA-Z0-9\\-]*$ \n",
				Required: true,
			},
			common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
				Description: "MaintenanceOption",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PeriodHour"): schema.StringAttribute{
						Description: "Period in hours \n" +
							"  - example: 1  \n",
						Optional: true,
					},
					common.ToSnakeCase("StartingDayOfWeek"): schema.StringAttribute{
						Description: "Starting day of week \n" +
							"  - example: 'MON' \n",
						Optional: true,
					},
					common.ToSnakeCase("StartingTime"): schema.StringAttribute{
						Description: "Starting time \n" +
							"  - example: '0000' \n",
						Optional: true,
					},
					common.ToSnakeCase("UseMaintenanceOption"): schema.BoolAttribute{
						Description: "Use maintenance option \n" +
							"  - example: False \n",
						Optional: true,
						Computed: true,
					},
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Cluster name \n" +
					"  - example: 'test'  \n" +
					"  - minLength: 3  \n" +
					"  - maxLength: 20  \n" +
					"  - pattern: ^[a-zA-Z]*$ \n",
				Required: true,
			},
			common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
				Description: "NAT availability \n" +
					"  - example: False \n",
				Required: true,
			},
			common.ToSnakeCase("ServiceState"): schema.StringAttribute{
				Description: "Service state \n" +
					"  - example : 'RUNNING' (Create,Start) / 'STOPPED' (Stop) \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID",
				Required:    true,
			},
			common.ToSnakeCase("Timezone"): schema.StringAttribute{
				Description: "Timezone \n" +
					"  - example: 'Asia/Seoul' \n",
				Required: true,
			},
			common.ToSnakeCase("ServiceWatchLogCollection"): schema.BoolAttribute{
				Description: "ServiceWatchLogCollection",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *eventstreamsClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Eventstreams
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *eventstreamsClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan eventstreams.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// cluster id 반환
	clusterId := data.Resource.Id

	// cluster 조회 func
	getFunc := func(id string) (*scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
		return r.client.GetCluster(ctx, id)
	}

	// wait for 구현
	getData, err := databaseUtils.AsyncRequestPollingWithState(ctx, clusterId, 500, 10*time.Second,
		"ServiceState", "RUNNING", "FAILED", getFunc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Cluster",
			"Could not read Cluster, unexpected error: "+err.Error(),
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "eventstreams", "event-streams", clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	if len(plan.Tags.Elements()) > 0 {
		getTags, err := r.AsyncPollingTags(ctx, clusterId, "eventstreams", "event-streams",
			100, 3*time.Second)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Tag",
				err.Error(),
			)
			return
		}
		tagsMap = getTags
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	//Metadata 처리
	state, err := r.MapGetResponseToState(ctx, getData, plan, tagsMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventstreamsClusterResource) AsyncPollingTags(ctx context.Context, clusterId string, serviceName string,
	resourceType string, maxAttempts int, internal time.Duration) (types.Map, error) {
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tagsMap, err := tag.GetTags(r.clients, serviceName, resourceType, clusterId)

		if err != nil {
			return types.Map{}, fmt.Errorf("attempt %d/%d failed: %w",
				attempt, maxAttempts, err)
		}

		if len(tagsMap.Elements()) > 0 {
			return tagsMap, nil
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return types.Map{}, fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return types.Map{}, fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func (r *eventstreamsClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, plan eventstreams.ClusterResource, tagsMap types.Map) (eventstreams.ClusterResource, error) {

	var allowableIpAddresses types.Set
	if len(resp.AllowableIpAddresses) == 0 {
		allowableIpAddresses, _ = types.SetValue(types.StringType, []attr.Value{})
	} else {
		ipAddresses := make([]attr.Value, len(resp.AllowableIpAddresses))
		for i, ipAddress := range resp.AllowableIpAddresses {
			ipAddresses[i] = types.StringValue(ipAddress)
		}
		allowableIpAddresses, _ = types.SetValue(types.StringType, ipAddresses)
	}

	var initConfigOption = eventstreams.InitConfigOption{
		AkhqId:                plan.InitConfigOption.AkhqId,
		AkhqPassword:          plan.InitConfigOption.AkhqPassword,
		BrokerPort:            types.Int32PointerValue(resp.InitConfigOption.BrokerPort),
		BrokerSaslId:          plan.InitConfigOption.BrokerSaslId,
		BrokerSaslPassword:    plan.InitConfigOption.BrokerSaslPassword,
		ZookeeperPort:         types.Int32PointerValue(resp.InitConfigOption.ZookeeperPort),
		ZookeeperSaslId:       plan.InitConfigOption.ZookeeperSaslId,
		ZookeeperSaslPassword: plan.InitConfigOption.ZookeeperSaslPassword,
	}

	var InstanceGroups []eventstreams.InstanceGroup
	for _, instanceGroup := range plan.InstanceGroups {

		cutIGs, updatedIG := mapInstanceGroup(resp.InstanceGroups, instanceGroup)

		InstanceGroups = append(InstanceGroups, updatedIG)

		resp.InstanceGroups = cutIGs

	}

	var maintenanceOption = eventstreams.MaintenanceOption{}
	if resp.MaintenanceOption.Get() != nil {
		maintenanceOption = eventstreams.MaintenanceOption{
			PeriodHour:           types.StringPointerValue(resp.MaintenanceOption.Get().PeriodHour.Get()),
			StartingDayOfWeek:    types.StringPointerValue((*string)(resp.MaintenanceOption.Get().StartingDayOfWeek.Get())),
			StartingTime:         types.StringPointerValue(resp.MaintenanceOption.Get().StartingTime.Get()),
			UseMaintenanceOption: types.BoolPointerValue(resp.MaintenanceOption.Get().UseMaintenanceOption),
		}
	}

	return eventstreams.ClusterResource{
		Id:                        types.StringValue(resp.Id),
		AkhqEnabled:               plan.AkhqEnabled,
		AllowableIpAddresses:      allowableIpAddresses,
		DbaasEngineVersionId:      plan.DbaasEngineVersionId,
		InitConfigOption:          initConfigOption,
		InstanceGroups:            InstanceGroups,
		InstanceNamePrefix:        plan.InstanceNamePrefix,
		IsCombined:                plan.IsCombined,
		MaintenanceOption:         maintenanceOption,
		Name:                      types.StringValue(resp.Name),
		NatEnabled:                plan.NatEnabled,
		ServiceState:              types.StringValue(string(resp.ServiceState)),
		SubnetId:                  types.StringValue(resp.SubnetId),
		Tags:                      tagsMap,
		Timezone:                  types.StringValue(resp.Timezone),
		ServiceWatchLogCollection: types.BoolPointerValue(resp.ServiceWatchLogCollection),
	}, nil
}

func (r *eventstreamsClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventstreams.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read Cluster name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "eventstreams", "event-streams", state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState, err := r.MapGetResponseToState(ctx, data, state, tagsMap)
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

func (r *eventstreamsClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*eventstreams.UpdateHandler{
		{
			Fields:  []string{"ServiceState"},
			Handler: r.handlerUpdateClusterState,
		},
		{
			Fields:  []string{"AllowableIpAddresses"},
			Handler: r.handlerUpdateClusterAllowableIpAddresses,
		},
		{
			Fields:  []string{"InstanceGroups"},
			Handler: r.handlerUpdateInstanceGroups,
		},
		{
			Fields:  []string{"Tags"},
			Handler: r.handlerUpdateTag,
		},
	}

	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settableFields []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFields = append(settableFields, databaseUtils.SnakeToPascal(attrName))
		}
	}

	changeFields, err := databaseUtils.GetChangedFields(plan, state, settableFields)
	if err != nil {
		return
	}

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "IsCombined", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "VipPublicIpId", "VirtualIpAddress", "ServiceWatchLogCollection"}

	if databaseUtils.IsOverlapFields(immutableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
		)
		return
	}

	// 변경 확인
	for _, h := range handlers {
		if databaseUtils.IsOverlapFields(h.Fields, changeFields) {
			if err := h.Handler(ctx, req, resp); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Cluster",
					"Could not update cluster, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	data, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading cluster",
			"Could not read cluster name "+state.Name.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// read Tag
	tagsMap, err := tag.GetTags(r.clients, "eventstreams", "event-streams", state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	newState, _ := r.MapGetResponseToState(ctx, data, plan, tagsMap)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *eventstreamsClusterResource) getStateTransitions() map[string]map[string]func(ctx context.Context, clusterId string) error {
	transitions := make(map[string]map[string]func(ctx context.Context, clusterId string) error)

	addState := func(from string, to string, callFunc func(ctx context.Context, clusterId string) error) {
		// from map 이 구성 되지 않았을때 초기화
		if transitions[from] == nil {
			transitions[from] = make(map[string]func(ctx context.Context, clusterId string) error)
		}
		transitions[from][to] = callFunc
	}

	// State Transition Map
	addState("STOPPED", "RUNNING", r.client.StartCluster)
	addState("RUNNING", "STOPPED", r.client.StopCluster)

	return transitions
}

func (r *eventstreamsClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	currentState := state.ServiceState.ValueString()
	desiredState := plan.ServiceState.ValueString()

	if currentState == desiredState {
		return nil
	}

	// state에 따라 start, stop 구분
	err := r.getStateTransitions()[currentState][desiredState](ctx, plan.Id.ValueString())
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", desiredState, "ERROR", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventstreamsClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	clusterId := plan.Id.ValueString()

	addedIPs, removedIps := databaseUtils.CompareIPAddresses(state.AllowableIpAddresses, plan.AllowableIpAddresses)

	err := r.client.SetSecurityGroupRules(ctx, clusterId, addedIPs, removedIps)
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", "RUNNING", "FAILED", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventstreamsClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	for i := 0; i < len(plan.InstanceGroups); i++ {
		currentInstanceGroup := state.InstanceGroups[i]
		desiredInstanceGroup := plan.InstanceGroups[i]

		instanceGroupFields := []string{"BlockStorageGroups", "Id", "Instances", "RoleType", "ServerTypeName"}

		changedFields, err := databaseUtils.GetChangedFields(desiredInstanceGroup, currentInstanceGroup, instanceGroupFields)
		if err != nil {
			return err
		}

		immutableFields := []string{"Id", "RoleType"}

		if databaseUtils.IsOverlapFields(immutableFields, changedFields) {
			resp.Diagnostics.AddError(
				"Error Updating Cluster",
				"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
			)
			return nil
		}

		if len(changedFields) > 0 {
			// ServerTypeName Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"ServerTypeName"}) {
				err := r.client.SetServerType(ctx, currentInstanceGroup.Id.ValueString(), desiredInstanceGroup.ServerTypeName.ValueString())
				if err != nil {
					return err
				}
			}

			// BlockStorageGroups Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"BlockStorageGroups"}) {
				if len(currentInstanceGroup.BlockStorageGroups) == len(desiredInstanceGroup.BlockStorageGroups) {
					// Resize Block Storage
					for i := 0; i < len(currentInstanceGroup.BlockStorageGroups); i++ {
						currentBlockStorage := currentInstanceGroup.BlockStorageGroups[i]
						desiredBlockStorage := desiredInstanceGroup.BlockStorageGroups[i]

						bsFields := []string{"Id", "Name", "RoleType", "SizeGb", "VolumeType"}
						changedBsFields, err := databaseUtils.GetChangedFields(currentBlockStorage, desiredBlockStorage, bsFields)
						if err != nil {
							return err
						}

						immutableBsFields := []string{"RoleType", "VolumeType"}

						if databaseUtils.IsOverlapFields(immutableBsFields, changedBsFields) {
							resp.Diagnostics.AddError(
								"Error Updating Cluster",
								"Immutable fields cannot be modified: "+strings.Join(immutableFields, ", "),
							)
							return nil
						}

						if databaseUtils.IsOverlapFields(changedBsFields, []string{"SizeGb"}) {
							//client
							err := r.client.SetBlockStorageSize(ctx, currentBlockStorage.Id.ValueString(), desiredBlockStorage.SizeGb.ValueInt32())
							if err != nil {
								return err
							}
						}
					}
				} else {
					resp.Diagnostics.AddError(
						"Operation not permitted for BLOCK_STORAGE_GROUP type",
						"the evnetstreams product does not support the addition of storage, so the addition of block storage in instances is restricted.",
					)
					return nil
				}
			}

			// Instances Update
			if databaseUtils.IsOverlapFields(changedFields, []string{"Instances"}) {
				// Kibana or DASHBOARDS
				t := currentInstanceGroup.RoleType.ValueString()
				if t == "KIBANA" || t == "DASHBOARDS" {
					resp.Diagnostics.AddError(
						"Invalid Instance Group Type",
						fmt.Sprintf("Instance group of type '%s' does not support  adding instance", t),
					)
					return nil
				}

				currentInstanceLen := len(currentInstanceGroup.Instances)
				desiredInstanceLen := len(desiredInstanceGroup.Instances)

				if desiredInstanceLen > currentInstanceLen {
					instanceCount := int32(desiredInstanceLen - currentInstanceLen)

					var serviceIPAddresses []string

					for _, instance := range desiredInstanceGroup.Instances[currentInstanceLen:] {
						if instance.ServiceIpAddress.IsNull() || instance.ServiceIpAddress.IsUnknown() {
							serviceIPAddresses = []string{}
							break
						}

						ip := instance.ServiceIpAddress.ValueString()
						serviceIPAddresses = append(serviceIPAddresses, ip)
					}

					err := r.client.AddInstances(ctx, state.Id.ValueString(), instanceCount, serviceIPAddresses)
					if err != nil {
						return err
					}
				}
			}

			// wait for 구현
			getFunc := func(id string) (*scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
				return r.client.GetCluster(ctx, id)
			}

			_, err := databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
				"ServiceState", "RUNNING", "ERROR", getFunc)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *eventstreamsClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan eventstreams.ClusterResource
	var state eventstreams.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	// Update
	_, err := tag.UpdateTags(r.clients, "eventstreams", "event-streams", plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		return err
	}

	return nil
}

func (r *eventstreamsClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventstreams.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// cluster id 반환
	clusterId := state.Id.ValueString()

	// Delete cluster
	err := r.client.DeleteCluster(ctx, clusterId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// cluster 조회 func
	getFunc := func(id string) (*scpEventstreams.EventStreamsClusterDetailResponseV1Dot1, error) {
		return r.client.GetCluster(ctx, id)
	}

	// wait for 구현
	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, clusterId, 200, 20*time.Second,
		"ServiceState", "TERMINATED", "FAILED", getFunc)
	if err != nil {
		if err.Error() != "404 Not Found" {
			resp.Diagnostics.AddError(
				"Error reading server",
				"Could not read server, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func mapInstanceGroup(instanceGroups []scpEventstreams.InstanceGroupResponse, def eventstreams.InstanceGroup) ([]scpEventstreams.InstanceGroupResponse, eventstreams.InstanceGroup) {

	for rm, instanceGroup := range instanceGroups {

		if !isEqualInstanceGroup(instanceGroup, def) {
			continue
		}

		var BlockStorage []eventstreams.BlockStorageGroup
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			BlockStorage = append(BlockStorage, eventstreams.BlockStorageGroup{
				Id:         types.StringValue(blockStorage.Id),
				Name:       types.StringValue(blockStorage.Name),
				RoleType:   types.StringValue(string(blockStorage.RoleType)),
				SizeGb:     types.Int32Value(blockStorage.SizeGb),
				VolumeType: types.StringValue(string(blockStorage.VolumeType)),
			})
		}

		var Instance []eventstreams.Instance
		for _, instance := range instanceGroup.Instances {
			Instance = append(Instance, eventstreams.Instance{
				Name:             types.StringValue(instance.Name),
				RoleType:         types.StringValue(string(instance.RoleType)),
				ServiceIpAddress: types.StringPointerValue(instance.ServiceIpAddress.Get()),
				PublicIpId:       types.StringPointerValue(instance.PublicIpId.Get()),
				//PublicIpAddress:  types.StringPointerValue(instance.PublicIpAddress.Get()),
				//ServiceState:     types.StringValue(string(instance.ServiceState)),
			})
		}

		return append(instanceGroups[:rm], instanceGroups[rm+1:]...), eventstreams.InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: BlockStorage,
			Instances:          Instance,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		}

	}

	return instanceGroups, def

}

func isEqualInstanceGroup(actual scpEventstreams.InstanceGroupResponse, expect eventstreams.InstanceGroup) bool {

	equal := string(actual.RoleType) == expect.RoleType.ValueString()
	equal = equal && actual.ServerTypeName == expect.ServerTypeName.ValueString()

	actualIt := actual.GetInstances()
	expectIt := expect.Instances
	equal = equal && len(actualIt) == len(expectIt)
	if equal {
		for pos := range len(expectIt) {
			equal = equal && expectIt[pos].RoleType.ValueString() == string(actualIt[pos].RoleType)
		}
	}

	actualBS := actual.GetBlockStorageGroups()
	expectBS := expect.BlockStorageGroups
	equal = equal && len(actualBS) == len(expectBS)
	if equal {
		for pos := range len(expectBS) {
			equal = equal && expectBS[pos].RoleType.ValueString() == string(actualBS[pos].RoleType)
			equal = equal && expectBS[pos].VolumeType.ValueString() == string(actualBS[pos].VolumeType)
			equal = equal && expectBS[pos].SizeGb.ValueInt32() == actualBS[pos].SizeGb
		}
	}

	return equal

}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *eventstreamsClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
