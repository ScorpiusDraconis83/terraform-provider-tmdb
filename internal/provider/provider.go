package provider

import (
	"context"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &tmdbProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &tmdbProvider{
			version: version,
		}
	}
}

type tmdbProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *tmdbProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tmdb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *tmdbProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// tmdbProviderModel maps provider schema data to a Go type.
type tmdbProviderModel struct {
	APIKey types.String `tfsdk:"key"`
}

// Configure prepares an API client for data sources and resources.
func (p *tmdbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring TMDB client")

	// Retrieve provider data from configuration
	var config tmdbProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apiKey"),
			"Unknown TMDB API Key",
			"The provider cannot create the TMDB API client as there is an unknown configuration value for the TMDB API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TMDB_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apiKey := os.Getenv("TMDB_KEY")

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apiKey"),
			"Missing TMDB API Key",
			"The provider cannot create the TMDB API client as there is a missing or empty value for the TMDB API Key. "+
				"Set the API key value in the configuration or use the TMDB_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "tmdb_apikey", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "tmdb_apikey")

	tflog.Debug(ctx, "Creating TMDB client")

	// Create a new TMDB client using the configuration values
	client, err := tmdb.Init(apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create TMDB API Client",
			"An unexpected error occurred when creating the TMDB API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"TMDB Client Error: "+err.Error(),
		)
		return
	}

	// Make the TMDB client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured TMDB client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *tmdbProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPopularMoviesDataSource,
		NewMovieDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *tmdbProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
