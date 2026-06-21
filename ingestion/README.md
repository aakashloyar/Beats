# 📥 Ingestion Service

The **Ingestion Service** is responsible for receiving and validating audio uploads from artists.

Its job ends when the original master audio file is safely stored and the **Encoding Service** is triggered.

---

# Responsibilities

The Ingestion Service handles:

- Upload initialization
- Chunked upload coordination
- Temporary upload tracking
- Upload completion verification
- Master file reconstruction
- File validation
- Permanent storage
- Triggering encoding pipeline

---

# What is Uploaded?

Artists upload the **master audio file**.

Supported formats:

```text
master.wav
master.flac
```

---

## `master.wav`

Original uncompressed studio-quality audio.

Characteristics:

- Highest quality
- Very large file size
- Source-of-truth recording

Example:

```text
5 min WAV → 50–60 MB
```

---

## `master.flac`

FLAC = **Free Lossless Audio Codec**

It compresses audio **without losing quality**.

How?

- Detects repeating patterns
- Stores them efficiently

Important:

```text
WAV and FLAC have same audio quality
FLAC is just smaller in size
```

---

# High-Level Ingestion Flow

```text
Artist Upload
      ↓
Client splits into chunks
      ↓
Chunks uploaded
      ↓
Temporary storage
      ↓
Upload verification
      ↓
Chunk merge
      ↓
Master file validation
      ↓
Permanent storage
      ↓
Trigger Encoding Service
      ↓
Track marked as ingested
```

---

# Detailed Lifecycle

---

## 1. Artist Starts Upload

Artist selects:

```text
master.wav
or
master.flac
```

This file becomes the **source of truth**.

All future audio variants will be generated from this file.

---

## 2. Client Splits File into Chunks

This is done on the frontend.

---

### Why chunking?

Without chunking:

```text
10 MB upload fails at 95%
→ restart whole upload
```

With chunking:

```text
10 chunks uploaded
chunk 7 fails
→ resend only chunk 7
```

Much more reliable.

---

### Chunk Payload

Each chunk contains:

```json
{
  "upload_id": "upl_123",
  "chunk_number": 1,
  "raw_bytes": "..."
}
```

---

## 3. Chunks Are Uploaded

The ingestion pipeline receives chunk data.

Possible API:

```http
POST /ingestion/upload-chunk
```

---

## 4. Temporary Chunk Storage

Chunks are stored temporarily.

Possible storage:

- Local disk
- Object storage (S3)

Status:

```text
incomplete
```

---

### Why temporary?

Upload may:

- Fail
- Be cancelled
- Timeout

We should not permanently store incomplete uploads.

---

## 5. Upload Completion Verification

Before proceeding, service checks:

- Are all chunks present?
- Are chunk numbers continuous?
- Does total size match expected size?

Result:

```text
Reject or Continue
```

---

## 6. Chunk Merge

All chunks are read in order and appended.

Example:

```text
chunk_1 + chunk_2 + chunk_3
```

Result:

```text
master.flac
```

---

## 7. Master File Validation

Validate reconstructed file.

Checks include:

- File type valid?
- Duration readable?
- Audio headers valid?
- File corrupted?

Result:

```text
Reject or Continue
```

---

## 8. Permanent Storage

Store validated master file in long-term storage.

Example:

```text
Private S3 bucket
```

Important:

```text
NOT CDN
NOT public
```

This is the **source of truth**.

---

## 9. Trigger Encoding Service

Once stored successfully:

Ingestion Service publishes an event to a message queue.

Example:

```text
Kafka
RabbitMQ
NATS
```

Encoding Service consumes this event.

---

## 10. Mark Track as Ingested

Track status updated.

Example:

```sql
track.status = 'ingested'
```

---

# Upload API Design

---

# Step 1: Initialize Upload

Client first requests upload session.

### Endpoint

```http
POST /ingestion/init-upload
```

---

## Request

```go
type InitUploadRequest struct {
    ArtistIDs []string  `json:"artist_ids"`
    FileName  string    `json:"file_name"`
    FileSize  int64     `json:"file_size"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

## Response

```json
{
  "upload_id": "upl_123",
  "max_chunk_size": 5242880,
  "chunks": [
    {
      "part_number": 1,
      "url": "https://signed-upload-url"
    }
  ]
}
```

---

# Upload Strategies

There are two common ways to upload.

---

# 1. Server Upload

Architecture:

```text
Client → Ingestion Service → Storage
```

---

## How it works

Client calls:

```text
/init-upload
```

Then sends every chunk to backend.

Backend:

- Receives chunk
- Validates chunk
- Uploads to storage

---

## Pros

### Full control

Backend can inspect every upload immediately.

---

### Simpler frontend

Client only talks to backend.

---

## Cons

### Scalability issue

Example:

```text
10,000 users uploading simultaneously
```

Backend can become overloaded.

---

### Double bandwidth cost

Data travels twice:

```text
Client → Backend
Backend → Storage
```

---

### Higher latency

Extra network hop.

---

## Best for

- Small apps
- Internal tools
- Low traffic systems

---

# 2. Direct Upload (Preferred)

Architecture:

```text
Client → Storage (S3)
Client → Backend (metadata only)
```

---

## How it works

### Step 1

Client calls:

```http
POST /ingestion/init-upload
```

---

### Step 2

Backend returns:

```text
Signed URLs
```

---

### Step 3

Client uploads directly to:

```text
S3 / Object Storage
```

---

### Step 4

Client notifies backend after upload.

---

## Pros

### Highly scalable

Storage handles upload traffic.

---

### Cost efficient

No double transfer.

---

### Faster

Direct path to storage.

---

## Cons

### More frontend complexity

Client must handle:

- Chunk uploads
- Retry logic
- Signed URLs

---

### Less immediate control

Backend cannot validate before upload.

Validation happens later.

---

### Security complexity

Need to manage:

- Signed URL expiry
- Bucket permissions
- Access scope

---

# Final Architecture (Hybrid Approach)

We will combine both approaches.

Architecture:

```text
1. init-upload        → Backend
2. upload chunks      → S3 (direct)
3. complete-upload    → Backend
4. backend validates
5. trigger encoding pipeline
```

This gives:

- Scalability of direct upload
- Control of backend validation

Best of both worlds.

---

# Final Upload Lifecycle

---

## Step 1: Client initializes upload

### API

```http
POST /init
```

---

## Request

```go
type InitUploadRequest struct {
    ArtistIDs []string  `json:"artist_ids"`
    FileName  string    `json:"file_name"`
    FileSize  int64     `json:"file_size"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

## Response

```json
{
  "upload_id": "upl_123",
  "max_chunk_size": 5242880,
  "chunks": [
    {
      "part_number": 1,
      "url": "..."
    }
  ]
}
```

---

## Step 2: Client uploads chunks directly to storage

```text
PUT chunk → Signed S3 URL
```

---

## Step 3: Storage returns ETag

Example:

```text
ETag = abc123
```

ETag identifies uploaded chunk.

---

## Step 4: Client marks chunk uploaded

### API

```http
POST /mark-chunk-uploaded
```

---

## Request

```json
{
  "upload_id": "upl_123",
  "chunk_number": 1,
  "etag": "abc123"
}
```

---

## Response

```json
{}
```

---

## Step 5: Repeat for all chunks

Upload all parts.

---

## Step 6: Complete Upload

### API

```http
POST /upload-complete
```

---

## Request

```json
{
  "upload_id": "upl_123"
}
```

---

## Response

```json
{}
```

---

# Summary

The Ingestion Service ensures:

- Reliable chunked uploads
- Upload resumability
- Validation of master files
- Secure long-term storage
- Event-driven encoding trigger

It is the entry point of the BEATS audio processing pipeline.