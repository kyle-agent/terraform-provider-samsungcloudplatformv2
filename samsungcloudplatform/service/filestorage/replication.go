package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/filestorage"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/filestorage/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

const ReplicationId = "replicationId"
const VolumeId = "volumeId"
const FieldName = "fieldName"
const ReplicationPolicy = "replicationPolicy"
const BackupRetentionCount = "backupRetentionCount"
const ReplicationFrequency = "replicationFrequency"
const ModifySchedule = "modify_schedule"
const Backup = "backup"
const Use = "use"
const Paused = "paused"

var (
	_ resource.Resource              = &fileStorageReplicationResource{}
	_ resource.ResourceWithConfigure = &fileStorageReplicationResource{}
)

func NewFileStorageReplicationResource() resource.Resource {
	return &fileStorageReplicationResource{}
}

type fileStorageReplicationResource struct {
	config  *scpsdk.Configuration
	client  *filestorage.Client
	clients *client.SCPClient
}

func (f *fileStorageReplicationResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	inst, ok := request.ProviderData.(client.Instance)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	f.client = inst.Client.FileStorage
	f.clients = inst.Client
}

func (f *fileStorageReplicationResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_filestorage_replication"
}

func (f *fileStorageReplicationResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Replication Response Data",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name \n" +
					"  - example : 'my_volume' \n" +
					"  - maxLength: 21  \n" +
					"  - minLength: 3  \n" +
					"  - pattern: `^[a-z]([a-z0-9_]){2,20}$` \n",
				Required: true,
			},
			"replication_frequency": schema.StringAttribute{
				Description: "Replication Frequency \n" +
					"  - example : '5min' \n" +
					"  - pattern: `^(5min|hourly|daily|weekly|monthly)$` \n",
				Required: true,
			},
			"region": schema.StringAttribute{
				Description: "Region \n" +
					"  - example : 'kr-west1' \n",
				Required: true,
			},
			"volume_id": schema.StringAttribute{
				Description: "Source Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Required: true,
			},
			"cifs_password": schema.StringAttribute{
				Description: "Cifs Password \n" +
					"  - example : 'cifspwd0!!' \n" +
					"  - maxLength: 20  \n" +
					"  - minLength: 6  \n" +
					"  - pattern: `^(?=.*[a-zA-Z])(?=.*\\d)(?=.*[!#&\\'*+,-.:;<=>?@^_`~/|])[a-zA-Z\\d!#&\\'*+,-.:;<=>?@^_`~/|]{6,20}$` \n",
				Optional:  true,
				WriteOnly: true,
			},
			"replication_id": schema.StringAttribute{
				Description: "Replication ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Computed: true,
			},
			"replication_policy": schema.StringAttribute{
				Description: "Replication Policy \n" +
					"  - example : 'use' \n" +
					"  - pattern: `^(use|paused)$` \n",
				Optional: true,
				Computed: true,
			},
			"replication_status": schema.StringAttribute{
				Description: "Replication Status \n" +
					"  - example : 'creating' \n",
				Computed: true,
			},
			"replication_volume_access_level": schema.StringAttribute{
				Description: "Target Access Level \n" +
					"  - example : 'ro' \n",
				Computed: true,
			},
			"replication_volume_id": schema.StringAttribute{
				Description: "Target Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Computed: true,
			},
			"replication_volume_name": schema.StringAttribute{
				Description: "Target Volume Name \n" +
					"  - example : 'my_volume' \n",
				Computed: true,
			},
			"replication_volume_region": schema.StringAttribute{
				Description: "Target Volume Region \n" +
					"  - example : 'kr-west1' \n",
				Computed: true,
			},
			"source_volume_access_level": schema.StringAttribute{
				Description: "Source Access Level \n" +
					"  - example : 'ro' \n",
				Computed: true,
			},
			"source_volume_id": schema.StringAttribute{
				Description: "Source Volume ID \n" +
					"  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438' \n",
				Computed: true,
			},
			"source_volume_name": schema.StringAttribute{
				Description: "Source Volume Name \n" +
					"  - example : 'my_volume' \n",
				Computed: true,
			},
			"source_volume_region": schema.StringAttribute{
				Description: "Source Volume Region \n" +
					"  - example : 'kr-west1' \n",
				Computed: true,
			},
			"replication_update_type": schema.StringAttribute{
				Description: "Replication Update Type \n" +
					"  - example : 'policy' \n" +
					"  - pattern: `^(policy|modify_schedule)$` \n",
				Optional: true,
			},
			"backup_retention_count": schema.Int32Attribute{
				Description: "Backup Retention Count \n" +
					"  - example : 'policy' \n" +
					"  - maximum : 128 \n" +
					"  - minimum : 1  \n",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 128),
				},
			},
			"replication_type": schema.StringAttribute{
				Description: "Replication Type \n" +
					"  - example : 'replication' \n" +
					"  - pattern: `^(replication|backup)$` \n",
				Required: true,
			},
		},
	}
}

func (f *fileStorageReplicationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Values from plan
	var plan filestorage.ReplicationResource
	diags := request.Config.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := f.client.CreateReplication(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error creating replication",
			"Could not create replication, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var params = map[string]string{
		ReplicationId: data.ReplicationId,
		VolumeId:      data.SourceVolumeId,
		FieldName:     ReplicationPolicy,
	}

	replication, err := waitForReplicationStatus(ctx, f.client, params, []string{}, []string{Use})
	if err != nil {
		response.Diagnostics.AddError(
			"Error waiting for replication to complete",
			err.Error(),
		)
		return
	}

	plan = mapReplicationToPlan(replication, plan)
	plan.VolumeId = types.StringValue(replication.SourceVolumeId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (f *fileStorageReplicationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state filestorage.ReplicationResource
	diags := request.State.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	replication, err := f.client.GetVolumeReplication(ctx, state.ReplicationId.ValueString(), state.VolumeId.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Unable to Read Replication Policy", err.Error())
		return
	}
	state = mapReplicationToPlan(replication, state)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
}

func (f *fileStorageReplicationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state filestorage.ReplicationResource
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	policy := filestorage.VolumeReplicationPolicy{
		ReplicationFrequency:  coalesceString(plan.ReplicationFrequency, state.ReplicationFrequency),
		ReplicationPolicy:     coalesceString(plan.ReplicationPolicy, state.ReplicationPolicy),
		ReplicationUpdateType: coalesceString(plan.ReplicationUpdateType, state.ReplicationUpdateType),
		BackupRetentionCount:  plan.BackupRetentionCount,
	}
	f.client.Config.Region = state.ReplicationVolumeRegion.ValueString()
	err := f.client.UpdateVolumeReplication(ctx, state.ReplicationId.ValueString(), state.ReplicationVolumeId.ValueString(), policy)
	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Updating Replication",
			"Could not update Replication, unexpected error: "+err.Error()+"\nReason: "+detail)
		return
	}

	var params, targetState = createUpdateParamMapAndTargetState(plan, state)

	replication, err := waitForReplicationStatus(ctx, f.client, params, []string{}, []string{targetState})
	if err != nil {
		response.Diagnostics.AddError("Error waiting for replication to complete", err.Error())
		return
	}

	if plan.ReplicationType.ValueString() == Backup && plan.ReplicationUpdateType.ValueString() == ModifySchedule && !plan.BackupRetentionCount.Equal(state.BackupRetentionCount) {
		params[FieldName] = BackupRetentionCount
		_, err = waitForReplicationStatus(ctx, f.client, params, []string{}, []string{plan.BackupRetentionCount.String()})
		if err != nil {
			response.Diagnostics.AddError("Error waiting for replication to complete", err.Error())
			return
		}
	}

	plan = mapReplicationToPlan(replication, plan)
	plan.VolumeId = types.StringValue(replication.SourceVolumeId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func createUpdateParamMapAndTargetState(plan filestorage.ReplicationResource, state filestorage.ReplicationResource) (map[string]string, string) {
	var params = map[string]string{
		ReplicationId: state.ReplicationId.ValueString(),
		VolumeId:      state.ReplicationVolumeId.ValueString(),
	}

	if plan.ReplicationUpdateType.ValueString() == ModifySchedule {
		params[FieldName] = ReplicationFrequency
		return params, plan.ReplicationFrequency.ValueString()
	}

	params[FieldName] = ReplicationPolicy
	return params, plan.ReplicationPolicy.ValueString()
}

func coalesceString(planVal, stateVal types.String) types.String {
	if !planVal.IsNull() && !planVal.IsUnknown() {
		return planVal
	}
	return stateVal
}

func (f *fileStorageReplicationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	var state filestorage.ReplicationResource
	diags := request.State.Get(ctx, &state)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	f.client.Config.Region = state.ReplicationVolumeRegion.ValueString()
	replicationId := state.ReplicationId.ValueString()
	volumeId := state.ReplicationVolumeId.ValueString()

	// The platform refuses to delete a replication (or its backing volume) while the
	// replication policy is still active: 400 "Cannot delete volume because replication
	// is in use. Delete Replication Policy. -Replication Policy : paused > delete".
	// Mirror the loadbalancer #77 fix style: set the policy to "paused"
	// (PUT /v1/replications/{id} with replication_update_type "policy"), poll until the
	// replication reports paused, then delete, then wait until it is fully gone.
	replication, getErr := f.client.GetVolumeReplication(ctx, replicationId, volumeId)
	if getErr == nil && replication.ReplicationPolicy != Paused {
		if err = f.client.PauseVolumeReplication(ctx, replicationId, volumeId); err != nil {
			detail := client.GetDetailFromError(err)
			response.Diagnostics.AddError(
				"Error Delete Replication",
				"Could not pause replication policy before delete, unexpected error : "+err.Error()+"\nReason: "+detail)
			return
		}

		var params = map[string]string{
			ReplicationId: replicationId,
			VolumeId:      volumeId,
			FieldName:     ReplicationPolicy,
		}
		if _, err = waitForReplicationStatus(ctx, f.client, params, []string{}, []string{Paused}); err != nil {
			response.Diagnostics.AddError(
				"Error Delete Replication",
				"Error waiting for replication policy to pause: "+err.Error())
			return
		}
	}

	if state.ReplicationType == types.StringValue("backup") {
		err = f.client.DeleteVolume(ctx, volumeId)

	} else {
		err = f.client.DeleteVolumeReplication(ctx, replicationId, volumeId)
	}

	if err != nil {
		detail := client.GetDetailFromError(err)
		response.Diagnostics.AddError(
			"Error Delete Replication",
			"Could not Delete Replication, unexpected error : "+err.Error()+"\nReason: "+detail)
		return
	}

	// The delete is asynchronous; wait until the replication no longer exists
	// (Get -> 404) so a dependent volume delete in the same `terraform destroy`
	// does not fail with 400 "Cannot delete volume because replication is in use".
	err = waitForReplicationGone(ctx, f.client, replicationId, volumeId)
	if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		response.Diagnostics.AddError(
			"Error Delete Replication",
			"Error waiting for replication to be deleted: "+err.Error())
		return
	}
}

func mapReplicationToPlan(replication *scpfilestorage.ReplicationShowResponse, plan filestorage.ReplicationResource) filestorage.ReplicationResource {
	plan.SourceVolumeId = types.StringValue(replication.SourceVolumeId)
	plan.SourceVolumeName = types.StringValue(replication.SourceVolumeName)
	plan.SourceVolumeRegion = types.StringValue(replication.SourceVolumeRegion)
	plan.SourceVolumeAccessLevel = types.StringValue(replication.SourceVolumeAccessLevel)
	plan.ReplicationVolumeId = types.StringValue(replication.ReplicationVolumeId)
	plan.ReplicationVolumeName = types.StringValue(replication.ReplicationVolumeName)
	plan.ReplicationVolumeRegion = types.StringValue(replication.ReplicationVolumeRegion)
	plan.ReplicationVolumeAccessLevel = types.StringValue(replication.ReplicationVolumeAccessLevel)
	plan.ReplicationFrequency = types.StringValue(replication.ReplicationFrequency)
	plan.ReplicationId = types.StringValue(replication.ReplicationId)
	plan.ReplicationPolicy = types.StringValue(replication.ReplicationPolicy)
	plan.ReplicationStatus = types.StringValue(replication.ReplicationStatus)
	plan.VolumeId = types.StringValue(replication.SourceVolumeId)

	return plan
}

func waitForReplicationStatus(ctx context.Context, fileStorageClient *filestorage.Client, params map[string]string, pendingStates []string, targetStates []string) (*scpfilestorage.ReplicationShowResponse, error) {
	var response *scpfilestorage.ReplicationShowResponse
	return response, client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		replication, err := fileStorageClient.GetVolumeReplication(ctx, params[ReplicationId], params[VolumeId])
		if err != nil {
			return nil, "", nil
		}
		response = replication
		return replication, getCurrentStatus(params, replication), nil
	})
}

// waitForReplicationGone polls the replication until the show call returns
// 404/not found (fully deleted). It mirrors waitForLoadbalancerStatus' DELETED
// handling: the Get error is surfaced so the caller can treat 404 as success.
func waitForReplicationGone(ctx context.Context, fileStorageClient *filestorage.Client, replicationId string, volumeId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"DELETED"}, func() (interface{}, string, error) {
		replication, err := fileStorageClient.GetVolumeReplication(ctx, replicationId, volumeId)
		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found") {
				return "gone", "DELETED", nil
			}
			return nil, "", err
		}
		return replication, replication.ReplicationStatus, nil
	})
}

func getCurrentStatus(params map[string]string, replication *scpfilestorage.ReplicationShowResponse) string {
	if params[FieldName] == BackupRetentionCount {
		return strconv.Itoa(int(*replication.BackupRetentionCount.Get()))
	} else if params[FieldName] == ReplicationFrequency {
		return replication.ReplicationFrequency
	} else {
		return replication.ReplicationPolicy
	}
}
