package cloudkarafka

import (
	"context"
	"fmt"
	"terraform-provider-cloudkarafka/api"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &topicResource{}
	_ resource.ResourceWithConfigure   = &topicResource{}
	_ resource.ResourceWithImportState = &topicResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewTopicResource() resource.Resource {
	return &topicResource{}
}

// topicResource is the resource implementation.
type topicResource struct {
	client *api.API
}

type topicResourceModel struct {
	InstanceID        types.Int64              `tfsdk:"instance_id"`
	Name              types.String             `tfsdk:"name"`
	Partitions        types.Int64              `tfsdk:"partitions"`
	ReplicationFactor types.Int64              `tfsdk:"replication_factor"`
	Config            topicConfigResourceModel `tfsdk:"config"`
}

type topicConfigResourceModel struct {
	CleanupPolicy     types.String `tfsdk:"cleanup_policy"`
	MinInsyncReplicas types.Int64  `tfsdk:"min_insync_replicas"`
	RetentionBytes    types.Int64  `tfsdk:"retention_bytes"`
	RetentionMs       types.Int64  `tfsdk:"retention_ms"`
	DeleteRetentionMs types.Int64  `tfsdk:"delete_retention_ms"`
	SegmentBytes      types.Int64  `tfsdk:"segment_bytes"`
}

func (me topicConfigResourceModel) AsHash() map[string]interface{} {
	config := make(api.Hash)
	if !me.CleanupPolicy.IsNull() {
		config["cleanup.policy"] = me.CleanupPolicy.ValueString()
	}
	if !me.MinInsyncReplicas.IsNull() {
		config["min.insync.replicas"] = me.MinInsyncReplicas.ValueInt64()
	}
	if !me.RetentionBytes.IsNull() {
		config["retention.bytes"] = me.RetentionBytes.ValueInt64()
	}
	if !me.RetentionMs.IsNull() {
		config["retention.ms"] = me.RetentionMs.ValueInt64()
	}
	if !me.DeleteRetentionMs.IsNull() {
		config["delete.retention.ms"] = me.DeleteRetentionMs.ValueInt64()
	}
	if !me.SegmentBytes.IsNull() {
		config["segment.bytes"] = me.SegmentBytes.ValueInt64()
	}
	return config

}

// Metadata returns the data source type name.
func (r *topicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_topic"
}

// Schema defines the schema for the data source.
func (r *topicResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a topic.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.Int64Attribute{
				Description: "Id of the instance where we want to manage the topic.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of topic.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"partitions": schema.Int64Attribute{
				Description: "Number of partitions for the topic.",
				Required:    true,
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"replication_factor": schema.Int64Attribute{
				Description: "Replication factor for the topic.",
				Required:    true,
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Topic configuration.",
				Attributes: map[string]schema.Attribute{
					"cleanup_policy": schema.StringAttribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("delete", "compact")},
					},
					"delete_retention_ms": schema.Int64Attribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
					},
					"min_insync_replicas": schema.Int64Attribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
						Validators:  []validator.Int64{int64validator.AtLeast(1)},
					},
					"retention_bytes": schema.Int64Attribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
					},
					"retention_ms": schema.Int64Attribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
					},
					"segment_bytes": schema.Int64Attribute{
						Description: "Delete or compact when records hit their retention.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *topicResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.API)
}

// Create creates the resource and sets the initial Terraform state.
func (r *topicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan topicResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createRequest := api.Topic{
		Name:       plan.Name.ValueString(),
		Partitions: plan.Partitions.ValueInt64(),
		Replicas:   plan.ReplicationFactor.ValueInt64(),
		Config:     plan.Config.AsHash(),
	}
	err := r.client.CreateTopic(plan.InstanceID.ValueInt64(), createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating topic", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("created diag failed"))
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *topicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state topicResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	topic, err := r.client.ReadTopic(state.InstanceID.ValueInt64(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read topic state", err.Error())
		return
	}
	state.Partitions = types.Int64Value(topic.Partitions)
	state.ReplicationFactor = types.Int64Value(topic.Replicas)
	if v, ok := topic.Config["cleanup.policy"]; ok {
		state.Config.CleanupPolicy = types.StringValue(v.(string))
	}
	if v, ok := topic.Config["retention.ms"]; ok {
		state.Config.RetentionMs = types.Int64Value(int64(v.(float64)))
	}
	if v, ok := topic.Config["retention.bytes"]; ok {
		state.Config.RetentionBytes = types.Int64Value(int64(v.(float64)))
	}
	if v, ok := topic.Config["segment.bytes"]; ok {
		state.Config.SegmentBytes = types.Int64Value(int64(v.(float64)))
	}
	if v, ok := topic.Config["min.insync.replicas"]; ok {
		state.Config.MinInsyncReplicas = types.Int64Value(int64(v.(float64)))
	}
	if v, ok := topic.Config["delete.retention.ms"]; ok {
		state.Config.DeleteRetentionMs = types.Int64Value(int64(v.(float64)))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *topicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan *topicResourceModel
	)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.UpdateTopic(plan.InstanceID.ValueInt64(), plan.Name.ValueString(), api.UpdateTopicRequest{
		Partitions: plan.Partitions.ValueInt64(),
		Config:     plan.Config.AsHash(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating topic", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *topicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state topicResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTopic(state.InstanceID.ValueInt64(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting topic", err.Error())
		return
	}

	return
}

func (r *topicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
