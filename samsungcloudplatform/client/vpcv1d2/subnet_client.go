package vpcv1d2

import (
	"context"

	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
)

//------------ Subnet -------------------//

func (client *Client) GetSubnetList(ctx context.Context, request SubnetDataSource) (*vpc.SubnetListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ListSubnets(ctx)

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
	if len(request.Type) > 0 {
		req = req.Type_(vpc.Type{
			ArrayOfSubnetType: &request.Type,
		})
	}
	if !request.State.IsNull() {
		req = req.State(vpc.SubnetState(request.State.ValueString()))
	}
	if !request.Cidr.IsNull() {
		req = req.Cidr(request.Cidr.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}

	resp, _, err := req.Execute()

	return resp, err

}

func (client *Client) CreateSubnet(ctx context.Context, request SubnetResource) (*vpc.SubnetShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.CreateSubnet(ctx)
	description := request.Description.ValueString()
	descriptionNS := vpc.NullableString{}
	descriptionNS.Set(&description)

	tags := convertToTags(request.Tags.Elements())

	req = req.SubnetCreateRequestV1Dot2(vpc.SubnetCreateRequestV1Dot2{
		Name:             request.Name.ValueString(),
		VpcId:            request.VpcId.ValueString(),
		Type:             vpc.SubnetType(request.Type.ValueString()),
		Cidr:             request.Cidr.ValueString(),
		Description:      descriptionNS,
		AllocationPools:  convertAllocationPoolsToInterface(request.AllocationPools),
		DnsNameservers:   dnsNameserversToStringSlice(request.DnsNameservers),
		HostRoutes:       convertHostRoutesToInterface(request.HostRoutes),
		Tags:             tags,
		GatewayIpAddress: *vpc.NewNullableString(request.GatewayIpAddress.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetSubnet(ctx context.Context, subnetId string) (*vpc.SubnetShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.ShowSubnet(ctx, subnetId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSubnet(ctx context.Context, vpcId string, request SubnetResource) (*vpc.SubnetShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1SubnetsApiAPI.SetSubnet(ctx, vpcId)

	description := request.Description.ValueString()
	dhcpIpAddress := request.DhcpIpAddress.ValueString()

	req = req.SubnetSetRequestV1Dot2(vpc.SubnetSetRequestV1Dot2{
		Description:   *vpc.NewNullableString(&description),
		DhcpIpAddress: *vpc.NewNullableString(&dhcpIpAddress),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnet(ctx context.Context, subnetId string) error {
	req := client.sdkClient.VpcV1SubnetsApiAPI.DeleteSubnet(ctx, subnetId)

	_, err := req.Execute()
	return err
}
