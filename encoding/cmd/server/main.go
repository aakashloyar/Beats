package main 

import (
	"log"
	"context"
	"github.com/aakashloyar/beats/encoding/internal/adapter/out/s3"
	decoder "github.com/aakashloyar/beats/encoding/internal/adapter/out/decoder"
	encoder "github.com/aakashloyar/beats/encoding/internal/adapter/out/encoder"
	kafkaProducer "github.com/aakashloyar/beats/encoding/internal/adapter/out/kafka_producer"
	"github.com/aakashloyar/beats/encoding/internal/application/service"
	kafkaConsumer "github.com/aakashloyar/beats/encoding/internal/adapter/in/kafka_consumer"
)

func main() {
	ctx := context.Background()

	s3Config := s3.Config {
		Region: "ap-south-1",
		Bucket: "your-bucket-name",
	}

	s3Client, err := s3Config.NewS3Client(ctx)

	s3Storage :=  s3.NewS3Storage(s3Client.Client,s3Config.Bucket)

	kafkaProducerConfig := kafkaProducer.Config{
		Brokers:  []string{"localhost:9092"},
		Topic:    "upload-completed",
		ClientID: "ingestion-service",
	}

	producer, err := kafkaProducer.NewFranzProducer(kafkaProducerConfig); if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	validator := service.NewValidator()
	decoderProbe := decoder.NewFFprobeAdapter()
	encoderProbe := encoder.NewEncoderAdapter()

	encodingService := service.NewEncodingService(s3Storage, decoderProbe, validator,encoderProbe, producer)


	kafkaConsumerConfig := kafkaConsumer.Config{
		Brokers:  []string{"localhost:9092"},
		Topic:    "upload-completed",
		ClientID: "ingestion-service",
	}

	consumer, err := kafkaConsumerConfig.NewConsumer(encodingService) 
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start(ctx)
	defer consumer.Close()
}