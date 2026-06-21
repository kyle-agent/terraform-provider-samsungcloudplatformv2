package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamGroupResource{}
	_ resource.ResourceWithConfigure   = &iamGroupResource{}
	_ resource.ResourceWithImportState = &iamGroupResource{}
)

// NewIamGroupResource is a helper function to simplify the provider implementation.
func NewIamGroupResource() resource.Resource {
	return &iamGroupResource{}
}

// iamGroupResource is the data source implementation.
type iamGroupResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

func (r *iamGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Group ID",
				MarkdownDescription: "Group ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Group Name",
				MarkdownDescription: "Group Name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Group Description",
				MarkdownDescription: "Group Description",
			},
			"tags": tag.ResourceSchema(),
			"policy_ids": schema.ListAttribute{
				Optional:    true,
				Description: "Policy IDs",
				ElementType: types.StringType,
			},
			"user_ids": schema.ListAttribute{
				Optional:    true,
				Description: "User IDs",
				ElementType: types.StringType,
			},
			"group": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "A detail of Group.",
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "생성 일시",
						MarkdownDescription: "생성 일시",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자",
						MarkdownDescription: "생성자",
					},
					"creator_email": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자 Email",
						MarkdownDescription: "생성자 Email",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자 성, 이름",
						MarkdownDescription: "생성자 성, 이름",
					},
					"description": schema.StringAttribute{
						Computed: true,
					},
					"domain_name": schema.StringAttribute{
						Computed:            true,
						Description:         "도메인 이름",
						MarkdownDescription: "도메인 이름",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
					},
					"members": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Members",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "생성 일시",
									MarkdownDescription: "생성 일시",
								},
								"created_by": schema.StringAttribute{
									Computed:            true,
									Description:         "생성자",
									MarkdownDescription: "생성자",
								},
								"creator_created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "생성 일시",
									MarkdownDescription: "생성 일시",
								},
								"creator_email": schema.StringAttribute{
									Computed:            true,
									Description:         "생성자 Email",
									MarkdownDescription: "생성자 Email",
								},
								"creator_last_login_at": schema.StringAttribute{
									Optional:            true,
									Description:         "생성자 마지막 로그인 일시",
									MarkdownDescription: "생성자 마지막 로그인 일시",
								},
								"creator_name": schema.StringAttribute{
									Computed:            true,
									Description:         "생성자 성, 이름",
									MarkdownDescription: "생성자 성, 이름",
								},
								"group_names": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Description:         "Group Names",
									MarkdownDescription: "Group Names",
								},
								"user_created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "생성 일시",
									MarkdownDescription: "생성 일시",
								},
								"user_email": schema.StringAttribute{
									Computed:            true,
									Description:         "User Email",
									MarkdownDescription: "User Email",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									Description:         "User ID",
									MarkdownDescription: "User ID",
								},
								"user_last_login_at": schema.StringAttribute{
									Optional:            true,
									Description:         "User 마지막 로그인 일시",
									MarkdownDescription: "User 마지막 로그인 일시",
								},
								"user_name": schema.StringAttribute{
									Computed:            true,
									Description:         "User 성, 이름",
									MarkdownDescription: "User 성, 이름",
								},
							},
						},
					},
					"policies": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
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

					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "수정 일시",
						MarkdownDescription: "수정 일시",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "수정자",
						MarkdownDescription: "수정자",
					},
					"modifier_email": schema.StringAttribute{
						Computed:            true,
						Description:         "수정자 Email",
						MarkdownDescription: "수정자 Email",
					},
					"modifier_name": schema.StringAttribute{
						Computed:            true,
						Description:         "수정자 성, 이름",
						MarkdownDescription: "수정자 성, 이름",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Group 이름",
						MarkdownDescription: "Group 이름",
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
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Group Type",
						MarkdownDescription: "Group Type",
					},
				},
			},
		},
	}
}

func (r *iamGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan iam.GroupResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateGroup(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// group
	group := data.Group

	// group members
	members := getGroupMembers(group.Members)

	// policies
	policies, hasError := getPolicies(ctx, group.Policies)
	if hasError {
		return
	}

	plan.Id = types.StringValue(group.Id)

	groupState := iam.Group{
		CreatedAt:     types.StringValue(group.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(group.CreatedBy),
		CreatorEmail:  types.StringPointerValue(group.CreatorEmail),
		CreatorName:   types.StringPointerValue(group.CreatorName),
		Description:   types.StringPointerValue(group.Description.Get()),
		DomainName:    types.StringValue(group.DomainName),
		Id:            types.StringValue(group.Id),
		Members:       members,
		Policies:      policies,
		ModifiedAt:    types.StringValue(group.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(group.ModifiedBy),
		ModifierEmail: types.StringPointerValue(group.ModifierEmail),
		ModifierName:  types.StringPointerValue(group.ModifierName),
		Name:          types.StringValue(group.Name),
		ResourceType:  types.StringPointerValue(group.ResourceType.Get()),
		ServiceName:   types.StringPointerValue(group.ServiceName.Get()),
		ServiceType:   types.StringPointerValue(group.ServiceType.Get()),
		Srn:           types.StringPointerValue(group.Srn.Get()),
		GroupType:     types.StringValue(group.Type),
	}
	groupObjectValue, diags := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	plan.Group = groupObjectValue
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *iamGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state iam.GroupResource
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing Group
	_, err := r.client.UpdateGroup(ctx, state.Id.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Group",
			"Could not update Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read Group",
			"Could not read group ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// group
	group := data.Group

	// group members
	members := getGroupMembers(group.Members)

	// policies
	policies, hasError := getPolicies(ctx, group.Policies)
	if hasError {
		return
	}

	groupState := iam.Group{
		CreatedAt:     types.StringValue(group.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(group.CreatedBy),
		CreatorEmail:  types.StringPointerValue(group.CreatorEmail),
		CreatorName:   types.StringPointerValue(group.CreatorName),
		Description:   types.StringPointerValue(group.Description.Get()),
		DomainName:    types.StringValue(group.DomainName),
		Id:            types.StringValue(group.Id),
		Members:       members,
		Policies:      policies,
		ModifiedAt:    types.StringValue(group.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(group.ModifiedBy),
		ModifierEmail: types.StringPointerValue(group.ModifierEmail),
		ModifierName:  types.StringPointerValue(group.ModifierName),
		Name:          types.StringValue(group.Name),
		ResourceType:  types.StringPointerValue(group.ResourceType.Get()),
		ServiceName:   types.StringPointerValue(group.ServiceName.Get()),
		ServiceType:   types.StringPointerValue(group.ServiceType.Get()),
		Srn:           types.StringPointerValue(group.Srn.Get()),
		GroupType:     types.StringValue(group.Type),
	}
	groupObjectValue, diags := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	state.Group = groupObjectValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state iam.GroupResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing iam group
	err := r.client.DeleteGroup(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group",
			"Could not delete Group, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *iamGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.GroupResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetGroup(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show Group",
			err.Error(),
		)
		return
	}

	// group
	group := data.Group

	// group members
	members := getGroupMembers(group.Members)

	// policies
	policies, hasError := getPolicies(ctx, group.Policies)
	if hasError {
		return
	}

	groupState := iam.Group{
		CreatedAt:     types.StringValue(group.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(group.CreatedBy),
		CreatorEmail:  types.StringPointerValue(group.CreatorEmail),
		CreatorName:   types.StringPointerValue(group.CreatorName),
		Description:   types.StringPointerValue(group.Description.Get()),
		DomainName:    types.StringValue(group.DomainName),
		Id:            types.StringValue(group.Id),
		Members:       members,
		Policies:      policies,
		ModifiedAt:    types.StringValue(group.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(group.ModifiedBy),
		ModifierEmail: types.StringPointerValue(group.ModifierEmail),
		ModifierName:  types.StringPointerValue(group.ModifierName),
		Name:          types.StringValue(group.Name),
		ResourceType:  types.StringPointerValue(group.ResourceType.Get()),
		ServiceName:   types.StringPointerValue(group.ServiceName.Get()),
		ServiceType:   types.StringPointerValue(group.ServiceType.Get()),
		Srn:           types.StringPointerValue(group.Srn.Get()),
		GroupType:     types.StringValue(group.Type),
	}

	groupObjectValue, _ := types.ObjectValueFrom(ctx, groupState.AttributeTypes(), groupState)
	state.Group = groupObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getGroupMembers(_members []scpsdkiam.GroupMember) []iam.Member {
	var members []iam.Member

	for _, member := range _members {
		var creatorLastLoginAt *string
		var userLastLoginAt *string

		if member.CreatorLastLoginAt.Get() != nil {
			t := member.CreatorLastLoginAt.Get().Format(time.RFC3339)
			creatorLastLoginAt = &t
		}
		if member.UserLastLoginAt.Get() != nil {
			t := member.UserLastLoginAt.Get().Format(time.RFC3339)
			userLastLoginAt = &t
		}

		groupNames := make([]types.String, 0, len(member.GroupNames))
		for _, groupName := range member.GroupNames {
			groupNames = append(groupNames, types.StringValue(groupName))
		}

		memberState := iam.Member{
			CreatedAt:          types.StringValue(member.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(member.CreatedBy),
			CreatorCreatedAt:   types.StringValue(member.CreatorCreatedAt.Format(time.RFC3339)),
			CreatorEmail:       types.StringPointerValue(member.CreatorEmail),
			CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
			CreatorName:        types.StringPointerValue(member.CreatorName),
			GroupNames:         groupNames,
			UserCreatedAt:      types.StringValue(member.UserCreatedAt.Format(time.RFC3339)),
			UserEmail:          types.StringPointerValue(member.UserEmail),
			UserId:             types.StringValue(member.UserId),
			UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
			UserName:           types.StringPointerValue(member.UserName),
		}

		members = append(members, memberState)
	}
	return members
}

func getPolicies(ctx context.Context, _policies interface{}) ([]iam.Policy, bool) {
	var policies []iam.Policy

	switch v := _policies.(type) {
	case []scpsdkiam.Policy:
		policies = _getPolicies(ctx, v, policies)
	}

	return policies, false
}

func _getPolicies(ctx context.Context, _policies []scpsdkiam.Policy, policies []iam.Policy) []iam.Policy {
	for _, policy := range _policies {

		var policyVersions []iam.PolicyVersion
		//policy versions
		for _, policyVersion := range policy.PolicyVersions {

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
			AccountId:        types.StringPointerValue(policy.AccountId.Get()),
			CreatedAt:        types.StringValue(policy.CreatedAt.Format(time.RFC3339)),
			CreatedBy:        types.StringValue(policy.CreatedBy),
			CreatorEmail:     types.StringPointerValue(policy.CreatorEmail.Get()),
			CreatorName:      types.StringPointerValue(policy.CreatorName.Get()),
			DefaultVersionId: types.StringValue(*policy.DefaultVersionId),
			Description:      types.StringPointerValue(policy.Description.Get()),
			DomainName:       types.StringValue(policy.DomainName),
			Id:               types.StringValue(*policy.Id),
			ModifiedAt:       types.StringValue(policy.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:       types.StringValue(policy.ModifiedBy),
			ModifierEmail:    types.StringPointerValue(policy.ModifierEmail.Get()),
			ModifierName:     types.StringPointerValue(policy.ModifierName.Get()),
			PolicyCategory:   types.StringValue(string(*policy.PolicyCategory)),
			PolicyName:       types.StringValue(*policy.PolicyName),
			PolicyType:       types.StringValue(string(*policy.PolicyType)),
			PolicyVersions:   policyVersions,
			ResourceType:     types.StringPointerValue(policy.ResourceType.Get()),
			ServiceName:      types.StringPointerValue(policy.ServiceName.Get()),
			ServiceType:      types.StringPointerValue(policy.ServiceType.Get()),
			Srn:              types.StringValue(policy.Srn),
			State:            types.StringValue(string(*policy.State)),
		}

		policies = append(policies, policyState)
	}
	return policies
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *iamGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
