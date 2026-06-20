package vpc

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &vpcVpcResource{}
	_ resource.ResourceWithConfigure   = &vpcVpcResource{}
	_ resource.ResourceWithImportState = &vpcVpcResource{}
)

// NewVpcVpcResource is a helper function to simplify the provider implementation.
func NewVpcVpcResource() resource.Resource {
	return &vpcVpcResource{}
}

// vpcVpcResource is the data source implementation.
type vpcVpcResource struct {
	config  *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *vpcVpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_vpc"
}

// Schema defines the schema for the data source.
func (r *vpcVpcResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "vpc",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "VPC CIDR\n" +
					"  - example : 192.167.0.0/18\n" +
					"  - maxMask : /24\n" +
					"  - minMask : /16",
				MarkdownDescription: "VPC CIDR\n" +
					"  - example : 192.167.0.0/18\n" +
					"  - maxMask : /24\n" +
					"  - minMask : /16",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : VPC description\n" +
					"  - maxLength : 50",
				MarkdownDescription: "Description\n" +
					"  - example : VPC description\n" +
					"  - maxLength : 50",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Description:         "Identifier of the resource.",
				MarkdownDescription: "Identifier of the resource.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "VPC Name \n" +
					"  - example : vpcName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				MarkdownDescription: "VPC Name \n" +
					"  - example : vpcName\n" +
					"  - maxLength : 20\n" +
					"  - minLength : 3\n" +
					"  - pattern : ^[a-zA-Z0-9-]+$",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9-]*$"), "Enter 3 -20 chars. (English, number, hyphen)"),
				},
				Required: true,
			},
			"tags": tag.ResourceSchema(),
			common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
				Description:         "Vpc",
				MarkdownDescription: "Vpc",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description:         "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						MarkdownDescription: "Account ID\n  - example: f1e6c81a2b054582878cb9724dc2ce9f",
						Computed:            true,
					},
					common.ToSnakeCase("CidrCount"): schema.Int32Attribute{
						Description:         "Cidr Count\n  - example: 20",
						MarkdownDescription: "Cidr Count\n  - example: 20",
						Computed:            true,
					},
					common.ToSnakeCase("cidrs"): schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"cidr": schema.StringAttribute{
									Computed:            true,
									Description:         "VPC Cidr\n  - example: 192.167.0.0/18",
									MarkdownDescription: "VPC Cidr\n  - example: 192.167.0.0/18",
								},
								"created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
									MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
								},
								"created_by": schema.StringAttribute{
									Computed:            true,
									Description:         "Created By\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
									MarkdownDescription: "Created By\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Cidr ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
									MarkdownDescription: "Cidr ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
								},
							},
						},
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
						Computed:            true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed:            true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description:         "Description\n  - maxLength: 50\n  - example: vpcDescription",
						MarkdownDescription: "Description\n  - maxLength: 50\n  - example: vpcDescription",
						Computed:            true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description:         "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
						MarkdownDescription: "VPC Id\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
						Computed:            true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed:            true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description:         "VPC Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: vpcName",
						MarkdownDescription: "VPC Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: `^[a-zA-Z0-9-]*$`\n  - example: vpcName",
						Computed:            true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description:         "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
						MarkdownDescription: "- enum: [\"CREATING\",\"ACTIVE\",\"DELETED\",\"ERROR\"]",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpcVpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vpcVpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpc.VpcResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpc
	data, err := r.client.CreateVpc(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpc",
			"Could not create vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(data.Vpc.Id)

	vpcModel := vpc.ResponseToVpcDSValue(data.Vpc)
	vpcObjectValueue, _ := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(ctx), vpcModel)
	plan.Vpc = vpcObjectValueue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *vpcVpcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpcVpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from vpc
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcModel := vpc.ResponseToVpcDSValue(data.Vpc)
	vpcObjectValueue, _ := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(ctx), vpcModel)
	state.Vpc = vpcObjectValueue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpcVpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state vpc.VpcResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateVpc(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpc",
			"Could not update vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetVpc as UpdateVpc items are not populated.
	data, err := r.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpc",
			"Could not read vpc ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpcModel := vpc.ResponseToVpcDSValue(data.Vpc)
	vpcObjectValueue, _ := types.ObjectValueFrom(ctx, vpcModel.AttributeTypes(ctx), vpcModel)
	state.Vpc = vpcObjectValueue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpcVpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpc.VpcResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpc
	err := r.client.DeleteVpc(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpc",
			"Could not delete vpc, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
