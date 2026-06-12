package filestorage

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpfilestorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/filestorage/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpfilestorage.APIClient // 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: scpfilestorage.NewAPIClient(config),
	}
}

// FileStorage
func (client *Client) GetVolumeList(ctx context.Context, request VolumeDataSourceIds) (*scpfilestorage.VolumeListResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeAPIsAPI.ListVolumes(ctx)
	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.TypeName.IsNull() {
		req = req.TypeName(request.TypeName.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateVolume(ctx context.Context, request VolumeResource) (*scpfilestorage.VolumeCreateResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeAPIsAPI.CreateVolume(ctx)

	if request.Id.ValueString() == "" {
		request.Id = types.StringNull()
	}

	// Tags
	var TagsObjects []scpfilestorage.Tag
	for k, v := range request.Tags.Elements() {
		tagObject := scpfilestorage.Tag{
			Key:   k,
			Value: *scpfilestorage.NewNullableString(v.(types.String).ValueStringPointer()),
		}
		TagsObjects = append(TagsObjects, tagObject)
	}

	req = req.VolumeCreateRequest(scpfilestorage.VolumeCreateRequest{
		CifsPassword: *scpfilestorage.NewNullableString(request.CifsPassword.ValueStringPointer()),
		Name:         request.Name.ValueString(),
		Protocol:     request.Protocol.ValueString(),
		TypeName:     request.TypeName.ValueString(),
		Tags:         TagsObjects,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVolume(ctx context.Context, id string) (*scpfilestorage.VolumeShowResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeAPIsAPI.ShowVolume(ctx, id)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVolume(ctx context.Context, id string, request VolumeResource) error {
	req := client.sdkClient.FilestorageV1VolumeAPIsAPI.SetVolume(ctx, id)

	req = req.VolumeSetRequest(scpfilestorage.VolumeSetRequest{
		FileUnitRecoveryEnabled: request.FileUnitRecoveryEnabled.ValueBool(),
	})

	_, _, err := req.Execute()
	return err
}

func (client *Client) CreateSnapshotSchedule(ctx context.Context, request SnapshotScheduleResource) (*scpfilestorage.SnapshotScheduleCreateResponse, error) {
	req := client.sdkClient.FilestorageV1SnapshotScheduleAPIsAPI.CreateSnapshotSchedule(ctx)

	snapshotScheduleElement := scpfilestorage.SnapshotSchedule{
		DayOfWeek: *scpfilestorage.NewNullableString(request.SnapshotSchedule.DayOfWeek.ValueStringPointer()),
		Hour:      request.SnapshotSchedule.Hour.ValueString(),
		Frequency: request.SnapshotSchedule.Frequency.ValueString(),
	}

	req = req.SnapshotScheduleCreateRequest(scpfilestorage.SnapshotScheduleCreateRequest{
		VolumeId:               request.VolumeId.ValueString(),
		SnapshotSchedule:       snapshotScheduleElement,
		SnapshotRetentionCount: *scpfilestorage.NewNullableInt32(request.SnapshotRetentionCount.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSnapshotSchedule(ctx context.Context, snapshotScheduleId string, request SnapshotScheduleResource) (*scpfilestorage.SnapshotScheduleSetResponse, error) {
	req := client.sdkClient.FilestorageV1SnapshotScheduleAPIsAPI.SetSnapshotSchedule(ctx, snapshotScheduleId)

	req = req.VolumeId(request.VolumeId.ValueString())

	snapshotScheduleElement := scpfilestorage.SnapshotSchedule{
		DayOfWeek: *scpfilestorage.NewNullableString(request.SnapshotSchedule.DayOfWeek.ValueStringPointer()),
		Hour:      request.SnapshotSchedule.Hour.ValueString(),
		Frequency: request.SnapshotSchedule.Frequency.ValueString(),
	}

	req = req.SnapshotScheduleSetRequest(scpfilestorage.SnapshotScheduleSetRequest{
		SnapshotSchedule:       snapshotScheduleElement,
		SnapshotRetentionCount: *scpfilestorage.NewNullableInt32(request.SnapshotRetentionCount.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSnapshotSchedule(ctx context.Context, snapshotScheduleId string, request SnapshotScheduleResource) error {
	req := client.sdkClient.FilestorageV1SnapshotScheduleAPIsAPI.DeleteSnapshotSchedule(ctx, snapshotScheduleId)

	req = req.VolumeId(request.VolumeId.ValueString())

	_, err := req.Execute()
	return err
}

func (client *Client) GetSnapshotScheduleList(ctx context.Context, id string) (*scpfilestorage.SnapshotScheduleListResponse, error) {
	req := client.sdkClient.FilestorageV1SnapshotScheduleAPIsAPI.ListSnapshotSchedule(ctx)

	req = req.VolumeId(id)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVolume(ctx context.Context, id string) error {
	req := client.sdkClient.FilestorageV1VolumeAPIsAPI.DeleteVolume(ctx, id)

	_, err := req.Execute()
	return err
}

// AccessRules
func (client *Client) GetVolumeAccessRules(ctx context.Context, id string) (*scpfilestorage.VolumeObjectAccessRuleListResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeAccessRulesAPIsAPI.ListAccessRules(ctx, id)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVolumeAccessRule(ctx context.Context, id string, request AccessRuleResource, action string) error {
	req := client.sdkClient.FilestorageV1VolumeAccessRulesAPIsAPI.SetAccessRule(ctx, id)

	req = req.AccessRuleRequest(scpfilestorage.AccessRuleRequest{
		Action:     action,
		ObjectId:   request.ObjectId.ValueString(),
		ObjectType: request.ObjectType.ValueString(),
	})
	_, _, err := req.Execute()
	return err
}

// Replication
func (client *Client) CreateReplication(ctx context.Context, request ReplicationResource) (*scpfilestorage.ReplicationCreateResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.CreateVolumeReplication(ctx)

	createReq := scpfilestorage.ReplicationCreateRequest{
		Name:                 request.Name.ValueString(),
		Region:               request.Region.ValueString(),
		ReplicationFrequency: request.ReplicationFrequency.ValueString(),
		VolumeId:             request.VolumeId.ValueString(),
		ReplicationType:      request.ReplicationType.ValueString(),
	}
	// NewNullableString/Int32 mark the field set even for a nil pointer, so an
	// unconfigured value is serialized as an explicit JSON null (the vpc-peering
	// description bug pattern). Only attach optionals the practitioner set.
	if !request.CifsPassword.IsNull() && !request.CifsPassword.IsUnknown() {
		createReq.CifsPassword = *scpfilestorage.NewNullableString(request.CifsPassword.ValueStringPointer())
	}
	if !request.BackupRetentionCount.IsNull() && !request.BackupRetentionCount.IsUnknown() {
		createReq.BackupRetentionCount = *scpfilestorage.NewNullableInt32(request.BackupRetentionCount.ValueInt32Pointer())
	}
	req = req.ReplicationCreateRequest(createReq)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVolumeReplication(ctx context.Context, id string, volumeId string, request VolumeReplicationPolicy) error {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.SetVolumeReplication(ctx, id).VolumeId(volumeId)
	freq := scpfilestorage.ReplicationUpdatePolicyEnum(request.ReplicationFrequency.ValueString())
	policy := scpfilestorage.ReplicationUpdateStatusEnum(request.ReplicationPolicy.ValueString())
	req = req.ReplicationUpdateRequest(scpfilestorage.ReplicationUpdateRequest{
		ReplicationFrequency:  *scpfilestorage.NewNullableReplicationUpdatePolicyEnum(&freq),
		ReplicationPolicy:     *scpfilestorage.NewNullableReplicationUpdateStatusEnum(&policy),
		ReplicationUpdateType: request.ReplicationUpdateType.ValueString(),
		BackupRetentionCount:  *scpfilestorage.NewNullableInt32(request.BackupRetentionCount.ValueInt32Pointer()),
	})

	_, _, err := req.Execute()
	return err
}

// PauseVolumeReplication sets the replication policy to "paused" via
// PUT /v1/replications/{replication_id}?volume_id={volume_id} with body
// {"replication_update_type": "policy", "replication_policy": "paused"}.
// The platform requires the policy to be paused before the replication can be
// deleted ("Replication Policy : paused > delete"), so Delete calls this first.
// Only the two fields are attached so unset optionals are not serialized as
// explicit JSON nulls (the vpc-peering description bug pattern).
func (client *Client) PauseVolumeReplication(ctx context.Context, id string, volumeId string) error {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.SetVolumeReplication(ctx, id).VolumeId(volumeId)

	updateReq := scpfilestorage.ReplicationUpdateRequest{
		ReplicationUpdateType: "policy",
	}
	updateReq.SetReplicationPolicy(scpfilestorage.REPLICATIONUPDATESTATUSENUM_PAUSED)
	req = req.ReplicationUpdateRequest(updateReq)

	_, _, err := req.Execute()
	return err
}

func (client *Client) GetVolumeReplicationList(ctx context.Context, id string) (*scpfilestorage.ReplicationListResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.ListVolumeReplications(ctx)
	req = req.VolumeId(id)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVolumeReplication(ctx context.Context, replicationId, volumeId string) (*scpfilestorage.ReplicationShowResponse, error) {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.ShowVolumeReplication(ctx, replicationId)
	req = req.VolumeId(volumeId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVolumeReplication(ctx context.Context, id string, request string) error {
	req := client.sdkClient.FilestorageV1VolumeReplicationAPIsAPI.DeleteVolumeReplication(ctx, id).VolumeId(request)

	_, err := req.Execute()
	return err
}
