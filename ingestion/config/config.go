package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)
type UploadConfig struct {
	PresignExpirty   int64 
	MaxChunkSize     int64
}

type S3Config struct {
	Region string
	Bucket string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	S3       S3Config
	Upload   UploadConfig
	Postgres PostgresConfig    
}

func load() Config{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env")
	}

	upload := UploadConfig{
		MaxChunkSize: 5 * 1024 * 1024, // 5MB
		PresignExpirty: 15, //15min
	}

	s3 := S3Config {
		Region: os.Getenv("AWS_REGION"),
		Bucket: os.Getenv("AWS_BUCKET"),
	}

	postgres := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DBNAME"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	return Config{
		S3: s3,
		Upload: upload,
		Postgres: postgres,
	}
}

var App = load()