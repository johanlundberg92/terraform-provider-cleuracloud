package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// OpenStackPrivileges represents the openstack privileges.
type ccpOpenstackPrivileges struct {
	Type types.String `tfsdk:"type"`
	Meta types.String `tfsdk:"meta"`
	// ProjectPrivileges []ccpProjectPrivileges `tfsdk:"project_privileges"`
}

// ProjectPrivilege represents individual project privilege.
// type ccpProjectPrivileges struct {
// 	ProjectID types.String `tfsdk:"project_id"`
// 	DomainID  types.String `tfsdk:"domain_id"`
// 	Type      types.String `tfsdk:"type"`
// }

// ------------------------- JSON
// User represents the structure of the JSON data.
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

// Privileges represents the nested privileges object.
type ccpPrivilegesJson struct {
	Users     ccpUsersPrivilegeJson      `json:"users"`
	OpenStack ccpOpenstackPrivilegesJson `json:"openstack"`
}

// InvoicePrivileges represents the invoice privileges.
type ccpUsersPrivilegeJson struct {
	Type string `json:"type"`
	Meta string `json:"meta"`
}

// OpenStackPrivileges represents the openstack privileges.
type ccpOpenstackPrivilegesJson struct {
	Type string `json:"type"`
	Meta string `json:"meta"`
	// ProjectPrivileges []ccpProjectPrivilegesJson `json:"project_privileges"`
}

// User represents the top-level user object in the schema.
type ccpUserResourceModelJson struct {
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	FirstName      string   `json:"firstname,omitempty"`
	LastName       string   `json:"lastname,omitempty"`
	Password       string   `json:"password,omitempty"`
	IpRestrictions []string `json:"ip_restrictions,omitempty"`
	// Privileges     Privileges `json:"privileges"`
}

type ccpUserResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Password  types.String `tfsdk:"password"`
	// IpRestrictions []string `tfsdk:"ip_restrictions"`
	// Privileges     Privileges `json:"privileges"`
}

// Privileges represents the nested privileges object within the schema.
// type ccpUserPrivilegesResourceJson struct {
// 	OpenStack   OpenStack   `json:"openstack,omitempty"`
// 	Invoice     Invoice     `json:"invoice,omitempty"`
// 	CityMonitor CityMonitor `json:"citymonitor,omitempty"`
// 	Users       Users       `json:"users,omitempty"`
// }

// // --------------------  CCP User
