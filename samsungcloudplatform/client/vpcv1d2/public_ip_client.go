package vpcv1d2

import (
	"context"

	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) ListPublicips(ctx context.Context, request PublicipDataSource) (*scpvpc.PublicipListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ListPublicip(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.IpAddress.IsNull() {
		req = req.IpAddress(request.IpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(request.State.ValueString())
	}
	if !request.AttachedResourceType.IsNull() {
		req = req.AttachedResourceType(request.AttachedResourceType.ValueString())
	}
	if !request.AttachedResourceId.IsNull() {
		req = req.AttachedResourceId(request.AttachedResourceId.ValueString())
	}
	if !request.AttachedResourceName.IsNull() {
		req = req.AttachedResourceName(request.AttachedResourceName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(scpvpc.PublicipType(request.Type.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

// GetPublicipWithStatus fetches a single public IP by ID using the v1.2 API
// (whose PublicipAttachedResourceType enum includes SUBNET, unlike v1.1) and
// returns the HTTP status code so callers can detect a real 404.
func (client *Client) GetPublicipWithStatus(ctx context.Context, publicipId string) (*scpvpc.PublicipShowResponse, int, error) {
	req := client.sdkClient.VpcV1PublicIpApiAPI.ShowPublicip(ctx, publicipId)

	resp, httpResp, err := req.Execute()
	statusCode := 0
	if httpResp != nil {
		statusCode = httpResp.StatusCode
	}
	return resp, statusCode, err
}
