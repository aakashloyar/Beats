# 🎛️ Encoding Service

The **Encoding Service** converts uploaded master audio files into multiple optimized streaming formats.

It consumes upload events from **Kafka**, downloads the source file, validates it, transcodes it into multiple bitrate variants, generates streaming manifests, and updates metadata in the database.

---

# Responsibilities

Encoding Service handles:

- Consuming ingestion-complete events
- Downloading master audio
- Audio validation
- Metadata extraction
- Transcoding into multiple qualities
- Packaging for streaming
- Uploading encoded assets
- Updating Track Service metadata

---

# High-Level Flow

```text
Kafka Event
    ↓
Download master file
    ↓
Temporary local storage
    ↓
Decode + validate
    ↓
Transcode with FFmpeg
    ↓
Generate HLS manifests
    ↓
Upload encoded assets
    ↓
Update database
```

---

# Step 1: Consume Upload Event from Kafka

The Ingestion Service publishes an event when upload is complete.

Example event:

```json
{
  "track_id": "trk_123",
  "upload_id": "upl_456",
  "master_file_url": "s3://beats/private/master.flac"
}
```

Encoding service consumes this event.

---

# Step 2: Download Master File

The master file is downloaded from object storage.

---

## Why download locally first?

FFmpeg is:

- CPU intensive
- Disk intensive
- Performs many random reads
- Needs stable high-throughput access

Reading directly from S3 repeatedly is inefficient.

Local disk is much faster.

---

## Temporary Storage Problem

As queue grows:

```text
More jobs → more downloaded files → disk fills up
```

Potential solutions:

### Vertical Scaling

Increase:

- CPU
- RAM
- Disk

Simple but limited.

---

### Horizontal Scaling (Preferred)

Run multiple encoding workers.

Example:

```text
encoding-worker-1
encoding-worker-2
encoding-worker-3
```

Kafka distributes jobs.

Benefits:

- Better scalability
- Fault isolation
- Parallel processing

---

# Step 3: Decode File and Extract Metadata

Before transcoding, the file must be decoded and inspected.

This is done using:

```text
ffprobe
```

---

## Example ffprobe output structure

```go
type FFProbeOutput struct {
    Streams []struct {
        CodecName string `json:"codec_name"`
        CodecType string `json:"codec_type"`
        Channels  int    `json:"channels"`
        SampleRate string `json:"sample_rate"`
    } `json:"streams"`

    Format struct {
        Duration string `json:"duration"`
        BitRate  string `json:"bit_rate"`
    } `json:"format"`
}
```

---

## Why `Streams[]` is an array?

A media file can contain multiple streams.

Examples:

- Audio stream
- Video stream
- Commentary track
- Subtitle stream

For now, we support only:

```text
Single simple audio file
```

Future support may include:

- Multiple audio tracks
- Commentary streams
- Video + audio

---

# Step 4: Validate Decoded Metadata

The extracted metadata is validated.

Checks include:

- Supported codec?
- Valid duration?
- Valid sample rate?
- Valid channel count?
- Corruption?

---

## Example validation rules

```text
codec ∈ {wav, flac}
sample_rate = 44100
channels ∈ {1, 2}
duration > 0
```

---

# Step 5: Encoding (Transcoding)

Master file is converted into multiple playable variants.

Example outputs:

```text
96 kbps
160 kbps
320 kbps
```

Generated using:

```text
FFmpeg
```

---

# Streaming Packaging Strategy

We are using:

```text
HLS (HTTP Live Streaming)
```

But there are multiple ways to package HLS assets.

---

# Option 1: Traditional HLS (Physical Segment Files)

Each segment is stored as a separate file.

---

## Output Structure

```text
playlist.m3u8
seg1.ts
seg2.ts
seg3.ts
```

---

## Example playlist

```text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:6
#EXT-X-MEDIA-SEQUENCE:0

#EXTINF:6.0,
seg1.ts

#EXTINF:6.0,
seg2.ts

#EXTINF:6.0,
seg3.ts

#EXT-X-ENDLIST
```

---

## What is inside `seg.ts`?

Each `.ts` file is an independent MPEG-TS segment.

Structure:

```text
PAT (Program Association Table)
PMT (Program Map Table)
Codec info
Timestamps
Audio frames
```

Visual:

```text
seg1.ts
 ├── headers
 ├── codec metadata
 ├── timestamps
 └── audio data
```

---

## Drawback

Metadata is repeated in every segment.

Storage inefficient.

---

# Option 2: HLS Byte-Range Mode

Single `.ts` file.

---

## Output Structure

```text
playlist.m3u8
audio.ts
```

---

## Internal layout

```text
|---segment1---|---segment2---|---segment3---|
0            120k           250k          380k
```

Single file, logical segments.

---

## Example playlist

```text
#EXTM3U

#EXTINF:6.0,
#EXT-X-BYTERANGE:120000@0
audio.ts

#EXTINF:6.0,
#EXT-X-BYTERANGE:130000@120000
audio.ts
```

Meaning:

```text
130000 = length
120000 = offset
```

Read:

```text
bytes 120000 → 249999
```

---

## Playback Flow

### Step 1

```http
GET /playlist.m3u8
```

---

### Step 2

```http
GET /audio.ts
Range: bytes=0-119999
```

Response:

```text
HTTP 206 Partial Content
```

---

### Step 3

Player continues requesting ranges.

---

# Option 3: CMAF (Preferred)

CMAF = **Common Media Application Format**

Modern fragmented MP4 format.

Preferred over `.ts`.

---

# CMAF Multi-File Structure

```text
playlist.m3u8
init.mp4
chunk1.m4s
chunk2.m4s
chunk3.m4s
```

---

## Playlist Example

```text
#EXTM3U
#EXT-X-MAP:URI="init.mp4"

#EXTINF:6.0,
chunk1.m4s

#EXTINF:6.0,
chunk2.m4s
```

---

# `init.mp4`

Contains only metadata.

Structure:

```text
[ftyp][moov]
```

---

## `ftyp`

Defines file type.

Example:

```text
format = mp4
brand = iso6
compatible = CMAF
```

---

## `moov`

Contains decoder instructions.

Example:

```text
codec = AAC
sample rate = 44100
channels = 2
track id = 1
duration info
```

---

# `chunk.m4s`

Contains media fragments.

Structure:

```text
[moof][mdat]
```

---

## `moof`

Movie fragment metadata:

- Chunk timing
- Sample count
- Offsets

---

## `mdat`

Actual encoded audio bytes.

---

# CMAF Playback Flow

### Step 1

```http
GET /playlist.m3u8
```

---

### Step 2

Player sees:

```text
#EXT-X-MAP:URI="init.mp4"
```

Requests:

```http
GET /init.mp4
```

Gets:

```text
[ftyp][moov]
```

---

### Step 3

Requests:

```http
GET /chunk1.m4s
```

Gets:

```text
[moof][mdat]
```

---

### Step 4

Player combines:

```text
(init.mp4) + (chunk1.m4s)
```

Playback begins.

---

# CMAF Single-File Mode (Chosen)

Instead of multiple chunk files:

```text
playlist.m3u8
audio.mp4
```

---

## Internal structure

```text
[ftyp][moov][moof][mdat][moof][mdat][moof][mdat]
```

Visual:

```text
audio.mp4
|-- init --|-- fragment1 --|-- fragment2 --|-- fragment3 --|
```

---

## Example playlist

```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:100000@5000
audio.mp4
```

---

## Meaning

Initialization section:

```text
bytes 0 → 4999
```

Segment 1:

```text
bytes 5000 → 104999
```

---

# Playback Flow

### Step 1

```http
GET /playlist.m3u8
```

---

### Step 2

Player requests init bytes:

```http
GET /audio.mp4
Range: bytes=0-4999
```

Receives:

```text
[ftyp][moov]
```

---

### Step 3

Requests segment:

```http
GET /audio.mp4
Range: bytes=5000-104999
```

Receives:

```text
[moof][mdat]
```

---

### Step 4

Player continues requesting fragments.

---

# Why CMAF over Byte-Range `.ts`?

---

## 1. Storage Efficiency

`.ts` repeats metadata.

CMAF stores metadata once.

Less storage.

---

## 2. Better Format

`.ts`:

- Old
- Designed for TV broadcasting

`.mp4`:

- Modern
- Cleaner
- Better tooling

---

## 3. Protocol Compatibility

CMAF works with both:

- **HLS**
- **DASH**

---

## 4. Lower Latency

Supports:

- Low-Latency HLS
- Low-Latency DASH

Smaller partial chunks.

---

## 5. Better CDN Behavior

Cache efficiency is generally better.

---

# HLS vs DASH

---

## HLS

HTTP Live Streaming.

Owned by Apple.

Manifest:

```text
playlist.m3u8
```

---

## DASH

Dynamic Adaptive Streaming over HTTP.

Open standard.

Manifest:

```text
manifest.mpd
```

(XML based)

---

## Both Provide

- Available qualities
- Segment locations
- Timing
- Switching rules

---

## Platform Preference

HLS:

- Preferred on Apple devices

DASH:

- Preferred elsewhere

---

## Current Decision

For now we support only:

```text
HLS (.m3u8)
```

DASH may be added later.

---

# Adaptive Bitrate Switching

Player first downloads:

```text
master.m3u8
```

---

## Master Playlist Example

```text
#EXTM3U

#EXT-X-STREAM-INF:BANDWIDTH=96000
96k/playlist.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=160000
160k/playlist.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=320000
320k/playlist.m3u8
```

---

Player chooses bitrate based on network speed.

---

## Media Playlist Example

```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:40000@5000
audio.mp4
```

---

# Final Output Structure

```text
outDir/
   master.m3u8
   96k/
      playlist.m3u8
      audio.mp4
   160k/
      playlist.m3u8
      audio.mp4
   320k/
      playlist.m3u8
      audio.mp4
```

---

# Final Responsibilities Summary

Encoding Service is responsible for:

- Consuming ingestion events
- Downloading source audio
- Validating media metadata
- Transcoding with FFmpeg
- Packaging CMAF/HLS assets
- Uploading streamable outputs
- Updating Track Service metadata

It transforms raw uploaded audio into production-ready streaming assets.