package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/dns"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dnsRecordResource{}
	_ resource.ResourceWithConfigure = &dnsRecordResource{}
)

// NewResourceManagerResourceGroupResource is a helper function to simplify the provider implementation.
func NewDnsRecordResource() resource.Resource {
	return &dnsRecordResource{}
}

// resourceManagerResourceGroupResource is the data source implementation.
type dnsRecordResource struct {
	config  *scpsdk.Configuration
	client  *dns.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *dnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *dnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Record.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ToSnakeCase("HostedZoneId"): schema.StringAttribute{
				Description: "Hosted zone ID.",
				Optional:    true,
			},
			common.ToSnakeCase("Record"): schema.SingleNestedAttribute{
				Description: "A detail of Record.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "Action",
						Optional:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Optional:    true,
					},
					common.ToSnakeCase("Links"): schema.SingleNestedAttribute{
						Description: "Links",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("Self"): schema.StringAttribute{
								Description: "Self",
								Optional:    true,
							},
						},
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("ProjectId"): schema.StringAttribute{
						Description: "ProjectId",
						Optional:    true,
					},
					common.ToSnakeCase("Records"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "Records",
						Optional:    true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "Status",
						Optional:    true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "Ttl",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
					common.ToSnakeCase("UpdatedAt"): schema.StringAttribute{
						Description: "UpdatedAt",
						Optional:    true,
					},
					common.ToSnakeCase("Version"): schema.Int32Attribute{
						Description: "Version",
						Optional:    true,
					},
					common.ToSnakeCase("ZoneId"): schema.StringAttribute{
						Description: "ZoneId",
						Optional:    true,
					},
					common.ToSnakeCase("ZoneName"): schema.StringAttribute{
						Description: "ZoneName",
						Optional:    true,
					},
				},
			},
			common.ToSnakeCase("RecordCreate"): schema.SingleNestedAttribute{
				Description: "Create Record.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Optional:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Optional:    true,
					},
					common.ToSnakeCase("Records"): schema.ListAttribute{
						ElementType: types.StringType,
						Description: "Records",
						Optional:    true,
					},
					common.ToSnakeCase("Ttl"): schema.Int32Attribute{
						Description: "Ttl",
						Optional:    true,
					},
					common.ToSnakeCase("Type"): schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { // 아직 정의하지 않은 Create 메서드를 추가한다.
	var plan dns.RecordResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateRecord(ctx, plan.HostedZoneId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Record",
			"Could not create Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	createErr := waitForRecordStatus(ctx, r.client, plan.HostedZoneId.ValueString(), *data.Id.Get(), []string{}, []string{"ACTIVE"})
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error creating record",
			"Error creating for record to become active: "+createErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetRecord(ctx, plan.HostedZoneId.ValueString(), *data.Id.Get())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(*data.Id.Get())

	recordModel := convertRecordDetail(*dataForShow)

	recordOjbectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	plan.Record = recordOjbectValue

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	// Get current state
	var state dns.RecordResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Gslb
	data, err := r.client.GetRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	recordModel := convertRecordDetail(*data)

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	state.Record = recordObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // 아직 정의하지 않은 Update 메서드를 추가한다.
	// Retrieve values from plan

	var state dns.RecordResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Record",
			"Could not update Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	updateErr := waitForRecordStatus(ctx, r.client, state.HostedZoneId.ValueString(), state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating record",
			"Error updating for record to become active: "+updateErr.Error(),
		)
		return
	}

	dataForShow, err := r.client.GetRecord(ctx, state.HostedZoneId.ValueString(), *data.Id.Get())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Record",
			"Could not read Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	recordModel := convertRecordDetail(*dataForShow)

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	state.Record = recordObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { // 아직 정의하지 않은 Delete 메서드를 추가한다.
	// Retrieve values from state
	var state dns.RecordResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.DeleteRecord(ctx, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Record",
			"Could not delete Record, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// DeleteRecord returns 202 Accepted and removes the record asynchronously.
	// Block until the Show 404s so the parent hosted zone Delete that follows in
	// the dependency graph does not 409 on a still-present record (which cascades
	// to the private DNS delete and leaks the bootstrap VPC). 404 is terminal.
	err = waitForRecordGone(ctx, r.client, state.HostedZoneId.ValueString(), state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting record",
			"Error waiting for record to become deleted: "+err.Error(),
		)
		return
	}

	recordModel := convertRecordDetail(convertRecordCreateResponseToRecord(*data))

	recordObjectValue, diags := types.ObjectValueFrom(ctx, recordModel.AttributeTypes(), recordModel)
	state.Record = recordObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func waitForRecordStatus(ctx context.Context, recordClient *dns.Client, hostedZoneId string, recordId string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := recordClient.GetRecord(ctx, hostedZoneId, recordId)
		if err != nil {
			return nil, "", err
		}
		return info, *info.Status.Get(), nil
	})
}

// waitForRecordGone polls the record Show until it 404s (the record has been
// fully removed). The refresh func surfaces that 404 as the returned error,
// which the Delete handler tolerates; DELETED is never actually observed because
// the Show stops returning the record once it is gone. It is nil-safe (does not
// dereference the response on error) unlike waitForRecordStatus.
func waitForRecordGone(ctx context.Context, recordClient *dns.Client, hostedZoneId string, recordId string) error {
	return client.WaitForStatus(ctx, nil, []string{"ACTIVE", "DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		info, err := recordClient.GetRecord(ctx, hostedZoneId, recordId)
		if err != nil {
			return nil, "", err
		}
		if info == nil || info.Status.Get() == nil {
			return nil, "", nil
		}
		return info, *info.Status.Get(), nil
	})
}
