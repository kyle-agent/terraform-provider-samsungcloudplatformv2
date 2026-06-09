package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scploadbalancer "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/loadbalancer/1.3"
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
	_ resource.Resource              = &loadbalancerLoadbalancerPublicNatIpResource{}
	_ resource.ResourceWithConfigure = &loadbalancerLoadbalancerPublicNatIpResource{}
)

// NewLoadbalancerLoadbalancerPublicNatIpResource is a helper function to simplify the provider implementation.
func NewLoadbalancerLoadbalancerPublicNatIpResource() resource.Resource {
	return &loadbalancerLoadbalancerPublicNatIpResource{}
}

// loadbalancerLoadbalancerPublicNatIpResource is the data source implementation.
type loadbalancerLoadbalancerPublicNatIpResource struct {
	config  *scpsdk.Configuration
	client  *loadbalancer.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_loadbalancer_public_nat_ip"
}

// Schema defines the schema for the data source.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Loadbalancer Public NAT.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("LoadbalancerId"): schema.StringAttribute{
				Description: "LoadbalancerId",
				Required:    true,
			},
			common.ToSnakeCase("LoadbalancerPublicNatIp"): schema.SingleNestedAttribute{
				Description: "A detail of public NAT.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Optional:    true,
					},
					common.ToSnakeCase("SubnetId"): schema.StringAttribute{
						Description: "SubnetId",
						Optional:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Optional:    true,
					},
					common.ToSnakeCase("ActionType"): schema.StringAttribute{
						Description: "ActionType",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("ExternalIpAddress"): schema.StringAttribute{
						Description: "ExternalIpAddress",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("InternalIpAddress"): schema.StringAttribute{
						Description: "InternalIpAddress",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerId"): schema.StringAttribute{
						Description: "OwnerId",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerName"): schema.StringAttribute{
						Description: "OwnerName",
						Optional:    true,
					},
					common.ToSnakeCase("OwnerType"): schema.StringAttribute{
						Description: "OwnerType",
						Optional:    true,
					},
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "PublicipId",
						Optional:    true,
					},
					common.ToSnakeCase("ServiceIpPortId"): schema.StringAttribute{
						Description: "ServiceIpPortId",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
					common.ToSnakeCase("vpc_id"): schema.StringAttribute{
						Description: "vpc_id",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("StaticNatCreate"): schema.SingleNestedAttribute{
				Description: "Create Loadbalancer static NAT.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("PublicipId"): schema.StringAttribute{
						Description: "PublicipId",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.LoadBalancer
}

// Create creates the resource and sets the initial Terraform state.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan loadbalancer.LoadbalancerPublicNatIpResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Lb Static NAT
	data, err := r.client.CreateLoadbalancerPublicNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Public NAT",
			"Could not create Public NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = virtualserverutil.ToNullableStringValue(data.StaticNat.Id.Get())

	// Map response body to schema and populate Computed attribute values
	staticNatModel := createLoadbalancerNatModel(data)
	staticNatObjectValue, diags := types.ObjectValueFrom(ctx, staticNatModel.AttributeTypes(), staticNatModel)
	plan.LoadbalancerPublicNatIp = staticNatObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *loadbalancerLoadbalancerPublicNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state loadbalancer.LoadbalancerPublicNatIpResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadbalancerId := state.LoadbalancerId.ValueString()

	// Delete (detach) the public NAT from the load balancer.
	err := r.client.DeleteLoadbalancerPublicNatIp(ctx, loadbalancerId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting LB Public NAT",
			"Could not delete LB Public NAT, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// The detach is asynchronous. If we return immediately, the public IP is still
	// ATTACHED and a later publicip/LB delete in the same `terraform destroy` fails
	// with 409 ("PublicIP state is not deletable (ATTACHED)" /
	// "Cannot delete the loadbalancer due to associated resources"). Wait until the
	// NAT is fully removed: poll Show until it 404s (gone) so the publicip becomes
	// deletable and the LB no longer has an associated NAT.
	err = waitForPublicNatIpRemoved(ctx, r.client, loadbalancerId)
	if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
		resp.Diagnostics.AddError(
			"Error Deleting LB Public NAT",
			"Error waiting for LB Public NAT to be detached: "+err.Error(),
		)
		return
	}
}

// waitForPublicNatIpRemoved polls the load balancer's public NAT until the Show
// call returns a not-found (404) error, indicating the NAT has been fully detached
// and the associated public IP is deletable again.
func waitForPublicNatIpRemoved(ctx context.Context, lbClient *loadbalancer.Client, loadbalancerId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"DELETED"}, func() (interface{}, string, error) {
		info, err := lbClient.GetLoadbalancerPublicNatIp(ctx, loadbalancerId)
		if err != nil {
			// 404 / not found means the NAT is gone, which is the terminal success
			// state. Surface the error so WaitForStatus stops; the caller treats a
			// 404 error as success.
			return nil, "", err
		}
		return info, info.StaticNat.State, nil
	})
}

func createLoadbalancerNatModel(data *scploadbalancer.StaticNatCreateResponse) loadbalancer.LoadbalancerPublicNatIpDetail {
	lbStaticNat := data.StaticNat
	return loadbalancer.LoadbalancerPublicNatIpDetail{
		AccountId:         virtualserverutil.ToNullableStringValue(lbStaticNat.AccountId.Get()),
		ActionType:        virtualserverutil.ToNullableStringValue(lbStaticNat.ActionType.Get()),
		CreatedAt:         types.StringValue(lbStaticNat.CreatedAt.Format(time.RFC3339)),
		CreatedBy:         types.StringValue(lbStaticNat.CreatedBy),
		Description:       virtualserverutil.ToNullableStringValue(lbStaticNat.Description.Get()),
		ExternalIpAddress: virtualserverutil.ToNullableStringValue(lbStaticNat.ExternalIpAddress.Get()),
		Id:                virtualserverutil.ToNullableStringValue(lbStaticNat.Id.Get()),
		InternalIpAddress: virtualserverutil.ToNullableStringValue(lbStaticNat.InternalIpAddress.Get()),
		ModifiedAt:        types.StringValue(lbStaticNat.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:        types.StringValue(lbStaticNat.ModifiedBy),
		Name:              virtualserverutil.ToNullableStringValue(lbStaticNat.Name.Get()),
		OwnerId:           virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerId.Get()),
		OwnerName:         virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerName.Get()),
		OwnerType:         virtualserverutil.ToNullableStringValue(lbStaticNat.OwnerType.Get()),
		PublicipId:        virtualserverutil.ToNullableStringValue(lbStaticNat.PublicipId.Get()),
		ServiceIpPortId:   virtualserverutil.ToNullableStringValue(lbStaticNat.ServiceIpPortId.Get()),
		State:             virtualserverutil.ToNullableStringValue(lbStaticNat.State.Get()),
		SubnetId:          virtualserverutil.ToNullableStringValue(lbStaticNat.AccountId.Get()),
		Type:              virtualserverutil.ToNullableStringValue(lbStaticNat.Type.Get()),
		VpcId:             virtualserverutil.ToNullableStringValue(lbStaticNat.VpcId.Get()),
	}
}
