package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamPolicyResource{}
	_ resource.ResourceWithConfigure   = &iamPolicyResource{}
	_ resource.ResourceWithImportState = &iamPolicyResource{}
)

// NewIamPolicyResource is a helper function to simplify the provider implementation.
func NewIamPolicyResource() resource.Resource {
	return &iamPolicyResource{}
}

// iamPolicyResource is the data source implementation.
type iamPolicyResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy"
}

func (r *iamPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Policy ID",
				MarkdownDescription: "Policy ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy Name",
				MarkdownDescription: "Policy Name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Policy Description",
				MarkdownDescription: "Policy Description",
			},
			"tags": tag.ResourceSchema(),
			"policy_version": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Policy Version",
				MarkdownDescription: "Policy Version",
				Attributes: map[string]schema.Attribute{
					"policy_document": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Policy Document",
						MarkdownDescription: "Policy Document",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								Optional:            true,
								Description:         "Statement",
								MarkdownDescription: "Statement",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.ListAttribute{
											Optional:            true,
											Description:         "Action",
											MarkdownDescription: "Action",
											ElementType:         types.StringType,
										},
										"not_action": schema.ListAttribute{
											Optional:            true,
											Description:         "Not Action",
											MarkdownDescription: "Not Action",
											ElementType:         types.StringType,
										},
										"effect": schema.StringAttribute{
											Optional:            true,
											Description:         "Effect",
											MarkdownDescription: "Effect",
										},
										"resource": schema.ListAttribute{
											Optional:            true,
											Description:         "Resource",
											MarkdownDescription: "Resource",
											ElementType:         types.StringType,
										},
										"sid": schema.StringAttribute{
											Optional:            true,
											Description:         "SID",
											MarkdownDescription: "SID",
										},
										"condition": schema.MapAttribute{
											ElementType: types.MapType{
												ElemType: types.ListType{
													ElemType: types.StringType,
												},
											},
											Optional: true,
										},
									},
								},
							},
							"version": schema.StringAttribute{
								Optional:            true,
								Description:         "Policy Version",
								MarkdownDescription: "Policy Version",
							},
						},
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Description: "A detail of Policy.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID",
						MarkdownDescription: "Account ID",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At",
						MarkdownDescription: "Created At",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By",
						MarkdownDescription: "Created By",
					},
					"creator_email": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Email",
						MarkdownDescription: "Creator Email",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Creator Name",
						MarkdownDescription: "Creator Name",
					},
					"default_version_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Default Version ID",
						MarkdownDescription: "Default Version ID",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Description",
						MarkdownDescription: "Description",
					},
					"domain_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Domain Name",
						MarkdownDescription: "Domain Name",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At",
						MarkdownDescription: "Modified At",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By",
						MarkdownDescription: "Modified By",
					},
					"modifier_email": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Email",
						MarkdownDescription: "Modifier Email",
					},
					"modifier_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier Name",
						MarkdownDescription: "Modifier Name",
					},
					"policy_category": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Category",
						MarkdownDescription: "Policy Category",
					},
					"policy_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Name",
						MarkdownDescription: "Policy Name",
					},
					"policy_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Policy Type",
						MarkdownDescription: "Policy Type",
					},
					"policy_versions": schema.ListNestedAttribute{
						Optional:            true,
						Description:         "Policy Versions",
						MarkdownDescription: "Policy Versions",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "Created At",
									MarkdownDescription: "Created At",
								},
								"created_by": schema.StringAttribute{
									Computed:            true,
									Description:         "Created By",
									MarkdownDescription: "Created By",
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "ID",
									MarkdownDescription: "ID",
								},
								"modified_at": schema.StringAttribute{
									Computed:            true,
									Description:         "Modified At",
									MarkdownDescription: "Modified At",
								},
								"modified_by": schema.StringAttribute{
									Computed:            true,
									Description:         "Modified By",
									MarkdownDescription: "Modified By",
								},
								"policy_document": schema.SingleNestedAttribute{
									Computed:            true,
									Description:         "Policy Document",
									MarkdownDescription: "Policy Document",
									Attributes: map[string]schema.Attribute{
										"statement": schema.ListNestedAttribute{
											Computed:            true,
											Description:         "Statement",
											MarkdownDescription: "Statement",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.ListAttribute{
														Optional:            true,
														Description:         "Action",
														MarkdownDescription: "Action",
														ElementType:         types.StringType,
													},
													"not_action": schema.ListAttribute{
														Optional:            true,
														Description:         "Not Action",
														MarkdownDescription: "Not Action",
														ElementType:         types.StringType,
													},
													"effect": schema.StringAttribute{
														Computed:            true,
														Description:         "Effect",
														MarkdownDescription: "Effect",
													},
													"resource": schema.ListAttribute{
														Optional:            true,
														Description:         "Resource",
														MarkdownDescription: "Resource",
														ElementType:         types.StringType,
													},
													"sid": schema.StringAttribute{
														Computed:            true,
														Description:         "SID",
														MarkdownDescription: "SID",
													},
													"condition": schema.MapAttribute{
														ElementType: types.MapType{
															ElemType: types.ListType{
																ElemType: types.StringType,
															},
														},
														Optional: true,
													},
													"principal": schema.SingleNestedAttribute{
														Optional:            true,
														Description:         "Principal",
														MarkdownDescription: "Principal",
														Attributes: map[string]schema.Attribute{
															"principal_string": schema.StringAttribute{
																Optional: true,
															},
															"principal_map": schema.MapAttribute{
																Optional: true,
																ElementType: types.ListType{
																	ElemType: types.StringType,
																},
															},
														},
													},
												},
											},
										},
										"version": schema.StringAttribute{
											Computed:            true,
											Description:         "Policy Version",
											MarkdownDescription: "Policy Version",
										},
									},
								},

								"policy_id": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy ID",
									MarkdownDescription: "Policy ID",
								},
								"policy_version_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy Version Name",
									MarkdownDescription: "Policy Version Name",
								},
							},
						},
					},
					"resource_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Resource Type",
						MarkdownDescription: "Resource Type",
					},
					"service_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Name",
						MarkdownDescription: "Service Name",
					},
					"service_type": schema.StringAttribute{
						Computed:            true,
						Description:         "Service Type",
						MarkdownDescription: "Service Type",
					},
					"srn": schema.StringAttribute{
						Computed:            true,
						Description:         "SRN",
						MarkdownDescription: "SRN",
					},
					"state": schema.StringAttribute{
						Computed:            true,
						Description:         "State",
						MarkdownDescription: "State",
					},
				},
			},
		},
	}
}

func (r *iamPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Iam
	r.clients = inst.Client
}

func (r *iamPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.PolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreatePolicy(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating policy",
			"Could not create policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var policyVersions []iam.PolicyVersion
	//policy versions
	for _, policyVersion := range data.PolicyVersions {

		var statements []iam.Statement
		for _, _statement := range policyVersion.PolicyDocument.Statement {

			// resource
			resources := make([]types.String, 0, len(_statement.Resource))
			for _, _resource := range _statement.Resource {
				resources = append(resources, types.StringValue(_resource))
			}

			// action
			actions := make([]types.String, 0, len(_statement.Action))
			for _, _action := range _statement.Action {
				actions = append(actions, types.StringValue(_action))
			}

			// not action
			notActions := make([]types.String, 0, len(_statement.NotAction))
			for _, _notAction := range _statement.NotAction {
				notActions = append(notActions, types.StringValue(_notAction))
			}

			// principal
			principal, _ := convertPrincipal(ctx, _statement.Principal)

			// condition
			condition, _ := convertCondition(ctx, _statement.Condition)

			statement := iam.Statement{
				Sid:       types.StringPointerValue(_statement.Sid),
				Effect:    types.StringValue(_statement.Effect),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: principal,
				Condition: condition,
			}

			statements = append(statements, statement)
		}

		policyDocument := iam.PolicyDocument{
			Version:   types.StringValue(policyVersion.PolicyDocument.Version),
			Statement: statements,
		}

		policyVersionState := iam.PolicyVersion{
			CreatedAt:         types.StringValue(policyVersion.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(policyVersion.CreatedBy),
			Id:                types.StringValue(*policyVersion.Id),
			ModifiedAt:        types.StringValue(policyVersion.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(policyVersion.ModifiedBy),
			PolicyDocument:    policyDocument,
			PolicyId:          types.StringValue(*policyVersion.PolicyId),
			PolicyVersionName: types.StringValue(*policyVersion.PolicyVersionName),
		}
		policyVersions = append(policyVersions, policyVersionState)

	}

	plan.Id = types.StringPointerValue(data.Id)

	policyState := iam.Policy{
		AccountId:        types.StringPointerValue(data.AccountId.Get()),
		CreatedAt:        types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.CreatedBy),
		CreatorEmail:     types.StringPointerValue(data.CreatorEmail.Get()),
		CreatorName:      types.StringPointerValue(data.CreatorName.Get()),
		DefaultVersionId: types.StringValue(*data.DefaultVersionId),
		Description:      types.StringPointerValue(data.Description.Get()),
		DomainName:       types.StringValue(data.DomainName),
		Id:               types.StringValue(*data.Id),
		ModifiedAt:       types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.ModifiedBy),
		ModifierEmail:    types.StringPointerValue(data.ModifierEmail.Get()),
		ModifierName:     types.StringPointerValue(data.ModifierName.Get()),
		PolicyCategory:   types.StringValue(string(*data.PolicyCategory)),
		PolicyName:       types.StringValue(*data.PolicyName),
		PolicyType:       types.StringValue(string(*data.PolicyType)),
		PolicyVersions:   policyVersions,
		ResourceType:     types.StringPointerValue(data.ResourceType.Get()),
		ServiceName:      types.StringPointerValue(data.ServiceName.Get()),
		ServiceType:      types.StringPointerValue(data.ServiceType.Get()),
		Srn:              types.StringValue(data.Srn),
		State:            types.StringValue(string(*data.State)),
	}

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	plan.Policy = policyObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.PolicyResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Policy",
			err.Error(),
		)
		return
	}

	var policyVersions []iam.PolicyVersion
	//policy versions
	for _, policyVersion := range data.PolicyVersions {

		var statements []iam.Statement
		for _, _statement := range policyVersion.PolicyDocument.Statement {

			// resource
			resources := make([]types.String, 0, len(_statement.Resource))
			for _, _resource := range _statement.Resource {
				resources = append(resources, types.StringValue(_resource))
			}

			// action
			actions := make([]types.String, 0, len(_statement.Action))
			for _, _action := range _statement.Action {
				actions = append(actions, types.StringValue(_action))
			}

			// not action
			notActions := make([]types.String, 0, len(_statement.NotAction))
			for _, _notAction := range _statement.NotAction {
				notActions = append(notActions, types.StringValue(_notAction))
			}

			// principal
			principal, _ := convertPrincipal(ctx, _statement.Principal)

			// condition
			condition, _ := convertCondition(ctx, _statement.Condition)

			statement := iam.Statement{
				Sid:       types.StringPointerValue(_statement.Sid),
				Effect:    types.StringValue(_statement.Effect),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: principal,
				Condition: condition,
			}

			statements = append(statements, statement)
		}

		policyDocument := iam.PolicyDocument{
			Version:   types.StringValue(policyVersion.PolicyDocument.Version),
			Statement: statements,
		}

		policyVersionState := iam.PolicyVersion{
			CreatedAt:         types.StringValue(policyVersion.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(policyVersion.CreatedBy),
			Id:                types.StringValue(*policyVersion.Id),
			ModifiedAt:        types.StringValue(policyVersion.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(policyVersion.ModifiedBy),
			PolicyDocument:    policyDocument,
			PolicyId:          types.StringValue(*policyVersion.PolicyId),
			PolicyVersionName: types.StringValue(*policyVersion.PolicyVersionName),
		}
		policyVersions = append(policyVersions, policyVersionState)

	}

	policyState := iam.Policy{
		AccountId:        types.StringPointerValue(data.AccountId.Get()),
		CreatedAt:        types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.CreatedBy),
		CreatorEmail:     types.StringPointerValue(data.CreatorEmail.Get()),
		CreatorName:      types.StringPointerValue(data.CreatorName.Get()),
		DefaultVersionId: types.StringValue(*data.DefaultVersionId),
		Description:      types.StringPointerValue(data.Description.Get()),
		DomainName:       types.StringValue(data.DomainName),
		Id:               types.StringValue(*data.Id),
		ModifiedAt:       types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.ModifiedBy),
		ModifierEmail:    types.StringPointerValue(data.ModifierEmail.Get()),
		ModifierName:     types.StringPointerValue(data.ModifierName.Get()),
		PolicyCategory:   types.StringValue(string(*data.PolicyCategory)),
		PolicyName:       types.StringValue(*data.PolicyName),
		PolicyType:       types.StringValue(string(*data.PolicyType)),
		PolicyVersions:   policyVersions,
		ResourceType:     types.StringPointerValue(data.ResourceType.Get()),
		ServiceName:      types.StringPointerValue(data.ServiceName.Get()),
		ServiceType:      types.StringPointerValue(data.ServiceType.Get()),
		Srn:              types.StringValue(data.Srn),
		State:            types.StringValue(string(*data.State)),
	}

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	state.Policy = policyObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state iam.PolicyResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdatePolicy(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Policy",
			"Could not update Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetPolicy(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Policy",
			"Could not read policy ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	var policyVersions []iam.PolicyVersion
	//policy versions
	for _, policyVersion := range data.PolicyVersions {

		var statements []iam.Statement
		for _, _statement := range policyVersion.PolicyDocument.Statement {

			// resource
			resources := make([]types.String, 0, len(_statement.Resource))
			for _, _resource := range _statement.Resource {
				resources = append(resources, types.StringValue(_resource))
			}

			// action
			actions := make([]types.String, 0, len(_statement.Action))
			for _, _action := range _statement.Action {
				actions = append(actions, types.StringValue(_action))
			}

			// not action
			notActions := make([]types.String, 0, len(_statement.NotAction))
			for _, _notAction := range _statement.NotAction {
				notActions = append(notActions, types.StringValue(_notAction))
			}

			// principal
			principal, _ := convertPrincipal(ctx, _statement.Principal)

			// condition
			condition, _ := convertCondition(ctx, _statement.Condition)

			statement := iam.Statement{
				Sid:       types.StringPointerValue(_statement.Sid),
				Effect:    types.StringValue(_statement.Effect),
				Resource:  resources,
				Action:    actions,
				NotAction: notActions,
				Principal: principal,
				Condition: condition,
			}

			statements = append(statements, statement)
		}

		policyDocument := iam.PolicyDocument{
			Version:   types.StringValue(policyVersion.PolicyDocument.Version),
			Statement: statements,
		}

		policyVersionState := iam.PolicyVersion{
			CreatedAt:         types.StringValue(policyVersion.CreatedAt.Format(time.RFC3339)),
			CreatedBy:         types.StringValue(policyVersion.CreatedBy),
			Id:                types.StringValue(*policyVersion.Id),
			ModifiedAt:        types.StringValue(policyVersion.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:        types.StringValue(policyVersion.ModifiedBy),
			PolicyDocument:    policyDocument,
			PolicyId:          types.StringValue(*policyVersion.PolicyId),
			PolicyVersionName: types.StringValue(*policyVersion.PolicyVersionName),
		}
		policyVersions = append(policyVersions, policyVersionState)

	}

	policyState := iam.Policy{
		AccountId:        types.StringPointerValue(data.AccountId.Get()),
		CreatedAt:        types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:        types.StringValue(data.CreatedBy),
		CreatorEmail:     types.StringPointerValue(data.CreatorEmail.Get()),
		CreatorName:      types.StringPointerValue(data.CreatorName.Get()),
		DefaultVersionId: types.StringValue(*data.DefaultVersionId),
		Description:      types.StringPointerValue(data.Description.Get()),
		DomainName:       types.StringValue(data.DomainName),
		Id:               types.StringValue(*data.Id),
		ModifiedAt:       types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:       types.StringValue(data.ModifiedBy),
		ModifierEmail:    types.StringPointerValue(data.ModifierEmail.Get()),
		ModifierName:     types.StringPointerValue(data.ModifierName.Get()),
		PolicyCategory:   types.StringValue(string(*data.PolicyCategory)),
		PolicyName:       types.StringValue(*data.PolicyName),
		PolicyType:       types.StringValue(string(*data.PolicyType)),
		PolicyVersions:   policyVersions,
		ResourceType:     types.StringPointerValue(data.ResourceType.Get()),
		ServiceName:      types.StringPointerValue(data.ServiceName.Get()),
		ServiceType:      types.StringPointerValue(data.ServiceType.Get()),
		Srn:              types.StringValue(data.Srn),
		State:            types.StringValue(string(*data.State)),
	}

	policyObjectValue, diags := types.ObjectValueFrom(ctx, policyState.Attributes(), policyState)
	state.Policy = policyObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.PolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicy(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam policy",
			"Could not delete Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertCondition(ctx context.Context, rawCondition map[string]map[string][]interface{}) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditionMapType := types.MapType{
		ElemType: types.ListType{
			ElemType: types.StringType,
		},
	}

	if rawCondition == nil {
		emptyOuterMap := map[string]attr.Value{}

		emptyConditionMap, emptyConditionDiags := types.MapValueFrom(ctx, conditionMapType, emptyOuterMap)
		diags.Append(emptyConditionDiags...)
		if emptyConditionDiags.HasError() {
			emptyConditionMap = types.MapUnknown(conditionMapType)
		}

		return emptyConditionMap, diags
	}

	outerMap := map[string]attr.Value{}
	for condType, innerMap := range rawCondition {
		if innerMap == nil {
			outerMap[condType] = types.MapNull(types.MapType{
				ElemType: types.ListType{ElemType: types.StringType},
			})
			continue
		}

		inner := map[string]attr.Value{}
		for key, values := range innerMap {
			stringValues := make([]attr.Value, len(values))
			for i, v := range values {
				if s, ok := v.(string); ok {
					stringValues[i] = types.StringValue(s)
				} else {
					stringValues[i] = types.StringNull()
					diags.AddAttributeWarning(
						path.Root("condition"),
						"Invalid Condition Value",
						"Value is not a string. Using null instead.",
					)
				}
			}

			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, stringValues)
			diags.Append(listDiags...)
			if listDiags.HasError() {
				listValue = types.ListNull(types.StringType)
			}

			inner[key] = listValue
		}

		innerMapType := types.ListType{
			ElemType: types.StringType,
		}

		mapValue, mapDiags := types.MapValueFrom(ctx, innerMapType, inner)
		if mapDiags.HasError() {
			mapValue = types.MapUnknown(innerMapType)
		}
		outerMap[condType] = mapValue
	}

	conditionMap, condDiags := types.MapValueFrom(ctx, conditionMapType, outerMap)
	diags.Append(condDiags...)
	if condDiags.HasError() {
		conditionMap = types.MapUnknown(conditionMapType)
	}

	return conditionMap, diags
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *iamPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
