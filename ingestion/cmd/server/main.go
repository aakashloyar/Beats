package main

import (
	"strconv"
	"net/http"
	"log"
	"context"
	"github.com/aakashloyar/beats/ingestion/config"
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

	port, err := strconv.Atoi(config.App.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid POSTGRES_PORT: %v", err)
	}
	postgresConfig := postgres.Config{
		Host:     config.App.Postgres.Host,
		Port:     port,
		User:     config.App.Postgres.User,
		Password: config.App.Postgres.Password,
		DBName:   config.App.Postgres.DBName,
		SSLMode:  config.App.Postgres.SSLMode,
	}

	db, err := postgresConfig.NewDB()

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	log.Printf("DB connection established: %v", db)

	s3Config := s3.Config{
		Region: config.App.S3.Region,
		Bucket: config.App.S3.Bucket,
	}

	s3Client, err := s3Config.NewS3Client(ctx)

	s3Storage :=  s3.NewS3Storage(s3Client.Client, s3Config.Bucket)

	log.Printf("S3 initialized: %v", s3Client.Client)

	kafkaConfig := kafkaProducer.Config{
		Brokers:  []string{config.App.Kafka.Brokers[0]},
		Topic:    config.App.Kafka.Topic,
		ClientID: config.App.Kafka.ClientID,
	}

	producer, err := kafkaProducer.NewFranzProducer(kafkaConfig); if err != nil {
		log.Fatal("error creating Kafka producer", err)
	}

	log.Printf("Kafka producer initialized: %v", producer)

	defer producer.Close()

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	uploadRepo := postgres.NewUploadRepository(db)
	chunkRepo := postgres.NewChunkRepository(db)

	log.Printf("Repositories initialized: uploadRepo=%v, chunkRepo=%v", uploadRepo, chunkRepo)

	initUploadService := uploadSvc.NewInitUploadService(uploadRepo,s3Storage,idGen,clock)
	completeUploadService := uploadSvc.NewCompleteUploadService(uploadRepo, s3Storage, producer, idGen, clock)
	markChunkService := chunkSvc.NewMarkChunkService(chunkRepo, idGen, clock)

	log.Printf("Services initialized: initUploadService=%v, completeUploadService=%v, markChunkService=%v", initUploadService, completeUploadService, markChunkService)

	uploadHandler := httpUpload.NewHandler(initUploadService, completeUploadService)
	chunkHandler := httpChunk.NewHandler(markChunkService)

	log.Printf("Handlers initialized: uploadHandler=%v, chunkHandler=%v", uploadHandler, chunkHandler)

	mux:= http.NewServeMux()

	httpUpload.RegisterRoutes(mux, uploadHandler)
	httpChunk.RegisterRoutes(mux, chunkHandler)
	
	log.Println("Routes registered successfully. Starting server on :8080...")
	http.ListenAndServe(":8080",mux)
}