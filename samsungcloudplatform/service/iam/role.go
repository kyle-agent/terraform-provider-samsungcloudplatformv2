package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpiam1d0 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamRoleResource{}
	_ resource.ResourceWithConfigure   = &iamRoleResource{}
	_ resource.ResourceWithImportState = &iamRoleResource{}
)

// NewIamRoleResource is a helper function to simplify the provider implementation.
func NewIamRoleResource() resource.Resource {
	return &iamRoleResource{}
}

// iamRoleResource is the data source implementation.
type iamRoleResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

func (r *iamRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Role.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Account ID",
				MarkdownDescription: "Account ID",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Description",
				MarkdownDescription: "Description",
			},
			"max_session_duration": schema.Int32Attribute{
				Optional:            true,
				Description:         "Max Session Duration",
				MarkdownDescription: "Max Session Duration",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Name",
				MarkdownDescription: "Name",
			},
			"policy_ids": schema.ListAttribute{
				Optional:            true,
				Description:         "Policy IDs",
				MarkdownDescription: "Policy IDs",
				ElementType:         types.StringType,
			},
			"principals": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Principals",
				MarkdownDescription: "Principals",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:            true,
							Description:         "Type of principal",
							MarkdownDescription: "Type of principal",
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Description:         "Value of principal",
							MarkdownDescription: "Value of principal",
						},
					},
				},
			},
			"tags": tag.ResourceSchema(),
			"assume_role_policy_document": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Assume Role Policy Document",
				MarkdownDescription: "Assume Role Policy Document",
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
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID",
				MarkdownDescription: "ID",
			},
			"role": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Role",
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
					"description": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Description",
						MarkdownDescription: "Description",
					},
					"account_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Account ID",
						MarkdownDescription: "Account ID",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
					},
					"max_session_duration": schema.Int32Attribute{
						Computed:            true,
						Description:         "Max Session Duration",
						MarkdownDescription: "Max Session Duration",
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
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Name",
						MarkdownDescription: "Name",
					},
					"assume_role_policy_document": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Assume Role Policy Document",
						MarkdownDescription: "Assume Role Policy Document",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								Optional:            true,
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

													Computed: true, ElementType: types.ListType{
														ElemType: types.StringType,
													},
												},
											},
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
					"policies": schema.ListNestedAttribute{
						Optional:            true,
						Description:         "Policies",
						MarkdownDescription: "Policies",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Optional:            true,
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
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Type",
						MarkdownDescription: "Type",
					},
				},
			},
		},
	}
}

func (r *iamRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.RoleResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateRole(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating role",
			"Could not create role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
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
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	plan.Id = types.StringValue(data.Role.Id)

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)
	plan.Role = roleObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.RoleResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetRole(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Role",
			err.Error(),
		)
		return
	}

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
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
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)
	state.Role = roleObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.RoleResource
	var state iam.RoleResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateRole(ctx, state.Id.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Role",
			"Could not update Role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetRole(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Role",
			"Could not read role ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// assume role policy document
	var assumeRolePolicyDocument iam.PolicyDocument

	var statements []iam.Statement
	for _, _statement := range data.Role.AssumeRolePolicyDocument.Statement {
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
			Sid:       types.StringValue(*_statement.Sid),
			Effect:    types.StringValue(_statement.Effect),
			Resource:  resources,
			Action:    actions,
			NotAction: notActions,
			Principal: principal,
			Condition: condition,
		}

		statements = append(statements, statement)
	}

	assumeRolePolicyDocument = iam.PolicyDocument{
		Version:   types.StringValue(data.Role.AssumeRolePolicyDocument.Version),
		Statement: statements,
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Role.Policies)
	if hasError {
		return
	}

	roleState := iam.Role{
		AccountId:                types.StringValue(*data.Role.AccountId.Get()),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		CreatedAt:                types.StringValue(data.Role.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                types.StringValue(data.Role.CreatedBy),
		CreatorEmail:             types.StringValue(*data.Role.CreatorEmail.Get()),
		CreatorName:              types.StringValue(*data.Role.CreatorName.Get()),
		Description:              types.StringPointerValue(data.Role.Description.Get()),
		Id:                       types.StringValue(data.Role.Id),
		MaxSessionDuration:       types.Int32Value(data.Role.MaxSessionDuration),
		ModifiedAt:               types.StringValue(data.Role.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:               types.StringValue(data.Role.ModifiedBy),
		ModifierEmail:            types.StringValue(*data.Role.ModifierEmail.Get()),
		ModifierName:             types.StringValue(*data.Role.ModifierName.Get()),
		Name:                     types.StringValue(data.Role.Name),
		Policies:                 policies,
		Type:                     types.StringValue(string(*data.Role.Type)),
	}

	roleObjectValue, diags := types.ObjectValueFrom(ctx, roleState.AttributeTypes(), roleState)

	plan.Role = roleObjectValue
	plan.Id = state.Id

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.RoleResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRole(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam Role",
			"Could not delete Role, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func convertPrincipal(ctx context.Context, principal interface{}) (iam.Principal, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch v := principal.(type) {

	case scpiam1d0.NullablePrincipal:

		if v.Get() == nil {
			_principal := iam.Principal{
				PrincipalString: types.StringNull(),
				PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
			}

			return _principal, diags
		}

		if v.Get().String != nil {
			_principal := iam.Principal{
				PrincipalString: types.StringValue(*v.Get().String),
				PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
			}
			return _principal, diags
		}

		if v.Get().MapmapOfStringarrayOfString != nil {
			tempMap := make(map[string]types.List, len(*v.Get().MapmapOfStringarrayOfString))
			for key, val := range *v.Get().MapmapOfStringarrayOfString {
				listVal, listDiags := types.ListValueFrom(ctx, types.StringType, val)
				diags.Append(listDiags...)
				tempMap[key] = listVal
			}

			mapVal, mapDiags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, tempMap)
			diags.Append(mapDiags...)

			if diags.HasError() {
				_principal := iam.Principal{
					PrincipalString: types.StringNull(),
					PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
				}
				return _principal, diags
			}

			_principal := iam.Principal{
				PrincipalString: types.StringNull(),
				PrincipalMap:    mapVal,
			}
			return _principal, diags
		}

	default:
		diags.AddError(
			"Error converting principal",
			fmt.Sprintf("Unsupported principal type: %T", v),
		)
		_principal := iam.Principal{
			PrincipalString: types.StringNull(),
			PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
		}
		return _principal, diags
	}

	_principal := iam.Principal{
		PrincipalString: types.StringNull(),
		PrincipalMap:    types.MapNull(types.ListType{ElemType: types.StringType}),
	}

	return _principal, diags
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *iamRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
