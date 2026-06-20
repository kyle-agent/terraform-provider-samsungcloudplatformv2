package resourcemanager

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/region"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceManagerResourceGroupResource{}
	_ resource.ResourceWithConfigure   = &resourceManagerResourceGroupResource{}
	_ resource.ResourceWithImportState = &resourceManagerResourceGroupResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewResourceManagerResourceGroupResource() resource.Resource {
	return &resourceManagerResourceGroupResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type resourceManagerResourceGroupResource struct {
	config  *scpsdk.Configuration
	client  *resourcemanager.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *resourceManagerResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resourcemanager_resource_group"
}

// Schema defines the schema for the data source.
func (r *resourceManagerResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource Group",
		Attributes: map[string]schema.Attribute{
			"region": region.ResourceSchema(),
			"tags":   tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the Resource Group",
				Computed:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description",
				Optional:    true,
			},
			common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "ResourceTypes",
				Optional:    true,
			},
			common.ToSnakeCase("GroupDefinitionTags"): schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Group Definition Tags",
			},
			common.ToSnakeCase("ResourceGroup"): schema.SingleNestedAttribute{
				Description: "Resource Group",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"tags": tag.ResourceSchema(),
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("Region"): schema.StringAttribute{
						Description: "Region",
						Computed:    true,
					},
					common.ToSnakeCase("Srn"): schema.StringAttribute{
						Description: "Srn",
						Computed:    true,
					},
					common.ToSnakeCase("ResourceTypes"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "ResourceTypes",
						Computed:    true,
					},
					common.ToSnakeCase("GroupDefinitionTags"): schema.ListNestedAttribute{
						Description: "A list of tag.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Key"): schema.StringAttribute{
									Description: "Key",
									Computed:    true,
								},
								common.ToSnakeCase("Value"): schema.StringAttribute{
									Description: "Value",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *resourceManagerResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.ResourceManager
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceManagerResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resourcemanager.ResourceGroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Region.IsNull() {
		r.client.Config.Region = plan.Region.ValueString()
	}

	// Create new Resource Group
	data, err := r.client.CreateResourceGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Resource Group",
			"Could not create Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup
	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringValue(resourceGroup.Id)

	tagElements := plan.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id, tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tags",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}
	resourceTypesObject, _ := types.ListValueFrom(ctx, types.StringType, resourceTypes)

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	plan.ResourceGroup = resourceGroupObjectValue
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceManagerResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state resourcemanager.ResourceGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Resource Group
	data, err := r.client.GetResourceGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			"Could not read Resource Group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup

	// Get Tags
	tagsMap, err := tag.GetTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Resource Group",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}
	resourceTypesObject, _ := types.ListValueFrom(ctx, types.StringType, resourceTypes)

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	state.ResourceGroup = resourceGroupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceManagerResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state resourcemanager.ResourceGroupResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Resource Group
	_, err := r.client.UpdateResourceGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Resource Group",
			"Could not update Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetResourceGroup as UpdateResourceGroup items are not populated.
	data, err := r.client.GetResourceGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading resourceGroup",
			"Could not read resourceGroup ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	resourceGroup := data.ResourceGroup

	tagElements := state.Tags.Elements()
	tagsMap, err := tag.UpdateTags(r.clients, "resourcemanager", "resource-group", resourceGroup.Id, tagElements)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tags",
			err.Error(),
		)
		return
	}

	var resourceTypes []string
	for _, resourceType := range resourceGroup.ResourceTypes {
		resourceTypes = append(resourceTypes, resourceType)
	}
	resourceTypesObject, _ := types.ListValueFrom(ctx, types.StringType, resourceTypes)

	var groupDefinitionTags []resourcemanager.Tag
	for _, t := range resourceGroup.Tags {
		tagState := resourcemanager.Tag{
			Key:   types.StringValue(t.Key),
			Value: types.StringPointerValue(t.Value.Get()),
		}
		groupDefinitionTags = append(groupDefinitionTags, tagState)
	}

	resourceGroupModel := resourcemanager.ResourceGroup{
		CreatedAt:           types.StringValue(resourceGroup.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(resourceGroup.CreatedBy),
		Description:         types.StringPointerValue(resourceGroup.Description.Get()),
		Id:                  types.StringValue(resourceGroup.Id),
		ModifiedAt:          types.StringValue(resourceGroup.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(resourceGroup.ModifiedBy),
		Name:                types.StringPointerValue(resourceGroup.Name.Get()),
		Region:              types.StringPointerValue(resourceGroup.Region.Get()),
		Srn:                 types.StringValue(resourceGroup.Srn),
		Tags:                tagsMap,
		ResourceTypes:       resourceTypesObject,
		GroupDefinitionTags: groupDefinitionTags,
	}
	resourceGroupObjectValue, diags := types.ObjectValueFrom(ctx, resourceGroupModel.AttributeTypes(), resourceGroupModel)
	state.ResourceGroup = resourceGroupObjectValue
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceManagerResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state resourcemanager.ResourceGroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag.UpdateTags(r.clients, "resourcemanager", "resource-group", state.Id.ValueString(), make(map[string]attr.Value))

	// Delete existing Resource Group
	err := r.client.DeleteResourceGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Resource Group",
			"Could not delete Resource Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *resourceManagerResourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
