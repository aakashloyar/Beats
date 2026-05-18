package kafkaconsumer

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
)

type Config struct {
	Brokers  []string
	Topic    string
	ClientID string
	GroupID  string
}

type Consumer struct {
	kafkaClient        *KafkaConsumerClient
	createTrackService in.CreateTrackService 
}
type KafkaConsumerClient struct {
	client *kgo.Client
	topic  string
}

func (cfg Config)NewConsumer(createTrackService in.CreateTrackService) (*Consumer, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.ClientID(cfg.ClientID),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		kafkaClient: &KafkaConsumerClient{
			client: client,
			topic: cfg.Topic,
		},
		createTrackService: createTrackService,
	}, nil
}