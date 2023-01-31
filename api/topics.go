package api

import (
	"fmt"
)

type ReadTopicResponse struct {
	Name       string `json:"name"`
	Partitions int    `json:"partitions"`
	Replicas   int    `json:"replicas"`
	Config     Json   `json:"config"`
}

func (api *API) ReadTopic(instanceId int, name string) (*ReadTopicResponse, error) {
	var (
		data   []ReadTopicResponse
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/topics", instanceId)
	_, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("topic %s not found", name)
}

type CreateTopicParams struct {
	Name       string `json:"name"`
	Partitions int    `json:"partitions"`
	Replicas   int    `json:"replicas"`
	Config     Json   `json:"config"`
}

func (api *API) CreateTopic(instanceId int, params CreateTopicParams) error {
	var (
		data   interface{}
		failed interface{}
	)

	path := fmt.Sprintf("/api/instances/%d/topics", instanceId)
	_, err := api.client.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	return err
}

type UpdateTopicParams struct {
	Partitions int  `json:"partitions"`
	Config     Json `json:"config"`
}

func (api *API) UpdateTopic(instanceId int, name string, params UpdateTopicParams) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/topics/%s", instanceId, name)
	_, err := api.client.New().Put(path).BodyJSON(params).Receive(&data, &failed)
	return err
}

func (api *API) DeleteTopic(instanceId int, name string) error {
	var (
		data   interface{}
		failed interface{}
	)
	path := fmt.Sprintf("/api/instances/%d/topics/%s", instanceId, name)
	_, err := api.client.New().Delete(path).Receive(&data, &failed)
	return err
}
