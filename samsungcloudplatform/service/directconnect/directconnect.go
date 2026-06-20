package directconnect

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpdirectconnect "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/direct-connect/1.0"
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
	_ resource.Resource                = &directConnectDirectConnectResource{}
	_ resource.ResourceWithConfigure   = &directConnectDirectConnectResource{}
	_ resource.ResourceWithImportState = &directConnectDirectConnectResource{}
)

// NewDirectConnectDirectConnectResource is a helper function to simplify the provider implementation.
func NewDirectConnectDirectConnectResource() resource.Resource {
	return &directConnectDirectConnectResource{}
}

// directConnectDirectConnectResource is the data source implementation.
type directConnectDirectConnectResource struct {
	config  *scpsdk.Configuration
	client  *directconnect.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *directConnectDirectConnectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directconnect_direct_connect"
}

// Schema defines the schema for the data source.
func (r *directConnectDirectConnectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "direct connect",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
				Description: "Type \n" +
					"  - example : 1 | 10 | 20 | 40",
				Required: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description\n" +
					"  - example : Direct Connect description\n" +
					"  - maxLength : 50\n" +
					"  - minLength : 1",
				Optional: true,
			},
			common.ToSnakeCase("FirewallEnabled"): schema.BoolAttribute{
				Description: "Firewall Enabled \n" +
					"  - example : true | false",
				Optional: true,
			},
			common.ToSnakeCase("FirewallLoggable"): schema.BoolAttribute{
				Description: "Firewall Loggable \n" +
					"  - example : true | false",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Direct Connect Name \n" +
					"  - example : directConnectName",
				Required: true,
			},
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "VPC ID \n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
			},
			common.ToSnakeCase("DirectConnect"): schema.SingleNestedAttribute{
				Description: "DirectConnect",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "name",
						Computed:    true,
					},
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "account id",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "description",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "vpc id",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "vpc name",
						Computed:    true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "bandwidth",
						Computed:    true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "firewall id",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "created at",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "created by",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "modified at",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "modified by",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "state",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *directConnectDirectConnectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.DirectConnect
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *directConnectDirectConnectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan directconnect.DirectConnectResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new direct connect
	data, err := r.client.CreateDirectConnect(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating direct connect",
			"Could not create direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.DirectConnect.Id)
	diags = resp.State.Set(ctx, plan)

	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	plan.DirectConnect = dconObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	err = waitForDirectConnectStatus(ctx, r.client, data.DirectConnect.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating direct connect",
			"Error waiting for direct connect to become active: "+err.Error(),
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
func (r *directConnectDirectConnectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state directconnect.DirectConnectResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from direct connect
	data, err := r.client.GetDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading direct connect",
			"Could not read direct connect ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	state.DirectConnect = dconObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *directConnectDirectConnectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state directconnect.DirectConnectResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateDirectConnect(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating direct connect",
			"Could not update direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetDirectConnect as UpdateDirectConnect items are not populated.
	data, err := r.client.GetDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading direct connect",
			"Could not read direct connect ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	dconModel := createDirectConnectModel(data)

	dconObjectValue, diags := types.ObjectValueFrom(ctx, dconModel.AttributeTypes(), dconModel)
	state.DirectConnect = dconObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *directConnectDirectConnectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state directconnect.DirectConnectResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing direct connect
	err := r.client.DeleteDirectConnect(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting direct connect",
			"Could not delete direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForDirectConnectStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting direct connect",
			"Error waiting for direct connect to become deleted: "+err.Error(),
		)
		return
	}
}

func createDirectConnectModel(data *scpdirectconnect.DirectConnectShowResponse) directconnect.DirectConnect {
	dcon := data.DirectConnect
	return directconnect.DirectConnect{
		Id:          types.StringValue(dcon.Id),
		Name:        types.StringValue(dcon.Name),
		AccountId:   types.StringValue(dcon.AccountId),
		Description: types.StringPointerValue(dcon.Description.Get()),
		VpcId:       types.StringValue(dcon.VpcId),
		VpcName:     types.StringValue(dcon.VpcName),
		Bandwidth:   types.Int32Value(dcon.Bandwidth),
		FirewallId:  types.StringPointerValue(dcon.FirewallId.Get()),
		CreatedAt:   types.StringValue(dcon.CreatedAt.Format(time.RFC3339)),
		CreatedBy:   types.StringValue(dcon.CreatedBy),
		ModifiedAt:  types.StringValue(dcon.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:  types.StringValue(dcon.ModifiedBy),
		State:       types.StringValue(string(dcon.State)),
	}
}

func waitForDirectConnectStatus(ctx context.Context, directConnectClient *directconnect.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := directConnectClient.GetDirectConnect(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.DirectConnect.State), nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *directConnectDirectConnectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
