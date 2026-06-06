package main

import (
	"context"
	"log"
	"net/http"
	"github.com/aakashloyar/beats/ingestion/config"

	httpChunk "github.com/aakashloyar/beats/ingestion/internal/adapter/in/http/chunk"
	httpUpload "github.com/aakashloyar/beats/ingestion/internal/adapter/in/http/upload"
	kafkaProducer "github.com/aakashloyar/beats/ingestion/internal/adapter/out/kafka-producer"
	"github.com/aakashloyar/beats/ingestion/internal/adapter/out/postgres"
	"github.com/aakashloyar/beats/ingestion/internal/adapter/out/s3"
	"github.com/aakashloyar/beats/ingestion/internal/application/ports/out/system"
	chunkSvc "github.com/aakashloyar/beats/ingestion/internal/application/service/chunk"
	uploadSvc "github.com/aakashloyar/beats/ingestion/internal/application/service/upload"
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

	s3Config := s3.Config{
		Region: config.App.S3.Region,
		Bucket: config.App.S3.Bucket,
	}

	s3Client, err := s3Config.NewS3Client(ctx)

	s3Storage :=  s3.NewS3Storage(s3Client.Client, s3Config.Bucket)

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