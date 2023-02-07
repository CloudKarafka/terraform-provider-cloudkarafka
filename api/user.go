package api

import (
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (api *API) ReadUser(instanceId int64, name string) (*User, error) {
	var (
		data   []User
		failed APIError
	)
	path := fmt.Sprintf("/api/instances/%d/users", instanceId)
	_, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("user %s not found", name)
}

func (api *API) CreateUser(instanceId int64, params User) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/users", instanceId)
	resp, err := api.client.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return failed
	}
	return nil
}

func (api *API) DeleteUser(instanceId int64, name string) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/users/%s", instanceId, name)
	resp, err := api.client.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return failed
	}
	return nil
}
