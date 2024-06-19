package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ==============
// RESOURCE MODEL
// ==============
type ipsecDataSourceModel struct {
	PsK            types.String `tfsdk:"psk"`
	Initiator      types.String `tfsdk:"initiator"`
	IpsecpolicyID  types.String `tfsdk:"ipsecpolicy_id"`
	AdminStateUp   types.Bool   `tfsdk:"admin_state_up"`
	Mtu            types.Int64  `tfsdk:"mtu"`
	PeerEpGroupID  types.String `tfsdk:"peer_ep_group_id"`
	IkepolicyID    types.String `tfsdk:"ikepolicy_id"`
	VpnserviceID   types.String `tfsdk:"vpnservice_id"`
	LocalEpGroupID types.String `tfsdk:"local_ep_group_id"`
	PeerAddress    types.String `tfsdk:"peer_address"`
	PeerID         types.String `tfsdk:"peer_id"`
	Name           types.String `tfsdk:"name"`
}

// ==============
// JSON
// ==============
type ipsecDataSourceModelJson struct {
	PsK            string `json:"psk"`
	Initiator      string `json:"initiator"`
	IpsecpolicyID  string `json:"ipsecpolicy_id"`
	AdminStateUp   bool   `json:"admin_state_up"`
	Mtu            int    `json:"mtu"`
	PeerEpGroupID  string `json:"peer_ep_group_id"`
	IkepolicyID    string `json:"ikepolicy_id"`
	VpnserviceID   string `json:"vpnservice_id"`
	LocalEpGroupID string `json:"local_ep_group_id"`
	PeerAddress    string `json:"peer_address"`
	PeerID         string `json:"peer_id"`
	Name           string `json:"name"`
}

// --------

type ipsecConnectionDataSource struct {
	Client *CleuraClient
}

func NewIpsecConnectionDataSource() datasource.DataSource {
	return &ipsecConnectionDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (c *ipsecConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (c *ipsecConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_connection"
}

// Read implements datasource.DataSource.
func (c *ipsecConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
func (c *ipsecConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for managing ipsec site connections",
		Attributes: map[string]schema.Attribute{
			"psk": schema.StringAttribute{
				Description: "Pre-Shared Key of the IPSec Site Connection.",
				Computed:    true,
			},
			"initiator": schema.StringAttribute{
				Description: "Initiator of the IPSec Site Connection.",
				Computed:    true,
			},
			"ipsecpolicy_id": schema.StringAttribute{
				Description: "ID of the IPSEC Policy associated with the connection.",
				Computed:    true,
			},
			"admin_state_up": schema.BoolAttribute{
				Description: "Admin state up of the IPSec Site Connection.",
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "Maximum Transmission Unit (MTU) for the IPSec Site Connection.",
				Computed:    true,
			},
			"peer_ep_group_id": schema.StringAttribute{
				Description: "ID of the Peer Endpoint Group associated with the connection.",
				Computed:    true,
			},
			"ikepolicy_id": schema.StringAttribute{
				Description: "ID of the IKE Policy associated with the connection.",
				Computed:    true,
			},
			"vpnservice_id": schema.StringAttribute{
				Description: "ID of the VPN Service associated with the connection.",
				Computed:    true,
			},
			"local_ep_group_id": schema.StringAttribute{
				Description: "ID of the Local Endpoint Group associated with the connection.",
				Computed:    true,
			},
			"peer_address": schema.StringAttribute{
				Description: "IP address of the peer endpoint.",
				Computed:    true,
			},
			"peer_id": schema.StringAttribute{
				Description: "ID of the peer endpoint.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the IPSec Site Connection.",
				Required:    true,
			},
		},
	}
}
