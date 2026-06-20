package gslb

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	gslb "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	virtualserverutil "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/virtualserver"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpgslb "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/gslb/1.1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const reasonPrefix = "\nReason: "

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &gslbGslbResource{}
	_ resource.ResourceWithConfigure   = &gslbGslbResource{}
	_ resource.ResourceWithImportState = &gslbGslbResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewGslbGslbResource() resource.Resource {
	return &gslbGslbResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type gslbGslbResource struct {
	config  *scpsdk.Configuration
	client  *gslb.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *gslbGslbResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gslb_gslb" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *gslbGslbResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Gslb.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("Tags"): tag.ResourceSchema(),
			common.ToSnakeCase("Gslb"): schema.SingleNestedAttribute{
				Description: "A detail of Gslb.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "Algorithm",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "Created at\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "Created by\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
						Description: "EnvUsage",
						Computed:    true,
					},
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "HealthCheck",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
								Description: "Created at\n" +
									"  - Example: 2024-05-17T00:23:17Z",
								Computed: true,
							},
							common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
								Description: "Created by\n" +
									"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "The GSLB Health Check Interval.\n" +
									"  - Range: 5 to 299",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Probe Timeout.\n" +
									"  - Range: 5 to 300",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "The GSLB Health Check User Name.\n" +
									"  - Max length: 60",
								Computed: true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "The GSLB Health Check Password.",
								Computed:    true,
							},
							common.ToSnakeCase("Id"): schema.StringAttribute{
								Description: "ID\n" +
									"  - Example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
								Computed: true,
							},
							common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
								Description: "Modified at\n" +
									"  - Example: 2024-05-17T00:23:17Z",
								Computed: true,
							},
							common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
								Description: "Modified by\n" +
									"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
								Computed: true,
							},
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The GSLB Health Check Protocol.",
								Computed:    true,
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "The GSLB Health Check Receive String.\n" +
									"  - Max length: 300",
								Computed: true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "The GSLB Health Check Send String.",
								Computed:    true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "The GSLB Health Check Service Port.\n" +
									"  - Range: 1 to 65535",
								Computed: true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Timeout.\n" +
									"  - Range: 6 to 300",
								Computed: true,
							},
						},
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "id",
						Computed:    true,
					},
					common.ToSnakeCase("LinkedResourceCount"): schema.Int32Attribute{
						Description: "LinkedResourceCount",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "Modified at\n" +
							"  - Example: 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "Modified by\n" +
							"  - Example: 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
				},
			},
			common.ToSnakeCase("GslbCreate"): schema.SingleNestedAttribute{
				Description: "Create Gslb.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Algorithm"): schema.StringAttribute{
						Description: "Algorithm",
						Required:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("EnvUsage"): schema.StringAttribute{
						Description: "EnvUsage",
						Required:    true,
					},
					common.ToSnakeCase("HealthCheck"): schema.SingleNestedAttribute{
						Description: "HealthCheck",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("HealthCheckInterval"): schema.Int32Attribute{
								Description: "The GSLB Health Check Interval.\n" +
									"  - Range: 5 to 299",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckProbeTimeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Probe Timeout.\n" +
									"  - Range: 5 to 300",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserId"): schema.StringAttribute{
								Description: "The GSLB Health Check User Name.\n" +
									"  - Max Length: 60",
								Optional: true,
							},
							common.ToSnakeCase("HealthCheckUserPassword"): schema.StringAttribute{
								Description: "The GSLB Health Check Password.",
								Optional:    true,
							},
							common.ToSnakeCase("Protocol"): schema.StringAttribute{
								Description: "The GSLB Health Check Protocol.",
								Required:    true,
							},
							common.ToSnakeCase("ReceiveString"): schema.StringAttribute{
								Description: "The GSLB Health Check Receive String.\n" +
									"  - Max Length: 300",
								Optional: true,
							},
							common.ToSnakeCase("SendString"): schema.StringAttribute{
								Description: "The GSLB Health Check Send String.",
								Optional:    true,
							},
							common.ToSnakeCase("ServicePort"): schema.Int32Attribute{
								Description: "The GSLB Health Check Service Port.\n" +
									"  - Range: 1 to 65535",
								Optional: true,
							},
							common.ToSnakeCase("Timeout"): schema.Int32Attribute{
								Description: "The GSLB Health Check Timeout.\n" +
									"  - Range: 6 to 300",
								Optional: true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Required:    true,
					},
					common.ToSnakeCase("Resources"): schema.ListNestedAttribute{
						Description: "Resources",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Description"): schema.StringAttribute{
									Description: "Description",
									Optional:    true,
								},
								common.ToSnakeCase("Destination"): schema.StringAttribute{
									Description: "Destination",
									Optional:    true,
								},
								common.ToSnakeCase("Region"): schema.StringAttribute{
									Description: "Region",
									Optional:    true,
								},
								common.ToSnakeCase("Weight"): schema.Int32Attribute{
									Description: "Weight",
									Optional:    true,
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
func (r *gslbGslbResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Gslb
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *gslbGslbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan gslb.GslbResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateGslb(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Gslb",
			"Could not create Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	err = waitForGslbStatus(ctx, r.client, data.Gslb.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating gslb",
			"Error waiting for gslb to become active: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(data.Gslb.Id)
	data, _ = r.client.GetGslb(ctx, data.Gslb.Id)

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	plan.Gslb = gslbObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *gslbGslbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state gslb.GslbResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetGslb(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Could not read Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	state.Gslb = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *gslbGslbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan
	var oldState gslb.GslbResource
	var state gslb.GslbResource
	req.State.Get(ctx, &oldState)
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if gslbResourceChanged(oldState, state) && gslbHealthCheckChanged(oldState, state) {
		resp.Diagnostics.AddError(
			"Error updating Gslb",
			"Could not change GSLB resources and health checks at the same time",
		)
		return
	}

	// Update existing order
	if gslbChanged(oldState, state) {
		_, err := r.client.UpdateGslb(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	if gslbResourceChanged(oldState, state) {
		_, err := r.client.UpdateGslbResource(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	if gslbHealthCheckChanged(oldState, state) {
		_, err := r.client.UpdateGslbHealthCheck(ctx, state.Id.ValueString(), state)
		if err != nil {
			detail := client.GetDetailFromError(err)
			resp.Diagnostics.AddError(
				"Error updating Gslb",
				"Could not update Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
			)
			return
		}
	}

	updateErr := waitForGslbStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating Gslb",
			"Error updating for Gslb to become active: "+updateErr.Error(),
		)
		return
	}

	data, err := r.client.GetGslb(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Gslb",
			"Could not read Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	state.Gslb = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *gslbGslbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state gslb.GslbResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Gslb
	data, err := r.client.DeleteGslb(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Gslb",
			"Could not delete Gslb, unexpected error: "+err.Error()+reasonPrefix+detail,
		)
		return
	}

	gslbModel := convertResponseToGslb(data)

	gslbObjectValue, diags := types.ObjectValueFrom(ctx, gslbModel.AttributeTypes(), gslbModel)
	state.Gslb = gslbObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertResponseToGslb(data *scpgslb.GslbShowResponse) gslb.GslbDetail {
	var healthCheck *gslb.HealthCheck
	if data.Gslb.HealthCheck.IsSet() {
		var healthCheckFromData = data.Gslb.HealthCheck.Get()
		healthCheck = &gslb.HealthCheck{
			CreatedAt:               types.StringValue(healthCheckFromData.CreatedAt.Format(time.RFC3339)),
			CreatedBy:               types.StringValue(healthCheckFromData.CreatedBy),
			HealthCheckInterval:     types.Int32Value(healthCheckFromData.GetHealthCheckInterval()),
			HealthCheckProbeTimeout: types.Int32Value(healthCheckFromData.GetHealthCheckProbeTimeout()),
			HealthCheckUserId:       types.StringValue(healthCheckFromData.GetHealthCheckUserId()),
			HealthCheckUserPassword: types.StringValue(healthCheckFromData.GetHealthCheckUserPassword()),
			Id:                      types.StringValue(healthCheckFromData.Id),
			ModifiedAt:              types.StringValue(healthCheckFromData.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:              types.StringValue(healthCheckFromData.ModifiedBy),
			Protocol:                types.StringValue(healthCheckFromData.Protocol),
			ReceiveString:           types.StringValue(healthCheckFromData.GetReceiveString()),
			SendString:              types.StringValue(healthCheckFromData.GetSendString()),
			ServicePort:             types.Int32Value(healthCheckFromData.GetServicePort()),
			Timeout:                 types.Int32Value(healthCheckFromData.GetTimeout()),
		}
	}
	return gslb.GslbDetail{
		Algorithm:           types.StringValue(data.Gslb.Algorithm),
		CreatedAt:           types.StringValue(data.Gslb.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(data.Gslb.CreatedBy),
		Description:         virtualserverutil.ToNullableStringValue(data.Gslb.Description.Get()),
		EnvUsage:            types.StringValue(data.Gslb.EnvUsage),
		HealthCheck:         healthCheck,
		Id:                  types.StringValue(data.Gslb.Id),
		LinkedResourceCount: types.Int32Value(data.Gslb.LinkedResourceCount),
		ModifiedAt:          types.StringValue(data.Gslb.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(data.Gslb.ModifiedBy),
		Name:                types.StringValue(data.Gslb.Name),
		State:               types.StringValue(data.Gslb.State),
	}
}

func gslbChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	if !oldState.GslbCreate.Algorithm.Equal(newState.GslbCreate.Algorithm) {
		return true
	}
	if !oldState.GslbCreate.Description.Equal(newState.GslbCreate.Description) {
		return true
	}
	return false
}

func gslbResourceChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	oldResources := oldState.GslbCreate.Resources
	newResources := newState.GslbCreate.Resources

	if len(oldResources) != len(newResources) {
		return true
	}

	for i := range oldResources {
		oldResource := oldResources[i]
		newResource := newResources[i]

		if oldResource.Description != newResource.Description {
			return true
		}

		if oldResource.Destination != newResource.Destination {
			return true
		}

		if oldResource.Region != newResource.Region {
			return true
		}

		if oldResource.Weight != newResource.Weight {
			return true
		}
	}

	return false
}

func gslbHealthCheckChanged(oldState gslb.GslbResource, newState gslb.GslbResource) bool {
	oldHealthCheck := oldState.GslbCreate.HealthCheck
	newHealthCheck := newState.GslbCreate.HealthCheck

	if oldHealthCheck == nil && newHealthCheck == nil {
		return false
	}
	if oldHealthCheck == nil || newHealthCheck == nil {
		return true
	}

	if oldHealthCheck.HealthCheckInterval != newHealthCheck.HealthCheckInterval {
		return true
	}
	if oldHealthCheck.HealthCheckProbeTimeout != newHealthCheck.HealthCheckProbeTimeout {
		return true
	}
	if oldHealthCheck.HealthCheckUserId != newHealthCheck.HealthCheckUserId {
		return true
	}
	if oldHealthCheck.HealthCheckUserPassword != newHealthCheck.HealthCheckUserPassword {
		return true
	}
	if oldHealthCheck.Protocol != newHealthCheck.Protocol {
		return true
	}
	if oldHealthCheck.ReceiveString != newHealthCheck.ReceiveString {
		return true
	}
	if oldHealthCheck.SendString != newHealthCheck.SendString {
		return true
	}
	if oldHealthCheck.ServicePort != newHealthCheck.ServicePort {
		return true
	}
	if oldHealthCheck.Timeout != newHealthCheck.Timeout {
		return true
	}

	return false
}

func waitForGslbStatus(ctx context.Context, gslbClient *gslb.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := gslbClient.GetGslb(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, info.Gslb.State, nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *gslbGslbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
