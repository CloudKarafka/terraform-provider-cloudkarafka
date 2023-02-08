package api

import (
	"errors"
	"fmt"
)

type AclRule struct {
	Id                  int64  `json:"id"`
	User                string `json:"name"`
	Operation           string `json:"operation"`
	Resource            string `json:"resource"`
	ResourcePattern     string `json:"resource_pattern"`
	CreatedAt           string `json:"created_at"`
	ResourcePatternType string `json:"resource_pattern_type"`
}

func (r *AclRule) Same(other *AclRule) bool {
	return r.User == other.User &&
		r.Operation == other.Operation && r.Resource == other.Resource &&
		r.ResourcePattern == other.ResourcePattern &&
		r.ResourcePatternType == r.ResourcePatternType
}

func (api *API) readAclRules(instanceId int64) ([]AclRule, error) {
	var (
		data   []AclRule
		failed APIError
	)
	path := fmt.Sprintf("/api/instances/%d/acls", instanceId)
	resp, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, failed
	}
	return data, nil
}

func (api *API) ReadAclRule(instanceId int64, id int64) (*AclRule, error) {
	data, err := api.readAclRules(instanceId)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("No rule found with id=%d", id)
}

func (api *API) CreateAclRule(instanceId int64, user string, rule AclRule) (int64, error) {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/acls", instanceId)
	body := map[string]interface{}{
		"user":  user,
		"rules": []AclRule{rule},
	}
	resp, err := api.client.New().Post(path).BodyJSON(body).Receive(nil, &failed)
	if err != nil {
		return -1, err
	}
	if resp.StatusCode != 201 {
		return -1, failed
	}
	data, err := api.readAclRules(instanceId)
	if err != nil {
		return -1, err
	}
	rule.User = user
	for _, v := range data {
		if v.Same(&rule) {
			return v.Id, nil
		}
	}
	return -1, errors.New("Failed to create rule")
}

func (api *API) DeleteAclRule(instanceId int64, id int64) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/acls/%d", instanceId, id)
	resp, err := api.client.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return failed
	}
	return nil
}
