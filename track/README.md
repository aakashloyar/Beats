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
