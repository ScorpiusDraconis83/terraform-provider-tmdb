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
	_ datasource.DataSource              = &searchDataSource{}
	_ datasource.DataSourceWithConfigure = &searchDataSource{}
)

func NewSearchDataSource() datasource.DataSource {
	return &searchDataSource{}
}

type searchDataSource struct {
	client *tmdb.Client
}

// Metadata returns the data source type name.
func (d *searchDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search"
}

// Schema defines the schema for the data source.
func (d *searchDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Required: true,
			},
			"movies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *searchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state searchDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	movies, err := d.client.GetSearchMovies(state.Query, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read TMDB Movies with query: %s", state.Query),
			err.Error(),
		)
		return
	}

	// Map response body to model.
	for _, movie := range movies.Results {
		movieState := moviesModel{
			ID:          types.Int64Value(movie.ID),
			Title:       types.StringValue(movie.Title),
			Overview:    types.StringValue(movie.Overview),
			ReleaseDate: types.StringValue(movie.ReleaseDate),
		}

		state.Movies = append(state.Movies, movieState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *searchDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type searchDataSourceModel struct {
	Movies []moviesModel `tfsdk:"movies"`
	Query  string        `tfsdk:"query"`
}
