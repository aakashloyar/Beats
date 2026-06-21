package main

import (
	"context"
	"log"

	"github.com/aakashloyar/beats/encoding/config"
	kafkaConsumer "github.com/aakashloyar/beats/encoding/internal/adapter/in/kafka_consumer"
	decoder "github.com/aakashloyar/beats/encoding/internal/adapter/out/decoder"
	encoder "github.com/aakashloyar/beats/encoding/internal/adapter/out/encoder"
	kafkaProducer "github.com/aakashloyar/beats/encoding/internal/adapter/out/kafka_producer"
	"github.com/aakashloyar/beats/encoding/internal/adapter/out/s3"
	"github.com/aakashloyar/beats/encoding/internal/application/service"
)

func main() {
	ctx := context.Background()

	s3Config := s3.Config{
		Region: config.App.S3.Region,
		Bucket: config.App.S3.Bucket,
	}

	s3Client, err := s3Config.NewS3Client(ctx)

	s3Storage := s3.NewS3Storage(s3Client.Client, s3Config.Bucket)

	kafkaProducerConfig := kafkaProducer.Config{
		Brokers:  config.App.KafkaProducer.Brokers,
		Topic:    config.App.KafkaProducer.Topic,
		ClientID: config.App.KafkaProducer.ClientID,
	}

	producer, err := kafkaProducer.NewFranzProducer(kafkaProducerConfig)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	validator := service.NewValidator()
	decoderProbe := decoder.NewFFprobeAdapter()
	encoderProbe := encoder.NewEncoderAdapter()

	encodingService := service.NewEncodingService(s3Storage, decoderProbe, validator, encoderProbe, producer)

	kafkaConsumerConfig := kafkaConsumer.Config{
		Brokers:  config.App.KafkaConsumer.Brokers,
		Topic:    config.App.KafkaConsumer.Topic,
		ClientID: config.App.KafkaConsumer.ClientID,
	}

	consumer, err := kafkaConsumerConfig.NewConsumer(encodingService)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start(ctx)
	defer consumer.Close()
}
