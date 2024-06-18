package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ccpUserDataSource struct {
	Client *CleuraClient
}

func NewCCPUserDataSource() datasource.DataSource {
	return &ccpUserDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (c *ccpUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (c *ccpUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ccp_user"
}

// Read implements datasource.DataSource.
func (c *ccpUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var userData ccpUserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &userData)...)
	if userData.Name.ValueString() == "" {
		resp.Diagnostics.AddError("name required", "name property must be specified")
		return
	}
	result, err := c.Client.GetCCPUser(ctx, userData.Name.ValueString())
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
func (c *ccpUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a CCP user in Cleura Cloud",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Required: false,
			},
			"name": schema.StringAttribute{
				Computed: false,
				Required: true,
			},
			"first_name": schema.StringAttribute{
				Computed: true,
			},
			"last_name": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"pending_email": schema.StringAttribute{
				Computed: true,
			},
			"language": schema.StringAttribute{
				Computed: true,
			},
			// "two_factor_login": schema.StringAttribute{
			// 	Computed: true,
			// },
			"admin": schema.BoolAttribute{
				Computed: true,
			},
			"auth_provider_id": schema.StringAttribute{
				Computed: true,
			},
			"currency": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"id":   types.StringType,
					"code": types.StringType,
					"name": types.StringType,
				},
			},
			// "ip_restrictions": schema.SetAttribute{
			// 	ElementType: tftype.StringType,
			// 	Computed:    true,
			// },
			"privileges": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed: true,
							},
							"meta": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"openstack": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed: true,
							},
							"meta": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}
