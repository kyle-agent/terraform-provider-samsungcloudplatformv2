package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/billing" // client 를 import 한다.
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &billingPlannedComputeResource{}
	_ resource.ResourceWithConfigure = &billingPlannedComputeResource{}
)

func NewBillingPlannedComputeResource() resource.Resource {
	return &billingPlannedComputeResource{}
}

type billingPlannedComputeResource struct {
	config  *scpsdk.Configuration
	client  *billing.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *billingPlannedComputeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_planned_compute" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *billingPlannedComputeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Planned compute",
		Attributes: map[string]schema.Attribute{
			"region": region.ResourceSchema(),
			"tags":   tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "AccountId",
				Optional:    true,
			},
			common.ToSnakeCase("ContractType"): schema.StringAttribute{
				Description: "Contract period code. Allowed API enum values: 01 (1-year), 03 (3-year), 05 (5-year).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("01", "03", "05"),
				},
			},
			common.ToSnakeCase("OsType"): schema.StringAttribute{
				Description: "OsType",
				Optional:    true,
			},
			common.ToSnakeCase("ServerType"): schema.StringAttribute{
				Description: "ServerType",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceId"): schema.StringAttribute{
				Description: "ServiceId",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceName"): schema.StringAttribute{
				Description: "ServiceName",
				Optional:    true,
			},
			common.ToSnakeCase("Action"): schema.StringAttribute{
				Description: "Action",
				Optional:    true,
			},
			common.ToSnakeCase("PlannedCompute"): schema.SingleNestedAttribute{
				Description: "PlannedCompute",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("ContractId"): schema.StringAttribute{
						Description: "Contract ID",
						Computed:    true,
					},
					common.ToSnakeCase("ContractType"): schema.StringAttribute{
						Description: "Contract Type",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created by",
						Computed:    true,
					},
					common.ToSnakeCase("DeleteYn"): schema.StringAttribute{
						Description: "Delete Y/N",
						Computed:    true,
					},
					common.ToSnakeCase("EndDate"): schema.StringAttribute{
						Description: "End date",
						Computed:    true,
					},
					common.ToSnakeCase("FirstContractStartAt"): schema.StringAttribute{
						Description: "First contract start at",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Planned compute ID",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified by",
						Computed:    true,
					},
					common.ToSnakeCase("NextContractType"): schema.StringAttribute{
						Description: "Next contract type",
						Computed:    true,
					},
					common.ToSnakeCase("NextEndDate"): schema.StringAttribute{
						Description: "Next end date",
						Computed:    true,
					},
					common.ToSnakeCase("NextStartDate"): schema.StringAttribute{
						Description: "Next end date",
						Computed:    true,
					},
					common.ToSnakeCase("OsName"): schema.StringAttribute{
						Description: "OS name",
						Computed:    true,
					},
					common.ToSnakeCase("OsType"): schema.StringAttribute{
						Description: "OS type",
						Computed:    true,
					},
					common.ToSnakeCase("Region"): schema.StringAttribute{
						Description: "Region",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceName"): schema.StringAttribute{
						Description: "Resource name",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceType"): schema.StringAttribute{
						Description: "Resource type",
						Computed:    true,
					},
					common.ToSnakeCase("ServerType"): schema.StringAttribute{
						Description: "Server type",
						Computed:    true,
					},
					common.ToSnakeCase("ServerTypeDescription"): schema.MapAttribute{
						Description: "Server type description",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("ServiceId"): schema.StringAttribute{
						Description: "Service ID",
						Computed:    true,
					},
					common.ToSnakeCase("ServiceName"): schema.StringAttribute{
						Description: "Service Name",
						Computed:    true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description: "srn",
						Computed:    true,
					},
					common.ToSnakeCase("StartDate"): schema.StringAttribute{
						Description: "Start date",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *billingPlannedComputeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Billing
	r.clients = inst.Client
}

func convertMapStringInterfaceToTypesMap(data map[string]interface{}) (types.Map, diag.Diagnostics) {
	tmp := make(map[string]attr.Value)
	for key, value := range data {
		strValue := fmt.Sprintf("%v", value)
		tmp[key] = basetypes.NewStringValue(strValue)
	}

	resultTypesMap, diags := types.MapValue(types.StringType, tmp)
	return resultTypesMap, diags
}

// Create creates the resource and sets the initial Terraform state.
func (r *billingPlannedComputeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan billing.PlannedComputeResource
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

	if !plan.Region.IsNull() {
		r.client.Config.Region = plan.Region.ValueString()
	}

	data, err := r.client.CreatePlannedCompute(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating planned compute",
			"Could not create planned compute, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plannedCompute := data.PlannedCompute
	var idString string
	if plannedCompute.Id.IsSet() && plannedCompute.Id.Get() != nil {
		idString = *plannedCompute.Id.Get()
	} else {
		idString = ""
	}
	plan.Id = types.StringValue(idString)

	serverTypeDesc, diags := convertMapStringInterfaceToTypesMap(plannedCompute.GetServerTypeDescription())

	plannedComputeModel := billing.PlannedCompute{
		AccountId:             types.StringPointerValue(plannedCompute.AccountId),
		ContractId:            types.StringPointerValue(plannedCompute.ContractId.Get()),
		ContractType:          types.StringPointerValue(plannedCompute.ContractType),
		CreatedAt:             types.StringValue(plannedCompute.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringPointerValue(plannedCompute.CreatedBy.Get()),
		DeleteYn:              types.StringPointerValue(plannedCompute.DeleteYn.Get()),
		EndDate:               types.StringValue(plannedCompute.GetEndDate()),
		FirstContractStartAt:  types.StringValue(plannedCompute.GetFirstContractStartAt()),
		Id:                    types.StringPointerValue(plannedCompute.Id.Get()),
		ModifiedAt:            types.StringValue(plannedCompute.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:            types.StringPointerValue(plannedCompute.ModifiedBy.Get()),
		NextContractType:      types.StringValue(plannedCompute.GetNextContractType()),
		NextEndDate:           types.StringValue(plannedCompute.GetNextEndDate()),
		NextStartDate:         types.StringValue(plannedCompute.GetNextStartDate()),
		OsName:                types.StringPointerValue(plannedCompute.OsName),
		OsType:                types.StringPointerValue(plannedCompute.OsType),
		Region:                types.StringPointerValue(plannedCompute.Region),
		ResourceName:          types.StringPointerValue(plannedCompute.ResourceName.Get()),
		ResourceType:          types.StringPointerValue(plannedCompute.ResourceType),
		ServerType:            types.StringPointerValue(plannedCompute.ServerType),
		ServerTypeDescription: serverTypeDesc,
		ServiceId:             types.StringPointerValue(plannedCompute.ServiceId),
		ServiceName:           types.StringPointerValue(plannedCompute.ServiceName),
		Srn:                   types.StringPointerValue(plannedCompute.Srn),
		StartDate:             types.StringPointerValue(plannedCompute.StartDate),
		State:                 types.StringPointerValue(plannedCompute.State),
	}
	plannedComputeObjectValue, diags := types.ObjectValueFrom(ctx, plannedComputeModel.AttributeTypes(), plannedComputeModel)
	plan.PlannedCompute = plannedComputeObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *billingPlannedComputeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state billing.PlannedComputeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetPlannedCompute(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			"Could not read Resource Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	plannedCompute := data.PlannedCompute

	serverTypeDesc, diags := convertMapStringInterfaceToTypesMap(plannedCompute.GetServerTypeDescription())

	plannedComputeModel := billing.PlannedCompute{
		AccountId:             types.StringValue(*plannedCompute.AccountId),
		ContractId:            types.StringPointerValue(plannedCompute.ContractId.Get()),
		ContractType:          types.StringPointerValue(plannedCompute.ContractType),
		CreatedAt:             types.StringValue(plannedCompute.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringPointerValue(plannedCompute.CreatedBy.Get()),
		DeleteYn:              types.StringPointerValue(plannedCompute.DeleteYn.Get()),
		EndDate:               types.StringValue(plannedCompute.GetEndDate()),
		FirstContractStartAt:  types.StringValue(plannedCompute.GetFirstContractStartAt()),
		Id:                    types.StringPointerValue(plannedCompute.Id.Get()),
		ModifiedAt:            types.StringValue(plannedCompute.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:            types.StringPointerValue(plannedCompute.ModifiedBy.Get()),
		NextContractType:      types.StringValue(plannedCompute.GetNextContractType()),
		NextEndDate:           types.StringValue(plannedCompute.GetNextEndDate()),
		NextStartDate:         types.StringValue(plannedCompute.GetNextStartDate()),
		OsName:                types.StringPointerValue(plannedCompute.OsName),
		OsType:                types.StringPointerValue(plannedCompute.OsType),
		Region:                types.StringPointerValue(plannedCompute.Region),
		ResourceName:          types.StringPointerValue(plannedCompute.ResourceName.Get()),
		ResourceType:          types.StringPointerValue(plannedCompute.ResourceType),
		ServerType:            types.StringPointerValue(plannedCompute.ServerType),
		ServerTypeDescription: serverTypeDesc,
		ServiceId:             types.StringPointerValue(plannedCompute.ServiceId),
		ServiceName:           types.StringPointerValue(plannedCompute.ServiceName),
		Srn:                   types.StringPointerValue(plannedCompute.Srn),
		StartDate:             types.StringPointerValue(plannedCompute.StartDate),
		State:                 types.StringPointerValue(plannedCompute.State),
	}
	plannedComputeObjectValue, diags := types.ObjectValueFrom(ctx, plannedComputeModel.AttributeTypes(), plannedComputeModel)
	state.PlannedCompute = plannedComputeObjectValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *billingPlannedComputeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state billing.PlannedComputeResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.UpdatePlannedCompute(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Planned Compute",
			"Could not update Planned Compute, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetPlannedCompute(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading plannedCompute",
			"Could not read PlannedCompute ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plannedCompute := data.PlannedCompute
	serverTypeDesc, diags := convertMapStringInterfaceToTypesMap(plannedCompute.GetServerTypeDescription())

	plannedComputeModel := billing.PlannedCompute{
		AccountId:             types.StringValue(*plannedCompute.AccountId),
		ContractId:            types.StringPointerValue(plannedCompute.ContractId.Get()),
		ContractType:          types.StringPointerValue(plannedCompute.ContractType),
		CreatedAt:             types.StringValue(plannedCompute.CreatedAt.Format(time.RFC3339)),
		CreatedBy:             types.StringPointerValue(plannedCompute.CreatedBy.Get()),
		DeleteYn:              types.StringPointerValue(plannedCompute.DeleteYn.Get()),
		EndDate:               types.StringValue(plannedCompute.GetEndDate()),
		FirstContractStartAt:  types.StringValue(plannedCompute.GetFirstContractStartAt()),
		Id:                    types.StringPointerValue(plannedCompute.Id.Get()),
		ModifiedAt:            types.StringValue(plannedCompute.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:            types.StringPointerValue(plannedCompute.ModifiedBy.Get()),
		NextContractType:      types.StringValue(plannedCompute.GetNextContractType()),
		NextEndDate:           types.StringValue(plannedCompute.GetNextEndDate()),
		NextStartDate:         types.StringValue(plannedCompute.GetNextStartDate()),
		OsName:                types.StringPointerValue(plannedCompute.OsName),
		OsType:                types.StringPointerValue(plannedCompute.OsType),
		Region:                types.StringPointerValue(plannedCompute.Region),
		ResourceName:          types.StringPointerValue(plannedCompute.ResourceName.Get()),
		ResourceType:          types.StringPointerValue(plannedCompute.ResourceType),
		ServerType:            types.StringPointerValue(plannedCompute.ServerType),
		ServerTypeDescription: serverTypeDesc,
		ServiceId:             types.StringPointerValue(plannedCompute.ServiceId),
		ServiceName:           types.StringPointerValue(plannedCompute.ServiceName),
		Srn:                   types.StringPointerValue(plannedCompute.Srn),
		StartDate:             types.StringPointerValue(plannedCompute.StartDate),
		State:                 types.StringPointerValue(plannedCompute.State),
	}
	plannedComputeObjectValue, diags := types.ObjectValueFrom(ctx, plannedComputeModel.AttributeTypes(), plannedComputeModel)
	state.PlannedCompute = plannedComputeObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *billingPlannedComputeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state billing.PlannedComputeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
