package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpn/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpnVpnTunnelResource{}
	_ resource.ResourceWithConfigure   = &vpnVpnTunnelResource{}
	_ resource.ResourceWithImportState = &vpnVpnTunnelResource{}
)

// NewVpnVpnTunnelResource is a helper function to simplify the provider implementation.
func NewVpnVpnTunnelResource() resource.Resource {
	return &vpnVpnTunnelResource{}
}

// vpnVpnTunnelResource is the data source implementation.
type vpnVpnTunnelResource struct {
	config *scpsdk.Configuration
	client *vpn.Client
}

// Metadata returns the data source type name.
func (r *vpnVpnTunnelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_tunnel"
}

// Schema defines the schema for the data source.
func (r *vpnVpnTunnelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Vpn tunnel",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : 0e3dffc50eb247a1adf4f2e5c82c4f99",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Description for VPN Tunnel",
				Optional: true,
			},
			"name": schema.StringAttribute{
				Description: "Name\n" +
					"  - example : ExampleVpnTunnel1",
				Required: true,
			},
			"vpn_gateway_id": schema.StringAttribute{
				Description: "VpnGatewayId\n" +
					"  - example : b156740b6335468d8354eb9ef8eddf5a",
				Required: true,
			},
			"phase1": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"dpd_retry_interval": schema.Int32Attribute{
						Description: "DpdRetryInterval\n" +
							"  - example : 60",
						Required: true,
					},
					"ike_version": schema.Int32Attribute{
						Description: "IkeVersion\n" +
							"  - example : 2",
						Required: true,
					},
					"peer_gateway_ip": schema.StringAttribute{
						Description: "PeerGatewayIp\n" +
							"  - example : 123.0.0.2",
						Required: true,
					},
					"phase1_diffie_hellman_groups": schema.ListAttribute{
						Description: "Phase1DiffieHellmanGroups\n" +
							"  - example : [30,31,32]",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"phase1_encryptions": schema.ListAttribute{
						Description: "Phase1Encryptions\n" +
							"  - example : ['des-md5', 'chacha20poly1305-prfsha256']",
						Required:    true,
						ElementType: types.StringType,
					},
					"phase1_life_time": schema.Int32Attribute{
						Description: "Phase1LifeTime\n" +
							"  - example : 86400",
						Required: true,
					},
					"pre_shared_key": schema.StringAttribute{
						Description: "PreSharedKey\n" +
							"  - example : PreSharedKey1",
						Required: true,
					},
				},
			},
			"phase2": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"perfect_forward_secrecy": schema.StringAttribute{
						Description: "PerfectForwardSecrecy\n" +
							"  - example : ENABLE",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ENABLE", "DISABLE"),
						},
					},
					"phase2_diffie_hellman_groups": schema.ListAttribute{
						Description: "Phase2DiffieHellmanGroups\n" +
							"  - example : [30,31,32]",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"phase2_encryptions": schema.ListAttribute{
						Description: "Phase2Encryptions\n" +
							"  - example : ['des-md5', 'chacha20poly1305-prfsha256']",
						Required:    true,
						ElementType: types.StringType,
					},
					"phase2_life_time": schema.Int32Attribute{
						Description: "Phase2LifeTime\n" +
							"  - example : 86400",
						Required: true,
					},
					"remote_subnets": schema.ListAttribute{
						Description: "RemoteSubnets\n" +
							"  - example : ['10.1.1.0/24', '10.1.2.0/24', '10.1.3.0/24']",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			common.ToSnakeCase("VpnTunnel"): schema.SingleNestedAttribute{
				Description: "Vpn tunnel",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
						Description: "VpnGatewayId",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
						Description: "VpnGatewayIpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
						Description: "VpnGatewayName",
						Computed:    true,
					},
					common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
						Description: "Phase1",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"dpd_retry_interval": schema.Int32Attribute{
								Description: "DpdRetryInterval \n - example: 60",
								Computed:    true,
							},
							"ike_version": schema.Int32Attribute{
								Description: "IkeVersion \n - example: 2",
								Computed:    true,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime \n - example: 86400 ",
								Computed:    true,
							},
							"peer_gateway_ip": schema.StringAttribute{
								Description: "PeerGatewayIp \n - example: 123.0.0.2",
								Computed:    true,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Diffie-Hellman Group 목록 \n - example : [\n   \"30\",\n    \"31\",\n   \"32\"\n  ]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Proposal 목록 \n - example : [\n   \"null-md5\",\n    \"aes128gcm\",\n   \"chacha20poly1305\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
					common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
						Description: "Phase2",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"life_time": schema.Int32Attribute{
								Description: "LifeTime \n - example: 86400 ",
								Computed:    true,
							},
							"perfect_forward_secrecy": schema.StringAttribute{
								Description: "PerfectForwardSecrecy \n - example: ENABLE",
								Computed:    true,
							},
							"remote_subnets": schema.ListAttribute{
								Description: "VPN Tunnel IPSec Remote Subnets \n - example : [\n   \"10.1.1.0/24\",\n    \"10.1.2.0/24\",\n   \"10.1.3.0/24\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Diffie-Hellman Group 목록 \n - example : [\n   \"30\",\n    \"31\",\n   \"32\"\n  ]",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "VPN Tunnel ISAKMP Proposal 목록 \n - example : [\n   \"null-md5\",\n    \"aes128gcm\",\n   \"chacha20poly1305\"\n  ]",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpnVpnTunnelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Vpn
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpnVpnTunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpn.VpnTunnel1d1Resource
	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpn tunnel
	data, err := r.client.CreateVpnTunnel1d1(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Could not create vpn tunnel, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpnTunnel := data.VpnTunnel
	plan.Id = types.StringValue(vpnTunnel.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForVpnTunnelStatus(ctx, r.client, vpnTunnel.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Error waiting for vpn tunnel to become active: "+err.Error(),
		)
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Read refreshes the Terraform state with the latest data.
func (r *vpnVpnTunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state vpn.VpnTunnel1d1Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetVpnTunnel(ctx, state.Id.ValueString())

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpn tunnel",
			"Could not read vpn tunnel ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vt := data.VpnTunnel
	vtModel := vpn.VpnTunnel{
		AccountId:           types.StringValue(vt.AccountId),
		CreatedAt:           types.StringValue(vt.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(vt.CreatedBy),
		Description:         types.StringPointerValue(vt.Description.Get()),
		Id:                  types.StringValue(vt.Id),
		Phase1:              mapPhase1Detail(vt.Phase1),
		Phase2:              mapPhase2Detail(vt.Phase2),
		ModifiedAt:          types.StringValue(vt.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(vt.ModifiedBy),
		Name:                types.StringValue(vt.Name),
		State:               types.StringValue(string(vt.State)),
		Status:              types.StringValue(string(vt.Status)),
		VpcId:               types.StringValue(vt.VpcId),
		VpcName:             types.StringValue(vt.VpcName),
		VpnGatewayId:        types.StringValue(vt.VpnGatewayId),
		VpnGatewayIpAddress: types.StringValue(vt.VpnGatewayIpAddress),
		VpnGatewayName:      types.StringValue(vt.VpnGatewayName),
	}

	vtObjectValue, diags := types.ObjectValueFrom(ctx, vtModel.AttributeTypes(), vtModel)
	state.VpnTunnel = vtObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpnVpnTunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, changedPlan, state vpn.VpnTunnel1d1Resource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	var nullString types.String
	changedPlan = plan

	if state.Phase1.PeerGatewayIp.Equal(plan.Phase1.PeerGatewayIp) {
		changedPlan.Phase1.PeerGatewayIp = nullString
	}

	// Comment this condition, since convert version v1.0 -> 1.1, RemoteSubnet -> RemoteSubnets[]
	//if state.Phase2.RemoteSubnet.Equal(plan.Phase2.RemoteSubnet) {
	//	changedPlan.Phase2.RemoteSubnet = nullString
	//}

	// Update existing order
	_, err := r.client.UpdateVpnTunnel(ctx, state.Id.ValueString(), changedPlan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpn tunnel",
			"Could not update vpn tunnel, unexpected error: "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnTunnelStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating vpn tunnel",
			"Error waiting for vpn tunnel to become active: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)

	resp.State = readResp.State
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpnVpnTunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpn.VpnTunnel1d1Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpn tunnel
	err := r.client.DeleteVpnTunnel1d1(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpn tunnel",
			"Could not delete vpn tunnel, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnTunnelStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpn tunnel",
			"Error waiting for vpn tunnel to become deleted: "+err.Error(),
		)
		return
	}
}

func mapPhase1Detail(phase1 scpvpn.VpnPhase1DetailV1Dot1) vpn.VpnPhase1v1Dot1Detail {
	return vpn.VpnPhase1v1Dot1Detail{
		DpdRetryInterval:    types.Int32Value(phase1.DpdRetryInterval),
		IkeVersion:          types.Int32Value(phase1.IkeVersion),
		LifeTime:            types.Int32Value(phase1.LifeTime),
		PeerGatewayIp:       types.StringValue(phase1.PeerGatewayIp),
		DiffieHellmanGroups: convertToTypesInt32Slice(phase1.DiffieHellmanGroups),
		Encryptions:         convertToTypesStringSlice(phase1.Encryptions),
	}
}

func mapPhase2Detail(phase2 scpvpn.VpnPhase2DetailV1Dot1) vpn.VpnPhase2v1Dot1Detail {
	return vpn.VpnPhase2v1Dot1Detail{
		DiffieHellmanGroups:   convertToTypesInt32Slice(phase2.DiffieHellmanGroups),
		Encryptions:           convertToTypesStringSlice(phase2.Encryptions),
		LifeTime:              types.Int32Value(phase2.LifeTime),
		PerfectForwardSecrecy: types.StringValue(phase2.PerfectForwardSecrecy),
		RemoteSubnets:         convertToTypesStringSlice(phase2.RemoteSubnets),
	}
}

func convertToTypesInt32Slice(intSlice []int32) []types.Int32 {
	result := make([]types.Int32, len(intSlice))
	for i, v := range intSlice {
		result[i] = types.Int32Value(v)
	}
	return result
}

func convertToTypesStringSlice(strSlice []string) []types.String {
	result := make([]types.String, len(strSlice))
	for i, v := range strSlice {
		result[i] = types.StringValue(v)
	}
	return result
}

func waitForVpnTunnelStatus(ctx context.Context, vpnClient *vpn.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpnClient.GetVpnTunnel(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpnTunnel.State), nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *vpnVpnTunnelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
