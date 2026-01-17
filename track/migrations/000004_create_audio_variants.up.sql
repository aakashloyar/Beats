CREATE TYPE codec_enum AS ENUM ('ogg', 'aac', 'mp3');

CREATE Table audio_varinats {
    id UUID PRIMARY KEY,
    track_id UUID NOT NULL,
    codec codec_enum NOT NULL,
    bitrate_kbps INT NOT NULL
    CHECK (bitrate_kbps IN (96, 160, 320)),
    sample_rate_hz INT NOT NULL,
    channels SMALLINT NOT NULL,
    duration_ms BIGINT NOT NULL,
    file_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT_NULL DEFAULT NOW(),
    CONSTRAINT fk_track
      FOREIGN KEY (track_id)
      REFERENCES tracks(id)
      ON DELETE CASCADE 
}