package api

import (
	"fmt"
	"time"
)

type ClusterStatus struct {
	Name       string `json:"name"`
	Ready      bool   `json:"ready"`
	Configured bool   `json:"configured"`
}

type VPC struct {
	Id     int    `json:"id"`
	Subnet string `json:"subnet"`
}

type InstanceResponse struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Plan         string   `json:"plan"`
	Region       string   `json:"region"`
	Tags         []string `json:"tags"`
	KafkaVersion string   `json:"kafka_version"`
	ApiKey       string   `json:"apikey"`
	BrokerUrl    string   `json:"brokers"`
	Password     string   `json:"password"`
	Username     string   `json:"username"`
	Vpc          VPC      `json:"vpc"`
}

type CreateInstanceRequest struct {
	Name         string   `json:"name"`
	Plan         string   `json:"plan"`
	Region       string   `json:"region"`
	Tags         []string `json:"tags,omitempty"`
	KafkaVersion string   `json:"kafka_version"`
	VpcId        int64    `json:"vpc_id,omitempty"`
	VpcSubnet    string   `json:"vpc_subnet,omitempty"`
}

type UpdateInstanceRequest struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Plan string   `json:"plan"`
}

func (api *API) waitUntilReady(id int64) error {
	var data ClusterStatus
	for {
		time.Sleep(10 * time.Second)
		path := fmt.Sprintf("api/instances/%d/cluster/status", id)
		_, err := api.client.New().Path(path).ReceiveSuccess(&data)
		if err != nil {
			return err
		}
		if data.Configured && data.Ready {
			return nil
		}
	}
}

func (api *API) readInstance(id int64) (InstanceResponse, error) {
	var data InstanceResponse
	var error APIError
	path := fmt.Sprintf("api/instances/%d", id)
	response, err := api.client.New().Path(path).Receive(&data, &error)
	if err != nil {
		return InstanceResponse{}, err
	}
	if response.StatusCode == 404 {
		return InstanceResponse{}, fmt.Errorf("Instance with id %d not found", id)
	}
	if response.StatusCode != 200 {
		return InstanceResponse{}, fmt.Errorf("failed to fetch info about instance: %s", error.Error())
	}
	return data, nil
}

func (api *API) CreateInstance(req CreateInstanceRequest) (InstanceResponse, error) {
	var (
		data   map[string]interface{}
		failed APIErrors
	)
	resp := InstanceResponse{}

	response, err := api.client.New().Post("/api/instances").BodyJSON(req).Receive(&data, &failed)
	if err != nil {
		return resp, err
	}
	if response.StatusCode == 401 {
		return resp, fmt.Errorf("Authentication error: %s", "invalid API key used")
	}
	if response.StatusCode == 400 {
		return resp, fmt.Errorf("Validation error: %s", failed.Error())
	}
	if response.StatusCode != 200 {
		return resp, fmt.Errorf("failed to create instance: %s", failed.Error())
	}
	instanceId := int64(data["id"].(float64))
	for {
		if err := api.waitUntilReady(instanceId); err != nil {
			return resp, err
		} else {
			break
		}
	}
	return api.readInstance(instanceId)
}

func (api *API) ReadInstance(id int64) (InstanceResponse, error) {
	return api.readInstance(id)
}

func (api *API) UpdateInstance(id int64, data UpdateInstanceRequest) error {
	var failed APIErrors
	path := fmt.Sprintf("api/instances/%d", id)
	response, err := api.client.New().Put(path).BodyJSON(data).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if response.StatusCode == 400 {
		return failed
	}
	if response.StatusCode == 401 {
		return fmt.Errorf("Authentication error: %s", "invalid API key used")
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("update instance failed: %s", failed.Error())
	}
	return api.waitUntilReady(id)
}

func (api *API) DeleteInstance(id int64, keep_vpc bool) error {
	var failed APIError
	path := fmt.Sprintf("api/instances/%d?keep_vpc=%v", id, keep_vpc)
	response, err := api.client.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("failed to delete instance: %s", failed.Error())
	}

	return nil
}
