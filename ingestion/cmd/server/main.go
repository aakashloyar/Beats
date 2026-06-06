package main

import (
	"net/http"
	"log"
	"context"
	"github.com/aakashloyar/beats/ingestion/internal/adapter/out/s3"
	"github.com/aakashloyar/beats/ingestion/internal/adapter/out/postgres"
	"github.com/aakashloyar/beats/ingestion/internal/application/ports/out/system"
	uploadSvc "github.com/aakashloyar/beats/ingestion/internal/application/service/upload"
	chunkSvc "github.com/aakashloyar/beats/ingestion/internal/application/service/chunk"
	httpUpload "github.com/aakashloyar/beats/ingestion/internal/adapter/in/http/upload"
	httpChunk "github.com/aakashloyar/beats/ingestion/internal/adapter/in/http/chunk"
	kafkaProducer "github.com/aakashloyar/beats/ingestion/internal/adapter/out/kafka-producer"
)
func main() {


	ctx := context.Background()
	postgresConfig := postgres.Config{
		Host:     "",
		Port:     123,
		User:     "",
		Password: "",
		DBName:   "",
		SSLMode:  "",
	}

	db, err := postgresConfig.NewDB()

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	s3Config := s3.Config {
		Region: "ap-south-1",
		Bucket: "your-bucket-name",
	}

	s3Client, err := s3Config.NewS3Client(ctx)

	s3Storage :=  s3.NewS3Storage(s3Client.Client,s3Config.Bucket)

	kafkaConfig := kafkaProducer.Config{
		Brokers:  []string{"localhost:9092"},
		Topic:    "upload-completed",
		ClientID: "ingestion-service",
	}

	producer, err := kafkaProducer.NewFranzProducer(kafkaConfig); if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	uploadRepo := postgres.NewUploadRepository(db)
	chunkRepo := postgres.NewChunkRepository(db)

	initUploadService := uploadSvc.NewInitUploadService(uploadRepo,s3Storage,idGen,clock)
	completeUploadService := uploadSvc.NewCompleteUploadService(uploadRepo, s3Storage, producer, idGen, clock)
	uploadHandler := httpUpload.NewHandler(initUploadService, completeUploadService)

	markChunkService := chunkSvc.NewMarkChunkService(chunkRepo, idGen, clock)
	chunkHandler := httpChunk.NewHandler(markChunkService)

	mux:= http.NewServeMux()
	httpUpload.RegisterRoutes(mux, uploadHandler)
	httpChunk.RegisterRoutes(mux, chunkHandler)
	http.ListenAndServe(":8080",mux)
}