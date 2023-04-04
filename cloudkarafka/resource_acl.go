package cloudkarafka

import (
	"context"
	"fmt"
	"terraform-provider-cloudkarafka/api"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
)

// NewAclResource is a helper function to simplify the provider implementation.
func NewAclResource() resource.Resource {
	return &aclResource{}
}

// aclResource is the resource implementation.
type aclResource struct {
	client *api.API
}

type aclResourceModel struct {
	InstanceID          types.Int64  `tfsdk:"instance_id"`
	User                types.String `tfsdk:"username"`
	ID                  types.Int64  `tfsdk:"id"`
	Operation           types.String `tfsdk:"operation"`
	Resource            types.String `tfsdk:"resource"`
	ResourcePattern     types.String `tfsdk:"resource_pattern"`
	ResourcePatternType types.String `tfsdk:"resource_pattern_type"`
}

// Metadata returns the data source type name.
func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aclrule"
}

// Schema defines the schema for the data source.
func (r *aclResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an ACL rule.",

		Attributes: map[string]schema.Attribute{
			"instance_id": schema.Int64Attribute{
				Description: "Id of the instance where we want to manage the rules.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Name of the user to apply the rules on.",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Rule ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"operation": schema.StringAttribute{
				Description: "Which operation to set the rule on.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive("read", "write",
					"create", "delete", "alter", "describe", "clusteraction", "describeconfigs", "alterconfigs",
					"idempotentwrite", "createtokens", "describetokens", "all")},
			},
			"resource": schema.StringAttribute{
				Description: "Which resource to set the rule on, cluster, topic or group are valid values.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive("cluster", "topic", "group")},
			},
			"resource_pattern": schema.StringAttribute{
				Description: "Which resource to match.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"resource_pattern_type": schema.StringAttribute{
				Description: "How to apply the resource_pattern, literal or prefixed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive("literal", "prefixed")},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *aclResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.API)
}

// Create creates the resource and sets the initial Terraform state.
func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule := api.AclRule{
		Operation:           plan.Operation.ValueString(),
		Resource:            plan.Resource.ValueString(),
		ResourcePattern:     plan.ResourcePattern.ValueString(),
		ResourcePatternType: plan.ResourcePatternType.ValueString(),
	}
	id, err := r.client.CreateAclRule(plan.InstanceID.ValueInt64(), plan.User.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error creating rules", err.Error())
		return
	}
	plan.ID = types.Int64Value(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("created diag failed"))
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aclResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.ReadAclRule(state.InstanceID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing rules", err.Error())
		return
	}
	state.Operation = types.StringValue(rule.Operation)
	state.Resource = types.StringValue(rule.Resource)
	state.ResourcePattern = types.StringValue(rule.ResourcePattern)
	state.ResourcePatternType = types.StringValue(rule.ResourcePatternType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *aclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No update

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aclResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAclRule(state.InstanceID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting rule", err.Error())
	}
	return
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
