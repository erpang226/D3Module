package kafka

import "main/pkg/config"

type KafkaClient struct {
	Config *config.KafkaClientConfig
}

// NewKafkaClient initializes a new kafka client instance
func NewKafkaClient(conf *config.KafkaClientConfig) *KafkaClient {
	return &KafkaClient{Config: conf}
}
