package vpcv1

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-vpc"

// ------------------- VPC PEERING -------------------//.

func (m VpcPeering) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                          types.StringType,
		"name":                        types.StringType,
		"account_type":                types.StringType,
		"approver_vpc_account_id":     types.StringType,
		"approver_vpc_id":             types.StringType,
		"approver_vpc_name":           types.StringType,
		"description":                 types.StringType,
		"created_at":                  types.StringType,
		"created_by":                  types.StringType,
		"modified_at":                 types.StringType,
		"modified_by":                 types.StringType,
		"requester_vpc_account_id":    types.StringType,
		"delete_requester_account_id": types.StringType,
		"requester_vpc_id":            types.StringType,
		"requester_vpc_name":          types.StringType,
		"state":                       types.StringType,
	}
}

type VpcPeeringResource struct {
	Id                   types.String `tfsdk:"id"`
	ApproverVpcAccountId types.String `tfsdk:"approver_vpc_account_id"`
	Name                 types.String `tfsdk:"name"`
	ApproverVpcId        types.String `tfsdk:"approver_vpc_id"`
	ApproverVpcName      types.String `tfsdk:"approver_vpc_name"`
	RequesterVpcId       types.String `tfsdk:"requester_vpc_id"`
	Description          types.String `tfsdk:"description"`
	VpcPeering           types.Object `tfsdk:"vpc_peering"`
	Tags                 types.Map    `tfsdk:"tags"`
}

type VpcPeeringListDataSource struct {
	Size             types.Int32  `tfsdk:"size"`
	Page             types.Int32  `tfsdk:"page"`
	Sort             types.String `tfsdk:"sort"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	RequesterVpcId   types.String `tfsdk:"requester_vpc_id"`
	RequesterVpcName types.String `tfsdk:"requester_vpc_name"`
	ApproverVpcId    types.String `tfsdk:"approver_vpc_id"`
	ApproverVpcName  types.String `tfsdk:"approver_vpc_name"`
	AccountType      types.String `tfsdk:"account_type"`
	State            types.String `tfsdk:"state"`
	VpcPeerings      []VpcPeering `tfsdk:"vpc_peerings"`
}

type VpcPeeringDataSource struct {
	Id         types.String `tfsdk:"id"`
	VpcPeering types.Object `tfsdk:"vpc_peering"`
}

type VpcPeeringAccountType string

type VpcPeering struct {
	AccountType              types.String `tfsdk:"account_type"`
	ApproverVpcAccountId     types.String `tfsdk:"approver_vpc_account_id"`
	ApproverVpcId            types.String `tfsdk:"approver_vpc_id"`
	ApproverVpcName          types.String `tfsdk:"approver_vpc_name"`
	CreatedAt                types.String `tfsdk:"created_at"`
	CreatedBy                types.String `tfsdk:"created_by"`
	Description              types.String `tfsdk:"description"`
	Id                       types.String `tfsdk:"id"`
	ModifiedAt               types.String `tfsdk:"modified_at"`
	ModifiedBy               types.String `tfsdk:"modified_by"`
	Name                     types.String `tfsdk:"name"`
	RequesterVpcAccountId    types.String `tfsdk:"requester_vpc_account_id"`
	RequesterVpcId           types.String `tfsdk:"requester_vpc_id"`
	RequesterVpcName         types.String `tfsdk:"requester_vpc_name"`
	DeleteRequesterAccountId types.String `tfsdk:"delete_requester_account_id"`
	State                    types.String `tfsdk:"state"`
}

// ------------------- VPC Peering Approval -------------------//
type VpcPeeringApprovalResource struct {
	// Input
	VpcPeeringID types.String `tfsdk:"vpc_peering_id"`
	Type         types.String `tfsdk:"type"`

	// Output
	VpcPeering types.Object `tfsdk:"vpc_peering"`
}
