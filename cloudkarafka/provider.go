package cloudkarafka

import (
	"context"
	"os"
	"terraform-provider-cloudkarafka/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &cloudkarafkaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &cloudkarafkaProvider{}
}

// cloudkarafkaProvider is the provider implementation.
type cloudkarafkaProvider struct{}

// cloudkarafkaProviderModel maps provider schema data to a Go type.
type cloudkarafkaProviderModel struct {
	APIKey types.String `tfsdk:"apikey"`
}

// Metadata returns the provider type name.
func (p *cloudkarafkaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudkarafka"
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudkarafkaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Cloudkarafka.",
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Description: "API key Cloudkarafka API.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Cloudkarafka API client for data sources and resources.
func (p *cloudkarafkaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Cloudkarafka client")
	var config cloudkarafkaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown Cloudkarafka APIKey",
			"The provider cannot create the Cloudkarafka API client as there is an unknown configuration value for the Cloudkarafka APIKey.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("CLOUDKARAFKA_HOST")
	if host == "" {
		host = "https://customer.cloudkarafka.com"
	}
	apikey := ""
	if !config.APIKey.IsNull() {
		apikey = config.APIKey.ValueString()
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing Cloudkarafka API key",
			"The provider cannot create the Cloudkarafka API client as there is a missing or empty value for the Cloudkarafka API key. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	client := api.New(host, apikey)
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *cloudkarafkaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *cloudkarafkaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInstanceResource,
		NewTopicResource,
		NewUserResource,
		NewAclResource,
	}
}
