# 🎵 Track Service

The **Track Service** is responsible for managing all metadata related to songs and their playable audio variants.

It acts as the **source of truth for music metadata** and provides information required by:

- Playback Service
- Discovery Service
- Recommendation Service
- Playlist Service

---

# Responsibilities

Track Service manages:

- Track metadata
- Artist metadata
- Album metadata
- Audio variant metadata
- Streaming file references

---

# Database Design

---

# Table 1: `tracks`

Stores the primary metadata of each song.

## Schema

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary identifier |
| `title` | TEXT | Song title |
| `artist_id` | UUID | FK → artists |
| `album_id` | UUID (nullable) | FK → albums |
| `cover_image_url` | TEXT (nullable) | Used when no album exists |
| `duration_ms` | BIGINT | Logical song duration |
| `language` | ENUM | Primary language |
| `release_date` | DATE | Song release date |
| `created_at` | TIMESTAMPTZ | Record creation timestamp |

---

## Field Explanations

---

### 1. `id`

Primary key of the track.

We use **UUID** instead of auto-increment integers.

### Why not auto-increment?

Problems in distributed systems:

- Hard to shard
- Database coordination needed
- Migration complexity
- Predictable IDs

---

### Why UUID?

UUID = **Universal Unique Identifier**

Properties:

- **Universal** → works across all systems
- **Unique** → extremely low collision probability
- **Identifier** → uniquely identifies a row/object

Uniqueness comes from combinations of:

- Timestamp
- Randomness
- Machine information

---

### UUID v4 vs UUID v7

#### UUID v4

```text
Pure random
```

Pros:

- Simple
- Highly unique

Cons:

- Poor index locality
- Random insert positions in B-tree

---

#### UUID v7 (Preferred)

```text
Timestamp + randomness
```

Pros:

- Time ordered
- Better indexing
- Faster inserts
- Better storage locality

Recommended for this project.

---

### 2. `title`

Human-readable track name.

Example:

```text
Shape of You
```

---

### 3. `artist_id`

Foreign key to:

```sql
artists.id
```

Must also be UUID.

---

### 4. `album_id`

Foreign key to:

```sql
albums.id
```

Nullable because not every song belongs to an album.

Example:

- Singles
- Independent releases
- Promotional tracks

---

### 5. `cover_image_url`

Fallback image URL when `album_id` is NULL.

Example:

```text
https://cdn.beats.com/covers/song123.jpg
```

---

### 6. `duration_ms`

Stores song duration in **milliseconds**.

Example:

```text
233456 ms
```

Why milliseconds?

- No floating-point precision issues
- Easier seeking calculations
- Better chunk calculations during playback

---

### 7. `language`

Currently storing **one dominant language only**.

Although some songs are multilingual, usually one language dominates.

Using ENUM for strict validation.

Supported values:

| Code | Language |
|------|----------|
| `hi` | Hindi |
| `en` | English |
| `hr` | Haryanvi |
| `pn` | Punjabi |

---

### 8. `release_date`

Stored as:

```sql
DATE
```

Only the date matters.

No need to store time.

---

### 9. `created_at`

Stored as:

```sql
TIMESTAMPTZ
```

Meaning:

```text
Timestamp with timezone
```

Stored internally in **UTC**.

Example:

```text
UTC = Coordinated Universal Time
IST = UTC + 5:30
```

---

# Table 2: `artists`

Stores artist metadata.

## Schema

| Column | Type |
|--------|------|
| `id` | UUID |
| `name` | TEXT |
| `bio` | TEXT |
| `profile_image_url` | TEXT |
| `created_at` | TIMESTAMPTZ |

---

## Important Notes

### `profile_image_url`

Stored in **CDN**.

CDN = **Content Delivery Network**

Why not database?

Because databases are not optimized for:

- Large binary files
- Fast global image delivery
- Caching

CDN provides:

- Low latency
- Global caching
- Reduced backend load

---

# Table 3: `albums`

Stores album metadata.

## Schema

| Column | Type |
|--------|------|
| `id` | UUID |
| `title` | TEXT |
| `cover_image_url` | TEXT |
| `release_date` | DATE |
| `created_at` | TIMESTAMPTZ |

---

# Table 4: `audio_variants`

Stores playable encoded versions of tracks.

Each track can have multiple variants.

Example:

```text
Track: Shape of You
Variants:
- OGG 96 kbps
- OGG 160 kbps
- OGG 320 kbps
```

---

## Schema

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

# Audio Variant Details

---

## 1. `id`

UUID primary key.

---

## 2. `track_id`

Foreign key to:

```sql
tracks.id
```

---

## 3. `codec`

Codec = **Compressor + Decompressor**

Raw audio is extremely large.

Example:

```text
5-minute WAV file → 50–60 MB
```

Codec compresses during storage and decompresses during playback.

---

### Supported codecs

#### MP3

- Oldest
- Widest device support

---

#### AAC

- Successor of MP3
- Better compression
- Common on iOS devices

---

#### OGG (Preferred)

Current choice.

Important:

```text
OGG = container
Opus = actual codec inside
```

Benefits:

- Excellent quality
- Better compression
- Open standard

---

## 4. `bitrate_kbps`

Controls audio quality.

Current supported bitrates:

```text
96 kbps
160 kbps
320 kbps
```

Meaning:

```text
Bits consumed per second during playback
```

Note:

```text
1 kbps = 1000 bits/sec
```

Not bytes.

---

## 5. `sample_rate_hz`

Current value:

```text
44,100 Hz
```

Represents how many audio samples are captured per second.

Human hearing range:

```text
20 Hz → 20,000 Hz
```

According to **Nyquist theorem**:

```text
Sample rate ≥ 2 × max audible frequency
```

Hence:

```text
2 × 20,000 = 40,000+
```

So 44.1 kHz is standard.

---

## 6. `channels`

Current:

```text
2
```

Possible values:

- `1` → Mono
- `2` → Stereo

Stereo creates spatial sound.

Example:

- Left channel → guitar louder
- Right channel → vocals louder

Gives more natural listening.

---

## 7. `duration_ms`

Different from `tracks.duration_ms`.

### `tracks.duration_ms`

Logical duration shown in:

- UI
- Metadata
- Search results

---

### `audio_variants.duration_ms`

Actual playable duration of encoded file.

Different codecs introduce tiny differences.

Example:

```text
Track metadata: 233 sec
MP3: 233.01 sec
OGG: 232.98 sec
AAC: 233.05 sec
```

Used during:

- Playback
- Seeking
- Buffer calculations

---

## 8. `file_url`

Stores streaming entry point.

We are using:

# Segment-Based Streaming ✅

---

# Streaming Approaches

---

## 1. Byte-Range Streaming

Traditional streaming method.

Client requests chunks of the same file.

Example:

```text
GET bytes 0-200KB
GET bytes 200-400KB
```

HTTP supports:

```text
Range Requests
```

Advantages:

- Easy implementation
- No preprocessing

Limitations:

- Harder adaptive bitrate switching
- Less CDN-friendly
- Resume/seek handling less flexible

---

## 2. Segment-Based Streaming (Chosen)

Audio is pre-cut into small segments.

Example:

```text
segment_001.ts
segment_002.ts
segment_003.ts
```

Each segment:

```text
2–6 seconds
```

Independent files stored in CDN.

---

### `file_url` points to Manifest File

It does **not** point directly to the audio.

Example:

```text
https://cdn.beats.com/audio/track123/playlist.m3u8
```

---

## What is a Manifest File?

A small text file containing:

- Total duration
- Codec
- Bitrate
- Segment locations

Example:

```text
#EXTM3U
segment1.ts
segment2.ts
segment3.ts
```

Benefits:

- Better buffering
- Adaptive streaming
- Easy bitrate switching
- Excellent CDN caching
- Better seek/resume support

---

## 9. `created_at`

Stored as:

```sql
TIMESTAMPTZ
```

---

# Track Service Request Flow

Application follows:

```text
Route
   ↓
Handler
   ↓
Service
   ↓
Repository
   ↓
Database
   ↓
Repository
   ↓
Service
   ↓
Handler
   ↓
Response
```

---

# Summary

Track Service is responsible for:

- Managing song metadata
- Managing artist and album metadata
- Managing playable audio variants
- Storing streaming manifest URLs
- Serving metadata to playback and discovery systems

It acts as the central metadata authority for the BEATS platform.