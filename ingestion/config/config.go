package config

type UploadConfig struct {
	PresignExpirty   int64 
	MaxChunkSize     int64
}

var Upload = UploadConfig{
	MaxChunkSize: 5 * 1024 * 1024, // 5MB
	PresignExpirty: 15, //15min
}