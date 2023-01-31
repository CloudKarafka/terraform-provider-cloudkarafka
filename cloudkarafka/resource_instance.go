package cloudkarafka

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-cloudkarafka/api"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &instanceResource{}
	_ resource.ResourceWithConfigure = &instanceResource{}
	// _ resource.ResourceWithImportState = &instanceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

// instanceResource is the resource implementation.
type instanceResource struct {
	client *api.API
}

type instanceResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Plan         types.String `tfsdk:"plan"`
	Region       types.String `tfsdk:"region"`
	KafkaVersion types.String `tfsdk:"kafka_version"`
	VPCSubnet    types.String `tfsdk:"vpc_subnet"`
	VPCId        types.Int64  `tfsdk:"vpc_id"`
}

// Metadata returns the data source type name.
func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema defines the schema for the data source.
func (r *instanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Instance ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of instance.",
				Required:    true,
			},
			"plan": schema.StringAttribute{
				Description: "What plan to use.",
				Required:    true,
			},
			"region": schema.StringAttribute{
				Description: "Which region to use.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(amazon-web-services|azure-arm|google-compute-engine)::[a-z0-9\-]+$`),
						"must be a valid region identifier",
					),
				},
			},
			"kafka_version": schema.StringAttribute{
				Description: "Which Apache Kafka version to use.",
				Required:    true,
			},
			"vpc_subnet": schema.StringAttribute{
				Description: "Subnet for the VPC.",
				Optional:    true,
			},
			"vpc_id": schema.Int64Attribute{
				Description: "ID for which subnet to use.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.API)
}

// Create creates the resource and sets the initial Terraform state.
func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := api.CreateInstanceRequest{
		Name:         plan.Name.ValueString(),
		Plan:         plan.Plan.ValueString(),
		Region:       plan.Region.ValueString(),
		KafkaVersion: plan.KafkaVersion.ValueString(),
	}
	if !plan.VPCId.IsNull() {
		createRequest.VpcId = plan.VPCId.ValueInt64()
	} else {
		if !plan.VPCSubnet.IsNull() {
			createRequest.VpcSubnet = plan.VPCSubnet.ValueString()
		}
	}
	instance, err := r.client.CreateInstance(createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating instance", err.Error())
		return
	}

	plan.ID = types.Int64Value(instance.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("created diag failed"))
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.ReadInstance(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read instance state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan *instanceResourceModel
	)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateInstance(plan.ID.ValueInt64(), api.UpdateInstanceRequest{
		Name: plan.Name.ValueString(),
		Plan: plan.Plan.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	return
}

func (r *instanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
