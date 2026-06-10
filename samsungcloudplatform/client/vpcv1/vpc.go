package vpcv1

import (
	"context"
	"fmt"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpvpc.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpvpc.NewAPIClient(config),
	}
}

//------------ VPC PEERING-------------------//

func (client *Client) GetListVpcPeering(ctx context.Context, request VpcPeeringListDataSource) (*scpvpc.VpcPeeringListResponse, error) {
	fmt.Printf("Start call GetListVpcPeering -------------------\n")
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.ListVpcPeerings(ctx)
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.RequesterVpcId.IsNull() {
		req = req.RequesterVpcId(request.RequesterVpcId.ValueString())
	}
	if !request.RequesterVpcName.IsNull() {
		req = req.RequesterVpcName(request.RequesterVpcName.ValueString())
	}
	if !request.ApproverVpcName.IsNull() {
		req = req.ApproverVpcName(request.ApproverVpcName.ValueString())
	}
	if !request.ApproverVpcId.IsNull() {
		req = req.ApproverVpcId(request.ApproverVpcId.ValueString())
	}
	if !request.AccountType.IsNull() {
		req = req.AccountType(scpvpc.VpcPeeringAccountType(request.AccountType.ValueString()))
	}

	if !request.State.IsNull() {
		req = req.State(scpvpc.VpcPeeringState(request.State.ValueString()))
	}

	resp, status, err := req.Execute()
	fmt.Printf("client err-------------------%v\n", err)
	fmt.Printf("client status-------------------%v\n", status)

	return resp, err
}

func (client *Client) CreateVpcPeering(ctx context.Context, request VpcPeeringResource) (*scpvpc.VpcPeeringShowResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.CreateVpcPeering(ctx)

	tags := convertToTags(request.Tags.Elements())

	createReq := scpvpc.VpcPeeringCreateRequest{
		ApproverVpcAccountId: request.ApproverVpcAccountId.ValueString(),
		ApproverVpcId:        request.ApproverVpcId.ValueString(),
		Name:                 request.Name.ValueString(),
		RequesterVpcId:       request.RequesterVpcId.ValueString(),
		Description:          *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
		Tags:                 tags,
	}
	// The API requires approver_vpc_name. Send it when known (derived in the
	// resource Create from approver_vpc_id, or provided directly by the user).
	if !request.ApproverVpcName.IsNull() && !request.ApproverVpcName.IsUnknown() && request.ApproverVpcName.ValueString() != "" {
		createReq.ApproverVpcName = request.ApproverVpcName.ValueStringPointer()
	}
	req = req.VpcPeeringCreateRequest(createReq)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateVpcPeering(ctx context.Context, vpcPeeringId string, request VpcPeeringResource) (*scpvpc.VpcPeeringShowResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.SetVpcPeering(ctx, vpcPeeringId)
	description := request.Description.ValueString()

	req = req.VpcPeeringSetRequest(scpvpc.VpcPeeringSetRequest{
		Description: &description,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpcPeering(ctx context.Context, vpcPeeringId string) (*scpvpc.VpcPeeringShowResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.ShowVpcPeering(ctx, vpcPeeringId)

	resp, _, err := req.Execute()
	return resp, err
}

// GetVpcPeeringWithStatus returns the peering along with the HTTP status code so
// callers (e.g. delete waiters) can distinguish a real 404 (gone) from other errors.
func (client *Client) GetVpcPeeringWithStatus(ctx context.Context, vpcPeeringId string) (*scpvpc.VpcPeeringShowResponse, int, error) {
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.ShowVpcPeering(ctx, vpcPeeringId)

	resp, httpResp, err := req.Execute()
	statusCode := 0
	if httpResp != nil {
		statusCode = httpResp.StatusCode
	}
	return resp, statusCode, err
}

func (client *Client) DeleteVpcPeering(ctx context.Context, vpcPeeringId string) error {
	req := client.sdkClient.VpcV1VpcPeeringApiAPI.DeleteVpcPeering(ctx, vpcPeeringId)

	_, err := req.Execute()
	return err
}

func convertToTags(elements map[string]attr.Value) []scpvpc.Tag {
	var tags []scpvpc.Tag
	for k, v := range elements {
		tagObject := scpvpc.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func (client *Client) ApprovalVpcPeering(ctx context.Context, vpcPeeringId string, approvalType string) (*scpvpc.VpcPeeringShowResponse, error) {
	req := client.sdkClient.VpcV1VpcPeeringApprovalApiAPI.ApprovalVpcPeering(ctx, vpcPeeringId)

	finalType, err := scpvpc.NewVpcPeeringApprovalTypeFromValue(approvalType)
	if err == nil {
		req = req.VpcPeeringApprovalRequest(scpvpc.VpcPeeringApprovalRequest{
			Type: *finalType,
		})
	} else {
		req = req.VpcPeeringApprovalRequest(scpvpc.VpcPeeringApprovalRequest{})
	}

	resp, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	return resp, err
}
