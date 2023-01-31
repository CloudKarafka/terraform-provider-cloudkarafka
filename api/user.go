package api

import (
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (api *API) ReadUser(instanceId int, name string) (*User, error) {
	var (
		data   []User
		failed interface{}
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

func (api *API) CreateUser(instanceId int, params User) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/users", instanceId)
	_, err := api.client.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	return err
}

func (api *API) DeleteUser(instanceId int, name string) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/users/%s", instanceId, name)
	_, err := api.client.New().Delete(path).Receive(&data, &failed)
	return err
}
