package cloudkarafka

import (
	"context"
	"fmt"
	"terraform-provider-cloudkarafka/api"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &configResource{}
	_ resource.ResourceWithConfigure   = &configResource{}
	_ resource.ResourceWithImportState = &configResource{}
)

// NewConfigrResource is a helper function to simplify the provider implementation.
func NewConfigResource() resource.Resource {
	return &configResource{}
}

// configResource is the resource implementation.
type configResource struct {
	client *api.API
}

type configResourceModel struct {
	InstanceID        types.Int64 `tfsdk:"instance_id"`
	AutoCreateTopics  types.Bool  `tfsdk:"auto_create_topics_enable"`
	MinInsyncReplicas types.Int64 `tfsdk:"min_insync_replicas"`
	LogRetentionBytes types.Int64 `tfsdk:"log_retention_bytes"`
	LogRetentionMs    types.Int64 `tfsdk:"log_retention_ms"`
	LogSegmentBytes   types.Int64 `tfsdk:"log_segment_bytes"`
	NetworkThreads    types.Int64 `tfsdk:"num_network_threads"`
	IOThreads         types.Int64 `tfsdk:"num_io_threads"`
	MessageMaxBytes   types.Int64 `tfsdk:"message_max_bytes"`
}

// Metadata returns the data source type name.
func (r *configResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafkaconfig"
}

// Schema defines the schema for the data source.
func (r *configResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the Kafka configuration.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.Int64Attribute{
				Description: "Id of the instance where we want to manage the topic.",
				Required:    true,
			},
			"auto_create_topics_enable": schema.BoolAttribute{
				Description: "Enable auto creation of topic on the server.",
				Optional:    true,
			},
			"min_insync_replicas": schema.Int64Attribute{
				Description: "Minimum insync replicas avaiable with ACKing.",
				Optional:    true,
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"log_retention_bytes": schema.Int64Attribute{
				Description: "The maximum size of the log before deleting it.",
				Optional:    true,
			},
			"log_retention_ms": schema.Int64Attribute{
				Description: "The number of milliseconds to keep a log file before deleting it.",
				Optional:    true,
			},
			"log_segment_bytes": schema.Int64Attribute{
				Description: "The maximum size of a single log file.",
				Optional:    true,
			},
			"num_io_threads": schema.Int64Attribute{
				Description: "The number of threads that the server uses for processing requests, which may include disk I/O.",
				Optional:    true,
			},
			"num_network_threads": schema.Int64Attribute{
				Description: "The number of threads that the server uses for receiving requests from the network and sending responses to the network.",
				Optional:    true,
			},
			"message_max_bytes": schema.Int64Attribute{
				Description: "Max size of message.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *configResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.API)
}

// Create creates the resource and sets the initial Terraform state.
func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan configResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg := api.NewKafkaConfig()
	if !plan.AutoCreateTopics.IsNull() {
		cfg.AutoCreateTopics = plan.AutoCreateTopics.ValueBool()
	}
	if !plan.MinInsyncReplicas.IsNull() {
		cfg.MinInsyncReplicas = plan.MinInsyncReplicas.ValueInt64()
	}
	if !plan.LogRetentionBytes.IsNull() {
		cfg.LogRetentionBytes = plan.LogRetentionBytes.ValueInt64()
	}
	if !plan.LogRetentionMs.IsNull() {
		cfg.LogRetentionMs = plan.LogRetentionBytes.ValueInt64()
	}
	if !plan.LogSegmentBytes.IsNull() {
		cfg.LogSegmentBytes = plan.LogSegmentBytes.ValueInt64()
	}
	if !plan.NetworkThreads.IsNull() {
		cfg.NetworkThreads = plan.NetworkThreads.ValueInt64()
	}
	if !plan.IOThreads.IsNull() {
		cfg.IOThreads = plan.IOThreads.ValueInt64()
	}

	err := r.client.WriteConfig(plan.InstanceID.ValueInt64(), cfg)
	if err != nil {
		resp.Diagnostics.AddError("Error updating kafka config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("created diag failed"))
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.ReadConfig(state.InstanceID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read kafka config", err.Error())
		return
	}
	state.AutoCreateTopics = types.BoolValue(config.AutoCreateTopics)
	if config.MinInsyncReplicas != -1 {
		state.MinInsyncReplicas = types.Int64Value(config.MinInsyncReplicas)
	}
	if config.IOThreads != -1 {
		state.IOThreads = types.Int64Value(config.IOThreads)
	}
	if config.NetworkThreads != -1 {
		state.NetworkThreads = types.Int64Value(config.NetworkThreads)
	}
	if config.LogRetentionBytes != -1 {
		state.LogRetentionBytes = types.Int64Value(config.LogRetentionBytes)
	}
	if config.LogRetentionMs != -1 {
		state.LogRetentionMs = types.Int64Value(config.LogRetentionMs)
	}
	if config.LogSegmentBytes != -1 {
		state.LogSegmentBytes = types.Int64Value(config.LogSegmentBytes)
	}
	if config.MessageMaxBytes != -1 {
		state.MessageMaxBytes = types.Int64Value(config.MessageMaxBytes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *configResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *configResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
