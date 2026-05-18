package out

type Validtor interface {
	ValidateAudio(meta AudioMetadata) error 
}