package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/firewall/1.0"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallFirewallRuleResource{}
	_ resource.ResourceWithConfigure   = &firewallFirewallRuleResource{}
	_ resource.ResourceWithImportState = &firewallFirewallRuleResource{}
)

// NewFirewallFirewallRuleResource is a helper function to simplify the provider implementation.
func NewFirewallFirewallRuleResource() resource.Resource {
	return &firewallFirewallRuleResource{}
}

// networkFirewallRuleResource is the data source implementation.
type firewallFirewallRuleResource struct {
	config *scpsdk.Configuration
	client *firewall.Client
}

// Metadata returns the firewallFirewallRuleResource source type name.
func (r *firewallFirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall_rule"
}

// Schema defines the schema for the data source.
func (r *firewallFirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Firewall rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "Firewall ID \n" +
					"  - example : 68db67f78abd405da98a6056a8ee42af",
				Required: true,
			},
			common.ToSnakeCase("FirewallRule"): schema.SingleNestedAttribute{
				Description: "Firewall rule",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "FirewallId",
						Computed:    true,
					},
					common.ToSnakeCase("Sequence"): schema.Int32Attribute{
						Description: "Sequence",
						Computed:    true,
					},
					common.ToSnakeCase("SourceInterface"): schema.StringAttribute{
						Description: "SourceInterface",
						Computed:    true,
					},
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "SourceAddress",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationInterface"): schema.StringAttribute{
						Description: "DestinationInterface",
						Computed:    true,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "DestinationAddress",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "Service",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "ServiceType",
									Computed:    true,
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "ServiceValue",
									Computed:    true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action",
						Computed:    true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "Direction",
						Computed:    true,
					},
					common.ToSnakeCase("VendorRuleId"): schema.StringAttribute{
						Description: "VendorRuleId",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
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
				},
			},
			common.ToSnakeCase("FirewallRuleCreate"): schema.SingleNestedAttribute{
				Description: "Firewall rule create object",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "Source Address \n" +
							"  - example : ['10.10.10.10', '20.20.20.20']",
						Required:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "Destination Address \n" +
							"  - example : ['10.10.10.10', '20.20.20.20']",
						Required:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "Service",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "Service Type \n" +
										"  - example : TCP | UDP | ICMP | IP | TCP_ALL | UDP_ALL | ICMP_ALL | ALL",
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("TCP", "UDP", "ICMP", "IP", "TCP_ALL", "UDP_ALL", "ICMP_ALL", "ALL"),
									},
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "Service Value \n" +
										"  - example : 443",
									Optional: true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action \n" +
							"  - example : ALLOW | DENY",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ALLOW", "DENY"),
						},
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "Direction \n" +
							"  - example : INBOUND | OUTBOUND",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("INBOUND", "OUTBOUND"),
						},
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status \n" +
							"  - example : ENABLE | DISABLE",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ENABLE", "DISABLE"),
						},
					},
					common.ToSnakeCase("OrderRuleId"): schema.StringAttribute{
						Description: "OrderRule ID \n" +
							"  - example : 7087c92d295445cda2785a94aab93c65",
						Optional: true,
					},
					common.ToSnakeCase("OrderDirection"): schema.StringAttribute{
						Description: "Order Direction \n" +
							"  - example :  BEFORE | AFTER | BOTTOM",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("BEFORE", "AFTER", "BOTTOM"),
						},
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description\n" +
							"  - example : VPC description\n" +
							"  - maxLength : 100\n" +
							"  - minLength : 1",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(100),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *firewallFirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Firewall
}

// Create creates the resource and sets the initial Terraform state.
func (r *firewallFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan firewall.FirewallRuleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new firewall rule
	data, err := r.client.CreateFirewallRule(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating firewall rule",
			"Could not create firewall rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.FirewallRule.Id)

	// Map response body to schema and populate Computed attribute values
	firewallRuleModel := createFirewallRuleModel(data)
	firewallRuleObjectValue, diags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	plan.FirewallRule = firewallRuleObjectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState adopts an existing resource via `terraform import <addr> <id>`
// using its opaque id; Read then refreshes the remaining state. (#81)
func (r *firewallFirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *firewallFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state firewall.FirewallRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from firewall rule
	data, err := r.client.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading firewall rule",
			"Could not read firewall rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	firewallRuleModel := createFirewallRuleModel(data)

	firewallRuleObjectValue, diags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	state.FirewallRule = firewallRuleObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *firewallFirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state firewall.FirewallRuleResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	_, err := r.client.UpdateFirewallRule(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating direct connect",
			"Could not update direct connect, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Fetch updated items from GetFirewallRule as UpdateFirewallRule items are not populated.
	data, err := r.client.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading firewall rule",
			"Could not read firewall rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	firewallRuleModel := createFirewallRuleModel(data)

	firewallRuleObjectValue, diags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	state.FirewallRule = firewallRuleObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state firewall.FirewallRuleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing firewall rule
	err := r.client.DeleteFirewallRule(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting firewall rule",
			"Could not delete firewall rule, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func createFirewallRuleModel(data *scpfirewall.FirewallRuleShowResponse) firewall.FirewallRule {
	fwRule := data.FirewallRule
	sourceAddresses := make([]string, 0, len(fwRule.SourceAddress))
	for _, address := range fwRule.SourceAddress {
		sourceAddresses = append(sourceAddresses, address)
	}
	destinationAddresses := make([]string, 0, len(fwRule.DestinationAddress))
	for _, address := range fwRule.DestinationAddress {
		destinationAddresses = append(destinationAddresses, address)
	}
	services := make([]firewall.FirewallPort, 0, len(fwRule.Service))
	for _, service := range fwRule.Service {
		services = append(services, firewall.FirewallPort{
			ServiceType:  types.StringValue(string(service.ServiceType)),
			ServiceValue: types.StringValue(*service.ServiceValue),
		})
	}

	return firewall.FirewallRule{
		Id:                   types.StringValue(fwRule.Id),
		Name:                 types.StringPointerValue(fwRule.Name.Get()),
		FirewallId:           types.StringValue(fwRule.FirewallId),
		Sequence:             types.Int32Value(fwRule.Sequence),
		SourceInterface:      types.StringValue(fwRule.SourceInterface),
		SourceAddress:        sourceAddresses,
		DestinationInterface: types.StringValue(fwRule.DestinationInterface),
		DestinationAddress:   destinationAddresses,
		Service:              services,
		Action:               types.StringValue(string(fwRule.Action)),
		Direction:            types.StringValue(string(fwRule.Direction)),
		VendorRuleId:         types.StringValue(fwRule.VendorRuleId),
		Description:          types.StringPointerValue(fwRule.Description.Get()),
		State:                types.StringValue(string(fwRule.State)),
		Status:               types.StringValue(string(fwRule.Status)),
		CreatedAt:            types.StringValue(fwRule.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(fwRule.CreatedBy),
		ModifiedAt:           types.StringValue(fwRule.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(fwRule.ModifiedBy),
	}
}
