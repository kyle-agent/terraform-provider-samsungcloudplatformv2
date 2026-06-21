package ske

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	scpske "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/ske/1.4"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeClusterResource{}
	_ resource.ResourceWithConfigure   = &skeClusterResource{}
	_ resource.ResourceWithImportState = &skeClusterResource{}
)

// NewSkeClusterResource is a helper function to simplify the provider implementation.
func NewSkeClusterResource() resource.Resource {
	return &skeClusterResource{}
}

// skeClusterResource is the data source implementation.
type skeClusterResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *skeClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_cluster"
}

// Schema defines the schema for the data source.
func (r *skeClusterResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterResourceSchema(ctx)
}

func ClusterResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "cluster",
		Attributes: map[string]schema.Attribute{
			"cloud_logging_enabled": schema.BoolAttribute{
				Required:            true,
				Description:         "Cloud Logging Enabled\n  - example: true",
				MarkdownDescription: "Cloud Logging Enabled\n  - example: true",
			},
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"cloud_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Cloud Logging Enabled\n  - example: true",
						MarkdownDescription: "Cloud Logging Enabled\n  - example: true",
					},
					"cluster_namespace": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Namespace\n  - example: sample-cluster-12345",
						MarkdownDescription: "Cluster Namespace\n  - example: sample-cluster-12345",
					},
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
						MarkdownDescription: "ID\n  - example: 0fdd87aab8cb46f59b7c1f81ed03fb3e",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Version\n  - example: v1.29.8",
						MarkdownDescription: "Cluster Version\n  - example: v1.29.8",
					},
					"managed_security_group": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource name\n  - example: sample-name",
								MarkdownDescription: "External Resource name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Managed Security Group",
						MarkdownDescription: "Managed Security Group",
					},
					"max_node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Max Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Max Node Count\n  - example: 5",
					},
					"modified_at": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
						MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
					},
					"modified_by": schema.StringAttribute{
						Computed:            true,
						Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Name\n  - example: sample-cluster",
						MarkdownDescription: "Cluster Name\n  - example: sample-cluster",
					},
					"node_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Cluster Node Count\n  - example: 5",
						MarkdownDescription: "Cluster Node Count\n  - example: 5",
					},
					"private_endpoint_access_control_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Name\n  - example: sample-name",
									MarkdownDescription: "Private Endpoint Access Control Resource Name\n  - example: sample-name",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Private Endpoint Access Control Resource Type\n  - example: vm",
									MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - example: vm",
								},
							},
						},
						Computed:            true,
						Description:         "Private Endpoint Access Control Resources",
						MarkdownDescription: "Private Endpoint Access Control Resources",
					},
					"private_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Kubeconfig Download Yn\n  - example: N",
						MarkdownDescription: "Private Kubeconfig Download Yn\n  - example: N",
					},
					"private_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Private Endpoint URL\n  - example: https://sample-cluster.ske.private.kr-west1.samsungsdscloud.com:6443",
					},
					"public_endpoint_access_control_ip": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
						MarkdownDescription: "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
					},
					"public_endpoint_url": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
						MarkdownDescription: "Public Endpoint URL\n  - example: https://sample-cluster.ske.kr-west1.samsungsdscloud.com:6443",
					},
					"public_kubeconfig_download_yn": schema.StringAttribute{
						Computed:            true,
						Description:         "Public Kubeconfig Download Yn\n  - example: N",
						MarkdownDescription: "Public Kubeconfig Download Yn\n  - example: N",
					},
					"security_group_list": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
									MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "External Resource name\n  - example: sample-name",
									MarkdownDescription: "External Resource name\n  - example: sample-name",
								},
							},
						},
						Computed:            true,
						Description:         "Connected Security Group List",
						MarkdownDescription: "Connected Security Group List",
					},
					"service_watch_logging_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Service Watch Enabled\n  - example: true",
						MarkdownDescription: "Service Watch Enabled\n  - example: true",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Cluster Status\n  - example: RUNNING",
						MarkdownDescription: "Cluster Status\n  - example: RUNNING",
					},
					"subnet": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource name\n  - example: sample-name",
								MarkdownDescription: "External Resource name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Subnet of Cluster",
						MarkdownDescription: "Subnet of Cluster",
					},
					"volume": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource name\n  - example: sample-name",
								MarkdownDescription: "External Resource name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "Connected File Storage",
						MarkdownDescription: "Connected File Storage",
					},
					"vpc": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
								MarkdownDescription: "External Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "External Resource name\n  - example: sample-name",
								MarkdownDescription: "External Resource name\n  - example: sample-name",
							},
						},
						Computed:            true,
						Description:         "VPC of Cluster",
						MarkdownDescription: "VPC of Cluster",
					},
				},
				Computed: true,
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Version\n  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n  - example: v1.29.8",
				MarkdownDescription: "Cluster Version\n  - pattern: ^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$\n  - example: v1.29.8",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^v[0-9]{1}\\.[0-9]{1,2}\\.[0-9]{1,2}$"), ""),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster Name\n  - maxLength: 30\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-cluster",
				MarkdownDescription: "Cluster Name\n  - maxLength: 30\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-cluster",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 30),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9\\-]*[a-z0-9]$"), ""),
				},
			},
			"private_endpoint_access_control_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
							MarkdownDescription: "Private Endpoint Access Control Resource ID\n  - example: 2a9be312-5d4b-4bc8-b2ae-35100fa9241f",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource Name\n  - example: sample-name",
							MarkdownDescription: "Private Endpoint Access Control Resource Name\n  - example: sample-name",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Private Endpoint Access Control Resource Type\n  - example: vm",
							MarkdownDescription: "Private Endpoint Access Control Resource Type\n  - example: vm",
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Private Endpoint Access Control Resources",
				MarkdownDescription: "Private Endpoint Access Control Resources",
			},
			"public_endpoint_access_control_ip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
				MarkdownDescription: "Public Endpoint Access Control IP\n  - example: 192.168.0.0",
			},
			"security_group_id_list": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "Security Group ID List\n  - example: [bdfda539-bd2e-4a5c-9021-ec6d52d1ca79]",
				MarkdownDescription: "Security Group ID List\n  - example: [bdfda539-bd2e-4a5c-9021-ec6d52d1ca79]",
			},
			"service_watch_logging_enabled": schema.BoolAttribute{
				Required:            true,
				Description:         "Service Watch Enabled\n  - example: true",
				MarkdownDescription: "Service Watch Enabled\n  - example: true",
			},
			"subnet_id": schema.StringAttribute{
				Required:            true,
				Description:         "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
				MarkdownDescription: "Subnet ID\n  - example: 023c57b14f11483689338d085e061492",
			},
			"volume_id": schema.StringAttribute{
				Required:            true,
				Description:         "Volume ID\n  - example: [bfdbabf2-04d9-4e8b-a205-020f8e6da438]",
				MarkdownDescription: "Volume ID\n  - example: [bfdbabf2-04d9-4e8b-a205-020f8e6da438]",
			},
			"vpc_id": schema.StringAttribute{
				Required:            true,
				Description:         "VPC ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
				MarkdownDescription: "VPC ID\n  - example: 7df8abb4912e4709b1cb237daccca7a8",
			},
			"tags": tag.ResourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *skeClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Ske
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *skeClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.ClusterResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	data, err := r.client.CreateCluster(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Could not create cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{"CREATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
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
func (r *skeClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from cluster
	data, _, err := r.client.GetCluster(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Cluster",
			"Could not read cluster ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	cluster := data.Cluster

	var securityGroups []ske.ExternalResource
	for _, securityGroup := range cluster.SecurityGroupList {
		securityGroups = append(securityGroups, r.makeExternalResourceModel((*scpske.ExternalResource)(&securityGroup)))
	}

	var privateEndpointAccessControlResources []ske.PrivateEndpointAccessControlResource
	for _, privateEndpointAccessControlResource := range cluster.PrivateEndpointAccessControlResources {
		privateEndpointAccessControlResources = append(privateEndpointAccessControlResources, r.makePrivateEndpointAccessControlResourceModel((*scpske.PrivateEndpointAccessControlResource)(&privateEndpointAccessControlResource)))
	}

	clusterModel := ske.Cluster{
		Id:                                    types.StringValue(cluster.Id),
		Name:                                  types.StringValue(cluster.Name),
		AccountId:                             types.StringValue(cluster.AccountId),
		CloudLoggingEnabled:                   types.BoolValue(cluster.CloudLoggingEnabled),
		KubernetesVersion:                     types.StringValue(cluster.KubernetesVersion),
		ClusterNamespace:                      types.StringValue(cluster.ClusterNamespace),
		MaxNodeCount:                          types.Int32PointerValue(cluster.MaxNodeCount.Get()),
		NodeCount:                             types.Int32PointerValue(cluster.NodeCount.Get()),
		PrivateEndpointUrl:                    types.StringValue(cluster.PrivateEndpointUrl),
		PrivateKubeconfigDownloadYn:           types.StringValue(cluster.PrivateKubeconfigDownloadYn),
		PrivateEndpointAccessControlResources: privateEndpointAccessControlResources,
		PublicEndpointUrl:                     types.StringValue(cluster.GetPublicEndpointUrl()),
		PublicKubeconfigDownloadYn:            types.StringValue(cluster.PublicKubeconfigDownloadYn),
		PublicEndpointAccessControlIp:         types.StringValue(cluster.GetPublicEndpointAccessControlIp()),
		Vpc:                                   r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		Subnet:                                r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Subnet.Get())),
		Volume:                                r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Volume.Get())),
		SecurityGroupList:                     securityGroups,
		ManagedSecurityGroup:                  r.makeExternalResourceModel((*scpske.ExternalResource)(cluster.Vpc.Get())),
		CreatedAt:                             types.StringValue(cluster.CreatedAt.Format(time.RFC3339)),
		CreatedBy:                             types.StringValue(cluster.CreatedBy),
		ModifiedAt:                            types.StringValue(cluster.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:                            types.StringValue(cluster.ModifiedBy),
		Status:                                types.StringValue(cluster.Status),
		ServiceWatchLoggingEnabled:            types.BoolValue(cluster.GetServiceWatchLoggingEnabled()), // v1.1
	}
	clusterObjectValue, _ := types.ObjectValueFrom(ctx, clusterModel.AttributeTypes(), clusterModel)
	state.Cluster = clusterObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *skeClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state ske.ClusterResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	err := r.syncCloudLoggingEnabled(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncSecurityGroupList(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncKubernetesVersion(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPrivateEndpointAccessControlResources(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncPublicEndpointAccessControlIp(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncServiceWatchLoggingEnabled(ctx, state, &plan, resp)
	if err != nil {
		return
	}
	err = r.syncTags(ctx, state, &plan, resp)
	if err != nil {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

func (r *skeClusterResource) syncCloudLoggingEnabled(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if state.CloudLoggingEnabled.Equal(plan.CloudLoggingEnabled) {
		return nil
	}
	data, err := r.client.UpdateClusterLogging(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating CloudLoggingEnabled",
			"Could not update cloud logging enabled, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncSecurityGroupList(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	securityGroupIdListPlan, _ := types.ListValueFrom(ctx, types.StringType, plan.SecurityGroupIdList)
	securityGroupIdListState, _ := types.ListValueFrom(ctx, types.StringType, state.SecurityGroupIdList)
	if securityGroupIdListPlan.Equal(securityGroupIdListState) {
		return nil
	}
	data, err := r.client.UpdateClusterSecurityGroups(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating SecurityGroupList",
			"Could not update security group list, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.Cluster.Id)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncKubernetesVersion(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.KubernetesVersion.Equal(state.KubernetesVersion) {
		return nil
	}
	data, err := r.client.UpgradeCluster(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating KubernetesVersion",
			"Could not update kubernetes version, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPrivateEndpointAccessControlResources(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if reflect.DeepEqual(plan.PrivateEndpointAccessControlResources, state.PrivateEndpointAccessControlResources) {
		return nil
	}
	data, err := r.client.UpdatePrivateEndpointAccessControlResources(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PrivateEndpointAccessControlResources",
			"Could not update cluster private endpoint access control resources, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncPublicEndpointAccessControlIp(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.PublicEndpointAccessControlIp.Equal(state.PublicEndpointAccessControlIp) {
		return nil
	}
	data, err := r.client.UpdatePublicEndpointAccessControlIps(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating PublicEndpointAccessControlIp",
			"Could not update public endpoint access control ip, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncServiceWatchLoggingEnabled(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.ServiceWatchLoggingEnabled.Equal(state.ServiceWatchLoggingEnabled) {
		return nil
	}
	data, err := r.client.UpdateServiceWatchLoggingEnabled(ctx, plan.Id.ValueString(), *plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating ServiceWatchLoggingEnabled",
			"Could not update service watch logging enabled, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(data.ResourceId)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

func (r *skeClusterResource) syncTags(ctx context.Context, state ske.ClusterResource, plan *ske.ClusterResource, resp *resource.UpdateResponse) error {
	if plan.Tags.Equal(state.Tags) {
		return nil
	}
	_, err := tag.UpdateTags(r.clients, "ske", "cluster", plan.Id.ValueString(), plan.Tags.Elements())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating Tags",
			"Could not update tags, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return err
	}
	err = waitForClusterStatus(ctx, r.client, plan.Id.ValueString(), []string{"UPDATING"}, []string{"RUNNING"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cluster",
			"Error waiting for cluster to become running: "+err.Error(),
		)
		return err
	}
	plan.Id = types.StringValue(plan.Id.ValueString())
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return nil
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *skeClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.ClusterResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing cluster
	data, err := r.client.DeleteCluster(ctx, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Could not delete cluster, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForClusterStatus(ctx, r.client, data.ResourceId, []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cluster",
			"Error waiting for cluster to become deleted: "+err.Error(),
		)
		return
	}
}

func (r *skeClusterResource) makeExternalResourceModel(externalResource *scpske.ExternalResource) ske.ExternalResource {
	return ske.ExternalResource{
		Id:   types.StringValue(externalResource.GetId()),
		Name: types.StringValue(externalResource.GetName()),
	}
}

func (r *skeClusterResource) makePrivateEndpointAccessControlResourceModel(privateEndpointAccessControlResource *scpske.PrivateEndpointAccessControlResource) ske.PrivateEndpointAccessControlResource {
	return ske.PrivateEndpointAccessControlResource{
		Id:   types.StringValue(privateEndpointAccessControlResource.GetId()),
		Name: types.StringValue(privateEndpointAccessControlResource.GetName()),
		Type: types.StringValue(privateEndpointAccessControlResource.GetType()),
	}
}

//

func waitForClusterStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetCluster(ctx, id)
		if httpStatus == 200 {
			return info, info.Cluster.Status, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETED", nil
		} else if err != nil {
			return nil, "", err
		}

		return info, info.Cluster.Status, nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *skeClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
