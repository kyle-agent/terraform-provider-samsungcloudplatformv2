package vertica

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vertica"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	databaseUtils "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/database"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpVertica "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vertica/1.0"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &verticaClusterResource{}
	_ resource.ResourceWithConfigure   = &verticaClusterResource{}
	_ resource.ResourceWithImportState = &verticaClusterResource{}
)

func NewVerticaClusterResource() resource.Resource {
	return &verticaClusterResource{}
}

type verticaClusterResource struct {
	config  *scpsdk.Configuration
	client  *vertica.Client
	clients *client.SCPClient
}

func (r *verticaClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vertica_cluster"
}

func (r *verticaClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "vertica",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AllowableIpAddresses"): schema.SetAttribute{
				Description: "Allowed IP addresses list  \n" +
					"  - example: ['192.168.10.1/32']",
				Required:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("DbaasEngineVersionId"): schema.StringAttribute{
				Description: "DBaaS engine version ID \n" +
					"  - example: '09c2fe88089040ffa035604e38f7e4e9' (Vertica ENTERPRISE 24.2.0-2)",
				Required: true,
			},
			common.ToSnakeCase("NatEnabled"): schema.BoolAttribute{
				Description: "NAT availability \n" +
					"  - example: False \n",
				Required: true,
			},
			common.ToSnakeCase("InitConfigOption"): schema.SingleNestedAttribute{
				Description: "Init config option",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("BackupOption"): schema.SingleNestedAttribute{
						Description: "Backup option",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("RetentionPeriodDay"): schema.StringAttribute{
								Description: "Backup retention period (day) \n" +
									"  - example: 7 \n" +
									"  - min: 7 \n" +
									"  - max: 35 \n",
								Optional: true,
							},
							common.ToSnakeCase("StartingTimeHour"): schema.StringAttribute{
								Description: "Backup starting time (hour) \n" +
									"  - example: 12 \n" +
									"  - min: 00 \n" +
									"  - max: 23 \n",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("DatabaseLocale"): schema.StringAttribute{
						Description: "Database locale information\n" +
							"  - example: 'ko_KR.utf8' \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseName"): schema.StringAttribute{
						Description: "Database name \n" +
							"  - example: 'test' \n" +
							"  - minLength: 3  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-zA-Z][a-zA-Z0-9]*$ \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserName"): schema.StringAttribute{
						Description: "Database user name \n" +
							"  - example: 'test' \n" +
							"  - minLength: 2  \n" +
							"  - maxLength: 20  \n" +
							"  - pattern: ^[a-z]*$ \n",
						Required: true,
					},
					common.ToSnakeCase("DatabaseUserPassword"): schema.StringAttribute{
						Description: "Database user password \n" +
							"  - minLength: 8  \n" +
							"  - maxLength: 30  \n" +
							"  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\\-[\\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\\S*[^\\w\\s]).{8,30} (\"'제외) \n",
						Required: true,
					},
					common.ToSnakeCase("DatabasePort"): schema.Int32Attribute{
						Description: "Database service port",
						Computed:    true,
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
					},
					common.ToSnakeCase("McPort"): schema.Int32Attribute{
						Description: "Mc port",
						Computed:    true,
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
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
											"  - example: 'CONSOLE' \n" +
											"  - pattern: CONSOLE / DATA \n",
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("CONSOLE", "DATA"),
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
								},
							},
						},
						common.ToSnakeCase("RoleType"): schema.StringAttribute{
							Description: "Role type \n" +
								"  - example: 'CONSOLE' \n" +
								"  - pattern: CONSOLE / DATA \n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("CONSOLE", "DATA"),
							},
						},
						common.ToSnakeCase("ServerTypeName"): schema.StringAttribute{
							Description: "Server type name \n" +
								"  - example: 'db1v1m2' \n",
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
			common.ToSnakeCase("License"): schema.StringAttribute{
				Description: "License",
				Required:    true,
			},
			common.ToSnakeCase("MaintenanceOption"): schema.SingleNestedAttribute{
				Description: "Maintenance option",
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
		},
	}
}

func (r *verticaClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vertica
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *verticaClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vertica.ClusterResource
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
	getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
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
	tagsMap, err := tag.GetTags(r.clients, "vertica", "vertica", clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}

	if len(plan.Tags.Elements()) > 0 {
		getTags, err := r.AsyncPollingTags(ctx, clusterId, "vertica", "vertica",
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

func (r *verticaClusterResource) AsyncPollingTags(ctx context.Context, clusterId string, serviceName string,
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

func (r *verticaClusterResource) MapGetResponseToState(ctx context.Context,
	resp *scpVertica.VerticaClusterDetailResponse, plan vertica.ClusterResource, tagsMap types.Map) (vertica.ClusterResource, error) {

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

	var backupOption = vertica.BackupOption{}
	if resp.InitConfigOption.BackupOption.Get() != nil {
		backupOption = vertica.BackupOption{
			RetentionPeriodDay: types.StringValue(resp.InitConfigOption.BackupOption.Get().RetentionPeriodDay),
			StartingTimeHour:   types.StringValue(resp.InitConfigOption.BackupOption.Get().StartingTimeHour),
		}
	}

	var initConfigOption = vertica.InitConfigOption{
		BackupOption:         backupOption,
		DatabaseLocale:       types.StringPointerValue(resp.InitConfigOption.DatabaseLocale.Get()),
		DatabaseName:         types.StringValue(resp.InitConfigOption.DatabaseName),
		DatabasePort:         types.Int32PointerValue(resp.InitConfigOption.DatabasePort.Get()),
		DatabaseUserName:     types.StringValue(resp.InitConfigOption.DatabaseUserName),
		DatabaseUserPassword: plan.InitConfigOption.DatabaseUserPassword,
		McPort:               types.Int32Value(resp.InitConfigOption.GetMcPort()),
	}

	var InstanceGroups []vertica.InstanceGroup
	for _, instanceGroup := range resp.InstanceGroups {
		var BlockStorage []vertica.BlockStorageGroup
		for _, blockStorage := range instanceGroup.BlockStorageGroups {
			BlockStorage = append(BlockStorage, vertica.BlockStorageGroup{
				Id:         types.StringValue(blockStorage.Id),
				Name:       types.StringValue(blockStorage.Name),
				RoleType:   types.StringValue(string(blockStorage.RoleType)),
				SizeGb:     types.Int32Value(blockStorage.SizeGb),
				VolumeType: types.StringValue(string(blockStorage.VolumeType)),
			})
		}

		var Instance []vertica.Instance
		for _, instance := range instanceGroup.Instances {
			Instance = append(Instance, vertica.Instance{
				Name:             types.StringValue(instance.Name),
				RoleType:         types.StringValue(string(instance.RoleType)),
				ServiceIpAddress: types.StringPointerValue(instance.ServiceIpAddress.Get()),
				PublicIpId:       types.StringPointerValue(instance.PublicIpId.Get()),
			})
		}

		InstanceGroups = append(InstanceGroups, vertica.InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: BlockStorage,
			Instances:          Instance,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		})
	}

	var maintenanceOption = vertica.MaintenanceOption{}
	if resp.MaintenanceOption.Get() != nil {
		maintenanceOption = vertica.MaintenanceOption{
			PeriodHour:           types.StringPointerValue(resp.MaintenanceOption.Get().PeriodHour.Get()),
			StartingDayOfWeek:    types.StringPointerValue((*string)(resp.MaintenanceOption.Get().StartingDayOfWeek.Get())),
			StartingTime:         types.StringPointerValue(resp.MaintenanceOption.Get().StartingTime.Get()),
			UseMaintenanceOption: types.BoolPointerValue(resp.MaintenanceOption.Get().UseMaintenanceOption),
		}
	}

	return vertica.ClusterResource{
		Id:                   types.StringValue(resp.Id),
		AllowableIpAddresses: allowableIpAddresses,
		DbaasEngineVersionId: plan.DbaasEngineVersionId,
		InitConfigOption:     initConfigOption,
		InstanceGroups:       InstanceGroups,
		InstanceNamePrefix:   plan.InstanceNamePrefix,
		MaintenanceOption:    maintenanceOption,
		License:              plan.License,
		Name:                 types.StringValue(resp.Name),
		NatEnabled:           plan.NatEnabled,
		ServiceState:         types.StringValue(string(resp.ServiceState)),
		SubnetId:             types.StringValue(resp.SubnetId),
		Tags:                 tagsMap,
		Timezone:             types.StringValue(resp.Timezone),
	}, nil
}

func (r *verticaClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vertica.ClusterResource
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
	tagsMap, err := tag.GetTags(r.clients, "vertica", "vertica", state.Id.ValueString())
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

func (r *verticaClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	handlers := []*vertica.UpdateHandler{
		{
			Fields:  []string{"ServiceState"},
			Handler: r.handlerUpdateClusterState,
		},
		{
			Fields:  []string{"InitConfigOption"},
			Handler: r.handlerUpdateClusterInitConfig,
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

	var plan vertica.ClusterResource
	var state vertica.ClusterResource
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

	immutableFields := []string{"id", "MaintenanceOption", "DbaasEngineVersionId", "NatEnabled", "InstanceNamePrefix", "Name", "SubnetId", "Timezone", "License"}

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
	tagsMap, err := tag.GetTags(r.clients, "vertica", "vertica", state.Id.ValueString())
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

func (r *verticaClusterResource) getStateTransitions() map[string]map[string]func(ctx context.Context, clusterId string) error {
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

func (r *verticaClusterResource) handlerUpdateClusterState(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan vertica.ClusterResource
	var state vertica.ClusterResource
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

	getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", desiredState, "ERROR", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *verticaClusterResource) handlerUpdateClusterInitConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan vertica.ClusterResource
	var state vertica.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	clusterId := plan.Id.ValueString()

	backupState := state.InitConfigOption.BackupOption
	backupPlan := plan.InitConfigOption.BackupOption

	// 1. backup 최초 설정
	if isEmpty(backupState) && !isEmpty(backupPlan) {
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, startingTimeHour, retentionPeriodDay)
		if err != nil {
			return err
		}
	}

	// 2. backup 설정 변경
	if !isEmpty(backupState) && !isEmpty(backupPlan) && !reflect.DeepEqual(backupState, backupPlan) {
		startingTimeHour := backupPlan.StartingTimeHour.ValueString()
		retentionPeriodDay := backupPlan.RetentionPeriodDay.ValueString()

		err := r.client.SetBackup(ctx, clusterId, startingTimeHour, retentionPeriodDay)
		if err != nil {
			return err
		}
	}

	// 3. bacup 설정 삭제
	if !isEmpty(backupState) && isEmpty(backupPlan) {
		err := r.client.UnSetBackup(ctx, clusterId)
		if err != nil {
			return err
		}
	}

	getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err := databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", "RUNNING", "ERROR", getFunc)
	if err != nil {
		return err
	}
	return nil
}

func isEmpty(sp vertica.BackupOption) bool {
	return sp.StartingTimeHour.IsNull() && sp.RetentionPeriodDay.IsNull()
}

func (r *verticaClusterResource) handlerUpdateClusterAllowableIpAddresses(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan vertica.ClusterResource
	var state vertica.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	clusterId := plan.Id.ValueString()

	addedIPs, removedIps := databaseUtils.CompareIPAddresses(state.AllowableIpAddresses, plan.AllowableIpAddresses)

	err := r.client.SetSecurityGroupRules(ctx, clusterId, addedIPs, removedIps)
	if err != nil {
		return err
	}

	getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
		return r.client.GetCluster(ctx, id)
	}

	_, err = databaseUtils.AsyncRequestPollingWithState(ctx, plan.Id.ValueString(), 200, 10*time.Second,
		"ServiceState", "RUNNING", "FAILED", getFunc)
	if err != nil {
		return err
	}

	return nil
}

func (r *verticaClusterResource) handlerUpdateInstanceGroups(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan vertica.ClusterResource
	var state vertica.ClusterResource
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
					// Add Block Storage
					addBlockStorage := desiredInstanceGroup.BlockStorageGroups[len(desiredInstanceGroup.BlockStorageGroups)-1]
					err := r.client.AddBlockStorages(ctx, currentInstanceGroup.Id.ValueString(), addBlockStorage.RoleType.ValueString(), addBlockStorage.SizeGb.ValueInt32(), addBlockStorage.VolumeType.ValueString())
					if err != nil {
						return err
					}
				}
			}

			// wait for 구현
			getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
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

func (r *verticaClusterResource) handlerUpdateTag(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) error {
	var plan vertica.ClusterResource
	var state vertica.ClusterResource
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	// Update
	_, err := tag.UpdateTags(r.clients, "vertica", "vertica", plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		return err
	}

	return nil
}

func (r *verticaClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vertica.ClusterResource
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
	getFunc := func(id string) (*scpVertica.VerticaClusterDetailResponse, error) {
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

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *verticaClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
