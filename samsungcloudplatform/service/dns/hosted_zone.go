package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dnsHostedZoneResource{}
	_ resource.ResourceWithConfigure = &dnsHostedZoneResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsHostedZoneResource() resource.Resource {
	return &dnsHostedZoneResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsHostedZoneResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsHostedZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_hosted_zone" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsHostedZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "HostedZone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Zone"): schema.SingleNestedAttribute{
				Description: "A detail of HostedZone.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("HostedZoneType"): schema.StringAttribute{
						Description: "HostedZoneType",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "PoolId",
						Optional:    true,
					},
					common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
						Description: "PrivateDnsId",
						Optional:    true,
					},
					common.ToSnakeCase("PrivateDnsName"): schema.StringAttribute{
						Description: "PrivateDnsName",
						Optional:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Optional:    true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "Ttl",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("HostedZoneCreate"): schema.SingleNestedAttribute{
				Description: "Create HostedZone.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("PrivateDnsId"): schema.StringAttribute{
						Description: "PrivateDnsId",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsHostedZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Dns
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *dnsHostedZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.HostedZoneResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateHostedZone(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating HostedZone",
			"Could not create HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	createErr := waitForHostedZoneStatus(ctx, r.client, data.Id, []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating hosted zone",
			"Error creating for hosted zone to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetHostedZone(ctx, data.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.Id)

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*dataForShow)

	hostedZoneOjbectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	plan.Zone = hostedZoneOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsHostedZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.HostedZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*data)

	hostedZoneObjectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	state.Zone = hostedZoneObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsHostedZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var oldState dns.HostedZoneResource
	req.State.Get(ctx, &oldState)
	var state dns.HostedZoneResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if checkModifiedFieldsExcludingDescription(oldState, state) {
		resp.Diagnostics.AddError(
			"Error updating HostedZone",
			"Hosted zones cannot be modified except for the description field.",
		)
		return
	}

	data, err := r.client.UpdateHostedZone(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating HostedZone",
			"Could not update HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForHostedZoneStatus(ctx, r.client, data.Id, []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating hosted zone",
			"Error updating for hosted zone to become active: "+updateErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading HostedZone",
			"Could not read HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	hostedZoneModel := convertHostedZoneShowResponseV1Dot3ToHostedZone(*dataForShow)

	hostedZoneObjectValue, diags := types.ObjectValueFrom(ctx, hostedZoneModel.AttributeTypes(), hostedZoneModel)
	state.Zone = hostedZoneObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsHostedZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.HostedZoneResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHostedZone(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting HostedZone",
			"Could not delete HostedZone, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// DeleteHostedZone returns 202 Accepted and tears down asynchronously. Block
	// until the Show 404s so the parent private DNS Delete that follows in the
	// dependency graph does not fire while this zone is still attached (which 409s
	// the parent delete and ultimately leaks the bootstrap VPC). 404 is terminal.
	err = waitForHostedZoneStatus(ctx, r.client, state.Id.ValueString(), []string{"ACTIVE", "DELETING"}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting hosted zone",
			"Error waiting for hosted zone to become deleted: "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func waitForHostedZoneStatus(ctx context.Context, hostedZoneClient *dns.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := hostedZoneClient.GetHostedZone(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Status), nil
	})
}

func checkModifiedFieldsExcludingDescription(oldState dns.HostedZoneResource, newState dns.HostedZoneResource) bool {
	oldHostedZone := oldState.HostedZoneCreate
	newHostedZone := newState.HostedZoneCreate

	if oldHostedZone.Type != newHostedZone.Type {
		return true
	}
	if oldHostedZone.Name != newHostedZone.Name {
		return true
	}
	if oldHostedZone.PrivateDnsId != newHostedZone.PrivateDnsId {
		return true
	}
	return false
}
