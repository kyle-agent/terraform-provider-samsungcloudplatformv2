package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	vpcV1Dot2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &VpcCidrResource{}
	_ resource.ResourceWithConfigure = &VpcCidrResource{}
)

// NewVpcCidrResource is a helper function to simplify the provider implementation.
func NewVpcCidrResource() resource.Resource {
	return &VpcCidrResource{}
}

// VpcCidrResource is the resource implementation.
type VpcCidrResource struct {
	_config *scpsdk.Configuration
	client  *vpcV1Dot2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VpcCidrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_cidr"
}

// Schema defines the schema for the resource.
func (r *VpcCidrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VPC CIDR",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "CIDR \n" +
					"  - example : 192.168.0.0/24",
				Required: true,
			},

			// Output
			common.ToSnakeCase("Vpc"): schema.SingleNestedAttribute{
				Description: "VPC detail after adding CIDR",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "VPC ID",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "Account ID",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("CidrCount"): schema.Int32Attribute{
						Description: "CIDR Count",
						Computed:    true,
					},
					common.ToSnakeCase("Cidrs"): schema.ListNestedAttribute{
						Description: "CIDRs",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Id"): schema.StringAttribute{
									Description: "CIDR ID",
									Computed:    true,
								},
								common.ToSnakeCase("Cidr"): schema.StringAttribute{
									Description: "CIDR",
									Computed:    true,
								},
								common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
									Description: "Created At",
									Computed:    true,
								},
								common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
									Description: "Created By",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created At",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created By",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified At",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified By",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VpcCidrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VpcCidrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcV1Dot2.VpcCidrResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.AddVpcCidr(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to add VPC CIDR",
			fmt.Sprintf("An error occurred while creating VPC CIDR: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	vpcCidr := &vpcV1Dot2.VpcCidrDetail{
		Id:         types.StringValue(data.Vpc.Id),
		Name:       types.StringValue(data.Vpc.Name),
		AccountId:  types.StringValue(data.Vpc.AccountId),
		State:      types.StringValue(string(data.Vpc.State)),
		CidrCount:  types.Int32Value(data.Vpc.CidrCount),
		CreatedAt:  types.StringValue(data.Vpc.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringValue(data.Vpc.CreatedBy),
		ModifiedAt: types.StringValue(data.Vpc.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringValue(data.Vpc.ModifiedBy),
	}

	if data.Vpc.Description.IsSet() {
		if desc := data.Vpc.Description.Get(); desc != nil {
			vpcCidr.Description = types.StringValue(*desc)
		}
	}

	if data.Vpc.Cidrs != nil {
		for _, cidr := range data.Vpc.Cidrs {
			vpcCidr.Cidrs = append(vpcCidr.Cidrs, vpcV1Dot2.VpcCidrInfo{
				Id:        types.StringValue(cidr.Id),
				Cidr:      types.StringValue(cidr.Cidr),
				CreatedAt: types.StringValue(cidr.CreatedAt.Format(time.RFC3339)),
				CreatedBy: types.StringValue(cidr.CreatedBy),
			})
		}
	} else {
		vpcCidr.Cidrs = []vpcV1Dot2.VpcCidrInfo{}
	}

	vpcCidrObjectValue, _ := types.ObjectValueFrom(ctx, vpcCidr.AttributeTypes(), vpcCidr)
	plan.Vpc = vpcCidrObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VpcCidrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpcV1Dot2.VpcCidrResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the VPC. A real 404 means the VPC (and thus its CIDRs) is gone,
	// so the resource is removed from state. Other errors are surfaced.
	data, statusCode, err := r.client.GetVpcWithStatus(ctx, state.VpcId.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading VPC CIDR",
			fmt.Sprintf("Could not read VPC CIDR for VPC ID %s: %s. Details: %s", state.VpcId.ValueString(), err.Error(), detail),
		)
		return
	}

	// If the specific CIDR this resource manages is no longer present on the VPC,
	// it has been removed out-of-band; drop it from state so it is re-created.
	cidrFound := false
	for _, cidr := range data.Vpc.Cidrs {
		if cidr.Cidr == state.Cidr.ValueString() {
			cidrFound = true
			break
		}
	}
	if !cidrFound {
		resp.State.RemoveResource(ctx)
		return
	}

	// Map API response to object (mirrors Create).
	vpcCidr := &vpcV1Dot2.VpcCidrDetail{
		Id:         types.StringValue(data.Vpc.Id),
		Name:       types.StringValue(data.Vpc.Name),
		AccountId:  types.StringValue(data.Vpc.AccountId),
		State:      types.StringValue(string(data.Vpc.State)),
		CidrCount:  types.Int32Value(data.Vpc.CidrCount),
		CreatedAt:  types.StringValue(data.Vpc.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  types.StringValue(data.Vpc.CreatedBy),
		ModifiedAt: types.StringValue(data.Vpc.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy: types.StringValue(data.Vpc.ModifiedBy),
	}

	if data.Vpc.Description.IsSet() {
		if desc := data.Vpc.Description.Get(); desc != nil {
			vpcCidr.Description = types.StringValue(*desc)
		}
	}

	if data.Vpc.Cidrs != nil {
		for _, cidr := range data.Vpc.Cidrs {
			vpcCidr.Cidrs = append(vpcCidr.Cidrs, vpcV1Dot2.VpcCidrInfo{
				Id:        types.StringValue(cidr.Id),
				Cidr:      types.StringValue(cidr.Cidr),
				CreatedAt: types.StringValue(cidr.CreatedAt.Format(time.RFC3339)),
				CreatedBy: types.StringValue(cidr.CreatedBy),
			})
		}
	} else {
		vpcCidr.Cidrs = []vpcV1Dot2.VpcCidrInfo{}
	}

	vpcCidrObjectValue, objDiags := types.ObjectValueFrom(ctx, vpcCidr.AttributeTypes(), vpcCidr)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Vpc = vpcCidrObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VpcCidrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: Implement Update function
	resp.Diagnostics.AddError(
		"Update Not Implemented",
		"VPC CIDR Update function is not yet implemented.",
	)
}

// Delete removes the managed secondary CIDR from the VPC. This is the inverse of
// Create (which ADDED the CIDR). There is no bulk "set CIDRs" endpoint, so we
// locate the CIDR entry's id on the VPC and issue a DELETE for it, then wait
// until the CIDR is no longer present. A 404 (VPC or CIDR already gone) is
// treated as already deleted.
func (r *VpcCidrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcV1Dot2.VpcCidrResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpcId := state.VpcId.ValueString()
	targetCidr := state.Cidr.ValueString()

	// Fetch the VPC to resolve the CIDR entry's id (needed for the DELETE path).
	data, statusCode, err := r.client.GetVpcWithStatus(ctx, vpcId)
	if err != nil {
		if statusCode == 404 {
			// VPC is gone; the CIDR is gone with it.
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting VPC CIDR",
			fmt.Sprintf("Could not read VPC %s before removing CIDR %s: %s. Details: %s", vpcId, targetCidr, err.Error(), detail),
		)
		return
	}

	cidrId := ""
	for _, cidr := range data.Vpc.Cidrs {
		if cidr.Cidr == targetCidr {
			cidrId = cidr.Id
			break
		}
	}
	if cidrId == "" {
		// CIDR already removed out-of-band; nothing to do.
		return
	}

	removeStatus, err := r.client.RemoveVpcCidr(ctx, vpcId, cidrId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to remove VPC CIDR",
			fmt.Sprintf("An error occurred while removing CIDR %s from VPC %s: %s. Details: %s", targetCidr, vpcId, err.Error(), detail),
		)
		return
	}
	if removeStatus == 404 {
		// Already gone.
		return
	}

	// Wait until the CIDR is no longer present on the VPC.
	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		vpcData, sc, err := r.client.GetVpcWithStatus(ctx, vpcId)
		if err != nil {
			if sc == 404 {
				return
			}
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error Waiting for VPC CIDR Removal",
				fmt.Sprintf("Could not read VPC %s while waiting for CIDR %s removal: %s. Details: %s", vpcId, targetCidr, err.Error(), detail),
			)
			return
		}

		stillPresent := false
		for _, cidr := range vpcData.Vpc.Cidrs {
			if cidr.Cidr == targetCidr {
				stillPresent = true
				break
			}
		}
		if !stillPresent {
			return
		}

		select {
		case <-timeout:
			resp.Diagnostics.AddError(
				"Timeout Waiting for VPC CIDR Removal",
				fmt.Sprintf("CIDR %s was still present on VPC %s after waiting for removal.", targetCidr, vpcId),
			)
			return
		case <-ticker.C:
			// poll again
		}
	}
}
