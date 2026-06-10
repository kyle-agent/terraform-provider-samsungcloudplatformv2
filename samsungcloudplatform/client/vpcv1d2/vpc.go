package vpcv1d2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Client struct definition
type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scpvpc.APIClient
}

// NewClient initializes the VPC client
func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scpvpc.NewAPIClient(config),
	}
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

// ------------ Subnet VIP Client Methods ------------

func (client *Client) ListSubnetVips(ctx context.Context, request SubnetVipDataSources) (*scpvpc.VipListResponse, error) {
	fmt.Printf("request.SubnetId.ValueString(): %v", request.SubnetId.ValueString())
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.ListSubnetVips(ctx, request.SubnetId.ValueString())

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.VirtualIpAddress.IsNull() {
		req = req.VirtualIpAddress(request.VirtualIpAddress.ValueString())
	}
	if !request.PublicIpAddress.IsNull() {
		req = req.PublicIpAddress(request.PublicIpAddress.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) ShowSubnetVip(ctx context.Context, SubnetId string, VipID string) (*scpvpc.VipShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.ShowSubnetVip(ctx, SubnetId, VipID)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateSubnetVIP(ctx context.Context, request SubnetVipResource) (*scpvpc.VipCreateResponse, error) {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.CreateSubnetVIP(ctx, request.SubnetId.ValueString())

	createReq := scpvpc.VipCreateRequest{
		VirtualIpAddress: *scpvpc.NewNullableString(request.VirtualIpAddress.ValueStringPointer()),
		Description:      *scpvpc.NewNullableString(request.Description.ValueStringPointer()),
	}

	req = req.VipCreateRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateSubnetVIP(ctx context.Context, SubnetId string, VipID string, Description string) (*scpvpc.VipShowResponse, error) {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.SetSubnetVip(ctx, SubnetId, VipID)

	setReq := scpvpc.VipSetRequest{
		Description: Description,
	}

	req = req.VipSetRequest(setReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnetVIP(ctx context.Context, SubnetId string, VipID string) error {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.DeleteSubnetVip(ctx, SubnetId, VipID)
	_, err := req.Execute()
	return err
}

// ------------ Subnet VIP NAT IP Client Methods ------------

func (client *Client) CreateSubnetVipNatIp(ctx context.Context, request SubnetVipNatIpResource) (*scpvpc.VipNatCreateResponse, error) {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.CreateSubnetVIPNATIp(ctx, request.SubnetId.ValueString(), request.VipId.ValueString())

	createReq := scpvpc.VipNatCreateRequest{
		NatType:    request.NatType.ValueString(),
		PublicipId: request.PublicipId.ValueString(),
	}

	req = req.VipNatCreateRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnetVipNatIp(ctx context.Context, state SubnetVipNatIpResource) error {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.DeleteSubnetVIPNATIp(ctx, state.SubnetId.ValueString(), state.VipId.ValueString(), state.Id.ValueString())
	_, err := req.Execute()
	return err
}

// ------------ VPC CIDR Client Methods ------------

func (client *Client) AddVpcCidr(ctx context.Context, request VpcCidrResource) (*scpvpc.VpcShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcsApiAPI.AddVpcCidr(ctx, request.VpcId.ValueString())

	createReq := scpvpc.VpcCidrCreateRequest{
		Cidr: request.Cidr.ValueString(),
	}

	req = req.VpcCidrCreateRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

// RemoveVpcCidr removes a secondary CIDR from a VPC. The generated SDK only
// exposes AddVpcCidr (POST /v1/vpcs/{vpc_id}/cidrs); the inverse is a DELETE on
// /v1/vpcs/{vpc_id}/cidrs/{cidr_id}. We build that signed request using the
// SDK's own exported helpers so it behaves identically to generated calls.
// It returns the HTTP status code so callers can treat 404 as already-deleted.
func (client *Client) RemoveVpcCidr(ctx context.Context, vpcId string, cidrId string) (int, error) {
	cfg := client.sdkClient.GetConfig()

	// Resolve the service base path exactly like the generated SDK does.
	basePath := cfg.Endpoint
	if basePath == "" {
		catalog := scpsdk.NewCatalog(
			cfg.AuthUrl,
			cfg.Credentials.AccessKey,
			cfg.Credentials.SecretKey,
			cfg.DefaultRegion,
		)
		ep, err := catalog.GetEndpoint(cfg.ServiceType, cfg.Region, cfg.AccountId)
		if err != nil {
			return 0, err
		}
		basePath = ep
	}

	fullPath := basePath + "/v1/vpcs/" + url.PathEscape(vpcId) + "/cidrs/" + url.PathEscape(cidrId)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullPath, nil)
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Scp-API-Version", "vpc 1.2")
	if cfg.Credentials.AuthToken != "" {
		httpReq.Header.Set("X-Auth-Token", cfg.Credentials.AuthToken)
	}

	// Sign the request (Scp-AccessKey / Scp-Signature / Scp-Timestamp / etc.).
	cfg.SetupRequestHeader(fullPath, http.MethodDelete, httpReq)

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode >= 300 && statusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return statusCode, fmt.Errorf("failed to remove VPC CIDR (status %d): %s", statusCode, string(body))
	}
	return statusCode, nil
}

// ------------ Subnet VIP Port Client Methods ------------

func (client *Client) CreateSubnetVipPort(ctx context.Context, request SubnetVipPortResource) (*scpvpc.VipPortResponse, error) {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.CreateSubnetVipPort(ctx, request.SubnetId.ValueString(), request.VipId.ValueString())

	createReq := scpvpc.VipPortRequest{
		PortId: request.PortId.ValueString(),
	}

	req = req.VipPortRequest(createReq)
	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteSubnetVipPort(ctx context.Context, state SubnetVipPortResource) error {
	req := client.sdkClient.VpcV1SubnetVipsApiAPI.DeleteSubnetVipConnectedPort(ctx, state.SubnetId.ValueString(), state.VipId.ValueString(), state.Id.ValueString())
	_, err := req.Execute()
	return err
}

// ------------ VPC Endpoint v1.2 Client Methods ------------
func (client *Client) GetVpcEndpointList(ctx context.Context, request VpcEndpointDataSource) (*scpvpc.VpcEndpointListResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1VpcEndpointApiAPI.ListVpcEndpoints(ctx)

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
	if !request.VpcName.IsNull() {
		req = req.VpcName(request.VpcName.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}
	if !request.SubnetId.IsNull() {
		req = req.SubnetId(request.SubnetId.ValueString())
	}
	if !request.ResourceType.IsNull() {
		req = req.ResourceType(scpvpc.VpcEndpointResourceType(request.ResourceType.ValueString()))
	}
	if !request.ResourceKey.IsNull() {
		req = req.ResourceKey(request.ResourceKey.ValueString())
	}
	if !request.EndpointIpAddress.IsNull() {
		req = req.EndpointIpAddress(request.EndpointIpAddress.ValueString())
	}
	if !request.State.IsNull() {
		req = req.State(scpvpc.VpcEndpointState(request.State.ValueString()))
	}

	resp, _, err := req.Execute()
	return resp, err
}

// ------------ Transit Gateway Firewall Client Methods ------------

func (client *Client) CreateTransitGatewayFirewall(ctx context.Context, transitGatewayId string, request TransitGatewayFireWallResource) (*scpvpc.TransitGatewayShowResponseV1Dot2, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGatewayFirewall(ctx, transitGatewayId)

	firewallCreateRequest := scpvpc.TransitGatewayFirewallCreateRequest{
		ProductType: scpvpc.TransitGatewayFirewallProductType(request.ProductType.ValueString()),
	}

	req = req.TransitGatewayFirewallCreateRequest(firewallCreateRequest)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteTransitGatewayFirewall(ctx context.Context, transitGatewayId string, firewallId string) (*http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGatewayFirewall(ctx, transitGatewayId, firewallId)
	resp, err := req.Execute()
	return resp, err
}

// ------------ Transit Gateway Firewall Connection Client Methods ------------

func (client *Client) CreateTransitGatewayFirewallConnection(ctx context.Context, transitGatewayId string) (*scpvpc.TransitGatewayShowResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.CreateTransitGatewayFirewallConnection(ctx, transitGatewayId)
	resp, status, err := req.Execute()
	return resp, status, err
}

func (client *Client) DeleteTransitGatewayFirewallConnection(ctx context.Context, transitGatewayId string) (*scpvpc.TransitGatewayShowResponseV1Dot2, *http.Response, error) {
	req := client.sdkClient.VpcV1TransitGatewayApiAPI.DeleteTransitGatewayFirewallConnection(ctx, transitGatewayId)
	resp, status, err := req.Execute()
	return resp, status, err
}
