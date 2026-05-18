# 🎵 BEATS Backend

A scalable **music streaming backend system** inspired by platforms like Spotify, JioSaavn, and Apple Music.

This project is designed using **Microservices + Hexagonal Architecture (Ports & Adapters)** with **Go (Golang)**.

---

# 📌 Stage 1 Scope

We are currently focusing on the core backend services required to upload, process, store, and stream songs.

### Services in Stage 1

1. **Track Service**
2. **Ingestion Service**
3. **Encoding Service**
4. **Playback Service**

---

# 🏗️ Overall Architecture

```text
Artist Upload
     |
     v
Ingestion Service
     |
     v
Encoding Service
     |
     v
Track Service
     |
     v
Playback Service
     |
     v
Client (Web / Mobile)
```

---

# 📁 Monorepo Structure

```text
BEATS-backend/
├── track-service/
├── playback-service/
├── user-service/
├── auth-service/
├── playlist-service/
├── recommendation-service/
├── discovery-service/
├── ingestion-service/
├── encoding-service/
├── analytics-service/
├── notification-service/
├── gateway-service/
├── proto/                 # Shared gRPC contracts
├── libs/                  # Shared libraries (auth, logging, utils)
├── deploy/                # Kubernetes, Terraform, Helm configs
└── README.md
```

---

# 🧱 Architectural Style

We follow **Hexagonal Architecture (Ports & Adapters)**.

### Why Hexagonal?

- Business logic stays independent.
- Easy to swap database.
- Easy to change transport (HTTP → gRPC).
- Better testability.
- Clear separation of concerns.

---

## Hexagonal Structure Example

```text
track-service/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── track.go
│   │   ├── artist.go
│   │   └── album.go
│   │
│   ├── application/
│   │   ├── ports/
│   │   │   ├── in/
│   │   │   └── out/
│   │   └── service/
│   │
│   ├── adapters/
│   │   ├── in/
│   │   │   ├── http/
│   │   │   ├── grpc/
│   │   │   └── kafka-consumer/
│   │   │
│   │   └── out/
│   │       ├── postgres/
│   │       ├── redis/
│   │       ├── kafka-producer/
│   │       └── s3-cdn/
│   │
│   ├── config/
│   └── observability/
│
├── migrations/
├── go.mod
└── go.sum
```

---

# 🎼 Service 1: Track Service

Track service manages all song metadata.

---

## Database Design

---

## Table: `tracks`

```sql
tracks
```

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID | Primary key |
| `title` | TEXT | Song title |
| `artist_id` | UUID | FK → artists |
| `album_id` | UUID (nullable) | FK → albums |
| `cover_image_url` | TEXT | Used if album is null |
| `duration_ms` | BIGINT | Duration in milliseconds |
| `language` | ENUM | `hi`, `en`, `hr`, `pn` |
| `release_date` | DATE | Date only |
| `created_at` | TIMESTAMPTZ | UTC timestamp |

---

### Why UUID?

Why not auto-increment?

Problems with auto-increment:

- Hard in distributed systems
- Sharding issues
- Migration pain
- Predictable IDs

UUID solves this.

### UUID v4 vs UUID v7

#### UUID v4

- Random only
- Collision probability near zero

#### UUID v7

- Timestamp + randomness
- Better indexing
- Insert locality
- Recommended

---

## Table: `artists`

| Column | Type |
|--------|------|
| `id` | UUID |
| `name` | TEXT |
| `bio` | TEXT |
| `profile_image_url` | TEXT |
| `created_at` | TIMESTAMPTZ |

### Notes

Profile images are stored in **CDN**.

Why not database?

- Databases are bad for large binary files
- CDN provides caching and fast delivery

---

## Table: `albums`

| Column | Type |
|--------|------|
| `id` | UUID |
| `title` | TEXT |
| `cover_image_url` | TEXT |
| `release_date` | DATE |
| `created_at` | TIMESTAMPTZ |

---

## Table: `audio_variants`

Stores playable versions of a track.

| Column | Type |
|--------|------|
| `id` | UUID |
| `track_id` | UUID |
| `codec` | ENUM |
| `bitrate_kbps` | INT |
| `sample_rate_hz` | INT |
| `channels` | INT |
| `duration_ms` | BIGINT |
| `file_url` | TEXT |
| `created_at` | TIMESTAMPTZ |

---

### Codec

Codec = **Compressor + Decompressor**

Raw audio is huge.

Example:

```text
5 min WAV file ≈ 50-60 MB
```

Supported codecs:

- **MP3**
- **AAC**
- **OGG (Opus)** ← preferred

---

### Bitrate

Current options:

```text
96 kbps
160 kbps
320 kbps
```

Meaning:

```text
Bits consumed per second during playback
```

---

### Sample Rate

Current:

```text
44,100 Hz
```

Reason:

Human hearing:

```text
20 Hz → 20,000 Hz
```

Nyquist theorem:

```text
Sample rate ≥ 2 × max frequency
```

---

### Channels

Possible values:

- `1` → Mono
- `2` → Stereo

Stereo gives better spatial experience.

Example:

- Left ear → guitar louder
- Right ear → vocals louder

---

### Why `track.duration_ms` and `audio_variant.duration_ms` both?

Because codecs slightly alter actual duration.

Example:

```text
Track metadata: 233 sec
MP3: 233.01 sec
OGG: 232.98 sec
AAC: 233.05 sec
```

Track duration = logical duration shown in UI.

Variant duration = exact playback duration.

---

## Streaming Model

Two main approaches:

---

### 1. Byte-Range Streaming

```text
GET bytes 0-200KB
GET bytes 200-400KB
```

Advantages:

- Simpler
- No preprocessing
- HTTP Range Requests

---

### 2. Segment-Based Streaming (Preferred)

Audio is pre-cut into segments.

```text
segment_1.ts
segment_2.ts
segment_3.ts
```

Each segment:

```text
2–6 seconds
```

`file_url` points to **manifest file**, not actual audio.

Example:

```text
playlist.m3u8
```

Manifest contains:

- Duration
- Codec
- Bitrate
- Segment URLs

Benefits:

- Adaptive bitrate streaming
- Better buffering
- Better CDN caching

---

### 3. True Real-Time Streaming

Protocols:

- WebRTC
- RTMP

Used for:

- Zoom
- Google Meet

Not suitable for music streaming.

---

# 📥 Service 2: Ingestion Service

Responsible for accepting uploaded master audio.

---

## Upload Flow

```text
1. Artist uploads master file
2. Client splits into chunks
3. Ingestion service receives chunks
4. Temporary storage
5. Completion verification
6. Chunk merge
7. Validation
8. Permanent storage
9. Trigger encoding
10. Mark ingested
```

---

## Step-by-Step

---

### 1. Artist Uploads Master File

Formats:

```text
master.wav
master.flac
```

These are source-of-truth files.

---

### 2. Client Splits into Chunks

Done on frontend.

Why?

Without chunking:

```text
10MB upload fails → restart whole upload
```

With chunking:

```text
chunk_7 failed → resend only chunk_7
```

Chunk payload:

```json
{
  "upload_id": "...",
  "chunk_number": 7,
  "raw_bytes": "..."
}
```

---

### 3. Upload Chunk Endpoint

```http
POST /ingestion/upload-chunk
```

---

### 4. Temporary Chunk Storage

Possible storage:

- Local disk
- S3 temporary bucket

Status:

```text
incomplete
```

---

### 5. Completion Verification

Checks:

- All chunks exist
- Chunks continuous
- Total size matches

---

### 6. Chunk Merge

```text
chunk_1 + chunk_2 + ... + chunk_n
```

Result:

```text
master.flac
```

---

### 7. Master File Validation

Checks:

- File type
- Audio headers
- Readable duration
- Corruption

---

### 8. Permanent Storage

Store master file in:

```text
Private Object Storage
```

Examples:

- Amazon S3
- GCS
- MinIO

Not public CDN.

---

### 9. Trigger Encoding Service

Use message queue.

Examples:

- Kafka
- RabbitMQ
- NATS

---

### 10. Mark Ingestion Complete

```text
track.status = ingested
```

---

# 🎛️ Service 3: Encoding Service

Converts master file into multiple streaming variants.

---

## Responsibilities

Generate:

- OGG 96 kbps
- OGG 160 kbps
- OGG 320 kbps

Future:

- AAC
- MP3

---

## Segmenting

Creates:

```text
playlist.m3u8
segment_001.ts
segment_002.ts
```

Uploads to CDN.

---

# ▶️ Service 4: Playback Service

Responsible for serving songs to listeners.

---

## Responsibilities

- Fetch available variants
- Select bitrate
- Return manifest URL
- Track playback session
- Authorization checks
- Handle buffering retries

---

# 🔧 Tech Stack

### Language

- Go (Golang)

### Communication

- HTTP
- gRPC
- Kafka events

### Storage

- PostgreSQL
- Redis
- S3 / Object Storage
- CDN

### Deployment

- Docker
- Kubernetes
- Terraform
- Helm

### Observability

- Prometheus
- Grafana
- OpenTelemetry

---

# 🚀 Future Services

Planned for later stages:

- User Service
- Auth Service
- Playlist Service
- Recommendation Service
- Discovery Service
- Analytics Service
- Notification Service
- Gateway Service

---

# Goal

Build a production-grade music streaming backend that is:

- Scalable
- Fault tolerant
- Distributed
- Observable
- Cloud native
- Easy to extend

---