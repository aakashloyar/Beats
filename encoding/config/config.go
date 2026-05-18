package config

var UploadLocalPathForEncoding string = ""

var SupportedCodecs map[string]bool = map[string]bool{
	"mp3":       true,
	"flac":      true,
	"pcm_s16le": true, // WAV
	"aac":       true,
	"opus":      true,
	"vorbis":    true,
}

// 3. Sample rate validation
var SupportedSampleRates map[int]bool = map[int]bool{
	44100: true,
	48000: true,
}


type EncodingProfile struct {
	FFmpegCodec string
	Bitrate     string
	Bandwidth   int
	///bandwidth is basically in actuall you need more internet speed than bitrate 
	///as the internet is used in other request also 
}

var Profiles = []EncodingProfile{
	{
		FFmpegCodec: "aac",
		Bitrate:     "96k",
		Bandwidth:   110000,
	},
	{
		FFmpegCodec: "aac",
		Bitrate:     "128k",
		Bandwidth:   140000,
	},
	{
		FFmpegCodec: "aac",
		Bitrate:     "320k",
		Bandwidth:   340000,
	},
}

var OutPutDirForTranscodedFiles = "streams"

var StoragePathForTranscodedFiles = "streams"