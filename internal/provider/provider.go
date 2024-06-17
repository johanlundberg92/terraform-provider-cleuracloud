package provider

import (
	"context"
	"net/http"
	"os"

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
	_ provider.Provider = &cleuraProvider{}
	// _ provider.ProviderWithFunctions = &cleuraProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cleuraProvider{
			version: version,
		}
	}
}

// cleuraProviderModel maps provider schema data to a Go type.
type cleuraProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Url      types.String `tfsdk:"api_url"`
	DomainId types.String `tfsdk:"domain_id"`
}

type cleuraProvider struct {
	version string
}

// Metadata returns the provider type name.
func (p *cleuraProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cleura"
	resp.Version = p.version
}

func (p *cleuraProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configure cleura provider was called")
	tflog.Info(ctx, "Creating Cleura client")

	// Retrieve provider data from configuration
	var config cleuraProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Cleura Username",
			"The provider cannot create the cleura API client as there is an unknown configuration value for the cleura API host. ")
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Cleura Password",
			"The provider cannot create the cleura API client as there is an unknown configuration value for the cleura API username. ")
	}

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown Cleura Url",
			"The provider cannot create the cleura API client as there is an unknown configuration value for the cleura API password. ")
	}
	if config.DomainId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("domain_id"),
			"Unknown Cleura DomainId",
			"The provider cannot create the cleura API client as there is an unknown configuration value for the cleura API password. ")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	username := os.Getenv("CLEURA_USER")
	password := os.Getenv("CLEURA_PW")
	api_url := os.Getenv("CLEURA_URL")
	domain_id := os.Getenv("CLEURA_DOMAIN_ID")

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Url.IsNull() {
		api_url = config.Url.ValueString()
	}

	if !config.DomainId.IsNull() {
		domain_id = config.DomainId.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing cleura API username",
			"The provider cannot create the cleura API client as there is a missing or empty value for the cleura API username. ")
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing cleura API password",
			"The provider cannot create the cleura API client as there is a missing or empty value for the cleura API username. ")
	}

	if api_url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing cleura API url",
			"The provider cannot create the cleura API client as there is a missing or empty value for the cleura API password. ")
	}

	if domain_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("domain_id"),
			"Missing cleura API domain_id",
			"The provider cannot create the cleura API client as there is a missing or empty value for the cleura API password. ")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cleura_user", username)
	ctx = tflog.SetField(ctx, "cleura_password", password)
	ctx = tflog.SetField(ctx, "cleura url", api_url)
	ctx = tflog.SetField(ctx, "cleura domain_id", domain_id)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cleura_password")

	tflog.Debug(ctx, "Creating Cleura client")

	client := &CleuraClient{}
	client.User = username
	client.Password = password
	client.Url = api_url
	client.DomainId = domain_id
	client.Client = &http.Client{}
	err := client.Login()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to login to Cleura cloud",
			"An unexpected error occurred when creating the CleuraClient. "+
				"Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured CleuraClient", map[string]any{"success": true})
}

// Schema defines the provider-level schema for configuration data.
func (p *cleuraProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Cleura.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "Username for Cleura API. May also be provided via CLEURA_USER environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Cleura API. May also be provided via CLEURA_PW environment variable.",
				Optional:    true,
			},
			"api_url": schema.StringAttribute{
				Description: "Url for Cleura API. May also be provided via CLEURA_URL environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"domain_id": schema.StringAttribute{
				Description: "DomainId for Cleura API. May also be provided via CLEURA_DOMAIN_ID environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *cleuraProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOpenstackUserDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *cleuraProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOpenstackUserResource,
	}
}

// func (p *cleuraProvider) Functions(_ context.Context) []func() function.Function {
// 	return nil
// 	// return []func() function.Function{
// 	// 	NewComputeTaxFunction,
// 	// }
// }
