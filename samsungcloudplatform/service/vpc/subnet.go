package vpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"regexp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcSubnetResource{}
	_ resource.ResourceWithConfigure   = &vpcSubnetResource{}
	_ resource.ResourceWithImportState = &vpcSubnetResource{}
)

// NewVpcSubnetResource is a helper function to simplify the provider implementation.
func NewVpcSubnetResource() resource.Resource {
	return &vpcSubnetResource{}
}

// vpcSubnetResource is the data source implementation.
type vpcSubnetResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcSubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet"
}

// Schema defines the schema for the data source.
func (r *vpcSubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "subnet",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("AllocationPools"): schema.ListNestedAttribute{
				Description: "Allocation Pools \n" +
					"  - example : [{ \"start\": \"10.0.0.2\", \"end\": \"10.0.0.254\" }]",
				MarkdownDescription: "Allocation Pools \n" +
					"  - example : [{ \"start\": \"10.0.0.2\", \"end\": \"10.0.0.254\" }]",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("end"): schema.StringAttribute{
							Description: "End",
							Required:    true,
						},
						common.ToSnakeCase("start"): schema.StringAttribute{
							Description: "Start",
							Required:    true,
						},
					},
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description:         "AccountId",
				MarkdownDescription: "AccountId",
				Computed:            true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "Suabnet CIDR\n" +
					"  - example : 192.167.1.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
				MarkdownDescription: "Suabnet CIDR\n" +
					"  - example : 192.167.1.0/24 \n" +
					"  - maxMask : /28\n" +
					"  - minMask : /16",
				Required: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description:         "CreatedAt",
				MarkdownDescription: "CreatedAt",
				Computed:            true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description:         "CreatedBy",
				MarkdownDescription: "CreatedBy",
				Computed:            true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description:         "Description\n  - maxLength: 50\n  - example: Subnet Description",
				MarkdownDescription: "Description\n  - maxLength: 50\n  - example: Subnet Description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			common.ToSnakeCase("dhcp_ip_address"): schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "DHCP automatically assigned IP\n  - example: 192.167.1.2",
				MarkdownDescription: "DHCP automatically assigned IP\n  - example: 192.167.1.2",
			},
			common.ToSnakeCase("DnsNameservers"): schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "DNS Name Servers\n  - example: [\"1.1.1.1\", \"2.2.2.2\"]",
				MarkdownDescription: "DNS Name Servers\n  - example: [\"1.1.1.1\", \"2.2.2.2\"]",
			},
			common.ToSnakeCase("GatewayIpAddress"): schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Gateway IP Address\n  - example: 192.167.1.1",
				MarkdownDescription: "Gateway IP Address\n  - example: 192.167.1.1",
			},
			common.ToSnakeCase("HostRoutes"): schema.ListNestedAttribute{
				Description: "HostRoutes \n" +
					"  - example : [{ \"destination\": \"192.168.24.0/24\", \"nexthop\": \"192.168.20.5\" }]",
				MarkdownDescription: "HostRoutes \n" +
					"  - example : [{ \"destination\": \"192.168.24.0/24\", \"nexthop\": \"192.168.20.5\" }]",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Destination"): schema.StringAttribute{
							Description: "Destination",
							Required:    true,
						},
						common.ToSnakeCase("Nexthop"): schema.StringAttribute{
							Description: "Nexthop",
							Required:    true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.",
				MarkdownDescription: "Identifier of the resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description:         "ModifiedAt",
				MarkdownDescription: "ModifiedAt",
				Computed:            true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description:         "ModifiedBy",
				MarkdownDescription: "ModifiedBy",
				Computed:            true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Subnet Name \n" +
					"  - example : subnetName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "Subnet Name \n" +
					"  - example : subnetName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9-]*$"), "Enter 3 -20 chars. (English, number, hyphen)"),
				},
				Required: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "State",
				MarkdownDescription: "State",
				Computed:            true,
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Subnet Type \n" +
					"  - example : GENERAL | LOCAL | VPC_ENDPOINT",
				MarkdownDescription: "Subnet Type \n" +
					"  - example : GENERAL | LOCAL | VPC_ENDPOINT",
				Required: true,
			},
			common.ToSnakeCase("VpcID"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "VPC ID \n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description:         "VpcName",
				MarkdownDescription: "VpcName",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcSubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VpcV1Dot2
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpcSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.SubnetResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subnet
	data, err := r.client.CreateSubnet(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating subnet",
			"Could not create subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	subnet := data.Subnet
	plan.Id = types.StringValue(subnet.Id)
	plan.DnsNameservers = vpc.DnsNameserversToList(subnet.DnsNameservers)
	diags = resp.State.Set(ctx, plan)

	err = waitForSubnetStatus(ctx, r.client, subnet.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subnet",
			"Error waiting for subnet to become active: "+err.Error(),
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
// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *vpcSubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpcSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.SubnetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetSubnet(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading subnet",
			"Could not read subnet ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Set refreshed state
	subnet := data.Subnet
	state.Id = types.StringValue(subnet.Id)
	state.AccountId = types.StringValue(subnet.AccountId)
	state.GatewayIpAddress = types.StringPointerValue(subnet.GatewayIpAddress.Get())
	state.VpcName = types.StringValue(subnet.VpcName)
	state.Description = types.StringPointerValue(subnet.Description.Get())
	state.State = types.StringValue(string(subnet.State))
	state.CreatedAt = types.StringValue(subnet.CreatedAt.Format(time.RFC3339))
	state.CreatedBy = types.StringValue(subnet.CreatedBy)
	state.ModifiedAt = types.StringValue(subnet.ModifiedAt.Format(time.RFC3339))
	state.ModifiedBy = types.StringValue(subnet.ModifiedBy)
	state.DhcpIpAddress = types.StringPointerValue(subnet.DhcpIpAddress.Get())
	state.DnsNameservers = vpc.DnsNameserversToList(subnet.DnsNameservers)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.SubnetResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateSubnet(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating subnet",
			"Could not update subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.SubnetResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subnet
	err := r.client.DeleteSubnet(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting subnet",
			"Could not delete subnet, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForSubnetStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting subnet",
			"Error waiting for subnet to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForSubnetStatus(ctx context.Context, vpcClient *vpc.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.GetSubnet(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.Subnet.State), nil
	})
}
