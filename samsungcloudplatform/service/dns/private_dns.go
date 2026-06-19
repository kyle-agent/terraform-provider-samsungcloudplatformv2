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
	scpdns "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/dns/1.3"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dnsPrivateDnsResource{}
	_ resource.ResourceWithConfigure = &dnsPrivateDnsResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsPrivateDnsResource() resource.Resource {
	return &dnsPrivateDnsResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsPrivateDnsResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsPrivateDnsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_private_dns" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsPrivateDnsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "PrivateDns.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("PrivateDns"): schema.SingleNestedAttribute{
				Description: "A detail of PrivateDns.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AuthDnsName"): schema.StringAttribute{
						Description: "AuthDnsName",
						Optional:    true,
					},
					common.ToSnakeCase("ConnectedVpcIds"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "ConnectedVpcIds",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("PoolId"): schema.StringAttribute{
						Description: "PoolId",
						Optional:    true,
					},
					common.ToSnakeCase("PoolName"): schema.StringAttribute{
						Description: "PoolName",
						Optional:    true,
					},
					common.ToSnakeCase("RegisteredRegion"): schema.StringAttribute{
						Description: "RegisteredRegion",
						Optional:    true,
					},
					common.ToSnakeCase("ResolverIp"): schema.StringAttribute{
						Description: "ResolverIp",
						Optional:    true,
					},
					common.ToSnakeCase("ResolverName"): schema.StringAttribute{
						Description: "ResolverName",
						Optional:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("PrivateDnsCreate"): schema.SingleNestedAttribute{
				Description: "Create PrivateDns.",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("ConnectedVpcIds"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "ConnectedVpcIds",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsPrivateDnsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dnsPrivateDnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.PrivateDnsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state dns.PrivateDnsDataSource
	listData, listErr := r.client.GetPrivateDnsList(ctx, state)

	if listErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Private Dns List",
			listErr.Error(),
		)
		return
	}

	var data *scpdns.PrivateDnsShowResponse
	var err error

	isActivated := false

	for _, privateDns := range listData.PrivateDns {
		if privateDns.State == "INACTIVE" && privateDns.Name == plan.PrivateDnsCreate.Name.ValueString() {
			// ACTIVATE
			data, err = r.client.ActivatePrivateDns(ctx, plan)
			isActivated = true
			break
		}
	}

	if !isActivated {
		// CREATE
		data, err = r.client.CreatePrivateDns(ctx, plan)
	}

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating(activating) Private Dns",
			"Could not create(activate) Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	createErr := waitForPrivateDnsStatus(ctx, r.client, data.PrivateDns.Id, []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating(activating) private dns",
			"Error creating(activating) for private dns to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetPrivateDns(ctx, data.PrivateDns.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(dataForShow.PrivateDns.Id)

	privateDnsModel := convertPrivateDns(dataForShow.PrivateDns)

	privateDnsOjbectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	plan.PrivateDns = privateDnsOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsPrivateDnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.PrivateDnsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetPrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateDnsModel := convertPrivateDns(data.PrivateDns)

	privateDnsObjectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	state.PrivateDns = privateDnsObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsPrivateDnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan

	var state dns.PrivateDnsResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdatePrivateDns(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Private Dns",
			"Could not update Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForPrivateDnsStatus(ctx, r.client, data.PrivateDns.Id, []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating private dns",
			"Error updating for private dns to become active: "+updateErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetPrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Private Dns",
			"Could not read Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	privateDnsModel := convertPrivateDns(dataForShow.PrivateDns)

	privateDnsObjectValue, diags := types.ObjectValueFrom(ctx, privateDnsModel.AttributeTypes(), privateDnsModel)
	state.PrivateDns = privateDnsObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsPrivateDnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.PrivateDnsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Private Dns
	err := r.client.DeletePrivateDns(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Private Dns",
			"Could not delete Private Dns, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// DeletePrivateDns returns 202 Accepted; the platform tears the zone down
	// asynchronously and only then detaches it from its connected VPC(s). Block
	// until the Show 404s (the zone is truly gone, hence unbound) so a downstream
	// VPC destroy does not race the detach and 409 "Cannot terminate due to
	// associated resources", leaking the VPC. A 404 is the expected terminal state.
	err = waitForPrivateDnsStatus(ctx, r.client, state.Id.ValueString(), []string{"ACTIVE", "DELETING"}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting private dns",
			"Error waiting for private dns to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForPrivateDnsStatus(ctx context.Context, privateDnsClient *dns.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := privateDnsClient.GetPrivateDns(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.PrivateDns.State), nil
	})
}
