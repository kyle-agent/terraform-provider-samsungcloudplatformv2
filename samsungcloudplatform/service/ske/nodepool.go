package ske

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/ske"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/service/ske/converter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &skeNodepoolResource{}
	_ resource.ResourceWithConfigure   = &skeNodepoolResource{}
	_ resource.ResourceWithImportState = &skeNodepoolResource{}
)

func NewSkeNodepoolResource() resource.Resource {
	return &skeNodepoolResource{}
}

type skeNodepoolResource struct {
	config  *scpsdk.Configuration
	client  *ske.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (r *skeNodepoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ske_nodepool" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

// Schema defines the schema for the data source.
func (r *skeNodepoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "nodepool",
		Attributes: map[string]schema.Attribute{
			"advanced_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"allowed_unsafe_sysctls": schema.StringAttribute{
						Optional:            true,
						Description:         "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
						MarkdownDescription: "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
					},
					"container_log_max_files": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
						MarkdownDescription: "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
						Validators: []validator.Int32{
							int32validator.Between(2, 10),
						},
					},
					"container_log_max_size": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
						MarkdownDescription: "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
						Validators: []validator.Int32{
							int32validator.Between(10, 100),
						},
					},
					"image_gc_high_threshold": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
						MarkdownDescription: "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
						Validators: []validator.Int32{
							int32validator.Between(10, 85),
						},
					},
					"image_gc_low_threshold": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
						MarkdownDescription: "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
						Validators: []validator.Int32{
							int32validator.Between(10, 85),
						},
					},
					"max_pods": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
						MarkdownDescription: "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
						Validators: []validator.Int32{
							int32validator.Between(10, 250),
						},
					},
					"pod_max_pids": schema.Int32Attribute{
						Required:            true,
						Description:         "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
						MarkdownDescription: "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
						Validators: []validator.Int32{
							int32validator.Between(1024, 4194304),
						},
					},
				},
				Optional:            true,
				Description:         "Node Pool Advanced Settings",
				MarkdownDescription: "Node Pool Advanced Settings",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
			},
			"custom_image_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Custom Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Custom Image ID\n  - example: 10a599e031e749b7b260868f441e862b",
			},
			"desired_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Desired node count (is_auto_scale = false)\n  - example: 2",
				MarkdownDescription: "Desired node count (is_auto_scale = false)\n  - example: 2",
			},
			"image_os": schema.StringAttribute{
				Required:            true,
				Description:         "Image OS\n  - example: ubuntu",
				MarkdownDescription: "Image OS\n  - example: ubuntu",
			},
			"image_os_version": schema.StringAttribute{
				Required:            true,
				Description:         "Image OS Version\n  - example: 22.04",
				MarkdownDescription: "Image OS Version\n  - example: 22.04",
			},
			"is_auto_recovery": schema.BoolAttribute{
				Required:            true,
				Description:         "Is Auto Recovery\n  - example: true",
				MarkdownDescription: "Is Auto Recovery\n  - example: true",
			},
			"is_auto_scale": schema.BoolAttribute{
				Required:            true,
				Description:         "Is Auto Scale\n  - example: true",
				MarkdownDescription: "Is Auto Scale\n  - example: true",
			},
			"keypair_name": schema.StringAttribute{
				Required:            true,
				Description:         "Keypair Name\n  - example: test_keypair",
				MarkdownDescription: "Keypair Name\n  - example: test_keypair",
			},
			"kubernetes_version": schema.StringAttribute{
				Required:            true,
				Description:         "Kubernetes Version\n  - example: v1.29.8",
				MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
			},
			"labels": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:            true,
							Description:         "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							MarkdownDescription: "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$"), ""),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							MarkdownDescription: "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							Validators: []validator.String{
								stringvalidator.LengthAtMost(63),
								stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
							},
							Default: stringdefault.StaticString(""),
						},
					},
				},
				Optional:            true,
				Description:         "Node Pool Labels",
				MarkdownDescription: "Node Pool Labels",
			},
			"linked_resources": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource ID\n  - example: res-12345678",
							MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Name\n  - example: my-resource",
							MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
						},
						"type": schema.StringAttribute{
							Required:            true,
							Description:         "Linked Resource Type (fs/obs)\n  - example: fs",
							MarkdownDescription: "Linked Resource Type (fs/obs)\n  - example: fs",
						},
					},
				},
				Optional: true,
			},
			"max_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Maximum node count (is_auto_scale = true)\n  - example: 5",
				MarkdownDescription: "Maximum node count (is_auto_scale = true)\n  - example: 5",
			},
			"min_node_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Minimum node count (is_auto_scale = true)\n  - example: 1",
				MarkdownDescription: "Minimum node count (is_auto_scale = true)\n  - example: 1",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Nodepool Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-nodepool",
				MarkdownDescription: "Nodepool Name\n  - maxLength: 20\n  - minLength: 3\n  - pattern: ^[a-z][a-z0-9\\-]*[a-z0-9]$\n  - example: sample-nodepool",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9\\-]*[a-z0-9]$"), ""),
				},
			},
			"nodepool": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
						MarkdownDescription: "Account ID\n  - example: 617b3d0e90c24a5fa1f65a3824861354",
					},
					"advanced_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"allowed_unsafe_sysctls": schema.StringAttribute{
								Computed:            true,
								Description:         "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
								MarkdownDescription: "Node Pool Allowed unsafe sysctls\n  - example: kernel.msg*,net.ipv4.route.min_pmtu",
								Default:             stringdefault.StaticString(""),
							},
							"container_log_max_files": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
								MarkdownDescription: "Node Pool container log max files\n  - maximum: 10\n  - minimum: 2\n  - example: 5",
							},
							"container_log_max_size": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
								MarkdownDescription: "Node Pool container log max size\n  - maximum: 100\n  - minimum: 10\n  - example: 10",
							},
							"image_gc_high_threshold": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
								MarkdownDescription: "Node Pool image GC high threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 85",
							},
							"image_gc_low_threshold": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
								MarkdownDescription: "Node Pool image GC low threshold percent\n  - maximum: 85\n  - minimum: 10\n  - example: 80",
							},
							"max_pods": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
								MarkdownDescription: "Node Pool max pod number\n  - maximum: 250\n  - minimum: 10\n  - example: 110",
							},
							"pod_max_pids": schema.Int32Attribute{
								Computed:            true,
								Description:         "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
								MarkdownDescription: "Node Pool Pod Max pids constraint\n  - maximum: 4.194304e+06\n  - minimum: 1024\n  - example: 4096",
							},
						},
						Computed:            true,
						Description:         "Node Pool Advanced Settings",
						MarkdownDescription: "Node Pool Advanced Settings",
					},
					"auto_recovery_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Recovery\n  - example: true",
						MarkdownDescription: "Is Auto Recovery\n  - example: true",
					},
					"auto_scale_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Auto Scale\n  - example: true",
						MarkdownDescription: "Is Auto Scale\n  - example: true",
					},
					"cluster": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Cluster ID\n  - example: 70a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Cluster",
						MarkdownDescription: "Cluster",
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
					"current_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Current Node Count\n  - example: 1",
						MarkdownDescription: "Current Node Count\n  - example: 1",
					},
					"desired_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Desired Node Count\n  - example: 2",
						MarkdownDescription: "Desired Node Count\n  - example: 2",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
						MarkdownDescription: "Nodepool ID\n  - example: bdfda539-bd2e-4a5c-9021-ec6d52d1ca79",
					},
					"image": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"custom_image_name": schema.StringAttribute{
								Computed:            true,
								Description:         "Custom Image Name\n  - example: custom-image",
								MarkdownDescription: "Custom Image Name\n  - example: custom-image",
							},
							"os": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS\n  - example: ubuntu",
								MarkdownDescription: "Image OS\n  - example: ubuntu",
							},
							"os_version": schema.StringAttribute{
								Computed:            true,
								Description:         "Image OS Version\n  - example: 22.04",
								MarkdownDescription: "Image OS Version\n  - example: 22.04",
							},
							"scp_gpu_driver": schema.StringAttribute{
								Computed: true,
							},
						},
						Computed:            true,
						Description:         "Image",
						MarkdownDescription: "Image",
					},
					"keypair": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Keypair Name\n  - example: test_keypair",
								MarkdownDescription: "Keypair Name\n  - example: test_keypair",
							},
						},
						Computed:            true,
						Description:         "Keypair Name",
						MarkdownDescription: "Keypair Name",
					},
					"kubernetes_version": schema.StringAttribute{
						Computed:            true,
						Description:         "Kubernetes Version\n  - example: v1.29.8",
						MarkdownDescription: "Kubernetes Version\n  - example: v1.29.8",
					},
					"labels": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Label Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Label Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									Default:             stringdefault.StaticString(""),
								},
							},
						},
						Computed:            true,
						Description:         "Node Pool Labels",
						MarkdownDescription: "Node Pool Labels",
					},
					"linked_resources": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource ID\n  - example: res-12345678",
									MarkdownDescription: "Linked Resource ID\n  - example: res-12345678",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Name\n  - example: my-resource",
									MarkdownDescription: "Linked Resource Name\n  - example: my-resource",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Linked Resource Type (fs/obs)\n  - example: fs",
									MarkdownDescription: "Linked Resource Type (fs/obs)\n  - example: fs",
								},
							},
						},
						Computed: true,
					},
					"max_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Max Node Count\n  - example: 5",
						MarkdownDescription: "Max Node Count\n  - example: 5",
					},
					"min_node_count": schema.Int32Attribute{
						Computed:            true,
						Description:         "Min Node Count\n  - example: 1",
						MarkdownDescription: "Min Node Count\n  - example: 1",
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
						Description:         "Nodepool Name\n  - example: sample-nodepool",
						MarkdownDescription: "Nodepool Name\n  - example: sample-nodepool",
					},
					"server_group_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
						MarkdownDescription: "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
					},
					"server_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"description": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type Description\n  - example: Standard",
								MarkdownDescription: "Server Type Description\n  - example: Standard",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
						},
						Computed:            true,
						Description:         "Server Type",
						MarkdownDescription: "Server Type",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						Description:         "Node Pool Status\n  - example: Running",
						MarkdownDescription: "Node Pool Status\n  - example: Running",
					},
					"taints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect": schema.StringAttribute{
									Computed:            true,
									Description:         "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
									MarkdownDescription: "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
								},
								"key": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
									MarkdownDescription: "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									Description:         "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									MarkdownDescription: "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
									Default:             stringdefault.StaticString(""),
								},
							},
						},
						Computed:            true,
						Description:         "Node Pool Taints",
						MarkdownDescription: "Node Pool Taints",
					},
					"volume_max_iops": schema.Int32Attribute{
						Computed: true,
					},
					"volume_max_throughput": schema.Int32Attribute{
						Computed: true,
					},
					"volume_size": schema.Int32Attribute{
						Computed:            true,
						Description:         "Volume Size\n  - example: 104",
						MarkdownDescription: "Volume Size\n  - example: 104",
					},
					"volume_type": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt": schema.BoolAttribute{
								Computed:            true,
								Description:         "Volume Type Encrypt\n  - example: true",
								MarkdownDescription: "Volume Type Encrypt\n  - example: true",
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
								MarkdownDescription: "Volume Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Description:         "Volume Type Name\n  - example: SSD",
								MarkdownDescription: "Volume Type Name\n  - example: SSD",
							},
						},
						Computed:            true,
						Description:         "Volume Type",
						MarkdownDescription: "Volume Type",
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
			"scp_gpu_driver": schema.StringAttribute{
				Optional: true,
			},
			"server_group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
				MarkdownDescription: "Server Group ID\n  - example: 2b8d33d5-4de5-40a5-a34c-7e30204133xc",
			},
			"server_type_id": schema.StringAttribute{
				Required:            true,
				Description:         "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
				MarkdownDescription: "Server Type ID\n  - example: 10a599e031e749b7b260868f441e862b",
			},
			"taints": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"effect": schema.StringAttribute{
							Required:            true,
							Description:         "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
							MarkdownDescription: "- enum: [\"NoSchedule\",\"NoExecute\",\"PreferNoSchedule\"]",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"NoSchedule",
									"NoExecute",
									"PreferNoSchedule",
								),
							},
						},
						"key": schema.StringAttribute{
							Required:            true,
							Description:         "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							MarkdownDescription: "Node Pool Taint Key\n  - pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$\n  - example: example.com/my-app",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$"), ""),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							MarkdownDescription: "Node Pool Taint Value\n  - maxLength: 63\n  - pattern: ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$\n  - example: bar",
							Validators: []validator.String{
								stringvalidator.LengthAtMost(63),
								stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
							},
							Default: stringdefault.StaticString(""),
						},
					},
				},
				Optional:            true,
				Description:         "Node Pool Taints",
				MarkdownDescription: "Node Pool Taints",
			},
			"volume_max_iops": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(5000, 20000),
				},
			},
			"volume_max_throughput": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(250, 1000),
				},
			},
			"volume_size": schema.Int32Attribute{
				Required:            true,
				Description:         "Volume Size\n  - example: 104",
				MarkdownDescription: "Volume Size\n  - example: 104",
			},
			"volume_type_name": schema.StringAttribute{
				Required:            true,
				Description:         "Volume Type Name\n  - example: SSD",
				MarkdownDescription: "Volume Type Name\n  - example: SSD",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *skeNodepoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *skeNodepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ske.NodepoolResource
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new nodepool
	data, err := r.client.CreateNodepool(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Creating Nodepool",
			"Could not create nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	plan.Id = types.StringValue(data.GetNodepool().Id)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	if plan.IsAutoScale.ValueBool() {
		plan.DesiredNodeCount = types.Int32Value(data.GetNodepool().MinNodeCount)
	} else {
		plan.MinNodeCount = types.Int32Value(data.GetNodepool().DesiredNodeCount)
		plan.MaxNodeCount = types.Int32Value(data.GetNodepool().DesiredNodeCount)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.GetNodepool().Id, []string{"ScalingUp"}, []string{"Running"}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Nodepool",
			"Error waiting for nodepool to become running: "+err.Error(),
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

func (r *skeNodepoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ske.NodepoolResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value from Nodepool
	data, _, err := r.client.GetNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := converter.NodepoolResponseToNodepoolModel(data)

	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)
	state.Nodepool = nodepoolObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *skeNodepoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state ske.NodepoolResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	if !state.KubernetesVersion.Equal(plan.KubernetesVersion) { // 1. 노드 풀 업그레이드
		/*
			업그레이드 버전 이외 설정 값이 변경 되어서는 안됨 (자동확장, 자동복구, 노드 수)
			자동 확장 사용시 최대/최소 노드 수 변경 불가
			자동 확장 미 사용시 노드 수 변경 불가
		*/
		if state.IsAutoRecovery.Equal(plan.IsAutoRecovery) && state.IsAutoScale.Equal(plan.IsAutoScale) {
			if plan.IsAutoScale.ValueBool() && state.MinNodeCount.Equal(plan.MinNodeCount) && state.MaxNodeCount.Equal(plan.MaxNodeCount) ||
				!plan.IsAutoScale.ValueBool() && state.DesiredNodeCount.Equal(plan.DesiredNodeCount) {
				_, err := r.client.UpgradeNodepool(ctx, plan)
				if err != nil {
					detail := client.GetDetailFromError(err)
					resp.Diagnostics.AddError(
						"Error Upgrade Nodepool",
						"Could not upgrade nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
					)
					return
				}

				err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"Updating"}, []string{"Running"}, true)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Updating Nodepool",
						"Error waiting for nodepool to become running: "+err.Error(),
					)
					return
				}
			} else {
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Version",
					"When nodepool version update, must not modify node count",
				)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Error Updating Nodepool Version",
				"When nodepool version update, must not modify auto recovery and auto scale",
			)
			return
		}
	} else { // label / taint 수정
		if !reflect.DeepEqual(plan.Labels, state.Labels) {
			_, err := r.client.UpdateNodepoolLabels(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Labels",
					"Could not update nodepool labels, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if !reflect.DeepEqual(plan.Taints, state.Taints) {
			_, err := r.client.UpdateNodepoolTaints(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Taints",
					"Could not update nodepool taints, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if !reflect.DeepEqual(plan.LinkedResources, state.LinkedResources) {
			_, err := r.client.UpdateNodepoolLinkedResources(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool Taints",
					"Could not update nodepool taints, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}
		}

		if (plan.IsAutoScale.ValueBool() && (!state.MinNodeCount.Equal(plan.MinNodeCount) || !state.MaxNodeCount.Equal(plan.MaxNodeCount))) ||
			(!plan.IsAutoScale.ValueBool() && !state.DesiredNodeCount.Equal(plan.DesiredNodeCount)) ||
			!state.IsAutoRecovery.Equal(plan.IsAutoRecovery) {
			_, err := r.client.UpdateNodepool(ctx, plan.Id.ValueString(), plan)
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Updating Nodepool",
					"Could not update nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			err = waitForNodepoolStatus(ctx, r.client, plan.Id.ValueString(), []string{"ScalingUp", "ScalingDown"}, []string{"Running"}, true)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Nodepool",
					"Error waiting for nodepool to become running: "+err.Error(),
				)
			}
		}
	}

	// Get refreshed value from Nodepool
	data, _, err := r.client.GetNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading Nodepool",
			"Could not read Nodepool ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	nodepoolModel := converter.NodepoolResponseToNodepoolModel(data)
	nodepoolObjectValue, diags := types.ObjectValueFrom(ctx, nodepoolModel.AttributeTypes(), nodepoolModel)

	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Nodepool = nodepoolObjectValue
	if plan.IsAutoScale.ValueBool() {
		plan.DesiredNodeCount = types.Int32Value(data.GetNodepool().MinNodeCount)
	} else {
		plan.MinNodeCount = types.Int32Value(data.GetNodepool().DesiredNodeCount)
		plan.MaxNodeCount = types.Int32Value(data.GetNodepool().DesiredNodeCount)
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

func (r *skeNodepoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ske.NodepoolResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Resource Group
	data, err := r.client.DeleteNodepool(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting Nodepool",
			"Could not delete nodepool, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForNodepoolStatus(ctx, r.client, data.GetResourceId(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting nodepool",
			"Error waiting for nodepool to become deleted: "+err.Error(),
		)
		return
	}
}

func waitForNodepoolStatus(ctx context.Context, skeClient *ske.Client, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, httpStatus, err := skeClient.GetNodepool(ctx, id)
		if httpStatus == 200 {
			return info, info.GetNodepool().Status, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETED", nil
		} else if httpStatus == 500 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("cluster with id=%s not found", id)
			}

			return info, "DELETING", nil
		} else if err != nil {
			return nil, "", err
		}

		return info, info.GetNodepool().Status, nil
	})
}

// ImportState adopts an existing resource via `terraform import <addr> <id>` using its
// opaque id; Read then refreshes the remaining state. (#81)
func (r *skeNodepoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
