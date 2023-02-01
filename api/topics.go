package api

import (
	"fmt"
	"time"
)

type Topic struct {
	Name       string `json:"name"`
	Partitions int64  `json:"partitions"`
	Replicas   int64  `json:"replicas"`
	Status     string `json:"status,omitempty"`
	Config     Hash   `json:"config,omitempty"`
}

type UpdateTopicRequest struct {
	Partitions int64 `json:"partitions"`
	Config     Hash  `json:"config,omitempty"`
}

func (api *API) waitUntilTopicReady(instanceId int64, topic string) error {
	times := 0
	for {
		time.Sleep(5 * time.Second)
		topic, err := api.readTopic(instanceId, topic)
		if err != nil {
			return err
		}
		if topic.Status == "ready" {
			return nil
		}
		times += 1
		if times > 36 {
			return fmt.Errorf("Something appears to be failing waiting on topic %s to be created, please contact support", topic)
		}
	}
}

func (api *API) readTopics(instanceId int64) ([]Topic, error) {
	var (
		data   []Topic
		failed APIError
	)
	path := fmt.Sprintf("/api/instances/%d/topics", instanceId)
	_, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	return data, failed
}

func (api *API) readTopic(instanceId int64, name string) (*Topic, error) {
	topics, err := api.readTopics(instanceId)
	if err != nil {
		return nil, err
	}
	for _, v := range topics {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("topic %s not found", name)
}

func (api *API) ReadTopic(instanceId int64, name string) (*Topic, error) {
	return api.readTopic(instanceId, name)
}

func (api *API) CreateTopic(instanceId int64, params Topic) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/topics", instanceId)
	resp, err := api.client.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return failed
	}
	if err := api.waitUntilTopicReady(instanceId, params.Name); err != nil {
		return err
	}
	return nil
}

func (api *API) UpdateTopic(instanceId int64, name string, params UpdateTopicRequest) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/topics/%s", instanceId, name)
	_, err := api.client.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}
	return failed
}

func (api *API) DeleteTopic(instanceId int64, name string) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/topics/%s", instanceId, name)
	_, err := api.client.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}
	return failed
}
