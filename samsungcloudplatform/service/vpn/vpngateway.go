package vpn

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/vpn/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpnVpnGatewayResource{}
	_ resource.ResourceWithConfigure   = &vpnVpnGatewayResource{}
	_ resource.ResourceWithImportState = &vpnVpnGatewayResource{}
)

// NewVpnVpnGatewayResource is a helper function to simplify the provider implementation.
func NewVpnVpnGatewayResource() resource.Resource {
	return &vpnVpnGatewayResource{}
}

// vpnVpnGatewayResource is the data source implementation.
type vpnVpnGatewayResource struct {
	config *scpsdk.Configuration
	client *vpn.Client
}

// Metadata returns the data source type name.
func (r *vpnVpnGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_gateway"
}

// Schema defines the schema for the data source.
func (r *vpnVpnGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Vpn gateway",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.\n" +
					"  - example : b156740b6335468d8354eb9ef8eddf5a",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Description for VPN Gateway",
				Optional: true,
			},
			common.ToSnakeCase("IpAddress"): schema.StringAttribute{
				Description: "Ip Address\n" +
					"  - example : 123.0.0.1",
				Required: true,
			},
			common.ToSnakeCase("IpId"): schema.StringAttribute{
				Description: "Identifier of the IP\n" +
					"  - example : fcde872f75c145a0893d656cc698f13e",
				Required: true,
			},
			common.ToSnakeCase("IpType"): schema.StringAttribute{
				Description: "Type of IP\n" +
					"  - example : PUBLIC",
				Required: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name\n" +
					"  - example : ExampleVpnGW1",
				Required: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "Identifier of the VPC\n" +
					"  - example : ceb44ea5ecb34a49b16495f9a63b0718",
				Required: true,
			},
			common.ToSnakeCase("VpnGateway"): schema.SingleNestedAttribute{
				Description: "Vpn gateway",
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
					common.ToSnakeCase("IpAddress"): schema.StringAttribute{
						Description: "IpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("IpId"): schema.StringAttribute{
						Description: "IpId",
						Computed:    true,
					},
					common.ToSnakeCase("IpType"): schema.StringAttribute{
						Description: "IpType",
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
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpnVpnGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpnVpnGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpn.VpnGatewayResource
	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpn gateway
	data, err := r.client.CreateVpnGateway(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpn gateway",
			"Could not create vpn gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpnGateway := data.VpnGateway
	plan.Id = types.StringValue(vpnGateway.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForVpnGatewayStatus(ctx, r.client, vpnGateway.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpn gateway",
			"Error waiting for vpn gateway to become active: "+err.Error(),
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
func (r *vpnVpnGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpn.VpnGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vgModel := createVpnGatewayModel(data)

	vgObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	state.VpnGateway = vgObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpnVpnGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpn.VpnGatewayResource
	diags := req.Plan.Get(ctx, &state) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpnGateway(ctx, state.Id.ValueString(), state) // client 를 호출한다.
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpnGateway as UpdateVpnGateway items are not populated.
	data, err := r.client.GetVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading vpn gateway",
			"Could not read vpn gateway ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	vgModel := createVpnGatewayModel(data)

	vgObjectValue, diags := types.ObjectValueFrom(ctx, vgModel.AttributeTypes(), vgModel)
	state.VpnGateway = vgObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpnVpnGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpn.VpnGatewayResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpn gateway
	err := r.client.DeleteVpnGateway(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpn gateway",
			"Could not delete vpn gateway, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnGatewayStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpn gateway",
			"Error waiting for vpn gateway to become deleted: "+err.Error(),
		)
		return
	}
}

func createVpnGatewayModel(data *scpvpn.VpnGatewayShowResponse) vpn.VpnGateway {
	vg := data.VpnGateway
	return vpn.VpnGateway{
		AccountId:   types.StringValue(vg.AccountId),
		CreatedAt:   types.StringValue(vg.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(vg.CreatedBy),
		Description: types.StringPointerValue(vg.Description.Get()),
		Id:          types.StringValue(vg.Id),
		IpAddress:   types.StringValue(vg.IpAddress),
		IpId:        types.StringValue(vg.IpId),
		IpType:      types.StringValue(vg.IpType),
		ModifiedAt:  types.StringValue(vg.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(vg.ModifiedBy),
		Name:        types.StringValue(vg.Name),
		State:       types.StringValue(string(vg.State)),
		VpcId:       types.StringValue(vg.VpcId),
		VpcName:     types.StringValue(vg.VpcName),
	}
}

func waitForVpnGatewayStatus(ctx context.Context, vpnClient *vpn.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpnClient.GetVpnGateway(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpnGateway.State), nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *vpnVpnGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
