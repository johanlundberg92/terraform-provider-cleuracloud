package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sethvargo/go-password/password"
)

type CleuraClient struct {
	Token    string
	User     string
	Password string
	Url      string
	Client   *http.Client
	DomainId string
}
type CleuraAuth struct {
	Auth CleuraAuthInfo `json:"auth"`
}
type CleuraAuthInfo struct {
	Username string `json:"login"`
	Password string `json:"password"`
}
type CleuraAuthResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

func (c *CleuraClient) Login(ctx context.Context) error {
	tflog.Trace(ctx, "Login was called", nil)
	login_marshalled, err := json.Marshal(CleuraAuth{CleuraAuthInfo{Username: c.User, Password: c.Password}})
	if err != nil {
		tflog.Error(ctx, "Marshalling of login info failed", nil)
		tflog.Error(ctx, err.Error())
		return err
	}

	buffer := bytes.NewBuffer(login_marshalled)
	response, err := http.Post(c.Url+"/auth/v1/tokens", "application/json", buffer)
	// req, err := http.NewRequest(http.MethodPost,c.Url, buffer)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to execute http post for login, error: %s", err.Error()), nil)
		return err
	}
	defer response.Body.Close()
	reader, err := io.ReadAll(response.Body)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to read response body after login request, error: %s", err.Error()), nil)
		return err
	}
	authToken := CleuraAuthResponse{}
	err = json.Unmarshal(reader, &authToken)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal response body into CleuraAuthResponse struct, error: %s", err.Error()), nil)
		return err
	}
	if authToken.Result != "login_ok" {
		tflog.Error(ctx, fmt.Sprintf("Response was not login_ok, response was: %s", authToken.Result), nil)
		return fmt.Errorf(fmt.Sprintf("Authentication result was not login_ok. Result was %s", authToken.Result))
	}
	c.Token = authToken.Token
	tflog.Trace(context.Background(), "Login complete!", nil)
	return nil
}
func (c *CleuraClient) GetUser(ctx context.Context, user string) (openstackUserDatasourceModel, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s", c.DomainId, user)
	cleuraUser := openstackUserDatasourceModelJson{}
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return openstackUserDatasourceModel{}, err
	}
	resultByteArray, err := io.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to read result into byte array, error: %s", err.Error()), nil)
		return openstackUserDatasourceModel{}, err
	}
	err = json.Unmarshal(resultByteArray, &cleuraUser)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal byte array into CleuraUser struct, error: %s", err.Error()), nil)
		return openstackUserDatasourceModel{}, err
	}

	response := openstackUserDatasourceModel{
		Id:               types.StringValue(cleuraUser.Id),
		Name:             types.StringValue(cleuraUser.Name),
		DomainId:         types.StringValue(cleuraUser.DomainId),
		DefaultProjectId: types.StringValue(cleuraUser.DefaultProjectId),
		Enabled:          types.BoolValue(cleuraUser.Enabled),
		Description:      types.StringValue(cleuraUser.Description),
	}
	for _, proj := range cleuraUser.Projects {
		var roles []openstackRole
		for _, role := range proj.Roles {
			roles = append(roles, openstackRole{
				Id:   types.StringValue(role.Id),
				Name: types.StringValue(role.Name),
			})
		}
		response.Projects = append(response.Projects, openstackProject{
			Id:       types.StringValue(proj.Id),
			Name:     types.StringValue(proj.Name),
			DomainId: types.StringValue(proj.DomainId),
			Roles:    roles,
		})
	}
	return response, nil

}
func (c *CleuraClient) DeleteUser(ctx context.Context, user string) error {
	apiPath := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s", c.DomainId, user)
	resp, err := c.delete(apiPath)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := &apiError{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(body, apiErr)
		errMsg := fmt.Sprintf("failed to delete user, error: %+v", body)
		tflog.Error(ctx, errMsg)
		return errors.New(errMsg)
	}
	return nil
}
func (c *CleuraClient) GetUserResource(ctx context.Context, user string) (openstackUserResourceModel, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s", c.DomainId, user)
	cleuraUser := openstackUserDatasourceModelJson{}
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return openstackUserResourceModel{}, err
	}
	resultByteArray, err := io.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to read result into byte array, error: %s", err.Error()), nil)
		return openstackUserResourceModel{}, err
	}
	err = json.Unmarshal(resultByteArray, &cleuraUser)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal byte array into CleuraUser struct, error: %s", err.Error()), nil)
		return openstackUserResourceModel{}, err
	}
	response := openstackUserResourceModel{
		Id:       types.StringValue(cleuraUser.Id),
		Name:     types.StringValue(cleuraUser.Name),
		DomainId: types.StringValue(cleuraUser.DomainId),
		// DefaultProjectId: types.StringValue(cleuraUser.DefaultProjectId),
		Enabled: types.BoolValue(cleuraUser.Enabled),
		// Description:      types.StringValue(cleuraUser.Description),
	}
	if len(cleuraUser.DefaultProjectId) == 0 {
		response.DefaultProjectId = types.StringNull()
	} else {
		response.DefaultProjectId = types.StringValue(cleuraUser.DefaultProjectId)
	}
	if len(cleuraUser.Description) == 0 {
		response.Description = types.StringNull()
	} else {
		response.Description = types.StringValue(cleuraUser.Description)
	}
	for _, proj := range cleuraUser.Projects {
		var roles []string
		for _, role := range proj.Roles {
			roles = append(roles, role.Name)
		}
		response.Projects = append(response.Projects, openstackUserCreateProject{
			Id: proj.Id,
			// Name:     types.StringValue(proj.Name),
			// DomainId: types.StringValue(proj.DomainId),
			Roles: roles,
		})
	}
	return response, nil
}
func (c *CleuraClient) DoesUserExist(ctx context.Context, user string) (bool, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s", c.DomainId, user)
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return false, err
	}
	if result.StatusCode != 200 {
		if result.StatusCode == 400 {
			return false, nil
		}
		return false, fmt.Errorf("return code from cleura API is: %d", result.StatusCode)
	}
	return true, nil
}
func (c *CleuraClient) CreateUser(ctx context.Context, model openstackUserResourceModel) (openstackUserCreatedModel, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/openstack/%s/users", model.DomainId.ValueString())
	payload := createOpenstackUser{}
	pw, err := password.Generate(12, 2, 0, false, true)
	if err != nil {
		return openstackUserCreatedModel{}, err
	}
	payload.User = createOpenstackUserInfo{Name: model.Name.ValueString(), Password: pw, Description: model.Description.ValueString()}
	projectList := make([]openstackUserCreateProject, 0)
	for _, p := range model.Projects {
		projectList = append(projectList, openstackUserCreateProject{Id: p.Id, Roles: p.Roles})
	}
	payload.Projects = projectList
	// jsonPayload, err := json.Marshal(payload)
	// if err != nil {
	// 	return err
	// }
	result, err := c.post(payload, apiPath)
	if err != nil {
		return openstackUserCreatedModel{}, err
	}
	msg, err := io.ReadAll(result.Body)
	if err != nil {
		return openstackUserCreatedModel{}, err
	}
	if result.StatusCode != 201 {
		apiErr := &apiError{}
		json.Unmarshal(msg, apiErr)
		tflog.Error(ctx, fmt.Sprintf("%+v", result))
		return openstackUserCreatedModel{}, errors.New("return code was not 200")
	}
	created := &openstackUserCreatedModel{}
	json.Unmarshal(msg, created)
	return *created, nil

}
func (c *CleuraClient) post(payload interface{}, apiPath string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.Url, apiPath)
	marshaled_payload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshaled_payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AUTH-LOGIN", c.User)
	req.Header.Set("X-AUTH-TOKEN", c.Token)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()
	return resp, nil
}
func (c *CleuraClient) get(apiPath string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.Url, apiPath)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-AUTH-LOGIN", c.User)
	req.Header.Set("X-AUTH-TOKEN", c.Token)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
func (c *CleuraClient) delete(apiPath string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.Url, apiPath)
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-AUTH-LOGIN", c.User)
	req.Header.Set("X-AUTH-TOKEN", c.Token)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *CleuraClient) put(payload interface{}, apiPath string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.Url, apiPath)
	marshaled_payload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(marshaled_payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AUTH-LOGIN", c.User)
	req.Header.Set("X-AUTH-TOKEN", c.Token)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()
	return resp, nil
}
func (c *CleuraClient) AddUserToProjectRole(ctx context.Context, user string, projectId string, projectRole string) error {
	apiUrl := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s/projects", c.DomainId, user)
	roles := []string{projectRole}
	ass := openstackProjectAssignment{ProjectId: projectId, Roles: roles}
	assignments := []openstackProjectAssignment{ass}
	resp, err := c.post(openstackProjectUpdate{Projects: assignments}, apiUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		tflog.Error(ctx, fmt.Sprintf("status code returned is: %d", resp.StatusCode))
		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		errMsg := &apiError{}
		err = json.Unmarshal(r, errMsg)
		if err != nil {
			return err
		}
		errStr := fmt.Sprintf("error message from api is: %+v", errMsg)
		tflog.Error(ctx, errStr)
		return errors.New(errStr)
	}
	return nil
}
func (c *CleuraClient) RemoveUserFromProjectRole(ctx context.Context, user string, projectId string, role string) error {
	apiUrl := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s/projects/%s/%s", c.DomainId, user, projectId, role)
	resp, err := c.delete(apiUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		tflog.Error(ctx, fmt.Sprintf("status code returned is: %d", resp.StatusCode))
		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		errMsg := &apiError{}
		err = json.Unmarshal(r, errMsg)
		if err != nil {
			return err
		}
		errStr := fmt.Sprintf("error message from api is: %+v", errMsg)
		tflog.Error(ctx, errStr)
		return errors.New(errStr)
	}
	return nil

}
func (c *CleuraClient) AddUserToProject(projects openstackProjectUpdate) {
}
func (c *CleuraClient) ToggleUserEnabled(ctx context.Context, user string, enabled bool) error {
	url := fmt.Sprintf("accesscontrol/v1/openstack/%s/users/%s", c.DomainId, user)

	resp, err := c.put(openstackUserUpdate{User: openstackUserUpdateProperties{Enabled: enabled}}, url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		tflog.Error(ctx, fmt.Sprintf("status code returned is: %d", resp.StatusCode))
		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		errMsg := &apiError{}
		err = json.Unmarshal(r, errMsg)
		if err != nil {
			return err
		}
		errStr := fmt.Sprintf("error message from api is: %+v", errMsg)
		tflog.Error(ctx, errStr)
		return errors.New(errStr)
	}
	return nil
}
func (c *CleuraClient) GetCCPUser(ctx context.Context, name string) (ccpUserDataSourceModel, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/users/%s", name)
	ccpUser := ccpUserJson{}
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return ccpUserDataSourceModel{}, err
	}
	if result.StatusCode != 200 {
		tflog.Error(ctx, fmt.Sprintf("return code was not 200, return code: %d when searching for user: %s", result.StatusCode, name))
		return ccpUserDataSourceModel{}, fmt.Errorf("return code was not 200, return code: %d when searching for user: %s", result.StatusCode, name)
	}
	resultByteArray, err := io.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to read result into byte array, error: %s", err.Error()), nil)
		return ccpUserDataSourceModel{}, err
	}
	err = json.Unmarshal(resultByteArray, &ccpUser)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal byte array into CCPUser struct, error: %s", err.Error()), nil)
		return ccpUserDataSourceModel{}, err
	}
	response := ccpUserDataSourceModel{
		Id:             types.StringValue(ccpUser.Id),
		Name:           types.StringValue(ccpUser.Name),
		Admin:          types.BoolValue(ccpUser.Admin),
		FirstName:      types.StringValue(ccpUser.FirstName),
		LastName:       types.StringValue(ccpUser.LastName),
		Email:          types.StringValue(ccpUser.Email),
		Language:       types.StringValue(ccpUser.Language),
		Currency:       &ccpCurrency{Id: types.StringValue(ccpUser.Currency.Id), Code: types.StringValue(ccpUser.Currency.Code), Name: types.StringValue(ccpUser.Currency.Name)},
		AuthProviderId: types.StringValue(ccpUser.AuthProviderId),
		// TwoFactorLogin: ccpUser.TwoFactorLogin,
		// IPRestrictions: ccpUser.IPRestrictions,
	}
	userPrivileges := ccpUsersPrivilege{
		Type: types.StringValue(ccpUser.Privileges.Users.Type),
		Meta: types.StringValue(ccpUser.Privileges.Users.Meta),
	}

	// Uncomment when adding support for project privileges
	// var osProjPrivileges []ccpProjectPrivileges
	// for _, osp := range ccpUser.Privileges.OpenStack.ProjectPrivileges {
	// 	projectPrivileges := ccpProjectPrivileges{
	// 		ProjectId: types.StringValue(osp.ProjectId),
	// 		DomainId:  types.StringValue(osp.DomainId),
	// 		Type:      types.StringValue(osp.Type),
	// 	}
	// 	osProjPrivileges = append(osProjPrivileges, projectPrivileges)
	// }
	// osPrivileges := ccpOpenstackPrivileges{Type: types.StringValue(ccpUser.Privileges.OpenStack.Type), Meta: types.StringValue(ccpUser.Privileges.OpenStack.Meta), ProjectPrivileges: osProjPrivileges}
	osPrivileges := ccpOpenstackPrivileges{Type: types.StringValue(ccpUser.Privileges.OpenStack.Type), Meta: types.StringValue(ccpUser.Privileges.OpenStack.Meta)}
	privileges := ccpPrivileges{
		Users:     userPrivileges,
		OpenStack: osPrivileges,
	}
	response.Privileges = &privileges
	return response, nil
}

func (c *CleuraClient) GetCCPUserResource(ctx context.Context, name string) (ccpUserResourceModel, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/users/%s", name)
	ccpUser := ccpUserJson{}
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return ccpUserResourceModel{}, err
	}
	if result.StatusCode != 200 {
		tflog.Error(ctx, fmt.Sprintf("return code was not 200, return code: %d when searching for user: %s", result.StatusCode, name))
		return ccpUserResourceModel{}, fmt.Errorf("return code was not 200, return code: %d when searching for user: %s", result.StatusCode, name)
	}
	resultByteArray, err := io.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to read result into byte array, error: %s", err.Error()), nil)
		return ccpUserResourceModel{}, err
	}
	err = json.Unmarshal(resultByteArray, &ccpUser)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal byte array into CCPUser struct, error: %s", err.Error()), nil)
		return ccpUserResourceModel{}, err
	}
	response := ccpUserResourceModel{
		Id:   types.StringValue(ccpUser.Id),
		Name: types.StringValue(ccpUser.Name),
		// Admin:          types.BoolValue(ccpUser.Admin),
		FirstName: types.StringValue(ccpUser.FirstName),
		LastName:  types.StringValue(ccpUser.LastName),
		Email:     types.StringValue(ccpUser.Email),
		// Language:  types.StringValue(ccpUser.Language),
		// Currency:       &ccpCurrency{Id: types.StringValue(ccpUser.Currency.Id), Code: types.StringValue(ccpUser.Currency.Code), Name: types.StringValue(ccpUser.Currency.Name)},
		// AuthProviderId: types.StringValue(ccpUser.AuthProviderId),
		// TwoFactorLogin: ccpUser.TwoFactorLogin,
		// IPRestrictions: ccpUser.IPRestrictions,
	}
	userPrivileges := ccpUserResourceUserPrivilege{
		Type: types.StringValue(ccpUser.Privileges.Users.Type),
	}
	// Uncomment when adding support for project privileges
	// var osProjPrivileges []ccpUserResourceProjectPrivileges
	// for _, osp := range ccpUser.Privileges.OpenStack.ProjectPrivileges {
	// 	projectPrivileges := ccpUserResourceProjectPrivileges{
	// 		ProjectId: types.StringValue(osp.ProjectId),
	// 		DomainId:  types.StringValue(osp.DomainId),
	// 		Type:      types.StringValue(osp.Type),
	// 	}
	// 	osProjPrivileges = append(osProjPrivileges, projectPrivileges)
	// }
	// osPrivileges := ccpUserResourceOpenstackPrivileges{Type: types.StringValue(ccpUser.Privileges.OpenStack.Type), ProjectPrivileges: osProjPrivileges}
	osPrivileges := ccpUserResourceOpenstackPrivileges{Type: types.StringValue(ccpUser.Privileges.OpenStack.Type)}
	privileges := ccpResourcePrivileges{
		Users:     userPrivileges,
		OpenStack: osPrivileges,
	}
	response.Privileges = &privileges
	return response, nil
}
func (c *CleuraClient) CreateCCPUser(ctx context.Context, model ccpUserResourceModel) (ccpUserResourceModel, error) {
	apiPath := "accesscontrol/v1/users"
	// payload := ccpUserResourceModelJson{}
	// pw, err := password.Generate(12, 2, 0, false, true)
	// if err != nil {
	// 	return ccpUserResourceModel{}, err
	// }
	// model.Password = types.StringValue(pw)
	modelJson := ccpUserResourceModelJson{
		Name:      model.Name.ValueString(),
		Email:     model.Email.ValueString(),
		FirstName: model.FirstName.ValueString(),
		LastName:  model.LastName.ValueString(),
	}
	if model.Privileges == nil {
		modelJson.Privileges = &ccpResourcePrivilegesJson{}
	} else {
		modelJson.Privileges = &ccpResourcePrivilegesJson{}
		modelJson.Privileges.Users = ccpUserResourcePrivilegeJson{
			Type: model.Privileges.Users.Type.ValueString(),
			// Meta: model.Privileges.Users.Meta.ValueString(),
		}
		// Uncomment when adding support for project privileges
		// var projectPrivileges []ccpUserResourceProjectPrivilegesJson
		// if model.Privileges.OpenStack.ProjectPrivileges != nil {
		// 	for _, p := range model.Privileges.OpenStack.ProjectPrivileges {
		// 		projectPrivileges = append(projectPrivileges, ccpUserResourceProjectPrivilegesJson{ProjectId: p.ProjectId.ValueString(), DomainId: p.DomainId.ValueString(), Type: p.Type.ValueString()})
		// 	}
		// 	modelJson.Privileges.OpenStack.ProjectPrivileges = projectPrivileges
		// }
		modelJson.Privileges.OpenStack = ccpUserOpenstackPrivilegesJson{
			Type: model.Privileges.OpenStack.Type.ValueString(),
			// Meta: model.Privileges.OpenStack.Meta.ValueString(),
		}
	}

	// var result map[string]interface{}
	// json.Unmarshal(payload, &result)
	// userObj, ok := result["user"].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("User object not found")
	// 	return ccpUserResourceModel{}, err
	// }

	// privileges, ok := userObj["privileges"].(map[string]interface{})
	// if !ok || len(privileges) == 0 {
	// 	fmt.Println("Privileges not found or empty")
	// } else {
	// 	for key := range privileges {
	// 		delete(privileges, key)
	// 	}
	// 	userObj["privileges"] = privileges // Update the privileges field in the user object
	// 	result["user"] = userObj           // Update the user object in the result map
	// }

	//"{\"user\":{\"name\":\"johan.testberg\",\"email\":\"johan.testberg@thernfrst.io\",\"firstname\":\"johan\",\"lastname\":\"testberg\",\"privileges\":{\"users\":{\"type\":\"\",\"meta\":\"\"},\"openstack\":{\"type\":\"\",\"meta\":\"\"}}}}"
	resp, err := c.post(ccpUserCreateJson{User: modelJson}, apiPath)
	if err != nil {
		return ccpUserResourceModel{}, err
	}
	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return ccpUserResourceModel{}, err
	}
	if resp.StatusCode == 409 {
		tflog.Error(ctx, "the specified user already exists")
		return ccpUserResourceModel{}, fmt.Errorf("the specified user already exists")
	}
	if resp.StatusCode != 200 {
		// if err != nil {
		// 	return ccpUserResourceModel{}, err
		// }
		apiErr := &apiError{}
		json.Unmarshal(msg, apiErr)
		tflog.Error(ctx, fmt.Sprintf("%+v", apiErr))
		return ccpUserResourceModel{}, errors.New("return code was not 200")
	}
	created := &openstackUserCreatedModel{}
	json.Unmarshal(msg, created)
	return model, nil

}
func (c *CleuraClient) DoesCCPUserExist(ctx context.Context, user string) (bool, error) {
	apiPath := fmt.Sprintf("accesscontrol/v1/users/%s", user)
	result, err := c.get(apiPath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occurred when executing get, error: %s", err.Error()))
		return false, err
	}
	if result.StatusCode != 200 {
		if result.StatusCode == 400 {
			return false, nil
		}
		return false, fmt.Errorf("return code from cleura API is: %d", result.StatusCode)
	}
	return true, nil
}
func (c *CleuraClient) UpdateCCPUser(ctx context.Context, resource ccpUserUpdate) error {
	apiPath := fmt.Sprintf("accesscontrol/v1/users/%s", resource.User.Name)
	result, err := c.put(resource, apiPath)
	if err != nil {
		return err
	}
	if result.StatusCode != 200 {
		defer result.Body.Close()
		tflog.Error(ctx, fmt.Sprintf("return code when updating CCP user was not 200, code was: %d", result.StatusCode))
		msg, err := io.ReadAll(result.Body)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode response body, error: %s", err.Error()))
			return fmt.Errorf("failed to decode response body, error: %s", err.Error())
		}
		apiErr := &apiError{}
		json.Unmarshal(msg, apiErr)
		tflog.Error(ctx, fmt.Sprintf("%+v", apiErr))
		return fmt.Errorf("error returned from API while updating user is: %+v", apiErr)
	}
	return nil
}
func (c *CleuraClient) CheckApiError(response *http.Response) {

}
func (c *CleuraClient) DeleteCCPUser(ctx context.Context, user string) error {
	apiPath := fmt.Sprintf("accesscontrol/v1/users/%s", user)
	resp, err := c.delete(apiPath)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := &apiError{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(body, apiErr)
		errMsg := fmt.Sprintf("failed to delete user, error: %+v", body)
		tflog.Error(ctx, errMsg)
		return errors.New(errMsg)
	}
	return nil
}
