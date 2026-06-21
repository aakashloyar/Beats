package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type S3Config struct {
	Region string
	Bucket string
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	ClientID string
}

type EncodingProfile struct {
	FFmpegCodec string
	Bitrate     string
	Bandwidth   int
}

type Config struct {
	S3                            S3Config
	KafkaProducer                 KafkaConfig
	KafkaConsumer                 KafkaConfig
	UploadLocalPathForEncoding    string
	SupportedCodecs               map[string]bool
	SupportedSampleRates          map[int]bool
	Profiles                      []EncodingProfile
	OutPutDirForTranscodedFiles   string
	StoragePathForTranscodedFiles string
}

func load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env", err)
	}

	s3 := S3Config{
		Region: os.Getenv("AWS_REGION"),
		Bucket: os.Getenv("AWS_BUCKET"),
	}

	kafkaProducer := KafkaConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		ClientID: os.Getenv("KAFKA_CLIENT_ID"),
	}

	kafkaConsumer := KafkaConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		ClientID: os.Getenv("KAFKA_CLIENT_ID"),
	}

	return Config{
		S3:                         s3,
		KafkaProducer:              kafkaProducer,
		KafkaConsumer:              kafkaConsumer,
		UploadLocalPathForEncoding: os.Getenv("UPLOAD_LOCAL_PATH_FOR_ENCODING"),
		SupportedCodecs: map[string]bool{
			"mp3":       true,
			"flac":      true,
			"pcm_s16le": true,
			"aac":       true,
			"opus":      true,
			"vorbis":    true,
		},
		SupportedSampleRates: map[int]bool{
			44100: true,
			48000: true,
		},
		Profiles: []EncodingProfile{
			{
				FFmpegCodec: "aac",
				Bitrate:     "96k",
				Bandwidth:   110000,
			},
			{
				FFmpegCodec: "aac",
				Bitrate:     "128k",
				Bandwidth:   140000,
			},
			{
				FFmpegCodec: "aac",
				Bitrate:     "320k",
				Bandwidth:   340000,
			},
		},
		OutPutDirForTranscodedFiles:   "streams",
		StoragePathForTranscodedFiles: "streams",
	}
}

var App = load()

var UploadLocalPathForEncoding = App.UploadLocalPathForEncoding

var SupportedCodecs = App.SupportedCodecs

var SupportedSampleRates = App.SupportedSampleRates

var Profiles = App.Profiles

var OutPutDirForTranscodedFiles = App.OutPutDirForTranscodedFiles

var StoragePathForTranscodedFiles = App.StoragePathForTranscodedFiles
