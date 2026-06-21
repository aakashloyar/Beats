package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/beats/track/config"
	httpalbum "github.com/aakashloyar/beats/track/internal/adapter/in/http/album"
	httpartist "github.com/aakashloyar/beats/track/internal/adapter/in/http/artist"
	httptrack "github.com/aakashloyar/beats/track/internal/adapter/in/http/track"
	kafkaConsumer "github.com/aakashloyar/beats/track/internal/adapter/in/kafka_consumer"
	postgres "github.com/aakashloyar/beats/track/internal/adapter/out/postgres"
	"github.com/aakashloyar/beats/track/internal/application/ports/out/system"
	albumsvc "github.com/aakashloyar/beats/track/internal/application/service/album"
	artistsvc "github.com/aakashloyar/beats/track/internal/application/service/artist"
	tracksvc "github.com/aakashloyar/beats/track/internal/application/service/track"
)

func main() {

	ctx := context.Background()

	port, err := strconv.Atoi(config.App.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid POSTGRES_PORT: %v", err)
	}

	dbConfig := postgres.Config{
		Host:     config.App.Postgres.Host,
		Port:     port,
		User:     config.App.Postgres.User,
		Password: config.App.Postgres.Password,
		DBName:   config.App.Postgres.DBName,
		SSLMode:  config.App.Postgres.SSLMode,
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	db, err := dbConfig.NewDB()

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	//all repositories
	trackRepo := postgres.NewTrackRepository(db)
	artistRepo := postgres.NewArtistRepository(db)
	albumRepo := postgres.NewAlbumRepository(db)

	//track services
	createtrackService := tracksvc.NewCreateTrackService(trackRepo, idGen, clock)
	gettrackService := tracksvc.NewGetTrackService(trackRepo)
	listtracksService := tracksvc.NewListTracksService(trackRepo)

	//track handler
	trackHandler := httptrack.NewHandler(createtrackService, gettrackService, listtracksService)

	//artist services
	createartistService := artistsvc.NewCreateTrackService(artistRepo, idGen, clock)
	getartistService := artistsvc.NewGetArtistService(artistRepo)

	//artist handler
	artistHandler := httpartist.NewHandler(createartistService, getartistService)

	//album services
	createablumService := albumsvc.NewCreateAlbumService(albumRepo, idGen, clock)
	getalbumService := albumsvc.NewGetAlbumService(albumRepo)
	listalbumsService := albumsvc.NewListAlbumsService(albumRepo)

	//album handler
	albumHandler := httpalbum.NewHandler(createablumService, getalbumService, listalbumsService)

	kafkaConsumerConfig := kafkaConsumer.Config{
		Brokers:  config.App.Kafka.Brokers,
		Topic:    config.App.Kafka.Topic,
		ClientID: config.App.Kafka.ClientID,
		GroupID:  config.App.Kafka.GroupID,
	}

	consumer, err := kafkaConsumerConfig.NewConsumer(createtrackService)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start(ctx)
	defer consumer.Close()

	mux := http.NewServeMux()
	httptrack.RegisterRoutes(mux, trackHandler)
	httpartist.RegisterRoutes(mux, artistHandler)
	httpalbum.RegisterRoutes(mux, albumHandler)

	http.ListenAndServe(":8080", mux)
}
