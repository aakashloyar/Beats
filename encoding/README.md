# Encoding Service 

# Flow
-> Encoding service takes upload information from kafka
-> Got the full video from storage
-> Encode upload using ffmpeg
-> update the database 

# steps:- 
1. Kafka Consumer hand over downloaded temporary file to encoding service 
# Note 
-> ffmpeg needs temporary storage for file 
-> bcz it is very computational heavy and retry
-> stable high-throughtput reads
-> so we will donwload the file first 
-> now here comes the problem 
-> if it will be going like this then the temporary storage will become heavy as queue grows
-> so we need solution for this 
-> we can do it without downloading also but not a good idea as dicussed above
-> so we can do something like vertical or horizontal scaling

2. Downloaded file is decoded
-> now the downloded file need to be validated before transcoding 
-> done by ffmpeg
->type ffprobeOutput struct {
	//array of stream because
	//suppose first 1 contains audio 
	//2nd one contain video or some other audio
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
# Note :-
-> now we are taking care about only simple audio file 
-> not taking care about other like commentry stream or video 
-> will see those in future

2. Vailating the decoded file
-> the meta data that we got after decoded file will get validated

3. Encoding 

* Note:-
-> way to store segment based streaming
-> Physical segment Files (Multiple files ) (Traditional HLS)
-> Byte range based storage (Single file  )
-> HLS (HTTP Live Streaming)

-> ok now let us talk about ways in segment based streaming
# Physical segment files (Normal HLS)
-> Multiple files 
-> playlist.m3u8 | chunk1.ts | chunk2.ts | .........
-> now let us look how our playlist.m3u8 looks like


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

-> let us look at seg.ts file 

```text
[MPEG-TS packet]
[MPEG-TS packet]
[MPEG-TS packet]
...

```

-> each packet looks like this

```text
PAT (program association table)
PMT (program map table)
AAC codec info
timestamps
audio frames
```
-> so seg.ts is tiny media independent file

-> visually
seg1.ts
 ├── headers
 ├── codec metadata
 ├── timestamps
 └── audio data

seg2.ts
 ├── headers
 ├── codec metadata
 ├── timestamps
 └── audio data


# Byte range
-> here donot create many files
-> pull all segments in one file
-> it looks like this way 
-> playlist.m3u8 | audio.ts
-> what audio.ts file looks like
audio.ts
|----segment1----|----segment2----|----segment3----|
0              120k             250k             380k
-> it is continuous file 
-> but logically it looks like 
But logically:
segment 1 = first 120k bytes
segment 2 = next 130k bytes
segment 3 = next 130k bytes

-> how playlist.m3u8 looks like 
```text
#EXTM3U

#EXTINF:6.0,
#EXT-X-BYTERANGE:120000@0
audio.ts

#EXTINF:6.0,
#EXT-X-BYTERANGE:130000@120000
audio.ts

#EXTINF:6.0,
#EXT-X-BYTERANGE:130000@250000
audio.ts
```

-> let us decode 

#EXT-X-BYTERANGE:130000@120000

Means:
130000 = length
120000 = starting byte offset

read bytes 120000 → 249999

-> steps
Step 1: Player requests playlist
-> req
GET /playlist.m3u8
-> res
playlist.m3u8 file 

Step 2: Player requests first segment

-> req
GET /audio.ts
Range: bytes=0-119999

-> res
HTTP 206 Partial Content
Content-Range: bytes 0-119999/380000

Step 3: Player requests next segment and continue

# CMAF/ fragmented mp4
-> now there also can be 2 ways multiple files/ single multifest file
-> let us talk about multiple files
-> now the file exist here
-> playlist.m3u8 | init.mp4 | chunk1.m4s | chunk2.m4s | chunk3.m4s | ...
-> here the difference is it will have all meta data stored in init.mp4 not in chunks 
-> so metadata is not repeated

-> let us look at playlist.m3u8 
```text
#EXTM3U
#EXT-X-MAP:URI="init.mp4"

#EXTINF:6.0,
chunk1.m4s

#EXTINF:6.0,
chunk2.m4s

#EXTINF:6.0,
chunk3.m4s
```

-> let us look at how init.mp4 looks like
-> it does not have media data only metadata needed to decode media
[ftyp][moov]
-> Box1 [ftyp]
format = mp4
brand = iso6
compatible = CMAF
-> tells i am an mp4 file 

-> Box2 [moov]
codec = AAC
sample rate = 44100
channels = 2
track id = 1
timescale = 48000
sample tables
duration info
-> decoder instructions


-> now let us look at chunk*.m4s 
[moof][mdat]
-> moof  (movie fragment metadata)
-> media data (actual encoded codec btyes)
-> 


Step1. 
-> req
GET /playlist.m3u8
-> res
playlist.m3u8

Step2. 
#EXT-X-MAP:URI="init.mp4"
-> req
GET /init.mp4
-> res
init.mp4

Step3. 
-> req
GET /chunk1.m4s
-> res
chunk1.m4s

Step4. Continue ask chunk*.m4s file

# Note:- 
-> why we are preferring CMAF over byte range 

1. Storage efficiency
-> because in byte range also each segment also contain metadata and every stuff
-> so it is not preferred for each segment metadata 

2. Format quality 
-> .ts(mpeg transport stream) was older made for tv broadcasting
-> .mp4 much modern

3. Protocol compatibility
-> Byte-range HLS
-> HLS
-> CMAF
-> HLS | DASH

4. Latency 
-> CMAF wins here
Designed to support:
Low-Latency HLS
Low-Latency DASH
Smaller partial chunks.

5. CDN caching 
-> nearly similar but CMAF win

* Note:- 
-> difference between HLS and DASH

-> HLS -> http Live streaming
-> owned by apple
-> player reads .m3u8

-> DASH Dynamic Adaptive streaming over HTTP 
-> xml based data 
-> used manifest.mpd
-> not owned by single company

-> HLS flow
Player
  ↓
GET playlist.m3u8
  ↓
Read chunk names
  ↓
GET chunk1.m4s
GET chunk2.m4s

-> DASH flow
Player
  ↓
GET manifest.mpd
  ↓
Read chunk info
  ↓
GET chunk1.m4s
GET chunk2.m4s

-> they do the same job 
-> both tell player 
-> the available quality 
-> segment location
-> duration
-> timing
-> switching rules

-> hls is preferred on apple devices
-> rest devices dash is preferred

# Note :- 
-> but for now as we only have playlist.m3u8
-> so we are only supporting hls not dash
-> in future we will see

# CMAF (Common media application format) single file
-> instead of multiple files it will have 
-> playlist.m3u8 | audio.mp4
-> audio.mp4
audio.mp4
|-- init section --|-- fragment1 --|-- fragment2 --|-- fragment3 --|
-> internally
[ftyp][moov][moof][mdat][moof][mdat][moof][mdat]
[ftyp]
-> i am mp4
[moov]
-> global metadata
[moof]
-> fragment metadata
Chunk timing | Sample count | Offsets
audio.mp4
   [ftyp][moov]      ← init part
   [moof][mdat]      ← segment 1
   [moof][mdat]      ← segment 2
   [moof][mdat]      ← segment 3

-> playlist.m3u8 
```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:100000@5000
audio.mp4

#EXTINF:6.0,
#EXT-X-BYTERANGE:98000@105000
audio.mp4

#EXTINF:6.0,
#EXT-X-BYTERANGE:102000@203000
audio.mp4
```   
-> here it means 
-> initialisation data lives in byte 0-4999
-> then byte start

-> now let us look a steps

Step 1: Player requests playlist
-> req
GET /playlist.m3u8
-> res
playlist.m3u8 file

Step 2: Player reads init byte-range
BYTERANGE="5000@0"
-> req
GET /audio.mp4
Range: bytes=0-4999
-> res
[ftyp][moov]
-> player now got codec and track structure

Step 3: Player requests segment 1
-> req
GET /audio.mp4
Range: bytes=5000-104999
-> res
[moof][mdat]
-> player combines in memory
(init bytes) + (segment bytes)
-> player start

Step 4: Next segments


# Note :- how switching happens
-> Master playlist
-> 1 for all quality
-> media playlist
-> each quality has it's own playlist file

1. Master playlist file 
```text
#EXTM3U

#EXT-X-STREAM-INF:BANDWIDTH=64000
64k/playlist.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=128000
128k/playlist.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=320000
320k/playlist.m3u8
```
-> this tell player
-> available version 
-> player download this first

2. Media playlist file

-> For 64k:

```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:40000@5000
audio.mp4

#EXTINF:6.0,
#EXT-X-BYTERANGE:42000@45000
audio.mp4
```

-> For 128k:

```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:80000@5000
audio.mp4

#EXTINF:6.0,
#EXT-X-BYTERANGE:82000@85000
audio.mp4
```
-> For 320k:

```text
#EXTM3U
#EXT-X-MAP:URI="audio.mp4",BYTERANGE="5000@0"

#EXTINF:6.0,
#EXT-X-BYTERANGE:200000@5000
audio.mp4
```

# Note:-
->so the step 1 is Get /masterplaylist
-> normal mp4 store metadata at last bad for streaming
-> mp4 streamable store that first

outDir/
   master.m3u8
   96k/
      playlist.m3u8
      audio.mp4
   128k/
      playlist.m3u8
      audio.mp4
   320k/
      playlist.m3u8
      audio.mp4