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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	InstanceID types.Int64            `tfsdk:"instance_id"`
	User       types.String           `tfsdk:"username"`
	Rules      []aclRuleResourceModel `tfsdk:"rules"`
}

type aclRuleResourceModel struct {
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
			"rules": schema.ListNestedAttribute{
				Description: "List of items in the order.",
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Instance ID.",
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"operation": schema.StringAttribute{
							Description: "Which operation to set the rule on.",
							Required:    true,
							Validators: []validator.String{stringvalidator.OneOfCaseInsensitive("read", "write",
								"create", "delete", "alter", "describe", "clusteraction", "describeconfigs", "alterconfigs",
								"idempotentwrite", "createtokens", "describetokens", "all")},
						},
						"resource": schema.StringAttribute{
							Description: "Which resource to set the rule on, cluster, topic or group are valid values.",
							Required:    true,
							Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("cluster", "topic", "group")},
						},
						"resource_pattern": schema.StringAttribute{
							Description: "Which resource to match.",
							Required:    true,
						},
						"resource_pattern_type": schema.StringAttribute{
							Description: "How to apply the resource_pattern, literal or prefixed.",
							Required:    true,
							Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("literal", "prefixed")},
						},
					},
				},
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

func (r *aclResource) refresh(ctx context.Context, model *aclResourceModel, instanceId int64, user string) error {
	rules, err := r.client.ReadAclRuleForUser(instanceId, user)
	if err != nil {
		return err
	}

	modelRules := make([]aclRuleResourceModel, len(rules))
	for i, r := range rules {
		modelRules[i] = aclRuleResourceModel{
			ID:                  types.Int64Value(r.Id),
			Operation:           types.StringValue(r.Operation),
			Resource:            types.StringValue(r.Resource),
			ResourcePattern:     types.StringValue(r.ResourcePattern),
			ResourcePatternType: types.StringValue(r.ResourcePatternType),
		}
	}
	model.Rules = modelRules
	return nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	rules := make([]api.AclRule, len(plan.Rules))
	for i, r := range plan.Rules {
		rules[i] = api.AclRule{
			Operation:           r.Operation.ValueString(),
			Resource:            r.Resource.ValueString(),
			ResourcePattern:     r.ResourcePattern.ValueString(),
			ResourcePatternType: r.ResourcePatternType.ValueString(),
		}
	}
	err := r.client.CreateAclRules(plan.InstanceID.ValueInt64(), plan.User.ValueString(), rules)
	if err != nil {
		resp.Diagnostics.AddError("Error creating rules", err.Error())
		return
	}
	if err := r.refresh(ctx, &plan, plan.InstanceID.ValueInt64(), plan.User.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating rules", err.Error())
		return
	}

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

	if err := r.refresh(ctx, &state, state.InstanceID.ValueInt64(), state.User.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error refreshing rules", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan *aclResourceModel
	)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	/*
		err := r.client.UpdateTopic(plan.InstanceID.ValueInt64(), plan.Name.ValueString(), api.UpdateTopicRequest{
			Partitions: plan.Partitions.ValueInt64(),
			Config:     plan.Config.AsHash(),
		})

		if err != nil {
			resp.Diagnostics.AddError("Error updating rules on user", err.Error())
			return
		}
	*/
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
	for _, rule := range state.Rules {
		err := r.client.DeleteAclRule(state.InstanceID.ValueInt64(), rule.ID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError("Error deleting rule", err.Error())
		}
	}
	return
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
