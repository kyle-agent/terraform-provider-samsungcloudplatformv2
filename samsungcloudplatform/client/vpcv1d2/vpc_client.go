package vpcv1d2

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

func (client *Client) GetVpcList(ctx context.Context, request VpcDataSource) (*vpc.VpcListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ListVpcs(ctx)
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
	if !request.State.IsNull() {
		req = req.State(vpc.VpcState(request.State.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	resp, _, err := req.Execute()
	return resp, err
}
func (client *Client) CreateVpc(ctx context.Context, request VpcResource) (*vpc.VpcShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.CreateVpc(ctx)

	tags := convertToTags(request.Tags.Elements())

	req = req.VpcCreateRequest(vpc.VpcCreateRequest{
		Cidr:        request.Cidr.ValueString(),
		Description: *vpc.NewNullableString(request.Description.ValueStringPointer()),
		Name:        request.Name.ValueString(),
		Tags:        tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetVpc(ctx context.Context, vpcId string) (*vpc.VpcShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ShowVpc(ctx, vpcId)

	resp, _, err := req.Execute()
	return resp, err
}

// GetVpcWithStatus returns the VPC along with the HTTP status code so callers
// (e.g. Read) can distinguish a real 404 (gone) from other errors.
func (client *Client) GetVpcWithStatus(ctx context.Context, vpcId string) (*vpc.VpcShowResponseV1Dot2, int, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.ShowVpc(ctx, vpcId)

	resp, httpResp, err := req.Execute()
	statusCode := 0
	if httpResp != nil {
		statusCode = httpResp.StatusCode
	}
	return resp, statusCode, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, request VpcResource) (*vpc.VpcShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.SetVpc(ctx, vpcId)

	req = req.VpcSetRequest(vpc.VpcSetRequest{
		Description: *vpc.NewNullableString(request.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteVpc(ctx context.Context, vpcId string) error {
	req := client.sdkClient.VpcV1VpcsApiAPI.DeleteVpc(ctx, vpcId)

	_, err := req.Execute()
	return err
}
