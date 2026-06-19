package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &iamUserResource{}
	_ resource.ResourceWithConfigure = &iamUserResource{}
)

// NewIamUserResource is a helper function to simplify the provider implementation.
func NewIamUserResource() resource.Resource {
	return &iamUserResource{}
}

// iamUserResource is the data source implementation.
type iamUserResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *iamUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user"
}

func (r *iamUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User.",
		Attributes: map[string]schema.Attribute{
			// account_id is the {account_id} PATH parameter of
			// POST /v1/accounts/{account_id}/users (server-side required). It was
			// previously Optional with no validator, so an unset value sent an empty
			// path segment (".../accounts//users"), producing a malformed signed URL
			// that the gateway rejected as a misleading "401 [HMAC] HMAC valid fail".
			// It is now Required with a non-empty validator so the user gets a clear
			// plan-time error instead of an opaque auth failure. See fork issue #74.
			"account_id": schema.StringAttribute{
				Required:            true,
				Description:         "Account ID (required: the owning account in which to create the user)",
				MarkdownDescription: "Account ID (required: the owning account in which to create the user)",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Description",
				MarkdownDescription: "Description",
			},
			"group_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Group IDs",
				MarkdownDescription: "Group IDs",
			},
			"policy_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Policy IDs",
				MarkdownDescription: "Policy IDs",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Description:         "Password",
				MarkdownDescription: "Password",
			},
			"tags": tag.ResourceSchema(),
			"temporary_password": schema.BoolAttribute{
				Optional:            true,
				Description:         "Temporary Password",
				MarkdownDescription: "Temporary Password",
			},
			"user_name": schema.StringAttribute{
				Optional:            true,
				Description:         "User Name",
				MarkdownDescription: "User Name",
			},
			"password_reuse_count": schema.Int32Attribute{
				Optional:            true,
				Description:         "Password Reuse Count",
				MarkdownDescription: "Password Reuse Count",
			},
			"user_id": schema.StringAttribute{
				Computed:            true,
				Description:         "User ID",
				MarkdownDescription: "User ID",
			},
			"user": schema.SingleNestedAttribute{
				Description: "A detail of User.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID",
						MarkdownDescription: "Account ID",
					},
					"company_name": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Company Name",
						MarkdownDescription: "Company Name",
					},
					"console_url": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Console URL",
						MarkdownDescription: "Console URL",
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
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Description",
						MarkdownDescription: "Description",
					},
					"dst_offset": schema.StringAttribute{
						Computed:            true,
						Description:         "Dst Offset",
						MarkdownDescription: "Dst Offset",
					},
					"email": schema.StringAttribute{
						Computed:            true,
						Description:         "Email",
						MarkdownDescription: "Email",
					},
					"email_authenticated": schema.BoolAttribute{
						Computed:            true,
						Description:         "Email Authenticated",
						MarkdownDescription: "Email Authenticated",
					},
					"first_name": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "First Name",
						MarkdownDescription: "First Name",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID",
						MarkdownDescription: "ID",
					},
					"last_login_at": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Last Login At",
						MarkdownDescription: "Last Login At",
					},
					"last_name": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Last Name",
						MarkdownDescription: "Last Name",
					},
					"last_password_update_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Last Password Update At",
						MarkdownDescription: "Last Password Update At",
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
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Name",
						MarkdownDescription: "Name",
					},
					"password": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Password",
						MarkdownDescription: "Password",
					},
					"password_reuse_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Password Reuse Count",
						MarkdownDescription: "Password Reuse Count",
					},
					"phone_authenticated": schema.BoolAttribute{
						Computed:            true,
						Description:         "Phone Authenticated",
						MarkdownDescription: "Phone Authenticated",
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
					"timezone": schema.StringAttribute{
						Computed:            true,
						Description:         "Timezone",
						MarkdownDescription: "Timezone",
					},
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Type",
						MarkdownDescription: "Type",
					},
					"tz_id": schema.StringAttribute{
						Computed:            true,
						Description:         "TZ ID",
						MarkdownDescription: "TZ ID",
					},
					"user_name": schema.StringAttribute{
						Computed:            true,
						Description:         "User Name",
						MarkdownDescription: "User Name",
					},
					"utc_offset": schema.StringAttribute{
						Computed:            true,
						Description:         "UTC Offset",
						MarkdownDescription: "UTC Offset",
					},
					"access_keys": schema.ListNestedAttribute{
						Computed:            true,
						Description:         "Access Keys",
						MarkdownDescription: "Access Keys",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_key": schema.StringAttribute{
									Description:         "Access Key",
									MarkdownDescription: "Access Key",
									Computed:            true,
								},
								"created_at": schema.StringAttribute{
									Description:         "Created At",
									MarkdownDescription: "Created At",
									Computed:            true,
								},
								"expiration_timestamp": schema.StringAttribute{
									Description:         "Expiration Timestmap",
									MarkdownDescription: "Expiration Timestmap",
									Computed:            true,
								},
								"id": schema.StringAttribute{
									Description:         "ID",
									MarkdownDescription: "ID",
									Computed:            true,
								},
								"is_enabled": schema.BoolAttribute{
									Description:         "Is Enabled",
									MarkdownDescription: "Is Enabled",
									Computed:            true,
								},
							},
						},
					},
					"groups": schema.ListNestedAttribute{
						Computed:            true,
						Description:         "Groups",
						MarkdownDescription: "Groups",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Group ID",
									MarkdownDescription: "Group ID",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Group Name",
									MarkdownDescription: "Group Name",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *iamUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.UserResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// account_id is the {account_id} path parameter of CreateIAMUser. If it is
	// empty the request targets ".../accounts//users", which is signed as a
	// malformed URL and surfaces as a misleading "401 [HMAC] HMAC valid fail"
	// (fork issue #74). Default it to the caller's own account id when omitted,
	// and fail with an actionable error rather than an opaque 401 if it cannot
	// be resolved.
	if plan.AccountId.ValueString() == "" {
		accountId, acctErr := r.client.GetAccountId()
		if acctErr != nil || accountId == "" {
			resp.Diagnostics.AddError(
				"Error creating user",
				"account_id is required to create an IAM user (it is the {account_id} path "+
					"segment of POST /v1/accounts/{account_id}/users). It was empty and the "+
					"caller's own account id could not be resolved automatically. Set "+
					"`account_id` on the samsungcloudplatformv2_iam_user resource.",
			)
			return
		}
		plan.AccountId = types.StringValue(accountId)
	}

	data, err := r.client.CreateUser(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	plan.UserId = types.StringValue(data.Id)

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		ConsoleUrl:           types.StringValue(data.ConsoleUrl),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		Password:             types.StringValue(data.Password),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	plan.User = userObjectValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.UserResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Show User",
			err.Error(),
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	state.User = userObjectValue
	state.UserId = types.StringValue(data.Id)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.UserResource
	var state iam.UserResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	_, err := r.client.UpdateUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating User",
			"Could not update User, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Unable to Read User",
			"Could not read User ID "+state.UserId.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// policies
	policies, hasError := getPolicies(ctx, data.Policies)
	if hasError {
		return
	}

	// user nil check
	userCompanyName := data.CompanyName.Get()
	if userCompanyName == nil {
		emptyStr := ""
		userCompanyName = &emptyStr
	}

	userFirstName := data.FirstName.Get()
	if userFirstName == nil {
		emptyStr := ""
		userFirstName = &emptyStr
	}

	userLastName := data.LastName.Get()
	if userLastName == nil {
		emptyStr := ""
		userLastName = &emptyStr
	}

	userLastLoginAt := data.LastLoginAt.Get()
	if userLastLoginAt == nil {
		emptyTime := time.Time{}
		userLastLoginAt = &emptyTime
	}

	userLastPasswordUpdateAt := data.LastPasswordUpdateAt.Get()
	if userLastPasswordUpdateAt == nil {
		emptyTime := time.Time{}
		userLastPasswordUpdateAt = &emptyTime
	}

	// empty list
	accessKeyInfos := make([]iam.AccessKeyV1Dot4, 0)
	groupInfos := make([]iam.GroupInfo, 0)

	userState := iam.User{
		AccountId:            types.StringValue(*data.AccountId.Get()),
		CompanyName:          types.StringValue(*userCompanyName),
		CreatedAt:            types.StringValue(data.CreatedAt.Format(time.RFC3339)),
		CreatedBy:            types.StringValue(data.CreatedBy),
		Description:          types.StringValue(*data.Description.Get()),
		DstOffset:            types.StringValue(*data.DstOffset.Get()),
		Email:                types.StringValue(*data.Email.Get()),
		EmailAuthenticated:   types.BoolValue(data.EmailAuthenticated),
		FirstName:            types.StringValue(*userFirstName),
		Id:                   types.StringValue(data.Id),
		LastLoginAt:          types.StringValue(userLastLoginAt.Format(time.RFC3339)),
		LastName:             types.StringValue(*userLastName),
		LastPasswordUpdateAt: types.StringValue(userLastPasswordUpdateAt.Format(time.RFC3339)),
		ModifiedAt:           types.StringValue(data.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:           types.StringValue(data.ModifiedBy),
		Name:                 types.StringValue(*data.Name.Get()),
		PasswordReuseCount:   types.Int32Value(data.PasswordReuseCount),
		PhoneAuthenticated:   types.BoolValue(data.PhoneAuthenticated),
		Policies:             policies,
		Timezone:             types.StringValue(*data.Timezone.Get()),
		Type:                 types.StringValue(data.Type),
		TzId:                 types.StringValue(*data.TzId.Get()),
		UserName:             types.StringValue(*data.UserName.Get()),
		UtcOffset:            types.StringValue(*data.UtcOffset.Get()),
		AccessKeys:           accessKeyInfos,
		Groups:               groupInfos,
	}

	userObjectValue, diags := types.ObjectValueFrom(ctx, userState.AttributeTypes(), userState)
	plan.User = userObjectValue
	plan.UserId = state.UserId

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.UserResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.AccountId.ValueString(), state.UserId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam User",
			"Could not delete User, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
