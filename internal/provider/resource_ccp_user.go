// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"

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

func NewCCPUserResource() resource.Resource {
	return &cleuraUserResource{}
}

type ccpUserResource struct {
	Client *CleuraClient
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

	result, err := c.Client.CreateCCPUser(plan)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to create user, error: %s", err.Error()))
		resp.Diagnostics.AddError("Failed to create user", fmt.Sprintf("error: %s", err.Error()))
		return
	}
	plan.Id = types.StringValue(result.Id)
	tflog.Trace(ctx, "created user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *ccpUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state openstackUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	exist, err := c.Client.DoesUserExist(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to check if user already exists", err.Error())
	}
	if !exist {
		// The user has been removed from outside Terraform, recreate it
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Cleura User resource has been deleted outside terraform", "New resource will be created")
		return
	}
	userResponse, err := c.Client.GetUserResource(state.Id.ValueString())
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
	var plan openstackUserResourceModel
	var currentState openstackUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if currentState.Enabled != plan.Enabled {
		err := c.Client.ToggleUserEnabled(currentState.Id.ValueString(), plan.Enabled.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Failed to update user", err.Error())
			return
		}
	}
	currentProjects := make([]openstackUserCreateProject, 0)
	currentProjects = append(currentProjects, currentState.Projects...)
	for _, p := range plan.Projects {
		// Check if the planned project is in the current state
		for _, st := range currentProjects {
			if st.Id == p.Id {
				// This planned project is in the current state
				// Compare what roles are to be ADDED
				for _, r := range p.Roles {
					if !slices.Contains(st.Roles, r) {
						// The planned role is not in the current state, therefore we must add it
						err := c.Client.AddUserToProjectRole(currentState.Id.ValueString(), p.Id, r)
						if err != nil {
							resp.Diagnostics.AddError("Failed to add user to project role", err.Error())
							return
						}
					}
				}
				// Compare what roles are to be DELETED
				for _, r := range st.Roles {
					if !slices.Contains(p.Roles, r) {
						// A role in the state has been removed from the plan, therefore we must REMOVE it
						err := c.Client.RemoveUserFromProjectRole(currentState.Id.ValueString(), p.Id, r)
						if err != nil {
							resp.Diagnostics.AddError("Failed to remove user from project role", err.Error())
							return
						}
					}
				}
			} else {
				// This planned project is NOT in the current state
				// We must add the user to it along with the specified roles
				projAssign := openstackProjectAssignment{ProjectId: p.Id, Roles: p.Roles}
				projSlice := []openstackProjectAssignment{projAssign}
				proj := openstackProjectUpdate{Projects: projSlice}
				c.Client.AddUserToProject(proj)
			}
		}
	}
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (c *ccpUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state openstackUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := c.Client.DeleteUser(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cleura user",
			"Could not delete Cleura user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (c *ccpUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
