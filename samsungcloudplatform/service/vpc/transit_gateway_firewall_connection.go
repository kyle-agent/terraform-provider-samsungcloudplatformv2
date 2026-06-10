package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpcv1d2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &tgwFirewallConnectionResource{}
	_ resource.ResourceWithConfigure = &tgwFirewallConnectionResource{}
)

// NewVpcTgwFirewallConnectionResource is a helper function to simplify the provider implementation.
func NewVpcTgwFirewallConnectionResource() resource.Resource {
	return &tgwFirewallConnectionResource{}
}

type tgwFirewallConnectionResource struct {
	config     *scpsdk.Configuration
	clientv1d2 *vpcv1d2.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *tgwFirewallConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_firewall_connection"
}

// Schema defines the schema for the data source.
func (r *tgwFirewallConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit Gateway Firewall Connection",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "Transit Gateway ID",
				Required:    true,
			},

			// Output
			common.ToSnakeCase("TransitGateway"): schema.SingleNestedAttribute{
				Description: "Transit Gateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "Transit Gateway Port Bandwidth\n" +
							"  - example: 1",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At \n" +
							"  - example : 2024-05-17T00:23:17Z \n",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Transit Gateway Description\n" +
							"  - example : TransitGateway Description",
						Computed: true,
					},
					common.ToSnakeCase("FirewallConnectionState"): schema.StringAttribute{
						Description: "Firewall Connection State\n" +
							"  - enum: ATTACHING | ACTIVE | DETACHING | DELETED | INACTIVE | ERROR\n" +
							"  - example: INACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
						Description: "Firewall ID\n" +
							"  - example: ['bbb93aca123f4bb2b2c0f206f4a86b2b']",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Transit Gateway ID\n" +
							"  - example: fe860e0af0c04dcd8182b84f907f31f4",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At \n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Transit Gateway Name\n" +
							"  - minLength: 3\n" +
							"  - maxLength: 20\n" +
							"  - pattern: ^[a-zA-Z0-9-]*$\n" +
							"  - example: TransitGatewayName",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State\n" +
							"  - enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
						Description: "Uplink Enabled?\n" +
							"  - default: false\n" +
							"  - example: false",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *tgwFirewallConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clientv1d2 = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

func (r *tgwFirewallConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.TransitGatewayFirewallConnectionResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := r.clientv1d2.CreateTransitGatewayFirewallConnection(ctx, plan.TransitGatewayId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Transit Gateway Firewall Connection",
			"Could not create Transit Gateway Firewall Connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map API response to object
	tgw := vpcv1d2.TransitGateway{
		Id:            types.StringValue(data.TransitGateway.Id),
		Name:          types.StringValue(data.TransitGateway.Name),
		AccountId:     types.StringValue(data.TransitGateway.AccountId),
		CreatedAt:     types.StringValue(data.TransitGateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(data.TransitGateway.CreatedBy),
		ModifiedAt:    types.StringValue(data.TransitGateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(data.TransitGateway.ModifiedBy),
		State:         types.StringValue(string(data.TransitGateway.State)),
		UplinkEnabled: types.BoolPointerValue(data.TransitGateway.UplinkEnabled),
	}
	if data.TransitGateway.Description.IsSet() {
		if val := data.TransitGateway.Description.Get(); val != nil {
			tgw.Description = types.StringValue(*val)
		}
	}
	if data.TransitGateway.FirewallIds.IsSet() {
		if val := data.TransitGateway.FirewallIds.Get(); val != nil {
			tgw.FirewallIds = types.StringValue(*val)
		}
	}
	if data.TransitGateway.Bandwidth.IsSet() {
		if val := data.TransitGateway.Bandwidth.Get(); val != nil {
			tgw.Bandwidth = types.Int32PointerValue(val)
		}
	}
	if data.TransitGateway.FirewallConnectionState.IsSet() {
		if desc := data.TransitGateway.FirewallConnectionState.Get(); desc != nil {
			tgw.FirewallConnectionState = types.StringValue(string(*desc))
		}
	}

	// Firewall connection states (SDK TransitGatewayFirewallConnectionState):
	// ATTACHING, ACTIVE, DETACHING, INACTIVE, DELETED, ERROR. Listing the
	// transitional state in Pending lets StateChangeConf short-circuit on a
	// parked/ERROR state instead of polling for the full timeout (issue #76).
	err = waitForFirewallConnectionState(ctx, r.clientv1d2, plan.TransitGatewayId.ValueString(),
		[]string{"ATTACHING"},
		[]string{common.ActiveState})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become active: "+err.Error(),
		)
		return
	}

	tgw.FirewallConnectionState = types.StringValue("ACTIVE")
	tgwObjectValue, _ := types.ObjectValueFrom(ctx, tgw.AttributeTypes(), tgw)
	plan.TransitGateway = tgwObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tgwFirewallConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.TransitGatewayFirewallConnectionResource

	diags := req.State.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tgwResp, err := r.clientv1d2.GetTransitGatewayInfo(ctx, state.TransitGatewayId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GetTransitGatewayInfo",
			err.Error(),
		)
		return
	}

	if tgwResp == nil {
		return
	}
	data := tgwResp.TransitGateway

	transitGateway := vpcv1d2.TransitGateway{
		Id:            types.StringValue(data.Id),
		Description:   types.StringPointerValue(data.Description.Get()),
		Name:          types.StringValue(data.Name),
		AccountId:     types.StringValue(data.AccountId),
		Bandwidth:     types.Int32PointerValue(data.Bandwidth.Get()),
		CreatedAt:     types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(data.CreatedBy),
		FirewallIds:   types.StringPointerValue(data.FirewallIds.Get()),
		ModifiedAt:    types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(data.ModifiedBy),
		State:         types.StringValue(string(data.State)),
		UplinkEnabled: types.BoolPointerValue(data.UplinkEnabled),
	}

	tgwObjectValue, _ := types.ObjectValueFrom(ctx, transitGateway.AttributeTypes(), transitGateway)
	state.TransitGateway = tgwObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tgwFirewallConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Implemented",
		"TGW Firewall Connection update function is not yet implemented.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tgwFirewallConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d2.TransitGatewayFirewallConnectionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.clientv1d2.DeleteTransitGatewayFirewallConnection(ctx, state.TransitGatewayId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting TGW Firewall Connection",
			"Could not delete TGW Firewall Connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForFirewallConnectionState(ctx, r.clientv1d2, state.TransitGatewayId.ValueString(),
		[]string{"DETACHING", common.ActiveState},
		[]string{common.InActiveState})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating transit gateway routing rule",
			"Error waiting for transit gateway routing rule to become active: "+err.Error(),
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

func waitForFirewallConnectionState(ctx context.Context, r *vpcv1d2.Client, transitGatewayId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := r.GetTransitGatewayInfo(ctx, transitGatewayId)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.TransitGateway.GetFirewallConnectionState()), nil
	})
}
