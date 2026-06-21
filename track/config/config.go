package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	ClientID string
	GroupID  string
}

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
}

func load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env", err)
	}

	postgres := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	kafka := KafkaConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		ClientID: os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:  os.Getenv("KAFKA_GROUP_ID"),
	}

	return Config{
		Postgres: postgres,
		Kafka:    kafka,
	}
}

var App = load()
