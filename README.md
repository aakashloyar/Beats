# Beats

# Functional requirements


Now let us look at some services

# *** Services ***

1. Track Service
2. Playback Service
3. User 
4. Auth
5. Playlist
6. Recommendation
7. Discovery
8. Ingestion (upload)
9. Encoding

*** Services ***



# *** Track Service ***

let us look at the database design of this service 


*** Tables ***



*** Table1 -> tracks ***
1. id
2. title
3. artist_id
4. album_id
5. cover_image_url
6. duration_ms
7. languages
8. release_date
9. created_at

Important Points

1. id
-> So this will be uuid
-> Why not autoincrement-> not good for distributed systems / also autoincrement is painful for migration
-> UUID -> Universal Unique Identifier
-> Universal-> works everywhere
-> Unique-> always unique
-> Identifire-> Used to identify any row or object

-> Why it is always unique/ Chance of Collision is zero
-> Because-> it depends on -> Time/ randomness/ Machineinfo

-> V4 vs V7 
-> V4 -> only random numbers
-> V7 -> Timestamp vs randomness
-> in V7 indexing also good 
-> as it uses timestamp so arranged according to that only



2. title



3. artist_id
-> It must be uuid only

4. album_id
-> uuid only
-> from here you will get your image url


5. cover_image_url
-> some song does not belong to album
-> so album_id is also nullable
-> in those cases you will get image from here

6. duration_ms
-> Must be in microsecond for a nice chunk size or not dealing with float


7. languages
-> now here come 1 main things
-> a song can be in english and hindi both 
-> so we are not looking at this case currently bcz this is only 1% chance
-> if there is mix also then most of the cases have one language dominating
-> so now we are considering to only 1 language
-> now language must be enum 
-> here we are giving 4 options
-> hindi-> hi 
-> english-> en
-> haryanvi-> hr
-> punjabi-> pn


8. release_date
-> stored as date only not as time
-> bcz you donot require any time 
-> just date is fine here


9. created_at
-> date with time 
-> TIMESTAMPTZ = timestamp with timezone
-> stored as UTC 
-> UTC = Coordinated Universal Time
-> IST (UTC+5:30)



*** Table1 end ***


*** Table2 -> artists  ***
1. id
2. name
3. bio
4. profile_image_url
5. created_at

1. id 
-> uuid

2. name

3. bio
-> about them

4. profile_image_url
-> stored in CDN
-> CDN -> content delivery network
-> why cdn 
-> databases are not designed for images

5. created_at
-> TIMESTAMPZ


*** Table2 end ***


*** Table3 -> albums ***
1. id
2. title
3. cover_image_url
4. release_date
5. created_at

1. id 
-> uuid

2. title

3. cover_image_url
-> stored in CDN

4. release_date 
-> date

5. created_at
-> TIMESTAMPZ

***  Table3 end  ***
*** Tables ***

*** Table4 -> audio_variants ***
1. id
2. track_id
3. codec
4. bitrate_kbps
5. sample_rate_hz
6. channels
7. duration_ms
8. file_url
9. created_at

1. id
-> uuid

2. track_id
-> uuid

3. codec 
*  currently moving with ogg
-> co+dec -> compressor + decompressor
-> raw audio is huge
-> 5 min song (raw WAV) 50-60MB
-> it is like a way to compress and decompress song
-> compresses during audio storage and decompresses during playback
-> currently we are considering only ogg
-> but we are planning to considering 3 codecs ogg, aac, mp3
-> mp3 is the oldest
-> most devices support mp3
-> aac is successor of mp3
-> it is mostly used for ios
-> ogg  with opus (ogg is contaniner opus is codec inside it)

4. bitrates
*  currently using only 3 (96| 160| 320)kb
-> how many bits do we need to play this song per sec
-> KB -> kilo byte
-> earlier 1KB -> 1024 Kilo Bytes
-> now 1KB = 1000 Kilo Bytes
-> 1kbps = 1000 bits/sec not bytes per seconds

5. sample_rate_hz
*  44,100Hz
-> refresh rate 
-> sample_rate must be 2*x of frequency you listen 
-> so 20Hz-20KHz  20KHz*2

6. channels
*  currently 2
-> in song there are 2 channels 
-> in other there can be 1 or 2 
-> we have left and right ear
-> so left-> guitar slightly louder
-> right -> singer slightly louder
-> this will give 3d feel more natural
-> mono(1) or stereo(2)

7. duration
-> track.duration -> logical duration of that song -> shown everywhere in the app (UI and metadata)
-> audio_varaiants.duration-> actual playable duration of that file -> used in playback and streaming
-> Track: Shape of You 
-> Duration: 233 seconds
-> MP3 version → 233.01 sec
-> OGG version → 232.98 sec
-> AAC version → 233.05 sec
-> this is due to codec behaviour


8. file_url
*  moving with segment bases streaming
-> this is the main part of our application 
-> in real system there is only two type of streaming
-> byte range streaming / segment based streaming

***  Byte range streaming ***
-> byte range streaming -> in this file is actual audio file
-> how this happens -> give me chunk from 0-200kb-> then 200-400kb
-> byte range = chunks
-> chunk stored in memory 
-> https supports range requests
-> no need to precut files

***  Segment Based streaming ***
-> in this audio is pre-cut into smaller pieces
-> each segment -> 2-6sec, independent file, stored in cdn
-> file_url does not point to audio
-> it points to manifest file

-> what is a manifest file
-> it is a text file 
-> contains -> song duration, codec, bitrate, where segments are



9. created_at
-> TIMESTAMPZ

Important points


*** Track Service ***


# *** Ingestion Service ***

-> This is service which take care of uploading the song
-> master.wav / master.flac | This is the main file that is uploaded by artist
-> master.wav is the original recorded file
-> master.flac is the compressed file with same quality || compressed by finding the same pattern
-> these both are same just flac is space saver




-> now let us look at all the steps including in ingestion service
1. artist upload master file
2. client splits into chunks
3. Ingestion service receives chunks
4. chunks are temporarily stored
5. upload completion is verified
6. chunks are merged into master file
7. master file is validated 
8. master file is stored permanently 
9. encoding is trigerred
10. track is mark as ingested

-> now we will start exploring the steps 
1. Artists start upload
-> usually master.wav or master.flac
-> high quality
-> in future only from this file all variants are made

2. client splits into chunks
-> this is done by frontend
-> why in chunks 
-> suppose 10mb song
-> some error then resend whole song
-> but suppose 10 chunks -> then network error -> only that chunk you need to resend
-> each chunks has upload_id,  chunk_number, raw_bytes

3. Ingestion service receiving chunks
-> POST /ingestion/upload-chunk


4. Temporary chunks storage
-> chunks are stored in temporary storage
-> marked as incomplete
-> stored in local storage, object storage like s3
-> why temporary -> upload may be fail or may be cancelled

5. Upload completion verification
-> ingestion service checks -> are all chunks present, are continuous, size matches the actual size
-> this will reject or proceed

6. chunk merge
-> read chunks in order
-> append bytes
-> create one file -> master.flac

7. master file validation 
=> validiates 
-> file types
-> duration readable
-> no corruption 
-> audio headers valid
-> reject or continue

8. permanent storage of master
-> master file is stored in long term storage
-> object storage bucked, not cdn, not public 
-> source of truth 

9. Encoding

10. mark ingestion complete
-> track.status= ingested


-> ingestion service is done till step 8
-> at step 9 encoding service come into picture
-> ingestion service triggers encoding via messaging queue

*** Ingestion Service ***



cheers
-> ui
-> backend
  -> ingestion
  -> encoding
  -> trackservice
  -> playback 









# Project planning 

1. Track Service
2. Ingestion Service 
3. Encoding Service
4. Playback Service 

-> Here it is sufficient 
-> this is stage1 of our project
-> then we will do other stages


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
├── proto/                 # gRPC contracts (shared)
├── libs/                  # shared infra libs (auth, logging)
├── deploy/                # k8s, terraform, helm
└── README.md

-> we will follow the hexagonal architecture 


track-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── 
│   │   └── 
│   ├── application/
│   │   ├── ports/
│   │   │   ├── in/
│   │   │   └── out/
│   │   └── service/
│   ├── adapters/
│   │   ├── in/
│   │   │   ├── http/
│   │   │   └── grpc/
│   │   └── out/
│   │       ├── postgres/
│   │       ├── redis/
│   │       └── kafka/
│   ├── config/
│   └── observability/
├── migrations/
├── go.mod
└── go.sum

