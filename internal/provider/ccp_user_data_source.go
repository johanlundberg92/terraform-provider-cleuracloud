package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------- Types
type ccpUserDataSourceModel struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Privileges   *ccpPrivileges `tfsdk:"privileges"`
	Admin        types.Bool     `tfsdk:"admin"`
	FirstName    types.String   `tfsdk:"first_name"`
	LastName     types.String   `tfsdk:"last_name"`
	Email        types.String   `tfsdk:"email"`
	PendingEmail types.String   `tfsdk:"pending_email"`
	Language     types.String   `tfsdk:"language"`
	// TwoFactorLogin interface{}    `tfsdk:"two_factor_login"`
	// IPRestrictions interface{}    `tfsdk:"ip_restrictions"`
	Currency       *ccpCurrency `tfsdk:"currency"`
	AuthProviderId types.String `tfsdk:"auth_provider_id"`
}

type ccpCurrency struct {
	Id   types.String `tfsdk:"id"`
	Code types.String `tfsdk:"code"`
	Name types.String `tfsdk:"name"`
}

// Privileges represents the nested privileges object.
type ccpPrivileges struct {
	Users     ccpUsersPrivilege      `tfsdk:"users"`
	OpenStack ccpOpenstackPrivileges `tfsdk:"openstack"`
}

// InvoicePrivileges represents the invoice privileges.
type ccpUsersPrivilege struct {
	Type types.String `tfsdk:"type"`
	Meta types.String `tfsdk:"meta"`
}

type ccpOpenstackPrivileges struct {
	Type              types.String           `tfsdk:"type"`
	Meta              types.String           `tfsdk:"meta"`
	ProjectPrivileges []ccpProjectPrivileges `tfsdk:"project_privileges"`
}

type ccpProjectPrivileges struct {
	ProjectId types.String `tfsdk:"project_id"`
	DomainId  types.String `tfsdk:"domain_id"`
	Type      types.String `tfsdk:"type"`
}

// ------------------------- JSON
type ccpUserJson struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Privileges     ccpPrivilegesJson `json:"privileges"`
	Admin          bool              `json:"admin"`
	FirstName      string            `json:"firstname"`
	LastName       string            `json:"lastname"`
	Email          string            `json:"email"`
	PendingEmail   string            `json:"pending_email,omitempty"`
	Language       string            `json:"language,omitempty"`
	TwoFactorLogin []string          `json:"twofactorLogin,omitempty"`  // Use interface{} for nullable fields
	IPRestrictions []string          `json:"ip_restrictions,omitempty"` // Assuming IP restrictions are strings; adjust as needed
	Currency       ccpCurrencyJson   `json:"currency,omitempty"`
	AuthProviderId string            `json:"auth_provider_id"`
}
type ccpCurrencyJson struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
type ccpPrivilegesJson struct {
	Users     ccpUsersPrivilegeJson      `json:"users"`
	OpenStack ccpOpenstackPrivilegesJson `json:"openstack"`
}
type ccpUsersPrivilegeJson struct {
	Type string `json:"type"`
	Meta string `json:"meta"`
}
type ccpOpenstackPrivilegesJson struct {
	Type              string                     `json:"type"`
	Meta              string                     `json:"meta"`
	ProjectPrivileges []ccpProjectPrivilegesJson `json:"project_privileges,omitempty"`
}
type ccpProjectPrivilegesJson struct {
	ProjectId string `json:"project_id"`
	DomainId  string `json:"domain_id"`
	Type      string `json:"type"`
}

// --------

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
							"project_privileges": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"project_id": schema.StringAttribute{
											Computed: true,
										},
										"domain_id": schema.StringAttribute{
											Computed: true,
										},
										"type": schema.StringAttribute{
											Computed: true,
										},
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
