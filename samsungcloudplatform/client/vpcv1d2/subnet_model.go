package vpcv1d2

import (
	vpc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpc/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SubnetDataSource struct {
	Cidr       types.String     `tfsdk:"cidr"`
	Id         types.String     `tfsdk:"id"`
	Name       types.String     `tfsdk:"name"`
	Page       types.Int32      `tfsdk:"page"`
	Size       types.Int32      `tfsdk:"size"`
	Sort       types.String     `tfsdk:"sort"`
	State      types.String     `tfsdk:"state"`
	Subnets    []Subnet         `tfsdk:"subnets"`
	TotalCount types.Int32      `tfsdk:"total_count"`
	Type       []vpc.SubnetType `tfsdk:"type"`
	VpcId      types.String     `tfsdk:"vpc_id"`
	VpcName    types.String     `tfsdk:"vpc_name"`
}

type Subnet struct {
	AccountId        types.String `tfsdk:"account_id"`
	Cidr             types.String `tfsdk:"cidr"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	GatewayIpAddress types.String `tfsdk:"gateway_ip_address"`
	Id               types.String `tfsdk:"id"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	Name             types.String `tfsdk:"name"`
	State            types.String `tfsdk:"state"`
	Type             types.String `tfsdk:"type"`
	VpcId            types.String `tfsdk:"vpc_id"`
	VpcName          types.String `tfsdk:"vpc_name"`
}

type SubnetResource struct {
	AccountId        types.String     `tfsdk:"account_id"`
	AllocationPools  []AllocationPool `tfsdk:"allocation_pools"`
	Cidr             types.String     `tfsdk:"cidr"`
	CreatedAt        types.String     `tfsdk:"created_at"`
	CreatedBy        types.String     `tfsdk:"created_by"`
	Description      types.String     `tfsdk:"description"`
	DhcpIpAddress    types.String     `tfsdk:"dhcp_ip_address"`
	DnsNameservers   types.List       `tfsdk:"dns_nameservers"`
	GatewayIpAddress types.String     `tfsdk:"gateway_ip_address"`
	HostRoutes       []HostRoute      `tfsdk:"host_routes"`
	Id               types.String     `tfsdk:"id"`
	ModifiedAt       types.String     `tfsdk:"modified_at"`
	ModifiedBy       types.String     `tfsdk:"modified_by"`
	Name             types.String     `tfsdk:"name"`
	State            types.String     `tfsdk:"state"`
	Type             types.String     `tfsdk:"type"`
	VpcId            types.String     `tfsdk:"vpc_id"`
	VpcName          types.String     `tfsdk:"vpc_name"`
	Tags             types.Map        `tfsdk:"tags"`
}

type AllocationPool struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type HostRoute struct {
	Destination types.String `tfsdk:"destination"`
	Nexthop     types.String `tfsdk:"nexthop"`
}

// dnsNameserversToStringSlice converts a types.List of strings into a plain
// []string for the SDK request. Null/unknown lists yield nil so the field is
// simply omitted (json:"...,omitempty").
func dnsNameserversToStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elems := list.Elements()
	if len(elems) == 0 {
		return nil
	}
	result := make([]string, 0, len(elems))
	for _, e := range elems {
		if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			result = append(result, s.ValueString())
		}
	}
	return result
}

// DnsNameserversToList builds a types.List(StringType) from a []string for state.
func DnsNameserversToList(values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	list, diags := types.ListValue(types.StringType, elems)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}

func convertAllocationPoolsToInterface(pools []AllocationPool) []interface{} {
	result := make([]interface{}, len(pools))
	for i, pool := range pools {
		result[i] = map[string]string{
			"start": pool.Start.ValueString(),
			"end":   pool.End.ValueString(),
		}
	}
	return result
}

func convertHostRoutesToInterface(routes []HostRoute) []interface{} {
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		result[i] = map[string]string{
			"destination": route.Destination.ValueString(),
			"nexthop":     route.Nexthop.ValueString(),
		}
	}
	return result
}
