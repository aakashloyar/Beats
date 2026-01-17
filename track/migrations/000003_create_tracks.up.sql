CREATE TYPE language_enum AS ENUM('hi','en','hr','pn')

CREATE TABLE TRACKS {
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    artist_id UUID NOT NULL,
    album_id UUID,
    cover_image_url TEXT,
    duration BIGINT NOT NULL,
    language language_enum NOT NULL,
    release_date DATE, 
    created_at TIMESTAMPTZ NOT NULL NOT(), 
    CONSTRAINT fk_artist
      FOREIGN KEY(artist_id)
      REFERENCES artists(id),
    CONSTRAINT fk_album
      FOREIGN KEY (album_id)
      REFERENCES albums(id)    
};
