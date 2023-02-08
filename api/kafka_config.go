package api

import (
	"fmt"
	"strings"
)

type KafkaConfig struct {
	AutoCreateTopics  bool
	MinInsyncReplicas int64
	LogRetentionBytes int64
	LogRetentionMs    int64
	LogSegmentBytes   int64
	NetworkThreads    int64
	IOThreads         int64
	MessageMaxBytes   int64
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		AutoCreateTopics:  false,
		MinInsyncReplicas: -1,
		LogRetentionBytes: -1,
		LogRetentionMs:    -1,
		LogSegmentBytes:   -1,
		NetworkThreads:    -1,
		IOThreads:         -1,
		MessageMaxBytes:   -1,
	}
}

func (me *KafkaConfig) AsProperties() string {
	var b strings.Builder
	if me.AutoCreateTopics {
		b.WriteString(fmt.Sprintf("auto.create.topics.enable=%v\n", me.AutoCreateTopics))
	}
	if me.MinInsyncReplicas != -1 {
		b.WriteString(fmt.Sprintf("min.insync.replicas=%d\n", me.MinInsyncReplicas))
	}
	if me.LogRetentionBytes != -1 {
		b.WriteString(fmt.Sprintf("log.retention.bytes=%d\n", me.LogRetentionBytes))
	}
	if me.LogRetentionMs != -1 {
		b.WriteString(fmt.Sprintf("log.retention.ms=%d\n", me.LogRetentionMs))
	}
	if me.LogSegmentBytes != -1 {
		b.WriteString(fmt.Sprintf("log.segment.bytes=%d\n", me.LogSegmentBytes))
	}
	if me.NetworkThreads != -1 {
		b.WriteString(fmt.Sprintf("num.network.threads=%d\n", me.NetworkThreads))
	}
	if me.IOThreads != -1 {
		b.WriteString(fmt.Sprintf("num.io.threads=%d\n", me.IOThreads))
	}
	if me.MessageMaxBytes != -1 {
		b.WriteString(fmt.Sprintf("message.max.bytes=%d\n", me.MessageMaxBytes))
	}
	return b.String()

}

func (api *API) ReadConfig(instanceId int64) (*KafkaConfig, error) {
	var (
		data   []Hash
		failed APIError
	)
	path := fmt.Sprintf("/api/instances/%d/config/kafka", instanceId)
	resp, err := api.client.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, failed
	}
	cfg := new(KafkaConfig)
	for _, r := range data {
		switch r["name"] {
		case "auto.create.topics.enable":
			cfg.AutoCreateTopics = r["value"].(bool)
		case "num.io.threads":
			cfg.IOThreads = int64(r["value"].(float64))
		case "num.network.threads":
			cfg.NetworkThreads = int64(r["value"].(float64))
		case "log.retention.ms":
			cfg.LogRetentionMs = int64(r["value"].(float64))
		case "log.retention.bytes":
			cfg.LogRetentionBytes = int64(r["value"].(float64))
		case "log.segment.bytes":
			cfg.LogSegmentBytes = int64(r["value"].(float64))
		case "min.insync.replicas":
			cfg.MinInsyncReplicas = int64(r["value"].(float64))
		case "message.max.bytes":
			cfg.MessageMaxBytes = int64(r["value"].(float64))
		default:
			return nil, fmt.Errorf("Unhandled config value %s", r["name"])
		}
	}
	return cfg, nil
}

func (api *API) WriteConfig(instanceId int64, config *KafkaConfig) error {
	var failed APIError
	path := fmt.Sprintf("/api/instances/%d/config/kafka", instanceId)
	body := strings.NewReader(config.AsProperties())
	resp, err := api.client.New().Post(path).Body(body).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return failed
	}
	return nil
}
