package virtualserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/virtualserver"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/virtualserver/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &virtualServerServerGroupResource{}
	_ resource.ResourceWithConfigure   = &virtualServerServerGroupResource{}
	_ resource.ResourceWithImportState = &virtualServerServerGroupResource{}
)

func NewVirtualServerServerGroupResource() resource.Resource {
	return &virtualServerServerGroupResource{}
}

type virtualServerServerGroupResource struct {
	config  *scpsdk.Configuration
	client  *virtualserver.Client
	clients *client.SCPClient
}

func (r *virtualServerServerGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualserver_server_group"
}

func (r *virtualServerServerGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Server",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			common.ToSnakeCase("Policy"): schema.StringAttribute{
				Description: "Policy",
				Required:    true,
			},
			common.ToSnakeCase("AccountId"): schema.StringAttribute{
				Description: "Account ID",
				Computed:    true,
			},
			common.ToSnakeCase("UserId"): schema.StringAttribute{
				Description: "User ID",
				Computed:    true,
			},
			common.ToSnakeCase("Members"): schema.ListAttribute{
				Description: "Members",
				Computed:    true,
				ElementType: types.StringType,
			},
			common.ToSnakeCase("PartitionSize"): schema.Int32Attribute{
				Description: "Partition Size",
				Computed:    true,
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

func (r *virtualServerServerGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.VirtualServer
	r.clients = inst.Client
}

func (r *virtualServerServerGroupResource) MapGetResponseToState(resp *scpvirtualserver.ServerGroup, tagsMap types.Map) virtualserver.ServerGroupResource {
	//Members
	members := make([]attr.Value, len(resp.Members))
	for i, member := range resp.Members {
		members[i] = types.StringValue(member)
	}

	return virtualserver.ServerGroupResource{
		Id:            types.StringValue(resp.Id),
		Name:          types.StringValue(resp.Name),
		Policy:        types.StringValue(resp.Policy),
		AccountId:     types.StringValue(resp.AccountId),
		UserId:        types.StringValue(resp.UserId),
		Members:       types.ListValueMust(types.StringType, members),
		PartitionSize: types.Int32PointerValue(resp.PartitionSize.Get()),
		Tags:          tagsMap,
	}
}

func (r *virtualServerServerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualserver.ServerGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateServerGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating server group",
			"Could not create server group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	getData, err := r.client.GetServerGroup(ctx, data.Id)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading server group",
			"Could not create server group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeServerGroup, data.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state := r.MapGetResponseToState(getData, tagsMap)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualserver.ServerGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetServerGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading server group",
			"Could not read server group id "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tagsMap, err := tag.GetTags(r.clients, ServiceNameVirtualServer, ResourceTypeServerGroup, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, state.Tags)

	newState := r.MapGetResponseToState(data, tagsMap)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan virtualserver.ServerGroupResource
	var state virtualserver.ServerGroupResource
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settableFileds []string
	for attrName, attribute := range req.Plan.Schema.GetAttributes() {
		if attribute.IsRequired() || attribute.IsOptional() {
			settableFileds = append(settableFileds, virtualserverutil.SnakeToPascal(attrName))
		}
	}

	changeFields, err := virtualserverutil.GetChangedFields(plan, state, settableFileds)
	if err != nil {
		return
	}

	immuntableFields := []string{"Name", "Policy"}

	if virtualserverutil.IsOverlapFields(immuntableFields, changeFields) {
		resp.Diagnostics.AddError(
			"Error Updating Server Group",
			"Immutable fields cannot be modified: "+strings.Join(immuntableFields, ", "),
		)
		return
	}

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, ServiceNameVirtualServer, ResourceTypeServerGroup, plan.Id.ValueString(), tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tags",
			err.Error(),
		)
		return
	}
	tagsMap = common.NullTagCheck(tagsMap, plan.Tags)

	state.Tags = tagsMap

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualServerServerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualserver.ServerGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServerGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting server group",
			"Could not delete server group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *virtualServerServerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
