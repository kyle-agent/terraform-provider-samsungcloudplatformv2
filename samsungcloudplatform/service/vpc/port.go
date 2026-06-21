package vpc

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpcPortResource{}
	_ resource.ResourceWithConfigure   = &vpcPortResource{}
	_ resource.ResourceWithImportState = &vpcPortResource{}
)

// NewVpcPortResource is a helper function to simplify the provider implementation.
func NewVpcPortResource() resource.Resource {
	return &vpcPortResource{}
}

// vpcPortResource is the data source implementation.
type vpcPortResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcPortResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_port"
}

// Schema defines the schema for the data source.
func (r *vpcPortResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "port",
		MarkdownDescription: "port",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.",
				MarkdownDescription: "Identifier of the resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description:         "AccountId",
				MarkdownDescription: "AccountId",
				Computed:            true,
			},
			common.ToSnakeCase("AttachedResourceId"): schema.StringAttribute{
				Description:         "AttachedResourceId",
				MarkdownDescription: "AttachedResourceId",
				Computed:            true,
			},
			common.ToSnakeCase("AttachedResourceType"): schema.StringAttribute{
				Description:         "AttachedResAttachedResourceType",
				MarkdownDescription: "AttachedResAttachedResourceType",
				Computed:            true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description:         "CreatedAt",
				MarkdownDescription: "CreatedAt",
				Computed:            true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Port description\n" +
					"  - maxLength : 50",
				MarkdownDescription: "Description\n" +
					"  - example : Port description\n" +
					"  - maxLength : 50",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Default: stringdefault.StaticString(""),
			},
			common.ToSnakeCase("FixedIpAddress"): schema.StringAttribute{
				Description: "Fixed IP Address \n" +
					"  - example : 172.24.4.2",
				MarkdownDescription: "Fixed IP Address \n" +
					"  - example : 172.24.4.2",
				Optional: true,
			},
			common.ToSnakeCase("MacAddress"): schema.StringAttribute{
				Description:         "MacAddress",
				MarkdownDescription: "MacAddress",
				Computed:            true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
				Description:         "ModifiedAt",
				MarkdownDescription: "ModifiedAt",
				Computed:            true,
			},
			common.ToSnakeCase("SecurityGroups"): schema.ListNestedAttribute{
				Description:         "Security groups",
				MarkdownDescription: "Security groups",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Port Name \n" +
					"  - example : portName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "Port Name \n" +
					"  - example : portName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Required: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description:         "State",
				MarkdownDescription: "State",
				Computed:            true,
			},
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "Subnet ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				MarkdownDescription: "Subnet ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("SubnetName"): schema.StringAttribute{
				Description:         "SubnetName",
				MarkdownDescription: "SubnetName",
				Computed:            true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description:         "VpcId",
				MarkdownDescription: "VpcId",
				Computed:            true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description:         "VpcName",
				MarkdownDescription: "VpcName",
				Computed:            true,
			},
			common.ToSnakeCase("VirtualIpAddresses"): schema.ListAttribute{
				Description:         "Virtual IP Addresses",
				MarkdownDescription: "Virtual IP Addresses",
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcPortResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.PortResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new port
	data, err := r.client.CreatePort(ctx, plan)
	if err != nil || !data.Port.IsSet() {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating port",
			"Could not create port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpc.MapPort(data.Port.Get(), &plan)
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *vpcPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpcPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.PortResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetPort(ctx, state.Id.ValueString())
	if err != nil || !data.Port.IsSet() {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading port",
			"Could not read port ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Set refreshed state
	vpc.MapPort(data.Port.Get(), &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.PortResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	data, err := r.client.UpdatePort(ctx, state.Id.ValueString(), state)
	if err != nil || !data.Port.IsSet() {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating port",
			"Could not update port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpc.MapPort(data.Port.Get(), &state)
	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.PortResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing port
	err := r.client.DeletePort(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting port",
			"Could not delete port, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
