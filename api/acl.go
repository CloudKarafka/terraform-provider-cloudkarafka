package api

import (
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

func (api *API) ReadAclRuleForUser(instanceId int64, user string) ([]AclRule, error) {
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
	var res []AclRule
	for _, v := range data {
		if v.User == user {
			res = append(res, v)
		}
	}
	return res, nil
}

func (api *API) CreateAclRules(instanceId int64, user string, params []AclRule) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/acls", instanceId)
	body := map[string]interface{}{
		"user":  user,
		"rules": params,
	}
	resp, err := api.client.New().Post(path).BodyJSON(body).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return failed
	}
	return nil

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
