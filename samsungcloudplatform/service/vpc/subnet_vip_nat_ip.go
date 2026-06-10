package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &VPCSubnetVipNatIpResource{}
	_ resource.ResourceWithConfigure = &VPCSubnetVipNatIpResource{}
)

// NewVPCSubnetVipNatIpResource is a helper function to simplify the provider implementation.
func NewVPCSubnetVipNatIpResource() resource.Resource {
	return &VPCSubnetVipNatIpResource{}
}

// VPCSubnetVipNatIpResource is the resource implementation.
type VPCSubnetVipNatIpResource struct {
	_config *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VPCSubnetVipNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vip_nat_ip"
}

// Schema defines the schema for the resource.
func (r *VPCSubnetVipNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC Subnet VIP NAT IP",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("VipId"): schema.StringAttribute{
				Description: "Subnet Vip Id \n" +
					"  - example : 0466a9448d9a4411a86055939e451c8f",
				Required: true,
			},
			common.ToSnakeCase("PublicipId"): schema.StringAttribute{
				Description: "Publicip ID \n" +
					"  - example : 12f56e27070248a6a240a497e43fbe18",
				Required: true,
			},
			common.ToSnakeCase("NatType"): schema.StringAttribute{
				Description: "NAT Type \n" +
					"  - example : PUBLIC",
				Required: true,
			},

			// Output
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Static Nat Id \n" +
					"  - example : 0009e49548154745948e9722adefbf40",
				Computed: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "Static Nat State \n" +
					"  - example : ACTIVE",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VPCSubnetVipNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *VPCSubnetVipNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.SubnetVipNatIpResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := r.client.CreateSubnetVipNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create VPC Subnet VIP NAT IP",
			fmt.Sprintf("An error occurred while creating VPC Subnet VIP NAT IP: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	plan.Id = types.StringValue(apiResponse.Id)
	plan.State = types.StringValue(apiResponse.State)

	waitForState := "ACTIVE"
	err = waitForVpcSubnetNatIpStatus(ctx, r.client, plan.SubnetId.ValueString(), plan.VipId.ValueString(), plan.PublicipId.ValueString(), []string{}, []string{waitForState})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VPC Subnet VIP NAT IP",
			"Error waiting for VPC Subnet VIP NAT IP to become active: "+err.Error(),
		)
		return
	}

	plan.State = types.StringValue(waitForState)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VPCSubnetVipNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.SubnetVipNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ShowSubnetVip(ctx, state.SubnetId.ValueString(), state.VipId.ValueString())
	if err != nil {
		// Subnet VIP was deleted => remove Subnet VIP NAT IP too
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Subnet VIP",
			"Could not read Subnet VIP Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map API response to object
	if data.SubnetVip.StaticNat.IsSet() {
		staticNat := data.SubnetVip.StaticNat.Get()
		if staticNat != nil && staticNat.PublicipId == state.PublicipId.ValueString() {
			state.Id = types.StringValue(data.SubnetVip.StaticNat.Get().Id)
			state.State = types.StringValue(data.SubnetVip.StaticNat.Get().State)
		} else {
			// Subnet VIP NAT IP was changed without us knowing so we are not managed this VIP NAT IP resource anymore
			resp.State.RemoveResource(ctx)
			return
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VPCSubnetVipNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d2.SubnetVipNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubnetVipNatIp(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting VPC subnet VIP NAT IP",
			"Could not delete VPC subnet VIP NAT IP with Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpcSubnetNatIpStatus(ctx, r.client, state.SubnetId.ValueString(), state.VipId.ValueString(), state.PublicipId.ValueString(), []string{}, []string{"DELETED"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VPC Subnet VIP NAT IP",
			"Error waiting for VPC Subnet VIP NAT IP to be deleted: "+err.Error(),
		)
		return
	}

	// Deleting the static NAT detaches the public IP from the VIP, but the
	// public IP's attachment (recorded with attached_resource_type=SUBNET) may
	// clear asynchronously. Wait until the public IP reports detached so that a
	// subsequent publicip refresh/destroy does not fail reading a SUBNET-typed
	// attachment ("SUBNET is not a valid PublicipAttachedResourceType").
	err = waitForPublicipDetached(ctx, r.client, state.PublicipId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VPC Subnet VIP NAT IP",
			"Error waiting for public IP to be detached from the VIP: "+err.Error(),
		)
		return
	}
}

// waitForPublicipDetached polls the public IP (via the v1.2 API, which supports
// the SUBNET attached-resource-type) until it is no longer attached to any
// resource. A 404 (public IP already gone) is treated as detached.
func waitForPublicipDetached(ctx context.Context, vpcClient *vpcv1d2.Client, publicipId string) error {
	return client.WaitForStatus(ctx, nil, []string{"ATTACHED"}, []string{"DETACHED"}, func() (interface{}, string, error) {
		info, statusCode, err := vpcClient.GetPublicipWithStatus(ctx, publicipId)
		if err != nil {
			if statusCode == 404 {
				return struct{}{}, "DETACHED", nil
			}
			return nil, "", err
		}
		publicip := info.GetPublicip()
		// Attachment is considered cleared when there is no attached resource id.
		if !publicip.AttachedResourceId.IsSet() || publicip.AttachedResourceId.Get() == nil || *publicip.AttachedResourceId.Get() == "" {
			return info, "DETACHED", nil
		}
		return info, "ATTACHED", nil
	})
}

func (r *VPCSubnetVipNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"VPC Subnet VIP NAT IP resource do not support update operations. The resource will not be updated.",
	)
}

func waitForVpcSubnetNatIpStatus(ctx context.Context, vpcClient *vpcv1d2.Client, subnetId string, vipId string, publicIpId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.ShowSubnetVip(ctx, subnetId, vipId)
		if err != nil {
			return nil, "", err
		}
		if info.SubnetVip.StaticNat.IsSet() {
			staticNat := info.SubnetVip.StaticNat.Get()
			if staticNat != nil && staticNat.PublicipId == publicIpId {
				return info, string(info.SubnetVip.StaticNat.Get().State), nil
			}
		}
		return info, "DELETED", nil
	})
}
