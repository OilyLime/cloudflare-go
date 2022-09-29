package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PermissionGroup struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Meta        map[string]string `json:"meta"`
	Permissions []Permission      `json:"permissions"`
}

type Permission struct {
	ID         string            `json:"id"`
	Key        string            `json:"key"`
	Attributes map[string]string `json:"attributes,omitempty"` // same as Meta in other structs
}

type PermissionGroupListResponse struct {
	Success  bool              `json:"success"`
	Errors   []string          `json:"errors"`
	Messages []string          `json:"messages"`
	Result   []PermissionGroup `json:"result"`
}

type PermissionGroupDetailResponse struct {
	Success  bool            `json:"success"`
	Errors   []string        `json:"errors"`
	Messages []string        `json:"messages"`
	Result   PermissionGroup `json:"result"`
}

// PermissionGroup returns a specific permission group from the API given the account ID and permission group ID
func (api *API) PermissionGroup(ctx context.Context, accountId string, permissionGroupId string) (PermissionGroup, error) {
	uri := fmt.Sprintf("/accounts/%s/iam/permission_groups/%s?depth=2", accountId, permissionGroupId)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return PermissionGroup{}, err
	}

	var permissionGroupResponse PermissionGroupDetailResponse
	err = json.Unmarshal(res, &permissionGroupResponse)
	if err != nil {
		return PermissionGroup{}, err
	}

	return permissionGroupResponse.Result, nil
}

// PermissionGroups returns all valid permission groups for a provided accountID
func (api *API) PermissionGroups(ctx context.Context, accountId string) ([]PermissionGroup, error) {
	uri := fmt.Sprintf("/accounts/%s/iam/permission_groups?depth=2", accountId)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []PermissionGroup{}, err
	}

	var permissionGroupResponse PermissionGroupListResponse
	err = json.Unmarshal(res, &permissionGroupResponse)
	if err != nil {
		return []PermissionGroup{}, err
	}

	return permissionGroupResponse.Result, nil
}

// FindPermissionGroupByName returns an array of permission groups which have been filtered by the name provided
func (api *API) FindPermissionGroupByName(ctx context.Context, accountId string, name string) ([]PermissionGroup, error) {
	uri := fmt.Sprintf("/accounts/%s/iam/permission_groups?name=%s&depth=2", accountId, name)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []PermissionGroup{}, err
	}

	var permissionGroupResponse PermissionGroupListResponse
	err = json.Unmarshal(res, &permissionGroupResponse)
	if err != nil {
		return []PermissionGroup{}, err
	}

	return permissionGroupResponse.Result, nil
}

// FindPermissionGroupForRole is a utility method to find permission groups with the role's name
func (api *API) FindPermissionGroupForRole(ctx context.Context, accountId string, role AccountRole) ([]PermissionGroup, error) {
	return api.FindPermissionGroupByName(ctx, accountId, role.Name)
}
