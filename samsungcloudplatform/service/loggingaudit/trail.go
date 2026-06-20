package loggingaudit

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &loggingauditTrailResource{}
	_ resource.ResourceWithConfigure   = &loggingauditTrailResource{}
	_ resource.ResourceWithImportState = &loggingauditTrailResource{}
)

// NewLoggingauditLogResource is a helper function to simplify the provider implementation.
func NewLoggingauditTrailResource() resource.Resource {
	return &loggingauditTrailResource{}
}

// loggingauditLogResource is the data source implementation.
type loggingauditTrailResource struct {
	config  *scpsdk.Configuration
	client  *loggingaudit.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loggingauditTrailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loggingaudit_trail"
}

func (r *loggingauditTrailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Trail",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the access key.",
				Computed:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "AccountId",
				Required:    true,
			},
			common.ToSnakeCase("BucketName"): schema.StringAttribute{
				Description: "BucketName",
				Optional:    true,
			},
			common.ToSnakeCase("BucketRegion"): schema.StringAttribute{
				Description: "BucketRegion",
				Optional:    true,
			},
			common.ToSnakeCase("LogTypeTotalYn"): schema.StringAttribute{
				Description: "LogTypeTotalYn",
				Required:    true,
			},
			common.ToSnakeCase("LogVerificationYn"): schema.StringAttribute{
				Description: "LogVerificationYn",
				Required:    true,
			},
			common.ToSnakeCase("RegionNames"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "RegionNames",
				Optional:    true,
			},
			common.ToSnakeCase("RegionTotalYn"): schema.StringAttribute{
				Description: "RegionTotalYn",
				Required:    true,
			},
			common.ToSnakeCase("ResourceTypeTotalYn"): schema.StringAttribute{
				Description: "ResourceTypeTotalYn",
				Required:    true,
			},
			common.ToSnakeCase("TagCreateRequests"): schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":   types.StringType,
						"value": types.StringType,
					},
				},
				Description: "TagCreateRequests",
				Optional:    true,
			},
			common.ToSnakeCase("TargetLogTypes"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "TargetLogTypes",
				Optional:    true,
			},
			common.ToSnakeCase("TargetResourceTypes"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "TargetResourceTypes",
				Optional:    true,
			},
			common.ToSnakeCase("TargetUsers"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "TargetUsers",
				Optional:    true,
			},
			common.ToSnakeCase("TrailDescription"): schema.StringAttribute{
				Description: "TrailDescription",
				Required:    true,
			},
			common.ToSnakeCase("TrailName"): schema.StringAttribute{
				Description: "TrailName",
				Optional:    true,
			},
			common.ToSnakeCase("TrailSaveType"): schema.StringAttribute{
				Description: "TrailSaveType",
				Optional:    true,
			},
			common.ToSnakeCase("UserTotalYn"): schema.StringAttribute{
				Description: "UserTotalYn",
				Required:    true,
			},
			common.ToSnakeCase("OrganizationTrailYn"): schema.StringAttribute{
				Description: "OrganizationTrailYn",
				Required:    true,
			},
			common.ToSnakeCase("LogArchiveAccountId"): schema.StringAttribute{
				Description: "LogArchiveAccountId",
				Required:    true,
			},
			common.ToSnakeCase("Trail"): schema.SingleNestedAttribute{
				Description: "Trail.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("AccountName"): schema.StringAttribute{
						Description: "AccountName",
						Computed:    true,
					},
					common.ToSnakeCase("BucketName"): schema.StringAttribute{
						Description: "BucketName",
						Computed:    true,
					},
					common.ToSnakeCase("BucketRegion"): schema.StringAttribute{
						Description: "BucketRegion",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedUserId"): schema.StringAttribute{
						Description: "CreatedUserId",
						Computed:    true,
					},
					common.ToSnakeCase("DelYn"): schema.StringAttribute{
						Description: "DelYn",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("LogTypeTotalYn"): schema.StringAttribute{
						Description: "LogTypeTotalYn",
						Computed:    true,
					},
					common.ToSnakeCase("LogVerificationYn"): schema.StringAttribute{
						Description: "LogVerificationYn",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("RegionNames"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "RegionNames",
						Computed:    true,
					},
					common.ToSnakeCase("RegionTotalYn"): schema.StringAttribute{
						Description: "RegionTotalYn",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceTypeTotalYn"): schema.StringAttribute{
						Description: "ResourceTypeTotalYn",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("TargetLogTypes"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "TargetLogTypes",
						Computed:    true,
					},
					common.ToSnakeCase("TargetResourceTypes"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "TargetResourceTypes",
						Computed:    true,
					},
					common.ToSnakeCase("TargetUsers"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "TargetUsers",
						Computed:    true,
					},
					common.ToSnakeCase("TrailBatchEndAt"): schema.StringAttribute{
						Description: "TrailBatchEndAt",
						Computed:    true,
					},
					common.ToSnakeCase("TrailBatchFirstStartAt"): schema.StringAttribute{
						Description: "TrailBatchFirstStartAt",
						Computed:    true,
					},
					common.ToSnakeCase("TrailBatchLastState"): schema.StringAttribute{
						Description: "TrailBatchLastState",
						Computed:    true,
					},
					common.ToSnakeCase("TrailBatchStartAt"): schema.StringAttribute{
						Description: "TrailBatchStartAt",
						Computed:    true,
					},
					common.ToSnakeCase("TrailBatchSuccessAt"): schema.StringAttribute{
						Description: "TrailBatchSuccessAt",
						Computed:    true,
					},
					common.ToSnakeCase("TrailDescription"): schema.StringAttribute{
						Description: "TrailDescription",
						Computed:    true,
					},
					common.ToSnakeCase("TrailName"): schema.StringAttribute{
						Description: "TrailName",
						Computed:    true,
					},
					common.ToSnakeCase("TrailSaveType"): schema.StringAttribute{
						Description: "TrailSaveType",
						Computed:    true,
					},
					common.ToSnakeCase("UserTotalYn"): schema.StringAttribute{
						Description: "UserTotalYn",
						Computed:    true,
					},
					common.ToSnakeCase("OrganizationTrailYn"): schema.StringAttribute{
						Description: "OrganizationTrailYn",
						Computed:    true,
					},
					common.ToSnakeCase("LogArchiveAccountId"): schema.StringAttribute{
						Description: "LogArchiveAccountId",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loggingauditTrailResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoggingAudit
	r.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (r *loggingauditTrailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loggingaudit.TrailResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetTrail(ctx, state.Id.ValueString())

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Trail",
			"Could not read Trail ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail

	trailModel := loggingaudit.Trail{
		AccountId:              types.StringPointerValue(trail.AccountId.Get()),
		AccountName:            types.StringPointerValue(trail.AccountName.Get()),
		BucketName:             types.StringValue(trail.BucketName),
		BucketRegion:           types.StringValue(trail.BucketName),
		CreatedAt:              types.StringValue(trail.CreatedAt.Format(time.RFC3339)),
		CreatedBy:              types.StringValue(trail.CreatedBy),
		CreatedUserId:          types.StringPointerValue(trail.CreatedUserId.Get()),
		DelYn:                  types.StringPointerValue(trail.DelYn.Get()),
		Id:                     types.StringValue(trail.Id),
		LogTypeTotalYn:         types.StringPointerValue(trail.LogTypeTotalYn.Get()),
		LogVerificationYn:      types.StringPointerValue(trail.LogVerificationYn.Get()),
		ModifiedAt:             types.StringValue(trail.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:             types.StringValue(trail.ModifiedBy),
		RegionNames:            ConvertInterfaceListToStringList(trail.RegionNames),
		RegionTotalYn:          types.StringPointerValue(trail.RegionTotalYn.Get()),
		ResourceTypeTotalYn:    types.StringPointerValue(trail.ResourceTypeTotalYn.Get()),
		State:                  types.StringPointerValue(trail.State.Get()),
		TargetLogTypes:         ConvertInterfaceListToStringList(trail.TargetLogTypes),
		TargetResourceTypes:    ConvertInterfaceListToStringList(trail.TargetResourceTypes),
		TargetUsers:            ConvertInterfaceListToStringList(trail.TargetUsers),
		TrailBatchEndAt:        TimePointValue(trail.TrailBatchEndAt.Get()),
		TrailBatchFirstStartAt: TimePointValue(trail.TrailBatchFirstStartAt.Get()),
		TrailBatchLastState:    types.StringPointerValue(trail.TrailBatchLastState.Get()),
		TrailBatchStartAt:      TimePointValue(trail.TrailBatchStartAt.Get()),
		TrailBatchSuccessAt:    TimePointValue(trail.TrailBatchSuccessAt.Get()),
		TrailDescription:       types.StringPointerValue(trail.TrailDescription.Get()),
		TrailName:              types.StringValue(trail.TrailName),
		TrailSaveType:          types.StringValue(trail.TrailSaveType),
		UserTotalYn:            types.StringPointerValue(trail.UserTotalYn.Get()),
		OrganizationTrailYn:    types.StringPointerValue(trail.OrganizationTrailYn.Get()),
		LogArchiveAccountId:    types.StringPointerValue(trail.LogArchiveAccountId.Get()),
	}

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)
	state.Trail = trailObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func TimePointValue(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func (r *loggingauditTrailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loggingaudit.TrailResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateTrail(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Trail",
			"Could not create Trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail
	plan.Id = types.StringValue(trail.Id)

	trailModel := loggingaudit.Trail{
		AccountId:              types.StringPointerValue(trail.AccountId.Get()),
		AccountName:            types.StringPointerValue(trail.AccountName.Get()),
		BucketName:             types.StringValue(trail.BucketName),
		BucketRegion:           types.StringValue(trail.BucketName),
		CreatedAt:              types.StringValue(trail.CreatedAt.Format(time.RFC3339)),
		CreatedBy:              types.StringValue(trail.CreatedBy),
		CreatedUserId:          types.StringPointerValue(trail.CreatedUserId.Get()),
		DelYn:                  types.StringPointerValue(trail.DelYn.Get()),
		Id:                     types.StringValue(trail.Id),
		LogTypeTotalYn:         types.StringPointerValue(trail.LogTypeTotalYn.Get()),
		LogVerificationYn:      types.StringPointerValue(trail.LogVerificationYn.Get()),
		ModifiedAt:             types.StringValue(trail.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:             types.StringValue(trail.ModifiedBy),
		RegionNames:            ConvertInterfaceListToStringList(trail.RegionNames),
		RegionTotalYn:          types.StringPointerValue(trail.RegionTotalYn.Get()),
		ResourceTypeTotalYn:    types.StringPointerValue(trail.ResourceTypeTotalYn.Get()),
		State:                  types.StringPointerValue(trail.State.Get()),
		TargetLogTypes:         ConvertInterfaceListToStringList(trail.TargetLogTypes),
		TargetResourceTypes:    ConvertInterfaceListToStringList(trail.TargetResourceTypes),
		TargetUsers:            ConvertInterfaceListToStringList(trail.TargetUsers),
		TrailBatchEndAt:        TimePointValue(trail.TrailBatchEndAt.Get()),
		TrailBatchFirstStartAt: TimePointValue(trail.TrailBatchFirstStartAt.Get()),
		TrailBatchLastState:    types.StringPointerValue(trail.TrailBatchLastState.Get()),
		TrailBatchStartAt:      TimePointValue(trail.TrailBatchStartAt.Get()),
		TrailBatchSuccessAt:    TimePointValue(trail.TrailBatchSuccessAt.Get()),
		TrailDescription:       types.StringPointerValue(trail.TrailDescription.Get()),
		TrailName:              types.StringValue(trail.TrailName),
		TrailSaveType:          types.StringValue(trail.TrailSaveType),
		UserTotalYn:            types.StringPointerValue(trail.UserTotalYn.Get()),
		OrganizationTrailYn:    types.StringPointerValue(trail.OrganizationTrailYn.Get()),
		LogArchiveAccountId:    types.StringPointerValue(trail.LogArchiveAccountId.Get()),
	}

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)

	plan.Trail = trailObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loggingauditTrailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state loggingaudit.TrailResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SetTrail(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Trail",
			"Could not update Trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetTrail(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading trail",
			"Could not read trail ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	trail := data.Trail

	trailModel := loggingaudit.Trail{
		AccountId:              types.StringPointerValue(trail.AccountId.Get()),
		AccountName:            types.StringPointerValue(trail.AccountName.Get()),
		BucketName:             types.StringValue(trail.BucketName),
		BucketRegion:           types.StringValue(trail.BucketName),
		CreatedAt:              types.StringValue(trail.CreatedAt.Format(time.RFC3339)),
		CreatedBy:              types.StringValue(trail.CreatedBy),
		CreatedUserId:          types.StringPointerValue(trail.CreatedUserId.Get()),
		DelYn:                  types.StringPointerValue(trail.DelYn.Get()),
		Id:                     types.StringValue(trail.Id),
		LogTypeTotalYn:         types.StringPointerValue(trail.LogTypeTotalYn.Get()),
		LogVerificationYn:      types.StringPointerValue(trail.LogVerificationYn.Get()),
		ModifiedAt:             types.StringValue(trail.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:             types.StringValue(trail.ModifiedBy),
		RegionNames:            ConvertInterfaceListToStringList(trail.RegionNames),
		RegionTotalYn:          types.StringPointerValue(trail.RegionTotalYn.Get()),
		ResourceTypeTotalYn:    types.StringPointerValue(trail.ResourceTypeTotalYn.Get()),
		State:                  types.StringPointerValue(trail.State.Get()),
		TargetLogTypes:         ConvertInterfaceListToStringList(trail.TargetLogTypes),
		TargetResourceTypes:    ConvertInterfaceListToStringList(trail.TargetResourceTypes),
		TargetUsers:            ConvertInterfaceListToStringList(trail.TargetUsers),
		TrailBatchEndAt:        TimePointValue(trail.TrailBatchEndAt.Get()),
		TrailBatchFirstStartAt: TimePointValue(trail.TrailBatchFirstStartAt.Get()),
		TrailBatchLastState:    types.StringPointerValue(trail.TrailBatchLastState.Get()),
		TrailBatchStartAt:      TimePointValue(trail.TrailBatchStartAt.Get()),
		TrailBatchSuccessAt:    TimePointValue(trail.TrailBatchSuccessAt.Get()),
		TrailDescription:       types.StringPointerValue(trail.TrailDescription.Get()),
		TrailName:              types.StringValue(trail.TrailName),
		TrailSaveType:          types.StringValue(trail.TrailSaveType),
		UserTotalYn:            types.StringPointerValue(trail.UserTotalYn.Get()),
		OrganizationTrailYn:    types.StringPointerValue(trail.OrganizationTrailYn.Get()),
		LogArchiveAccountId:    types.StringPointerValue(trail.LogArchiveAccountId.Get()),
	}

	trailObjectValue, diags := types.ObjectValueFrom(ctx, trailModel.AttributeTypes(), trailModel)
	state.Trail = trailObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *loggingauditTrailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loggingaudit.TrailResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing iam
	err := r.client.DeleteTrailKey(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting trail",
			"Could not delete trail, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func ConvertInterfaceListToStringList(rawList []interface{}) []types.String {
	result := make([]types.String, 0, len(rawList))
	for _, v := range rawList {
		if v == nil {
			result = append(result, types.StringNull())
		} else {
			strVal, _ := v.(string)
			result = append(result, types.StringValue(strVal))
		}
	}
	return result
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *loggingauditTrailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
