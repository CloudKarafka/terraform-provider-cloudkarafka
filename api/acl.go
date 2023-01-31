package api

import (
	"fmt"
)

type AclRule struct {
	Id              int    `json:"id"`
	User            string `json:"user"`
	Operation       string `json:"operation"`
	Resource        string `json:"resource"`
	ResourcePattern string `json:"resource_pattern"`
}

func (api *API) ReadAclRule(instanceId int, user, operation, resource, resourcePattern string) (*AclRule, error) {
	var (
		data   []AclRule
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/acls", instanceId)
	_, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		if v.User == user && v.Operation == operation &&
			v.Resource == resource && v.ResourcePattern == resourcePattern {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("rule not found")
}

type CreateAclRule struct {
	Operation       string `json:"operation"`
	Resource        string `json:"resource"`
	ResourcePattern string `json:"resource_pattern"`
}
type CreateAclRuleParams struct {
	User  string          `json:"user"`
	Rules []CreateAclRule `json:"rules"`
}

func (api *API) CreateAclRule(instanceId int, params CreateAclRuleParams) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/acls", instanceId)
	_, err := api.client.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	return err
}

func (api *API) DeleteAclRule(instanceId int, id int) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/acls/%d", instanceId, id)
	_, err := api.client.New().Delete(path).Receive(&data, &failed)
	return err
}
