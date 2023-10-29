package provider

import (
	"context"
	"fmt"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &movieDataSource{}
	_ datasource.DataSourceWithConfigure = &movieDataSource{}
)

func NewMovieDataSource() datasource.DataSource {
	return &movieDataSource{}
}

type movieDataSource struct {
	client *tmdb.Client
}

// Metadata returns the data source type name.
func (d *movieDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_movie"
}

// Schema defines the schema for the data source.
func (d *movieDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Computed: true,
			},
			"overview": schema.StringAttribute{
				Computed: true,
			},
			"releasedate": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *movieDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state movieDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	movie, err := d.client.GetMovieDetails(int(state.ID.ValueInt64()), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read TMDB Movie with ID: %v", state.ID.ValueInt64()),
			err.Error(),
		)
		return
	}

	state.ID = types.Int64Value(movie.ID)
	state.Title = types.StringValue(movie.Title)
	state.Overview = types.StringValue(movie.Overview)
	state.ReleaseDate = types.StringValue(movie.ReleaseDate)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *movieDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tmdb.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

type movieDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Overview    types.String `tfsdk:"overview"`
	ReleaseDate types.String `tfsdk:"releasedate"`
}
