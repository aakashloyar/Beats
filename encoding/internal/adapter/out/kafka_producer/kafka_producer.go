package kafkaproducer

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Config struct {
	Brokers []string
	Topic   string
	ClientID string
}

type FranzProducer struct {
	client *kgo.Client
	topic  string
}


func NewFranzProducer(cfg Config) (*FranzProducer, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.DefaultProduceTopic(cfg.Topic),
		kgo.ClientID(cfg.ClientID),
	)
	if err != nil {
		return nil, err
	}

	return &FranzProducer{
		client: client,
		topic:  cfg.Topic,
	}, nil
}