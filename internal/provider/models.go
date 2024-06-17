package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type apiError struct {
	Error apiErrorDetails `json:"error"`
}
type apiErrorDetails struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// ++++++++++++++++++++ User
type openstackUserDatasourceModel struct {
	Id               types.String       `json:"id" tfsdk:"id"`
	Name             types.String       `json:"name" tfsdk:"name"`
	DomainId         types.String       `json:"domain_id,omitempty" tfsdk:"domain_id"`
	DefaultProjectId types.String       `json:"default_project_id,omitempty" tfsdk:"default_project_id"`
	Enabled          types.Bool         `json:"enabled" tfsdk:"enabled"`
	Description      types.String       `json:"description,omitempty" tfsdk:"description"`
	Projects         []openstackProject `json:"projects,omitempty" tfsdk:"projects"`
	// Client           *CleuraClient
}
type openstackProject struct {
	Id       types.String    `json:"id" tfsdk:"id"`
	Name     types.String    `json:"name" tfsdk:"name"`
	DomainId types.String    `json:"domain_id" tfsdk:"domain_id"`
	Roles    []openstackRole `json:"roles" tfsdk:"roles"`
}
type openstackRole struct {
	Id   types.String `json:"id" tfsdk:"id"`
	Name types.String `json:"name" tfsdk:"name"`
}

// JSON UNMARSHALING
type openstackUserDatasourceModelJson struct {
	Id               string                 `json:"id" tfsdk:"id"`
	Name             string                 `json:"name" tfsdk:"name"`
	DomainId         string                 `json:"domain_id,omitempty" tfsdk:"domain_id"`
	DefaultProjectId string                 `json:"default_project_id,omitempty" tfsdk:"default_project_id"`
	Enabled          bool                   `json:"enabled" tfsdk:"enabled"`
	Description      string                 `json:"description,omitempty" tfsdk:"description"`
	Projects         []openstackProjectJson `json:"projects,omitempty" tfsdk:"projects"`
	// Client *CleuraClient

}

type openstackProjectJson struct {
	Id       string              `json:"id" tfsdk:"id"`
	Name     string              `json:"name" tfsdk:"name"`
	DomainId string              `json:"domain_id" tfsdk:"domain_id"`
	Roles    []openstackRoleJson `json:"roles" tfsdk:"roles"`
}

type openstackRoleJson struct {
	Id   string `json:"id" tfsdk:"id"`
	Name string `json:"name" tfsdk:"name"`
}

type openstackUserResourceModel struct {
	Id               types.String                 `json:"id" tfsdk:"id"`
	Name             types.String                 `json:"name" tfsdk:"name"`
	DomainId         types.String                 `json:"domain_id,omitempty" tfsdk:"domain_id"`
	DefaultProjectId types.String                 `json:"default_project_id,omitempty" tfsdk:"default_project_id"`
	Enabled          types.Bool                   `json:"enabled" tfsdk:"enabled"`
	Description      types.String                 `json:"description,omitempty" tfsdk:"description"`
	Projects         []openstackUserCreateProject `json:"projects,omitempty" tfsdk:"projects"`
	// Client           *CleuraClient
}

type createOpenstackUser struct {
	User     createOpenstackUserInfo      `json:"user"`
	Projects []openstackUserCreateProject `json:"projects,omitempty"`
}

type createOpenstackUserInfo struct {
	Name             string `json:"name"`
	Password         string `json:"password"`
	Description      string `json:"description,omitempty"`
	DefaultProjectId string `json:"default_project_id,omitempty" tfsdk:"default_project_id"`
}

// // Projects represents the projects array within the JSON structure.
//
//	type openstackUserCreateProjectList struct {
//		Type        string                   `json:"type"`
//		Required    bool                     `json:"required"`
//		Example     string                   `json:"example"`
//		ProjectList []map[string]interface{} `json:"projectList"` // Using map[string]interface{} for dynamic keys
//	}
type openstackUserCreateProject struct {
	Id    string   `json:"project_id" tfsdk:"id"`
	Roles []string `json:"roles" tfsdk:"roles"`
}

// type openstackUserCreateProjectRole struct {
// }
type openstackUserCreatedModel struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	DomainId         string `json:"domain_id"`
	DefaultProjectId string `json:"default_project_id"`
	Enabled          bool   `json:"enabled"`
	Description      bool   `json:"description"`
}

type openstackUserUpdate struct {
	User openstackUserUpdateProperties `json:"user"`
}
type openstackUserUpdateProperties struct {
	Enabled bool `json:"enabled"`
}
type openstackProjectUpdate struct {
	Projects []openstackProjectAssignment `json:"projects"`
}
type openstackProjectAssignment struct {
	ProjectId string   `json:"project_id"`
	Roles     []string `json:"roles"`
}

// --------------------  User
