// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ccpUserResource{}
var _ resource.ResourceWithImportState = &ccpUserResource{}

// ==============
// RESOURCE MODEL
// ==============
type ccpUserResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name" json:"name"`
	Email     types.String `tfsdk:"email" json:"email"`
	FirstName types.String `tfsdk:"first_name" json:"firstname"`
	LastName  types.String `tfsdk:"last_name" json:"lastname"`
	// Password  types.String `tfsdk:"password" json:"password"`
	// IpRestrictions []string `tfsdk:"ip_restrictions"`
	Privileges *ccpResourcePrivileges `tfsdk:"privileges" json:"privileges"`
}
type ccpResourcePrivileges struct {
	Users     ccpUserResourceUserPrivilege       `tfsdk:"users"`
	OpenStack ccpUserResourceOpenstackPrivileges `tfsdk:"openstack"`
}
type ccpUserResourceUserPrivilege struct {
	Type types.String `tfsdk:"type"`
	// Meta types.String `tfsdk:"meta"`
}
type ccpUserResourceOpenstackPrivileges struct {
	Type types.String `tfsdk:"type"`
	// Meta              types.String           `tfsdk:"meta"`
	// Uncomment when adding support for project privileges
	// ProjectPrivileges []ccpUserResourceProjectPrivileges `tfsdk:"project_privileges"`
}
type ccpUserResourceProjectPrivileges struct {
	ProjectId types.String `tfsdk:"project_id"`
	DomainId  types.String `tfsdk:"domain_id"`
	Type      types.String `tfsdk:"type"`
}

// ==============
// JSON MODEL
// ==============
type ccpUserCreateJson struct {
	User ccpUserResourceModelJson `json:"user"`
}
type ccpUserResourceModelJson struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Password  string `json:"password,omitempty"`
	// IpRestrictions []string `json:"ip_restrictions,omitempty"`
	Privileges *ccpResourcePrivilegesJson `json:"privileges"`
}
type ccpResourcePrivilegesJson struct {
	Users     ccpUserResourcePrivilegeJson   `json:"users"`
	OpenStack ccpUserOpenstackPrivilegesJson `json:"openstack"`
}
type ccpUserResourcePrivilegeJson struct {
	Type string `json:"type"`
	// Meta string `json:"meta"`
}
type ccpUserOpenstackPrivilegesJson struct {
	Type              string                                 `json:"type"`
	ProjectPrivileges []ccpUserResourceProjectPrivilegesJson `json:"project_privileges,omitempty"`
}

type ccpUserResourceProjectPrivilegesJson struct {
	ProjectId string `json:"project_id"`
	DomainId  string `json:"domain_id"`
	Type      string `json:"type"`
}
type ccpUserResourceCreateResponseJson struct {
}
type ccpUserCreate struct {
	User ccpUserResourceModel `json:"user"`
}
type ccpUserUpdate struct {
	User ccpUserResourceModelJson `json:"user"`
}

func NewCCPUserResource() resource.Resource {
	return &ccpUserResource{}
}

type ccpUserResource struct {
	Client *CleuraClient
}

func (c *ccpUserResource) GetJsonModel(obj ccpUserResourceModel) ccpUserResourceModelJson {
	result := ccpUserResourceModelJson{
		Name:      obj.Name.ValueString(),
		Email:     obj.Email.ValueString(),
		FirstName: obj.FirstName.ValueString(),
		LastName:  obj.LastName.ValueString(),
		Privileges: &ccpResourcePrivilegesJson{
			Users: ccpUserResourcePrivilegeJson{
				Type: obj.Privileges.Users.Type.ValueString(),
			},
			OpenStack: ccpUserOpenstackPrivilegesJson{
				Type: obj.Privileges.OpenStack.Type.ValueString(),
			},
		},
	}
	return result
}

func (c *ccpUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ccp_user"
}

func (c *ccpUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a CCP user in Cleura Cloud",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				// Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
			},
			"privileges": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"users": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required: true,
							},
							// "meta": schema.StringAttribute{
							// 	Optional: true,
							// },
						},
					},
					"openstack": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required: true,
							},
							// "meta": schema.StringAttribute{
							// 	Optional: true,
							// },
							// Uncomment when adding support for project privileges
							// "project_privileges": schema.ListNestedAttribute{
							// 	Optional: true,
							// 	NestedObject: schema.NestedAttributeObject{
							// 		Attributes: map[string]schema.Attribute{
							// 			"project_id": schema.StringAttribute{
							// 				Required: true,
							// 			},
							// 			"domain_id": schema.StringAttribute{
							// 				Required: true,
							// 			},
							// 			"type": schema.StringAttribute{
							// 				Required: true,
							// 			},
							// 		},
							// 	},
							// },
						},
					},
				},
			},
		},
	}
}

func (c *ccpUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (c *ccpUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ccpUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := c.Client.CreateCCPUser(ctx, plan)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to create user, error: %s", err.Error()))
		resp.Diagnostics.AddError("Failed to create user", fmt.Sprintf("error: %s", err.Error()))
		return
	}
	plan.Id = types.StringValue(result.Id.ValueString())
	tflog.Trace(ctx, "created user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *ccpUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ccpUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	exist, err := c.Client.DoesCCPUserExist(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to check if CCP user already exists", err.Error())
	}
	if !exist {
		// The user has been removed from outside Terraform, recreate it
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Cleura CCP User resource has been deleted outside terraform", "New resource will be created")
		return
	}
	userResponse, err := c.Client.GetCCPUserResource(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading user resource",
			"Could not read user resource named: "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("userResponse: %+v", userResponse))

	// Set refreshed state
	diags = resp.State.Set(ctx, &userResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ccpUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ccpUserResourceModel
	var currentState ccpUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateModel := c.GetJsonModel(plan)
	err := c.Client.UpdateCCPUser(ctx, ccpUserUpdate{User: updateModel})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update CCP user", err.Error())
		return
	}
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (c *ccpUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ccpUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := c.Client.DeleteCCPUser(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cleura CCP user",
			"Could not delete Cleura CCP user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (c *ccpUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
