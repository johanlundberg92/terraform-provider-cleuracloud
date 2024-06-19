package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type cleuraUserDataSource struct {
	Client *CleuraClient
}

func NewOpenstackUserDataSource() datasource.DataSource {
	return &cleuraUserDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (c *cleuraUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*CleuraClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unable to cast ProviderData to *CleuraClient",
			fmt.Sprintf("Expected *CleuraClient, got: %T", req.ProviderData),
		)
		return
	}
	c.Client = client
}

// Metadata implements datasource.DataSource.
func (c *cleuraUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openstack_user"
}

// Read implements datasource.DataSource.
func (c *cleuraUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var userData openstackUserDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &userData)...)
	result, err := c.Client.GetUser(ctx, userData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get user",
			err.Error(),
		)
		return
	}
	diags := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSource.
func (c *cleuraUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a user in Cleura Cloud",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: false,
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"domain_id": schema.StringAttribute{
				Computed: true,
			},
			"default_project_id": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"domain_id": schema.StringAttribute{
							Computed: true,
						},
						"roles": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

}
