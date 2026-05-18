module github.com/aakashloyar/beats/track

go 1.25.0

require (
	github.com/aakashloyar/beats v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.12.3
	github.com/twmb/franz-go v1.21.2
)

require (
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.13.1 // indirect
)

replace github.com/aakashloyar/beats => ../
